package managers

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/engine"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type StatsUpload struct {
	CollegeTeamStats   []structs.CollegeTeamGameStats
	CollegePlayerStats []structs.CollegePlayerGameStats
	ProTeamStats       []structs.ProfessionalTeamGameStats
	ProPlayerStats     []structs.ProfessionalPlayerGameStats
	CollegePlayByPlay  []structs.CollegePlayByPlay
	ProPlayByPlay      []structs.ProPlayByPlay
}

func RunGames() {
	// Get GameDTOs
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	gameDay := ts.GetGameDay()
	collegeGames := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay)
	proGames := []structs.ProfessionalGame{}
	collegeGameMap := MakeCollegeGameMap(collegeGames)
	proGameMap := MakeProGameMap(proGames)
	// proGames := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay)
	gameDTOs := PrepareGames(collegeGames, proGames)
	// RUN THE GAMES!
	results := engine.RunGames(gameDTOs)
	// collegeTeamMap := GetCollegeTeamMap()
	// proTeamMap := GetProTeamMap()
	// collegePlayerMap := GetCollegePlayersMap()
	// proPlayersMap := GetProPlayersMap()
	upload := NewStatsUpload()
	for _, r := range results {
		// Iterate through all lines, players, accumulate stats to upload
		// WriteBoxScoreFile(r, "test_results/test_twelve/box_score/"+r.HomeTeam+"_vs_"+r.AwayTeam+".csv")

		// Iterate through Play By Plays and record them to a CSV
		// if r.IsCollegeGame {
		// 	WritePlayByPlayCSVFile(pbps, "test_results/test_twelve/play_by_play/"+r.HomeTeam+"_vs_"+r.AwayTeam+".csv", collegePlayerMap, collegeTeamMap)
		// } else {
		// 	WriteProPlayByPlayCSVFile(pbps, "test_results/test_twelve/play_by_play/"+r.HomeTeam+"_vs_"+r.AwayTeam+".csv", proPlayersMap, proTeamMap)
		// }
		upload.Collect(r, ts.SeasonID)

		if r.IsCollegeGame {
			game := collegeGameMap[r.GameID]
			game.UpdateScore(uint(r.HomeTeamScore), uint(r.AwayTeamScore), uint(r.HomeTeamShootoutScore), uint(r.AwayTeamShootoutScore), r.IsOvertime, r.IsOvertimeShootout)
			repository.SaveCollegeGameRecord(game, db)
		} else {
			game := proGameMap[r.GameID]
			game.UpdateScore(uint(r.HomeTeamScore), uint(r.AwayTeamScore), uint(r.HomeTeamShootoutScore), uint(r.AwayTeamShootoutScore), r.IsOvertime, r.IsOvertimeShootout)
			repository.SaveProfessionalGameRecord(game, db)
		}
	}
	upload.Flush(db)
}

func NewStatsUpload() *StatsUpload {
	return &StatsUpload{}
}

func (u *StatsUpload) Collect(state engine.GameState, seasonID uint) {
	// Team stats
	u.collectTeamStats(state, seasonID)

	// Player stats for both teams
	u.collectPlayerStats(state.HomeStrategy, state.WeekID, state.GameID, state.IsCollegeGame)
	u.collectPlayerStats(state.AwayStrategy, state.WeekID, state.GameID, state.IsCollegeGame)

	// Play-by-play
	u.collectPbP(state.Collector.PlayByPlays, state.IsCollegeGame)
}

func (u *StatsUpload) collectTeamStats(state engine.GameState, seasonID uint) {
	if state.IsCollegeGame {
		u.CollegeTeamStats = append(u.CollegeTeamStats,
			makeCollegeTeamStatsObject(state.WeekID, state.GameID, seasonID, state.HomeTeamStats),
			makeCollegeTeamStatsObject(state.WeekID, state.GameID, seasonID, state.AwayTeamStats),
		)
	} else {
		u.ProTeamStats = append(u.ProTeamStats,
			makeProTeamStatsObject(state.WeekID, state.GameID, seasonID, state.HomeTeamStats),
			makeProTeamStatsObject(state.WeekID, state.GameID, seasonID, state.AwayTeamStats),
		)
	}
}

func (u *StatsUpload) collectPlayerStats(pl engine.GamePlaybook, week, gameID uint, isCollege bool) {
	types := [][]engine.LineStrategy{pl.Forwards, pl.Defenders, pl.Goalies}
	for _, group := range types {
		for _, line := range group {
			for _, p := range line.Players {
				if isCollege {
					u.CollegePlayerStats = append(u.CollegePlayerStats,
						makeCollegePlayerStatsObject(week, gameID, p.Stats),
					)
				} else {
					u.ProPlayerStats = append(u.ProPlayerStats,
						makeProPlayerStatsObject(week, gameID, p.Stats),
					)
				}
			}
		}
	}
}

