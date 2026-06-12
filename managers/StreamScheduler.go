package managers

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	fbsvc "github.com/CalebRose/SimHockey/firebase"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

const maxStreamSlots = 8

// PendingGame is a lightweight descriptor for a game waiting to enter a slot.
type PendingGame struct {
	GameID       uint
	HomeTeamID   uint
	AwayTeamID   uint
	HomeTeam     string
	AwayTeam     string
	IsUserGame   bool // true if either team is user-coached / user-owned
}

// GameStream represents one active streaming slot.
type GameStream struct {
	GameID    uint
	StartTime time.Time
	EndTime   time.Time
	League    string
}

// StreamScheduler manages up to maxStreamSlots concurrent game streams and an
// ordered queue of pending games for a single league.
type StreamScheduler struct {
	mu          sync.Mutex
	ActiveSlots [maxStreamSlots]*GameStream
	Queue       []PendingGame
	League      string // "chl" or "phl"
	isCollege   bool
}

// computeStreamTimes sums SecondsConsumed across a play-by-play slice and
// returns a start time of now, the corresponding end time, and the total seconds.
func computeStreamTimes(totalSecs int) (start, end time.Time) {
	start = time.Now().UTC()
	end = start.Add(time.Duration(totalSecs) * time.Second)
	return
}

// loadTotalSeconds queries the PbP table for gameID and sums SecondsConsumed.
// Returns 0 if no records are found.
func loadTotalSeconds(gameID uint, isCollege bool) int {
	gameIDStr := strconv.FormatUint(uint64(gameID), 10)
	total := 0
	if isCollege {
		plays := repository.FindCHLPlayByPlaysRecordsByGameID(gameIDStr)
		for _, p := range plays {
			total += int(p.SecondsConsumed)
		}
	} else {
		plays := repository.FindPHLPlayByPlaysRecordsByGameID(gameIDStr)
		for _, p := range plays {
			total += int(p.SecondsConsumed)
		}
	}
	return total
}

// loadTotalPlays returns the number of play-by-play records for a game.
func loadTotalPlays(gameID uint, isCollege bool) int {
	gameIDStr := strconv.FormatUint(uint64(gameID), 10)
	if isCollege {
		return len(repository.FindCHLPlayByPlaysRecordsByGameID(gameIDStr))
	}
	return len(repository.FindPHLPlayByPlaysRecordsByGameID(gameIDStr))
}

// dequeue pops the first item from the queue.
func (s *StreamScheduler) dequeue() (PendingGame, bool) {
	if len(s.Queue) == 0 {
		return PendingGame{}, false
	}
	next := s.Queue[0]
	s.Queue = s.Queue[1:]
	return next, true
}

// InitQueue loads all complete, unrevealed games for the current matchup and
// sorts them with user-coached/owned games first, then by GameID.
func (s *StreamScheduler) InitQueue(weekID, seasonID, gameDay string, isPreseason bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userGames, aiGames []PendingGame

	if s.isCollege {
		games := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
		teamMap := GetCollegeTeamMap()
		for _, g := range games {
			if !g.GameComplete || g.IsRevealed {
				continue
			}
			homeTeam := teamMap[g.HomeTeamID]
			awayTeam := teamMap[g.AwayTeamID]
			pg := PendingGame{
				GameID:     g.ID,
				HomeTeamID: g.HomeTeamID,
				AwayTeamID: g.AwayTeamID,
				HomeTeam:   homeTeam.Abbreviation,
				AwayTeam:   awayTeam.Abbreviation,
				IsUserGame: homeTeam.IsUserCoached || awayTeam.IsUserCoached,
			}
			if pg.IsUserGame {
				userGames = append(userGames, pg)
			} else {
				aiGames = append(aiGames, pg)
			}
		}
	} else {
		games := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
		teamMap := GetProTeamMap()
		for _, g := range games {
			if !g.GameComplete || g.IsRevealed {
				continue
			}
			homeTeam := teamMap[g.HomeTeamID]
			awayTeam := teamMap[g.AwayTeamID]
			isUser := homeTeam.Owner != "" || awayTeam.Owner != "" ||
				homeTeam.GM != "" || awayTeam.GM != ""
			pg := PendingGame{
				GameID:     g.ID,
				HomeTeamID: g.HomeTeamID,
				AwayTeamID: g.AwayTeamID,
				HomeTeam:   homeTeam.Abbreviation,
				AwayTeam:   awayTeam.Abbreviation,
				IsUserGame: isUser,
			}
			if pg.IsUserGame {
				userGames = append(userGames, pg)
			} else {
				aiGames = append(aiGames, pg)
			}
		}
	}

	// User-coached games fill the front of the queue; AI games follow.
	s.Queue = append(userGames, aiGames...)
	log.Printf("StreamScheduler(%s): queued %d games (%d user, %d AI)",
		s.League, len(s.Queue), len(userGames), len(aiGames))
}