func (u *StatsUpload) collectPbP(pbps []structs.PbP, isCollege bool) {
	if isCollege {
		for _, pbp := range pbps {
			u.CollegePlayByPlay = append(u.CollegePlayByPlay, structs.CollegePlayByPlay{PbP: pbp})
		}
	} else {
		for _, pbp := range pbps {
			u.ProPlayByPlay = append(u.ProPlayByPlay, structs.ProPlayByPlay{PbP: pbp})
		}
	}
}

func (u *StatsUpload) Flush(db *gorm.DB) error {
	const batchSize = 200
	const bigBatchSize = 500
	if err := repository.CreateCHLPlayByPlayRecordBatch(db, u.CollegePlayByPlay, bigBatchSize); err != nil {
		return err
	}
	if err := repository.CreatePHLPlayByPlayRecordBatch(db, u.ProPlayByPlay, bigBatchSize); err != nil {
		return err
	}
	if err := repository.CreateCHLPlayerGameStatsRecordBatch(db, u.CollegePlayerStats, batchSize); err != nil {
		return err
	}
	if err := repository.CreatePHLPlayerGameStatsRecordBatch(db, u.ProPlayerStats, batchSize); err != nil {
		return err
	}
	if err := repository.CreateCHLTeamGameStatsRecordBatch(db, u.CollegeTeamStats, batchSize); err != nil {
		return err
	}
	if err := repository.CreatePHLTeamGameStatsRecordBatch(db, u.ProTeamStats, batchSize); err != nil {
		return err
	}
	return nil
}

func PrepareGames(collegeGames []structs.CollegeGame, proGames []structs.ProfessionalGame) []structs.GameDTO {
	fmt.Println("Loading Games...")

	// Wait Groups
	var collegeGamesWg sync.WaitGroup
	// Mutex Lock
	var mutex sync.Mutex

	// College Only
	// collegeTeamMap := GetCollegeTeamMap()
	collegeTeamRosterMap := GetAllCollegePlayersMapByTeam()
	collegeLineupMap := GetCollegeLineupsMap()
	collegeShootoutLineupMap := GetCollegeShootoutLineups()
	arenaMap := GetArenaMap()
	// collegeGames := GetCollegeGamesForTesting(collegeTeamMap)
	collegeGamesWg.Add(len(collegeGames))
	gameDTOList := make([]structs.GameDTO, 0, len(collegeGames))
	sem := make(chan struct{}, 20)

	for _, c := range collegeGames {
		sem <- struct{}{}
		localC := c
		go func(c structs.CollegeGame) {
			defer func() { <-sem }()
			defer collegeGamesWg.Done()
			if c.GameComplete {
				return
			}
			mutex.Lock()
			htr := collegeTeamRosterMap[c.HomeTeamID]
			atr := collegeTeamRosterMap[c.AwayTeamID]
			htl := collegeLineupMap[c.HomeTeamID]
			atl := collegeLineupMap[c.AwayTeamID]
			htsl := collegeShootoutLineupMap[c.HomeTeamID]
			atsl := collegeShootoutLineupMap[c.AwayTeamID]
			hp := getCollegePlaybookDTO(htl, htr, htsl)
			ap := getCollegePlaybookDTO(atl, atr, atsl)
			capacity := 0

			arena := arenaMap[c.ArenaID]
			if arena.ID == 0 {
				capacity = 6000
			} else {
				capacity = int(arena.Capacity)
			}
			mutex.Unlock()

			match := structs.GameDTO{
				GameID:        c.ID,
				GameInfo:      c.BaseGame,
				HomeStrategy:  hp,
				AwayStrategy:  ap,
				IsCollegeGame: true,
				Attendance:    uint32(capacity),
			}

			mutex.Lock()
			gameDTOList = append(gameDTOList, match)
			mutex.Unlock()
		}(localC)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	var proGamesWg sync.WaitGroup
	// professionalTeamMap := GetProTeamMap()
	// proGames := GetProGamesForTesting(professionalTeamMap)
	// proGames := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay)
	proTeamRosterMap := GetAllProPlayersMapByTeam()
	proLineupMap := GetProLineupsMap()
	proShootoutLineupMap := GetProShootoutLineups()
	proGamesWg.Add(len(proGames))
	proSem := make(chan struct{}, 20)
	for _, g := range proGames {
		proSem <- struct{}{}
		localC := g
		go func(g structs.ProfessionalGame) {
			defer func() { <-proSem }()
			defer proGamesWg.Done()
			if g.GameComplete {
				return
			}
			mutex.Lock()
			htr := proTeamRosterMap[g.HomeTeamID]
			atr := proTeamRosterMap[g.AwayTeamID]
			htl := proLineupMap[g.HomeTeamID]
			atl := proLineupMap[g.AwayTeamID]
			htsl := proShootoutLineupMap[g.HomeTeamID]
			atsl := proShootoutLineupMap[g.AwayTeamID]
			hp := getProfessionalPlaybookDTO(htl, htr, htsl)
			ap := getProfessionalPlaybookDTO(atl, atr, atsl)
			capacity := 0

			arena := arenaMap[g.ArenaID]
			if arena.ID == 0 {
				capacity = 6000
			} else {
				capacity = int(arena.Capacity)
			}
			mutex.Unlock()

			match := structs.GameDTO{
				GameID:        g.ID,
				GameInfo:      g.BaseGame,
				HomeStrategy:  hp,
				AwayStrategy:  ap,
				IsCollegeGame: false,
				Attendance:    uint32(capacity),
			}

			mutex.Lock()
			gameDTOList = append(gameDTOList, match)
			mutex.Unlock()
		}(localC)
	}
	collegeGamesWg.Wait()
	proGamesWg.Wait()
	for i := 0; i < cap(proSem); i++ {
		proSem <- struct{}{}
	}
	return gameDTOList
}

func GeneratePreseasonGames() {
	db := dbprovider.GetInstance().GetDB()

	collegeTeamMap := GetCollegeTeamMap()
	proTeamMap := GetProTeamMap()

	collegeGames := GetCollegeGamesForPreseason(collegeTeamMap)
	proGames := GetProGamesForPreseason(proTeamMap)

	repository.CreateCHLGamesRecordsBatch(db, collegeGames, 20)
	repository.CreatePHLGamesRecordsBatch(db, proGames, 20)
}

func GetCollegeGamesForPreseason(teamMap map[uint]structs.CollegeTeam) []structs.CollegeGame {
	games := []structs.CollegeGame{}
	teamIDs := make([]uint, 0, len(teamMap))
	playedGameReference := make(map[uint]map[uint]bool)
	for id := range teamMap {
		teamIDs = append(teamIDs, id)
		playedGameReference[id] = make(map[uint]bool)
	}
	gameDay := "A"

	for round := 1; round <= 3; round++ {
		var pairings [][2]uint
		const maxTries = 500
		for tries := 0; tries < maxTries; tries++ {
			// shuffle
			rand.Shuffle(len(teamIDs), func(i, j int) {
				teamIDs[i], teamIDs[j] = teamIDs[j], teamIDs[i]
			})

			// attempt to pair
			var err error
			pairings, err = createCollegeGamePairings(teamIDs, teamMap, playedGameReference)
			if err == nil {
				// success!
				break
			}
		}
		if len(pairings)*2 != len(teamIDs) {
			log.Fatalf("couldn’t find a full pairing on round %d after %d tries", round, maxTries)
		}

		// generate the Game objects
		for _, pair := range pairings {
			game := generateCollegeGame(
				1,    // leagueID
				2501, // seasonID
				1,    // weekID
				pair[0], pair[1],
				gameDay, teamMap, true,
			)
			games = append(games, game)
		}

		// rotate gameDay A→B→C
		switch gameDay {
		case "A":
			gameDay = "B"
		case "B":
			gameDay = "C"
		}
	}

	return games
}

func GetProGamesForPreseason(teamMap map[uint]structs.ProfessionalTeam) []structs.ProfessionalGame {
	playedGameReference := make(map[uint]map[uint]bool)
	games := []structs.ProfessionalGame{}
	gameDay := "A"
	teamIDs := make([]uint, 0, len(teamMap))
	for id := range teamMap {
		teamIDs = append(teamIDs, id)
		playedGameReference[id] = make(map[uint]bool)
	}

	for round := 1; round <= 3; round++ {
		var pairings [][2]uint
		const maxTries = 500
		for tries := 0; tries < maxTries; tries++ {
			// shuffle
			rand.Shuffle(len(teamIDs), func(i, j int) {
				teamIDs[i], teamIDs[j] = teamIDs[j], teamIDs[i]
			})

			// attempt to pair
			var err error
			pairings, err = createProGamePairings(teamIDs, teamMap, playedGameReference)
			if err == nil {
				// success!
				break
			}
		}
		if len(pairings)*2 != len(teamIDs) {
			log.Fatalf("couldn’t find a full pairing on round %d after %d tries", round, maxTries)
		}

		// generate the Game objects
		for _, pair := range pairings {
			game := generateProfessionalGame(
				1,    // leagueID
				2501, // seasonID
				1,    // weekID
				pair[0], pair[1],
				gameDay, teamMap, true,
			)
			games = append(games, game)
		}

		// rotate gameDay A→B→C
		switch gameDay {
		case "A":
			gameDay = "B"
		case "B":
			gameDay = "C"
		}
	}
	return games
}

func createCollegeGamePairings(
	teamIDs []uint,
	teamMap map[uint]structs.CollegeTeam,
	playedGameReference map[uint]map[uint]bool,
) ([][2]uint, error) {
	n := len(teamIDs)
	if n%2 != 0 {
		return nil, fmt.Errorf("need an even number of teams; got %d", n)
	}

	paired := make(map[uint]bool)
	var pairings [][2]uint

	// Walk the shuffled slice of teamIDs
	for _, t1 := range teamIDs {
		if paired[t1] {
			continue
		}
		c1 := teamMap[t1].ConferenceID

		// Try to find a haven for t1
		found := false
		for _, t2 := range teamIDs {
			if t1 == t2 || paired[t2] {
				continue
			}
			c2 := teamMap[t2].ConferenceID

			// same‐conference? only allow if both are independent (7)
			if c1 != 7 && c1 == c2 {
				continue
			}
			// already played?
			if playedGameReference[t1][t2] {
				continue
			}

			// commit the pairing
			pairings = append(pairings, [2]uint{t1, t2})
			paired[t1], paired[t2] = true, true
			playedGameReference[t1][t2] = true
			playedGameReference[t2][t1] = true
			found = true
			break
		}

		// if we can’t find a partner for t1, bail out
		if !found {
			break
		}
	}

	// did we cover _all_ teams?
	if len(pairings)*2 != n {
		return nil, fmt.Errorf(
			"incomplete pairing: only paired %d of %d teams",
			len(pairings)*2, n,
		)
	}

	return pairings, nil
}

func createProGamePairings(teamIDs []uint, teamMap map[uint]structs.ProfessionalTeam, playedGameReference map[uint]map[uint]bool) ([][2]uint, error) {
	n := len(teamIDs)
	if n%2 != 0 {
		return nil, fmt.Errorf("need an even number of teams; got %d", n)
	}

	paired := make(map[uint]bool)
	var pairings [][2]uint

	// Walk the shuffled slice of teamIDs
	for _, t1 := range teamIDs {
		if paired[t1] {
			continue
		}
		c1 := teamMap[t1].DivisionID

		// Try to find a haven for t1
		found := false
		for _, t2 := range teamIDs {
			if t1 == t2 || paired[t2] {
				continue
			}
			c2 := teamMap[t2].DivisionID

			// same‐conference? only allow if both are independent (7)
			if c1 != 7 && c1 == c2 {
				continue
			}
			// already played?
			if playedGameReference[t1][t2] {
				continue
			}

			// commit the pairing
			pairings = append(pairings, [2]uint{t1, t2})
			paired[t1], paired[t2] = true, true
			playedGameReference[t1][t2] = true
			playedGameReference[t2][t1] = true
			found = true
			break
		}

		// if we can’t find a partner for t1, bail out
		if !found {
			break
		}
	}

	// did we cover _all_ teams?
	if len(pairings)*2 != n {
		return nil, fmt.Errorf(
			"incomplete pairing: only paired %d of %d teams",
			len(pairings)*2, n,
		)
	}

	return pairings, nil
}

func GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay string) []structs.CollegeGame {
	return repository.FindCollegeGamesByCurrentMatchup(weekID, seasonID, gameDay)
}

func GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay string) []structs.ProfessionalGame {
	return repository.FindProfessionalGamesByCurrentMatchup(weekID, seasonID, gameDay)
}

func GetCollegeGamesByTeamIDAndSeasonID(teamID, seasonID string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGames(seasonID, teamID, isPreseason)
}

func GetProfessionalGamesByTeamIDAndSeasonID(teamID, seasonID string, isPreseason bool) []structs.ProfessionalGame {
	return repository.FindProfessionalGames(seasonID, teamID, isPreseason)
}

func GetCollegeGamesBySeasonID(seasonID string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGames(seasonID, "", isPreseason)
}

func GetProfessionalGamesBySeasonID(seasonID string, isPreseason bool) []structs.ProfessionalGame {
	return repository.FindProfessionalGames(seasonID, "", isPreseason)
}

func GetCollegeGameByID(id string) structs.CollegeGame {
	return repository.FindCollegeGameRecord(id)
}

func GetProfessionalGameByID(id string) structs.ProfessionalGame {
	return repository.FindProfessionalGameRecord(id)
}

func GetArenaMap() map[uint]structs.Arena {
	arenas := repository.FindAllArenas()
	return MakeArenaMap(arenas)
}

func getCollegePlaybookDTO(lineups []structs.CollegeLineup, roster []structs.CollegePlayer, shootoutLineup structs.CollegeShootoutLineup) structs.PlayBookDTO {
	forwards, defenders, goalies := getCollegeForwardDefenderGoalieLineups(lineups)
	return structs.PlayBookDTO{
		Forwards:       forwards,
		Defenders:      defenders,
		Goalies:        goalies,
		CollegeRoster:  roster,
		ShootoutLineup: shootoutLineup.ShootoutPlayerIDs,
	}
}