// Tick is called by the cron on every interval.  It marks completed game slots
// as revealed in Firebase, then promotes pending games into freed slots.
func (s *StreamScheduler) Tick(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()

	// 1. Mark completed game slots revealed and free them.
	for i, slot := range s.ActiveSlots {
		if slot == nil || now.Before(slot.EndTime) {
			continue
		}
		// Slot has elapsed — mark revealed and clear.
		go func(gameID uint, league string) {
			if err := fbsvc.SetGameRevealed(ctx, gameID, league); err != nil {
				log.Printf("StreamScheduler: SetGameRevealed(gameID=%d, league=%s): %v", gameID, league, err)
			}
		}(slot.GameID, slot.League)
		s.ActiveSlots[i] = nil
	}

	// 2. Fill freed slots from the queue.
	for i, slot := range s.ActiveSlots {
		if slot != nil || len(s.Queue) == 0 {
			continue
		}
		next, ok := s.dequeue()
		if !ok {
			break
		}

		totalSecs := loadTotalSeconds(next.GameID, s.isCollege)
		if totalSecs == 0 {
			// No plays in DB yet — skip this game and try the next.
			log.Printf("StreamScheduler(%s): skipping game %d — no PbP records found", s.League, next.GameID)
			continue
		}
		totalPlays := loadTotalPlays(next.GameID, s.isCollege)

		start, end := computeStreamTimes(totalSecs)
		record := fbsvc.LiveGameRecord{
			GameID:          next.GameID,
			HomeTeamID:      next.HomeTeamID,
			AwayTeamID:      next.AwayTeamID,
			HomeTeam:        next.HomeTeam,
			AwayTeam:        next.AwayTeam,
			League:          s.League,
			StreamStartTime: start,
			StreamEndTime:   end,
			TotalPlays:      totalPlays,
			IsRevealed:      false,
		}
		go func(rec fbsvc.LiveGameRecord, league string) {
			if err := fbsvc.UploadLiveGame(ctx, rec, league); err != nil {
				log.Printf("StreamScheduler: UploadLiveGame(gameID=%d, league=%s): %v", rec.GameID, league, err)
			}
		}(record, s.League)

		s.ActiveSlots[i] = &GameStream{
			GameID:    next.GameID,
			StartTime: start,
			EndTime:   end,
			League:    s.League,
		}
		log.Printf("StreamScheduler(%s): activated game %d (ends at %s)", s.League, next.GameID, end.Format(time.RFC3339))
	}
}

// IsIdle returns true when all slots are empty and the queue is exhausted.
func (s *StreamScheduler) IsIdle() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Queue) > 0 {
		return false
	}
	for _, slot := range s.ActiveSlots {
		if slot != nil {
			return false
		}
	}
	return true
}

// StartCHLLiveStreamingCron initialises a CHL StreamScheduler, fills its queue,
// and runs Tick on a 5-second interval until all games are revealed.
func StartCHLLiveStreamingCron() {
	ts := GetTimestamp()
	scheduler := &StreamScheduler{League: "chl", isCollege: true}
	scheduler.InitQueue(
		strconv.Itoa(int(ts.WeekID)),
		strconv.Itoa(int(ts.SeasonID)),
		ts.GetGameDay(),
		ts.IsPreseason,
	)
	if len(scheduler.Queue) == 0 {
		log.Println("StreamScheduler(chl): no games to stream")
		return
	}

	ctx := context.Background()
	scheduler.Tick(ctx) // fill initial slots immediately

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			scheduler.Tick(ctx)
			if scheduler.IsIdle() {
				log.Println("StreamScheduler(chl): all games complete, stopping")
				return
			}
		}
	}()
}

// StartPHLLiveStreamingCron initialises a PHL StreamScheduler and runs it.
func StartPHLLiveStreamingCron() {
	ts := GetTimestamp()
	scheduler := &StreamScheduler{League: "phl", isCollege: false}
	scheduler.InitQueue(
		strconv.Itoa(int(ts.WeekID)),
		strconv.Itoa(int(ts.SeasonID)),
		ts.GetGameDay(),
		ts.IsPreseason,
	)
	if len(scheduler.Queue) == 0 {
		log.Println("StreamScheduler(phl): no games to stream")
		return
	}

	ctx := context.Background()
	scheduler.Tick(ctx)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			scheduler.Tick(ctx)
			if scheduler.IsIdle() {
				log.Println("StreamScheduler(phl): all games complete, stopping")
				return
			}
		}
	}()
}

// GetCHLLivePlays returns the ordered play-by-play slice for a single CHL game
// as a PlayByPlayResponse slice, suitable for the live-plays API endpoint.
// No Firebase reads occur.
func GetCHLLivePlays(gameID string) []structs.PlayByPlayResponse {
	plays := repository.FindCHLPlayByPlaysRecordsByGameID(gameID)
	if len(plays) == 0 {
		return []structs.PlayByPlayResponse{}
	}
	game := GetCollegeGameByID(gameID)
	teamMap := GetCollegeTeamMap()
	playerMap := GetCollegePlayersMap()
	return GenerateCHLPlayByPlayResponse(plays, teamMap, playerMap, true, game.HomeTeamID, game.AwayTeamID)
}

// GetPHLLivePlays returns the ordered play-by-play slice for a single PHL game
// as a PlayByPlayResponse slice, suitable for the live-plays API endpoint.
// No Firebase reads occur.
func GetPHLLivePlays(gameID string) []structs.PlayByPlayResponse {
	plays := repository.FindPHLPlayByPlaysRecordsByGameID(gameID)
	if len(plays) == 0 {
		return []structs.PlayByPlayResponse{}
	}
	game := GetProfessionalGameByID(gameID)
	teamMap := GetProTeamMap()
	playerMap := GetProPlayersMap()
	return GeneratePHLPlayByPlayResponse(plays, teamMap, playerMap, true, game.HomeTeamID, game.AwayTeamID)
}