func getProfessionalPlaybookDTO(lineups []structs.ProfessionalLineup, roster []structs.ProfessionalPlayer, shootoutLineup structs.ProfessionalShootoutLineup) structs.PlayBookDTO {
	forwards, defenders, goalies := getProfessionalForwardDefenderGoalieLineups(lineups)
	return structs.PlayBookDTO{
		Forwards:           forwards,
		Defenders:          defenders,
		Goalies:            goalies,
		ProfessionalRoster: roster,
		ShootoutLineup:     shootoutLineup.ShootoutPlayerIDs,
	}
}

// func getBaseRoster(roster []structs.CollegePlayer) []structs.BasePlayer {
// 	basePlayers := []structs.BasePlayer{}

// 	for _, p := range roster {
// 		basePlayers = append(basePlayers, p.BasePlayer)
// 	}
// 	return basePlayers
// }

func getCollegeForwardDefenderGoalieLineups(lineups []structs.CollegeLineup) ([]structs.BaseLineup, []structs.BaseLineup, []structs.BaseLineup) {
	forwards := []structs.BaseLineup{}
	defenders := []structs.BaseLineup{}
	goalies := []structs.BaseLineup{}
	for _, l := range lineups {
		if l.LineType == 1 {
			forwards = append(forwards, l.BaseLineup)
		} else if l.LineType == 2 {
			defenders = append(defenders, l.BaseLineup)
		} else {
			goalies = append(goalies, l.BaseLineup)
		}
	}
	return forwards, defenders, goalies
}

func getProfessionalForwardDefenderGoalieLineups(lineups []structs.ProfessionalLineup) ([]structs.BaseLineup, []structs.BaseLineup, []structs.BaseLineup) {
	forwards := []structs.BaseLineup{}
	defenders := []structs.BaseLineup{}
	goalies := []structs.BaseLineup{}
	for _, l := range lineups {
		if l.LineType == 1 {
			forwards = append(forwards, l.BaseLineup)
		} else if l.LineType == 2 {
			defenders = append(defenders, l.BaseLineup)
		} else {
			goalies = append(goalies, l.BaseLineup)
		}
	}
	return forwards, defenders, goalies
}

func generateCollegeGame(seasonID, weekID, week, hid, aid uint, gameDay string, teamMap map[uint]structs.CollegeTeam, isPreseason bool) structs.CollegeGame {
	return structs.CollegeGame{
		BaseGame: structs.BaseGame{
			WeekID:      weekID,
			Week:        int(week),
			GameDay:     gameDay,
			SeasonID:    seasonID,
			HomeTeamID:  hid,
			HomeTeam:    teamMap[hid].TeamName,
			AwayTeamID:  aid,
			AwayTeam:    teamMap[aid].TeamName,
			ArenaID:     uint(teamMap[hid].ArenaID),
			IsPreseason: isPreseason,
		},
	}
}

func generateProfessionalGame(seasonID, weekID, week, hid, aid uint, gameDay string, teamMap map[uint]structs.ProfessionalTeam, isPreseason bool) structs.ProfessionalGame {
	return structs.ProfessionalGame{
		BaseGame: structs.BaseGame{
			WeekID:      weekID,
			Week:        int(week),
			GameDay:     gameDay,
			SeasonID:    seasonID,
			HomeTeamID:  hid,
			HomeTeam:    teamMap[hid].Abbreviation,
			AwayTeamID:  aid,
			AwayTeam:    teamMap[aid].Abbreviation,
			ArenaID:     uint(teamMap[hid].ArenaID),
			IsPreseason: isPreseason,
		},
	}
}

func GetPlayoffSeriesBySeriesID(seriesID string) structs.PlayoffSeries {
	return repository.FindPlayoffSeriesByID(seriesID)
}
