package managers

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"

	util "github.com/CalebRose/SimHockey/_util"
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
	if ts.IsTesting {
		// Generate and run random college games for testing without database operations
		generateAndRunTestGames(ts, db)
		return
	}
	collegeGames := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)
	proGames := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)

	collegeStandingsMap := GetCollegeStandingsMap(seasonID)
	proStandingsMap := GetProStandingsMap(seasonID)
	gameDTOs := PrepareGames(collegeGames, proGames, collegeStandingsMap, proStandingsMap)
	// RUN THE GAMES!
	results := engine.RunGames(gameDTOs)

	// for _, r := range results {
	// 	homestats := r.HomeTeamStats
	// 	awaystats := r.AwayTeamStats

	// 	fmt.Printf("%s : Shots: %d, Goals: %d \n", r.HomeTeam, homestats.Shots, homestats.GoalsFor)
	// 	fmt.Printf("%s : Shots: %d, Goals: %d \n", r.AwayTeam, awaystats.Shots, awaystats.GoalsFor)
	// }
	collegeGameMap := MakeCollegeGameMap(collegeGames)
	proGameMap := MakeProGameMap(proGames)

	// collegeTeamMap := GetCollegeTeamMap()
	// proTeamMap := GetProTeamMap()
	collegePlayerMap := GetCollegePlayersMap()
	proPlayersMap := GetProPlayersMap()
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
		gameType, _ := ts.GetCurrentGameType(r.IsCollegeGame)
		upload.Collect(r, ts.SeasonID, uint(gameType))
		stars := GenerateThreeStars(r, ts.SeasonID)
		if r.IsCollegeGame {
			upload.ApplyGoalieStaminaChangesCollege(db, r, collegePlayerMap)
			game := collegeGameMap[r.GameID]
			game.UpdateScore(uint(r.HomeTeamScore), uint(r.AwayTeamScore), uint(r.HomeTeamShootoutScore), uint(r.AwayTeamShootoutScore), uint(r.Attendance), r.IsOvertime, r.IsOvertimeShootout)
			game.UpdateThreeStars(stars)
			repository.SaveCollegeGameRecord(game, db)
		} else {
			upload.ApplyGoalieStaminaChangesPro(db, r, proPlayersMap)
			game := proGameMap[r.GameID]
			game.UpdateScore(uint(r.HomeTeamScore), uint(r.AwayTeamScore), uint(r.HomeTeamShootoutScore), uint(r.AwayTeamShootoutScore), uint(r.Attendance), r.IsOvertime, r.IsOvertimeShootout)
			game.UpdateThreeStars(stars)
			repository.SaveProfessionalGameRecord(game, db)
		}
	}
	if !ts.IsTesting {
		upload.Flush(db)
	} else {

	}
}

func NewStatsUpload() *StatsUpload {
	return &StatsUpload{}
}

func (u *StatsUpload) Collect(state engine.GameState, seasonID, gameType uint) {
	// Team stats
	u.collectTeamStats(state, seasonID, gameType)

	// Player stats for both teams
	u.collectPlayerStats(state.HomeStrategy, state.WeekID, state.GameID, gameType, state.IsCollegeGame)
	u.collectPlayerStats(state.AwayStrategy, state.WeekID, state.GameID, gameType, state.IsCollegeGame)

	// Play-by-play
	u.collectPbP(state.Collector.PlayByPlays, state.IsCollegeGame)
}

func (u *StatsUpload) collectTeamStats(state engine.GameState, seasonID, gameType uint) {
	if state.IsCollegeGame {
		u.CollegeTeamStats = append(u.CollegeTeamStats,
			makeCollegeTeamStatsObject(state.WeekID, state.GameID, seasonID, gameType, state.HomeTeamStats),
			makeCollegeTeamStatsObject(state.WeekID, state.GameID, seasonID, gameType, state.AwayTeamStats),
		)
	} else {
		u.ProTeamStats = append(u.ProTeamStats,
			makeProTeamStatsObject(state.WeekID, state.GameID, seasonID, gameType, state.HomeTeamStats),
			makeProTeamStatsObject(state.WeekID, state.GameID, seasonID, gameType, state.AwayTeamStats),
		)
	}
}

func (u *StatsUpload) collectPlayerStats(pl engine.GamePlaybook, week, gameID, gameType uint, isCollege bool) {
	types := [][]engine.LineStrategy{pl.Forwards, pl.Defenders, pl.Goalies}
	for _, group := range types {
		for _, line := range group {
			for _, p := range line.Players {
				if isCollege {
					u.CollegePlayerStats = append(u.CollegePlayerStats,
						makeCollegePlayerStatsObject(week, gameID, gameType, p.Stats),
					)
				} else {
					u.ProPlayerStats = append(u.ProPlayerStats,
						makeProPlayerStatsObject(week, gameID, gameType, p.Stats),
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

func (u *StatsUpload) ApplyGoalieStaminaChangesCollege(db *gorm.DB, state engine.GameState, playerMap map[uint]structs.CollegePlayer) {
	homeStrategy := state.HomeStrategy
	awayStrategy := state.AwayStrategy

	homeGoalies := homeStrategy.Goalies
	awayGoalies := awayStrategy.Goalies
	homeBench := homeStrategy.BenchPlayers
	awayBench := awayStrategy.BenchPlayers

	for _, g := range homeGoalies {
		for _, p := range g.Players {
			triggerSave := false
			player := playerMap[p.ID]
			if player.ID <= 0 {
				continue
			}
			if p.Stats.TimeOnIce > 0 {
				triggerSave = true
				player.ApplyGoalieStaminaDrain()
			} else {
				if player.GoalieStamina < 100 {
					triggerSave = true
					player.RecoverGoalieStamina()
				}
			}
			if triggerSave {
				repository.SaveCollegeHockeyPlayerRecord(player, db)
			}
		}
	}

	for _, g := range awayGoalies {
		for _, p := range g.Players {
			triggerSave := false
			player := playerMap[p.ID]
			if player.ID <= 0 {
				continue
			}
			if p.Stats.TimeOnIce > 0 {
				triggerSave = true
				player.ApplyGoalieStaminaDrain()
			} else {
				if player.GoalieStamina < 100 {
					triggerSave = true
					player.RecoverGoalieStamina()
				}
			}
			if triggerSave {
				repository.SaveCollegeHockeyPlayerRecord(player, db)
			}
		}
	}

	for _, p := range homeBench {
		player := playerMap[p.ID]
		if player.ID <= 0 || player.Position != Goalie {
			continue
		}
		triggerSave := false
		if p.Stats.TimeOnIce > 0 {
			triggerSave = true
			player.ApplyGoalieStaminaDrain()
		} else {
			if player.GoalieStamina < util.MaxGoalieStamina {
				triggerSave = true
			}
			player.RecoverGoalieStamina()
		}

		if triggerSave {
			repository.SaveCollegeHockeyPlayerRecord(player, db)
		}
	}
	for _, p := range awayBench {
		player := playerMap[p.ID]
		if player.ID <= 0 || player.Position != Goalie {
			continue
		}
		triggerSave := false
		if p.Stats.TimeOnIce > 0 {
			triggerSave = true
			player.ApplyGoalieStaminaDrain()
		} else {
			if player.GoalieStamina < util.MaxGoalieStamina {
				triggerSave = true
			}
			player.RecoverGoalieStamina()
		}

		if triggerSave {
			repository.SaveCollegeHockeyPlayerRecord(player, db)
		}
	}

}
func (u *StatsUpload) ApplyGoalieStaminaChangesPro(db *gorm.DB, state engine.GameState, playerMap map[uint]structs.ProfessionalPlayer) {
	homeStrategy := state.HomeStrategy
	awayStrategy := state.AwayStrategy

	homeGoalies := homeStrategy.Goalies
	awayGoalies := awayStrategy.Goalies
	homeBench := homeStrategy.BenchPlayers
	awayBench := awayStrategy.BenchPlayers

	for _, g := range homeGoalies {
		for _, p := range g.Players {
			player := playerMap[p.ID]
			if player.ID <= 0 {
				continue
			}
			triggerSave := false
			if p.Stats.TimeOnIce > 0 {
				triggerSave = true
				player.ApplyGoalieStaminaDrain()
			} else {
				if player.GoalieStamina < util.MaxGoalieStamina {
					triggerSave = true
				}
				player.RecoverGoalieStamina()
			}

			if triggerSave {
				repository.SaveProPlayerRecord(player, db)
			}
		}
	}

	for _, g := range awayGoalies {
		for _, p := range g.Players {
			player := playerMap[p.ID]
			if player.ID <= 0 {
				continue
			}
			triggerSave := false
			if p.Stats.TimeOnIce > 0 {
				triggerSave = true
				player.ApplyGoalieStaminaDrain()
			} else {
				if player.GoalieStamina < util.MaxGoalieStamina {
					triggerSave = true
				}
				player.RecoverGoalieStamina()
			}

			if triggerSave {
				repository.SaveProPlayerRecord(player, db)
			}
		}
	}

	for _, p := range homeBench {
		player := playerMap[p.ID]
		if player.ID <= 0 || player.Position != Goalie {
			continue
		}
		triggerSave := false
		if p.Stats.TimeOnIce > 0 {
			triggerSave = true
			player.ApplyGoalieStaminaDrain()
		} else {
			if player.GoalieStamina < util.MaxGoalieStamina {
				triggerSave = true
			}
			player.RecoverGoalieStamina()
		}

		if triggerSave {
			repository.SaveProPlayerRecord(player, db)
		}
	}
	for _, p := range awayBench {
		player := playerMap[p.ID]
		if player.ID <= 0 || player.Position != Goalie {
			continue
		}
		triggerSave := false
		if p.Stats.TimeOnIce > 0 {
			triggerSave = true
			player.ApplyGoalieStaminaDrain()
		} else {
			if player.GoalieStamina < util.MaxGoalieStamina {
				triggerSave = true
			}
			player.RecoverGoalieStamina()
		}

		if triggerSave {
			repository.SaveProPlayerRecord(player, db)
		}
	}

}

func PrepareGames(collegeGames []structs.CollegeGame, proGames []structs.ProfessionalGame, collegeStandingsMap map[uint]structs.CollegeStandings, proStandingsMap map[uint]structs.ProfessionalStandings) []structs.GameDTO {
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

			// Enhanced panic recovery with detailed error reporting
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("PANIC in PrepareGames goroutine for game %s vs %s (IDs: %d vs %d):\n",
						c.HomeTeam, c.AwayTeam, c.HomeTeamID, c.AwayTeamID)
					fmt.Printf("Error: %v\n", r)

					// Check what data exists for debugging
					fmt.Printf("Debug info:\n")
					fmt.Printf("  - Home team roster exists: %t (length: %d)\n",
						collegeTeamRosterMap[c.HomeTeamID] != nil, len(collegeTeamRosterMap[c.HomeTeamID]))
					fmt.Printf("  - Away team roster exists: %t (length: %d)\n",
						collegeTeamRosterMap[c.AwayTeamID] != nil, len(collegeTeamRosterMap[c.AwayTeamID]))
					fmt.Printf("  - Home team lineups exist: %t (length: %d)\n",
						collegeLineupMap[c.HomeTeamID] != nil, len(collegeLineupMap[c.HomeTeamID]))
					fmt.Printf("  - Away team lineups exist: %t (length: %d)\n",
						collegeLineupMap[c.AwayTeamID] != nil, len(collegeLineupMap[c.AwayTeamID]))
					fmt.Printf("  - Home shootout lineup exists: %t\n",
						collegeShootoutLineupMap[c.HomeTeamID].TeamID != 0)
					fmt.Printf("  - Away shootout lineup exists: %t\n",
						collegeShootoutLineupMap[c.AwayTeamID].TeamID != 0)
					fmt.Printf("  - Arena exists: %t (ArenaID: %d)\n",
						arenaMap[c.ArenaID].ID != 0, c.ArenaID)
					fmt.Printf("  - Home team standings exist: %t\n",
						collegeStandingsMap[c.HomeTeamID].TeamID != 0)
				}
			}()

			if c.GameComplete {
				return
			}

			mutex.Lock()

			// Step-by-step data retrieval with nil checks and detailed logging
			fmt.Printf("Processing game: %s vs %s\n", c.HomeTeam, c.AwayTeam)

			// Check home team roster
			htr := collegeTeamRosterMap[c.HomeTeamID]
			if htr == nil {
				fmt.Printf("ERROR: No roster found for home team %s (ID: %d)\n", c.HomeTeam, c.HomeTeamID)
				mutex.Unlock()
				return
			}

			// Check away team roster
			atr := collegeTeamRosterMap[c.AwayTeamID]
			if atr == nil {
				fmt.Printf("ERROR: No roster found for away team %s (ID: %d)\n", c.AwayTeam, c.AwayTeamID)
				mutex.Unlock()
				return
			}

			// Check home team lineups
			htl := collegeLineupMap[c.HomeTeamID]
			if htl == nil {
				fmt.Printf("ERROR: No lineups found for home team %s (ID: %d)\n", c.HomeTeam, c.HomeTeamID)
				mutex.Unlock()
				return
			}

			// Check away team lineups
			atl := collegeLineupMap[c.AwayTeamID]
			if atl == nil {
				fmt.Printf("ERROR: No lineups found for away team %s (ID: %d)\n", c.AwayTeam, c.AwayTeamID)
				mutex.Unlock()
				return
			}

			// Check shootout lineups
			htsl := collegeShootoutLineupMap[c.HomeTeamID]
			if htsl.TeamID == 0 {
				fmt.Printf("ERROR: No shootout lineup found for home team %s (ID: %d)\n", c.HomeTeam, c.HomeTeamID)
				mutex.Unlock()
				return
			}

			atsl := collegeShootoutLineupMap[c.AwayTeamID]
			if atsl.TeamID == 0 {
				fmt.Printf("ERROR: No shootout lineup found for away team %s (ID: %d)\n", c.AwayTeam, c.AwayTeamID)
				mutex.Unlock()
				return
			}

			// Check arena
			arena := arenaMap[c.ArenaID]
			if arena.ID == 0 {
				fmt.Printf("ERROR: No arena found for ArenaID: %d\n", c.ArenaID)
				mutex.Unlock()
				return
			}

			// Generate playbooks with error checking
			fmt.Printf("Generating playbooks for %s vs %s...\n", c.HomeTeam, c.AwayTeam)

			var hp, ap structs.PlayBookDTO

			// Generate home playbook with error handling
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("PANIC generating home playbook for %s: %v\n", c.HomeTeam, r)
						fmt.Printf("  Home lineups count: %d\n", len(htl))
						fmt.Printf("  Home roster count: %d\n", len(htr))
						panic(r) // Re-panic to be caught by outer handler
					}
				}()
				hp = getCollegePlaybookDTO(htl, htr, htsl)
			}()

			// Generate away playbook with error handling
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("PANIC generating away playbook for %s: %v\n", c.AwayTeam, r)
						fmt.Printf("  Away lineups count: %d\n", len(atl))
						fmt.Printf("  Away roster count: %d\n", len(atr))
						panic(r) // Re-panic to be caught by outer handler
					}
				}()
				ap = getCollegePlaybookDTO(atl, atr, atsl)
			}()

			capacity := arena.Capacity
			currentStandings := collegeStandingsMap[c.HomeTeamID]
			attendancePercent := getAttendancePercent(int(currentStandings.TotalWins)+int(currentStandings.TotalOTWins), int(currentStandings.TotalLosses))
			if c.IsPreseason {
				attendancePercent = 1.0
			}
			fanCount := uint32(float64(capacity) * attendancePercent)
			mutex.Unlock()

			match := structs.GameDTO{
				GameID:        c.ID,
				GameInfo:      c.BaseGame,
				HomeStrategy:  hp,
				AwayStrategy:  ap,
				IsCollegeGame: true,
				Attendance:    fanCount,
				Capacity:      uint32(arena.Capacity),
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
			arena := arenaMap[g.ArenaID]
			capacity := arena.Capacity
			currentStandings := proStandingsMap[g.HomeTeamID]
			attendancePercent := getAttendancePercent(int(currentStandings.TotalWins)+int(currentStandings.TotalOTWins), int(currentStandings.TotalLosses))
			if g.IsPreseason {
				attendancePercent = 1.0
			}
			fanCount := uint32(float64(capacity) * attendancePercent)
			mutex.Unlock()

			match := structs.GameDTO{
				GameID:        g.ID,
				GameInfo:      g.BaseGame,
				HomeStrategy:  hp,
				AwayStrategy:  ap,
				IsCollegeGame: false,
				Attendance:    fanCount,
				Capacity:      uint32(arena.Capacity),
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

	sort.Slice(gameDTOList, func(i, j int) bool {
		return gameDTOList[i].IsCollegeGame
	})
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

func PrepareCollegeTournamentGamesFormat(db *gorm.DB, ts structs.Timestamp) {
	seasonID := ts.SeasonID
	nextGameID := repository.FindLatestCHLGameID() + 1
	collegeTeams := repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "1"})
	teamMap := MakeCollegeTeamMap(collegeTeams)
	standingsMap := GetCollegeStandingsMap(strconv.Itoa(int(seasonID)))
	conferenceMap := map[uint8][]*structs.CollegeStandings{}
	conferenceIDList := []uint8{1, 2, 3, 4, 5, 6, 7}
	quarterfinalsSeries := []structs.CollegeSeries{}
	semiFinalsAndFinalsGames := []structs.CollegeGame{}

	for _, t := range collegeTeams {
		standings := standingsMap[t.ID]
		if standings.ID == 0 {
			continue
		}
		if len(conferenceMap[t.ConferenceID]) > 0 {
			conferenceMap[t.ConferenceID] = append(conferenceMap[t.ConferenceID], &standings)
		} else {
			conferenceMap[t.ConferenceID] = []*structs.CollegeStandings{&standings}
		}
	}

	for _, cid := range conferenceIDList {
		conference := conferenceMap[cid]
		for _, s := range conference {
			s.CalculateConferencePoints()
		}

		sort.Slice(conferenceMap[cid], func(i, j int) bool {
			if conferenceMap[cid][i].Points == conferenceMap[cid][j].Points {
				return conferenceMap[cid][i].GoalsFor > conferenceMap[cid][j].GoalsFor
			}
			return conferenceMap[cid][i].Points > conferenceMap[cid][j].Points
		})

		// If CID == 2, conduct different tournament structure for Big Ten.
		// Else, standard 8 team tournament. Series are best of 3, followed by one semifinal game and one finals game
		if cid == 2 {
			seven := TopN(conferenceMap[cid], 7)
			pairs := [][2]*structs.CollegeStandings{
				{seven[1], seven[6]}, // 2v7  -> Semi #1 AWAY (vs #1)
				{seven[2], seven[5]}, // 3v6  -> Semi #2 (TBD HOA or later)
				{seven[3], seven[4]}, // 4v5  -> Semi #2 (TBD HOA or later)
			}
			semiFinalID1 := nextGameID     // 1 vs 2/7
			semiFinalID2 := nextGameID + 1 // 3/6 vs 4/5
			finalsID := nextGameID + 2     // winner of nextGameID & nextGameID2 == Conference Finals
			conference := ""

			for idx, p := range pairs {
				a, b := p[0], p[1]
				homeTeam := teamMap[a.TeamID]
				conference = homeTeam.Conference
				// Route: index 0 is (2/7) -> Semi #1, AWAY; others -> Semi #2
				ngID := semiFinalID2
				nextHOA := "H" // neutral placeholder; can be "" if you’ll reseed later
				if idx == 0 {
					ngID = semiFinalID1
					nextHOA = "A" // winner is away vs #1
				}

				series := structs.CollegeSeries{
					BaseSeries: structs.BaseSeries{
						SeasonID:    seasonID,
						SeriesName:  fmt.Sprintf("%s Conference Quarterfinals", conference),
						BestOfCount: 3,
						HomeTeamID:  a.TeamID, HomeTeam: a.TeamName, HomeTeamRank: 2 + uint(idx), // 2,3,4
						AwayTeamID: b.TeamID, AwayTeam: b.TeamName, AwayTeamRank: uint(7 - idx), // 7,6,5
						GameCount:     0,
						IsPlayoffGame: true,
					},
					NextGameID:   ngID,
					NextGameHOA:  nextHOA,
					ConferenceID: uint8(cid),
				}
				quarterfinalsSeries = append(quarterfinalsSeries, series)
			}
			// Semifinal Game 1
			top1 := seven[0]
			top1Team := teamMap[top1.TeamID]
			semifinalGame1 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID1},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conference),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					HomeTeamID: top1.TeamID, HomeTeam: top1.TeamName, HomeTeamRank: 1,
					Arena: top1Team.Arena, NextGameID: finalsID, NextGameHOA: "H",
					GameDay: "A",
				},
				IsConferenceTournament: true,
			}

			// Semifinal Game 2
			semifinalGame2 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID2},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conference),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					NextGameID: finalsID, NextGameHOA: "A",
					GameDay: "A",
				},
				IsConferenceTournament: true,
			}

			finalsGame := structs.CollegeGame{
				Model: gorm.Model{ID: finalsID},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Finals", conference),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					GameDay: "B",
				},
				IsConferenceTournament: true,
			}

			semiFinalsAndFinalsGames = append(semiFinalsAndFinalsGames, semifinalGame1, semifinalGame2, finalsGame)
			nextGameID += 3
		} else {
			eight := TopN(conferenceMap[cid], 8)
			if len(eight) < 8 {
				// Should not happen since each conference is at least 8 and above.
				// Not enough teams for 8—either skip or fall back to a smaller bracket.
				// For now: continue (or handle your smaller-bracket flow here)
				continue
			}
			pairs := SeededPairs(eight, 4) // (1v8),(2v7),(3v6),(4v5)
			semiFinalID1 := nextGameID     // 1 vs 4
			semiFinalID2 := nextGameID + 1 // 2 vs 3
			finalsID := nextGameID + 2     // winner of nextGameID & nextGameID2 == Conference Finals
			conference := ""
			for qfIdx, p := range pairs {
				a, b := p[0], p[1]
				homeTeam := teamMap[a.TeamID]
				conference = homeTeam.Conference

				// QF1 and QF4 feed Semi #1; QF2 and QF3 feed Semi #2
				ngID := semiFinalID2
				if qfIdx == 0 || qfIdx == 3 {
					ngID = semiFinalID1
				}

				// Since pairings in order are (1v8), (2v7), (3v6), (4v5);
				// the first two should point to H as their nextHOA. the rest will be A.
				nextHOA := "H"
				if qfIdx > 1 {
					nextHOA = "A"
				}

				series := structs.CollegeSeries{
					BaseSeries: structs.BaseSeries{
						SeasonID:    seasonID,
						SeriesName:  fmt.Sprintf("%s Conference Quarterfinals", conference),
						BestOfCount: 3,
						HomeTeamID:  a.TeamID, HomeTeam: a.TeamName, HomeTeamRank: uint(qfIdx + 1),
						AwayTeamID: b.TeamID, AwayTeam: b.TeamName, AwayTeamRank: uint(8 - qfIdx),
						GameCount:     1,
						IsPlayoffGame: true,
					},
					NextGameID:   ngID,
					NextGameHOA:  nextHOA,
					ConferenceID: uint8(cid),
				}
				quarterfinalsSeries = append(quarterfinalsSeries, series)
			}
			// Semifinal Game 1
			semifinalGame1 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID1},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conference),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					NextGameID: finalsID, NextGameHOA: "H",
					GameDay: "A", IsPlayoffGame: true,
				},
				IsConferenceTournament: true,
			}

			// Semifinal Game 2
			semifinalGame2 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID2},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conference),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					NextGameID: finalsID, NextGameHOA: "A",
					GameDay: "A", IsPlayoffGame: true,
				},
				IsConferenceTournament: true,
			}

			// Finals Game
			finalsGame := structs.CollegeGame{
				Model: gorm.Model{ID: finalsID},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Finals", conference),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					GameDay: "B",
				},
				IsConferenceTournament: true,
			}

			semiFinalsAndFinalsGames = append(semiFinalsAndFinalsGames, semifinalGame1, semifinalGame2, finalsGame)
			nextGameID += 3
		}
	}
	// Create College Series in batch
	// Create college games in batch
	repository.CreateCHLSeriesRecordsBatch(db, quarterfinalsSeries, 20)
	repository.CreateCHLGamesRecordsBatch(db, semiFinalsAndFinalsGames, 50)
}

func PrepareCHLPostSeasonGamesFormat(db *gorm.DB, ts structs.Timestamp) {
	seasonID := ts.SeasonID
	baseID := repository.FindLatestCHLGameID() + 1
	collegeTeams := repository.FindAllCollegeTeams(repository.TeamClauses{LeagueID: "1"})
	collegeStandings := repository.FindAllCollegeStandings(repository.StandingsQuery{SeasonID: strconv.Itoa(int(seasonID))})
	stMap := MakeCollegeStandingsMap(collegeStandings)
	pool := []*structs.CollegeStandings{}
	qualified := map[uint]bool{}

	// Conference Tournament Winners
	for _, t := range collegeTeams {
		standings := stMap[t.ID]
		if standings.ID == 0 {
			continue
		}
		if standings.IsConferenceTournamentChampion {
			sCopy := standings // be safe
			pool = append(pool, &sCopy)
			qualified[t.ID] = true
		}
	}

	// Sort collegeStandings by Points, Goals For.
	// Note: This current format may be biased towards conference tournament teams considering that these games are added to the standings as well.
	// May need to think of iterative approach where standings are updated in real time based on game and then we conduct a sort
	// Thoughts?
	sort.Slice(collegeStandings, func(i, j int) bool {
		if collegeStandings[i].Points == collegeStandings[j].Points {
			return collegeStandings[i].GoalsFor > collegeStandings[j].GoalsFor
		}
		return collegeStandings[i].Points > collegeStandings[j].Points
	})

	// Iterate by collegeStandings, checkif standings.TeamID is in isQualified map
	// If not, add to standingsList until length is 16
	for _, s := range collegeStandings {
		if len(pool) == 16 {
			break
		}
		if !qualified[s.TeamID] {
			sCopy := s // take address of stable copy
			pool = append(pool, &sCopy)
			qualified[s.TeamID] = true
		}
	}

	// Then sort standingsList by points, goals for, and seed 1-16
	// Then create individual games going by 1v16, 2v15, etc.
	sort.Slice(pool, func(i, j int) bool {
		if pool[i].Points == pool[j].Points {
			return pool[i].GoalsFor > pool[j].GoalsFor
		}
		return pool[i].Points > pool[j].Points
	})

	// Update Seed Ranks (1,1,1,1,2,2,2,2,3,3,3,3,4,4,4,4)
	for idx, s := range pool {
		s.AssignRank(idx/4 + 1)
	}

	order := []int{0, 15, 7, 8, 3, 12, 4, 11, 1, 14, 6, 9, 2, 13, 5, 10}
	games := make([]structs.CollegeGame, 0, 15)

	mk := func(id uint, title string, week int) structs.CollegeGame {
		return structs.CollegeGame{
			Model: gorm.Model{ID: id},
			BaseGame: structs.BaseGame{
				GameTitle:     title,
				SeasonID:      seasonID,
				WeekID:        util.GetWeekID(seasonID, uint(week)),
				Week:          week,
				IsNeutralSite: true,
				IsPlayoffGame: true,
				Arena:         "TBD", City: "TBD", State: "TBD", Country: "TBD",
			},
		}
	}

	arenas := repository.FindAllArenas(repository.ArenaQuery{Country: "USA", GreaterThanID: "66", LessThanID: "123"})

	// ---------- Round of 16 (ids baseID..baseID+7) ----------
	for i := 0; i < 8; i++ {
		a := pool[order[2*i+0]]
		b := pool[order[2*i+1]]
		arenaIdx := util.GenerateIntFromRange(0, len(arenas)-1)
		arena := arenas[arenaIdx]

		g := mk(baseID+uint(i), fmt.Sprintf("%d SimCHL Round of 16", ts.Season), 20)
		g.HomeTeamID, g.HomeTeam, g.HomeTeamRank = a.TeamID, a.TeamName, a.Rank
		g.AwayTeamID, g.AwayTeam, g.AwayTeamRank = b.TeamID, b.TeamName, b.Rank
		g.ArenaID = arena.ID
		g.Arena = arena.Name
		g.City = arena.City
		g.State = arena.State
		g.Country = arena.Country

		// parent = quarterfinal, HOA: upper child (even i) is H, lower (odd i) is A
		g.NextGameID = baseID + 8 + uint(i/2)
		if i%2 == 0 {
			g.NextGameHOA = "H"
		} else {
			g.NextGameHOA = "A"
		}

		g.GameDay = "A"
		games = append(games, g)
	}

	// ---------- Quarterfinals (ids baseID+8..baseID+11) ----------
	for q := 0; q < 4; q++ {
		arenaIdx := util.GenerateIntFromRange(0, len(arenas)-1)
		arena := arenas[arenaIdx]
		g := mk(baseID+8+uint(q), fmt.Sprintf("%d SimCHL Quarterfinal", ts.Season), 20)
		g.NextGameID = baseID + 12 + uint(q/2)
		if q%2 == 0 {
			g.NextGameHOA = "H"
		} else {
			g.NextGameHOA = "A"
		}
		g.GameDay = "B"
		g.ArenaID = arena.ID
		g.Arena = arena.Name
		g.City = arena.City
		g.State = arena.State
		g.Country = arena.Country
		games = append(games, g)
	}

	// ---------- Frozen Four (Semifinals) (ids baseID+12..baseID+13) ----------
	// Generate Frozen Four & National Championship Location
	arenaIdx := util.GenerateIntFromRange(0, len(arenas)-1)
	arena := arenas[arenaIdx]
	for s := 0; s < 2; s++ {
		g := mk(baseID+12+uint(s), fmt.Sprintf("%d SimCHL Frozen Four Semifinal", ts.Season), 21)
		g.NextGameID = baseID + 14
		if s == 0 {
			g.NextGameHOA = "H"
		} else {
			g.NextGameHOA = "A"
		}
		g.GameDay = "A"
		g.ArenaID = arena.ID
		g.Arena = arena.Name
		g.City = arena.City
		g.State = arena.State
		g.Country = arena.Country
		games = append(games, g)
	}

	// ---------- National Championship (id baseID+14) ----------
	final := mk(baseID+14, fmt.Sprintf("%d SimCHL National Championship", ts.Season), 21)
	final.IsNationalChampionship = true
	final.GameDay = "C"
	final.ArenaID = arena.ID
	final.Arena = arena.Name
	final.City = arena.City
	final.State = arena.State
	final.Country = arena.Country
	games = append(games, final)

	// Generate a Toilet Bowl Matchup between the two lowest ranked teams in the entire standings.
	// This is an exhibition game and does not affect any standings or stats.
	secondWorstTeam := collegeStandings[len(collegeStandings)-2]
	worstTeam := collegeStandings[len(collegeStandings)-1]
	toiletBowl := mk(baseID+15, fmt.Sprintf("%d SimCHL Toilet Bowl", ts.Season), 21)
	toiletBowl.GameDay = "C"
	toiletBowl.HomeTeamID = secondWorstTeam.TeamID
	toiletBowl.HomeTeam = secondWorstTeam.TeamName
	toiletBowl.HomeTeamRank = secondWorstTeam.Rank

	toiletBowl.AwayTeamID = worstTeam.TeamID
	toiletBowl.AwayTeam = worstTeam.TeamName
	toiletBowl.AwayTeamRank = worstTeam.Rank
	// The Toilet Bowl will always be held at the Hart Center in Worcester, MA (Holy Cross)
	// since it is a neutral site and a known location.
	toiletBowl.IsNeutralSite = true
	// Hart Center ArenaID = 24
	toiletBowl.ArenaID = 24
	toiletBowl.Arena = "The Hart Center"
	toiletBowl.City = "Worcester"
	toiletBowl.State = "MA"
	toiletBowl.Country = "USA"
	games = append(games, toiletBowl)
	repository.CreateCHLGamesRecordsBatch(db, games, 50)
}

func PreparePHLPostSeasonGamesFormat(db *gorm.DB, ts structs.Timestamp) {
	seasonID := ts.SeasonID
	nextSeriesID := repository.FindLatestPHLSeriesID() + 1
	proTeams := repository.FindAllProTeams(repository.TeamClauses{LeagueID: "1"})
	teamMap := MakeProTeamMap(proTeams)
	standingsMap := GetProStandingsMap(strconv.Itoa(int(seasonID)))
	divisionMap := map[uint8][]*structs.ProfessionalStandings{}
	divisionIDList := []uint8{1, 2, 3, 4}
	qualifyingTeams := []*structs.ProfessionalStandings{}
	postSeasonSeriesList := []structs.ProSeries{}

	for _, t := range proTeams {
		standings := standingsMap[t.ID]
		if standings.ID == 0 {
			continue
		}
		if len(divisionMap[t.DivisionID]) > 0 {
			divisionMap[t.DivisionID] = append(divisionMap[t.DivisionID], &standings)
		} else {
			divisionMap[t.DivisionID] = []*structs.ProfessionalStandings{&standings}
		}
	}

	// Get Qualifying Top 2 teams in each division
	for _, did := range divisionIDList {
		division := divisionMap[did]

		// Sort by Points, Goals For
		sort.Slice(division, func(i, j int) bool {
			if division[i].Points == division[j].Points {
				return division[i].GoalsFor > division[j].GoalsFor
			}
			return division[i].Points > division[j].Points
		})

		// Then get top two teams from the division
		qualifyingTeams = append(qualifyingTeams, division[:2]...)
	}

	// So, Divisions 1 and 2 should be part of the same conference. 3 and 4 the same.
	// We will need to pair the top team from 1 division to face the 2nd best team from the opposite division in the same conference. And then vice versa. This will serve as the quarterfinals series, best of 7.
	pairs := [][2]*structs.ProfessionalStandings{
		{qualifyingTeams[0], qualifyingTeams[3]}, // Div 1 #1 vs Div 2 #2
		{qualifyingTeams[1], qualifyingTeams[2]}, // Div 1 #2 vs Div 2 #1
		{qualifyingTeams[4], qualifyingTeams[7]}, // Div 3 #1 vs Div 4 #2
		{qualifyingTeams[5], qualifyingTeams[6]}, // Div 3 #2 vs Div 4 #1
	}

	// Quarterfinal Series
	for qfIdx, p := range pairs {
		a, b := p[0], p[1]
		homeTeam := teamMap[a.TeamID]
		awayTeam := teamMap[b.TeamID]
		nextSeriesHoa := "H"
		if qfIdx%2 == 1 {
			nextSeriesHoa = "A"
		}

		quarterFinalsSeries := structs.ProSeries{
			Model: gorm.Model{ID: nextSeriesID + uint(qfIdx)},
			BaseSeries: structs.BaseSeries{
				SeasonID:   seasonID,
				HomeTeamID: a.TeamID, HomeTeam: a.TeamName, HomeTeamRank: 1,
				HomeTeamCoach:   homeTeam.Coach,
				AwayTeamCoach:   awayTeam.Coach,
				IsInternational: homeTeam.LeagueID != 1 && awayTeam.LeagueID != 1,
				AwayTeamID:      b.TeamID, AwayTeam: b.TeamName, AwayTeamRank: 2,
				SeriesName:    fmt.Sprintf("%d %s vs %s Quarterfinals", ts.Season, homeTeam.Division, awayTeam.Division),
				BestOfCount:   7,
				GameCount:     1,
				IsPlayoffGame: true,
				NextSeriesHOA: nextSeriesHoa,                    // Higher seed is home
				NextSeriesID:  nextSeriesID + 4 + uint(qfIdx/2), // Semifinal Series ID
			},
		}
		postSeasonSeriesList = append(postSeasonSeriesList, quarterFinalsSeries)
	}

	// Make two Semifinals Series records
	for sfIdx := 0; sfIdx < 2; sfIdx++ {
		nextSeriesHoa := "H"
		if sfIdx == 1 {
			nextSeriesHoa = "A"
		}
		semiFinalsSeries := structs.ProSeries{
			Model: gorm.Model{ID: nextSeriesID + 4 + uint(sfIdx)},
			BaseSeries: structs.BaseSeries{
				SeasonID:   seasonID,
				HomeTeamID: 0, HomeTeam: "", HomeTeamRank: 0,
				AwayTeamID: 0, AwayTeam: "", AwayTeamRank: 0,
				HomeTeamCoach:   "",
				AwayTeamCoach:   "",
				IsInternational: false,
				NextSeriesHOA:   nextSeriesHoa,    // Higher seed is home
				NextSeriesID:    nextSeriesID + 6, // Finals Series ID
				SeriesName:      fmt.Sprintf("%d SimPHL Semifinals", ts.Season),
				BestOfCount:     7,
				GameCount:       1,
			},
		}
		postSeasonSeriesList = append(postSeasonSeriesList, semiFinalsSeries)
	}

	// Finals Series
	finalsSeries := structs.ProSeries{
		Model: gorm.Model{ID: nextSeriesID + 6},
		BaseSeries: structs.BaseSeries{
			SeasonID:   seasonID,
			HomeTeamID: 0, HomeTeam: "", HomeTeamRank: 0,
			AwayTeamID: 0, AwayTeam: "", AwayTeamRank: 0,
			HomeTeamCoach:   "",
			AwayTeamCoach:   "",
			IsInternational: false,
			NextSeriesHOA:   "", // Higher seed is home
			NextSeriesID:    0,
			IsTheFinals:     true,
			SeriesName:      fmt.Sprintf("%d SimPHL Finals", ts.Season),
			BestOfCount:     7,
			GameCount:       1,
		},
	}

	postSeasonSeriesList = append(postSeasonSeriesList, finalsSeries)

	repository.CreatePHLSeriesRecordsBatch(db, postSeasonSeriesList, 20)

}

func GenerateCollegeTournamentQuarterfinalsGames(db *gorm.DB, ts structs.Timestamp) {
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	teamMap := GetCollegeTeamMap()
	collegeGames := repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, WeekID: weekID, GameCompleted: "N"})
	if len(collegeGames) > 0 {
		return
	}
	collegeGamesUpload := []structs.CollegeGame{}
	activeCHLSeries := repository.FindCollegeSeriesRecords(seasonID)

	for _, s := range activeCHLSeries {
		if s.HomeTeamID == 0 || s.AwayTeamID == 0 {
			continue
		}
		gameCount := strconv.Itoa(int(s.GameCount))
		arena := ""
		city := ""
		state := ""
		country := ""
		seriesName := s.SeriesName
		matchName := seriesName + " Game: " + gameCount
		ht := teamMap[s.HomeTeamID]
		arenaID := ht.ArenaID
		arena = ht.Arena
		city = ht.City
		state = ht.State
		country = ht.Country
		weekID := util.GetWeekID(ts.SeasonID, 18)
		week := 18

		collegeGame := structs.CollegeGame{
			BaseGame: structs.BaseGame{
				GameTitle: matchName,
				SeasonID:  ts.SeasonID, WeekID: weekID, Week: week,
				HomeTeamID: s.HomeTeamID, HomeTeam: s.HomeTeam, HomeTeamRank: s.HomeTeamRank,
				HomeTeamCoach: s.HomeTeamCoach,
				AwayTeamID:    s.AwayTeamID, AwayTeam: s.AwayTeam, AwayTeamRank: s.AwayTeamRank,
				AwayTeamCoach: s.AwayTeamCoach,
				Arena:         arena,
				ArenaID:       uint(arenaID),
				City:          city,
				State:         state,
				Country:       country,
				GameDay:       "A",
				SeriesID:      s.ID,
			},
			IsConferenceTournament: true,
		}
		collegeGamesUpload = append(collegeGamesUpload, collegeGame)
	}

	repository.CreateCHLGamesRecordsBatch(db, collegeGamesUpload, 50)
}

func GenerateProPlayoffGames(db *gorm.DB, ts structs.Timestamp) {
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	teamMap := GetProTeamMap()
	professionalGames := repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID, WeekID: weekID})

	// Check if all series have been completed. If it's the middle of the week, then we should check & see if there are still active series.
	// If every series is complete, we can proceed to the next series.
	incompleteGames := []structs.ProfessionalGame{}
	for _, g := range professionalGames {
		if !g.GameComplete {
			incompleteGames = append(incompleteGames, g)
		}
	}
	// If game still exist, do NOT generate new games
	if len(incompleteGames) > 0 {
		return
	}

	gameDay := ts.GetGameDay()

	// Get Active Pro Series
	proSeries := repository.FindProSeriesRecords(strconv.Itoa(int(ts.SeasonID)))
	proGamesUpload := []structs.ProfessionalGame{}

	for _, s := range proSeries {
		if s.HomeTeamID == 0 || s.AwayTeamID == 0 || s.SeriesComplete {
			continue
		}

		if s.IsTheFinals && s.SeriesComplete {
			ts.EndTheProfessionalSeason()
			repository.SaveTimestamp(ts, db)
			break
		}
		gameCount := strconv.Itoa(int(s.GameCount))
		homeTeam := ""
		homeTeamID := 0
		homeTeamCoach := ""
		homeTeamRank := 0
		awayTeam := ""
		awayTeamID := 0
		awayTeamCoach := ""
		awayTeamRank := 0
		arenaID := 0
		arena := ""
		city := ""
		state := ""
		country := ""
		seriesName := s.SeriesName
		matchName := seriesName + " Game: " + gameCount
		// Game 1, 2, 5, or 7 => Higher Seed is Home
		// Game 3, 4, or 6 => Lower Seed is Home
		// If Game 7 does not exist, it means the series ended in 4, 5, or 6 games.
		if gameCount == "1" || gameCount == "2" || gameCount == "5" || gameCount == "7" {
			ht := teamMap[s.HomeTeamID]
			homeTeam = s.HomeTeam
			homeTeamID = int(s.HomeTeamID)
			homeTeamCoach = s.HomeTeamCoach
			homeTeamRank = int(s.HomeTeamRank)
			awayTeam = s.AwayTeam
			awayTeamID = int(s.AwayTeamID)
			awayTeamCoach = s.AwayTeamCoach
			awayTeamRank = int(s.AwayTeamRank)
			arenaID = int(ht.ArenaID)
			arena = ht.Arena
			city = ht.City
			state = ht.State
			country = ht.Country
		} else {
			ht := teamMap[s.AwayTeamID]
			homeTeam = s.AwayTeam
			homeTeamID = int(s.AwayTeamID)
			homeTeamCoach = s.AwayTeamCoach
			homeTeamRank = int(s.AwayTeamRank)
			awayTeam = s.HomeTeam
			awayTeamID = int(s.HomeTeamID)
			awayTeamCoach = s.HomeTeamCoach
			awayTeamRank = int(s.HomeTeamRank)
			arenaID = int(ht.ArenaID)
			arena = ht.Arena
			city = ht.City
			state = ht.State
			country = ht.Country
		}

		proGame := structs.ProfessionalGame{
			BaseGame: structs.BaseGame{
				GameTitle: matchName,
				SeasonID:  ts.SeasonID, WeekID: ts.WeekID, Week: int(ts.Week),
				HomeTeamID: uint(homeTeamID), HomeTeam: homeTeam, HomeTeamRank: uint(homeTeamRank),
				HomeTeamCoach: homeTeamCoach,
				AwayTeamID:    uint(awayTeamID), AwayTeam: awayTeam, AwayTeamRank: uint(awayTeamRank),
				AwayTeamCoach: awayTeamCoach,
				Arena:         arena,
				ArenaID:       uint(arenaID),
				City:          city,
				State:         state,
				Country:       country,
				GameDay:       gameDay,
				SeriesID:      s.ID,
				IsPlayoffGame: s.IsPlayoffGame,
			},
			IsStanleyCup: s.IsTheFinals,
			SeriesID:     s.ID,
		}
		proGamesUpload = append(proGamesUpload, proGame)
	}

	repository.CreatePHLGamesRecordsBatch(db, proGamesUpload, 50)
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
				gameDay, "", teamMap, true,
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

func GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGamesByCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
}

func GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay string, isPreseason bool) []structs.ProfessionalGame {
	return repository.FindProfessionalGamesByCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
}

func GetCollegeGamesByTeamIDAndSeasonID(teamID, seasonID string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, TeamID: teamID, IsPreseason: isPreseason})
}

func GetProfessionalGamesByTeamIDAndSeasonID(teamID, seasonID string, isPreseason bool) []structs.ProfessionalGame {
	return repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID, TeamID: teamID, IsPreseason: isPreseason})
}

func GetCollegeGamesBySeasonID(seasonID string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, IsPreseason: isPreseason})
}

func GetProfessionalGamesBySeasonID(seasonID string, isPreseason bool) []structs.ProfessionalGame {
	return repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID, IsPreseason: isPreseason})
}

func GetCollegeGameByID(id string) structs.CollegeGame {
	return repository.FindCollegeGameRecord(id)
}

func GetProfessionalGameByID(id string) structs.ProfessionalGame {
	return repository.FindProfessionalGameRecord(id)
}

func GetArenaMap() map[uint]structs.Arena {
	arenas := repository.FindAllArenas(repository.ArenaQuery{})
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
		switch l.LineType {
		case 1:
			forwards = append(forwards, l.BaseLineup)
		case 2:
			defenders = append(defenders, l.BaseLineup)
		default:
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
		switch l.LineType {
		case 1:
			forwards = append(forwards, l.BaseLineup)
		case 2:
			defenders = append(defenders, l.BaseLineup)
		default:
			goalies = append(goalies, l.BaseLineup)
		}
	}
	return forwards, defenders, goalies
}

func generateCollegeGame(seasonID, weekID, week, hid, aid uint, gameDay, gameTitle string, teamMap map[uint]structs.CollegeTeam, isPreseason bool) structs.CollegeGame {
	return structs.CollegeGame{
		BaseGame: structs.BaseGame{
			WeekID:      weekID,
			Week:        int(week),
			GameDay:     gameDay,
			GameTitle:   gameTitle,
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

func GetPlayoffSeriesBySeriesID(seriesID string) structs.ProSeries {
	return repository.FindPlayoffSeriesByID(seriesID)
}

func GenerateThreeStars(state engine.GameState, seasonID uint) structs.ThreeStars {
	types := [][]engine.LineStrategy{state.HomeStrategy.Forwards, state.HomeStrategy.Defenders, state.HomeStrategy.Goalies, state.AwayStrategy.Forwards, state.AwayStrategy.Defenders, state.AwayStrategy.Goalies}
	threeStars := []structs.ThreeStarsObj{}
	winningTeamID := state.HomeTeamID
	if state.AwayTeamWin {
		winningTeamID = state.AwayTeamID
	}
	winningTeamCount := 0
	totalCount := 0
	for _, group := range types {
		for _, line := range group {
			for _, p := range line.Players {
				wonGame := (p.TeamID == uint16(state.HomeTeamID) && state.HomeTeamWin) || (p.TeamID == uint16(state.AwayTeamID) && state.AwayTeamWin)
				if state.IsCollegeGame {
					statsObj := makeCollegePlayerStatsObject(state.WeekID, state.GameID, 0, p.Stats)
					star := structs.ThreeStarsObj{GameID: state.GameID, PlayerID: p.ID, TeamID: uint(p.TeamID)}
					star.MapPoints(statsObj.BasePlayerStats, wonGame)
					threeStars = append(threeStars, star)
				} else {
					statsObj := makeProPlayerStatsObject(state.WeekID, state.GameID, 0, p.Stats)
					star := structs.ThreeStarsObj{GameID: state.GameID, PlayerID: p.ID, TeamID: uint(p.TeamID)}
					star.MapPoints(statsObj.BasePlayerStats, wonGame)
					threeStars = append(threeStars, star)
				}
			}
		}
	}

	sort.Slice(threeStars, func(i, j int) bool {
		return threeStars[i].Points > threeStars[j].Points
	})
	starOne := 0
	starTwo := 0
	starThree := 0
	for _, star := range threeStars {
		if starOne > 0 && starTwo > 0 && starThree > 0 {
			break
		}
		if winningTeamCount > 1 && star.TeamID == winningTeamID {
			continue
		}
		if starOne == 0 {
			starOne = int(star.PlayerID)
		} else if starTwo == 0 {
			starTwo = int(star.PlayerID)
		} else if starThree == 0 {
			starThree = int(star.PlayerID)
		}
		if star.TeamID == winningTeamID {
			winningTeamCount++
		}
		totalCount++
	}
	return structs.ThreeStars{
		StarOne:   uint(starOne),
		StarTwo:   uint(starTwo),
		StarThree: uint(starThree),
	}
}

func getAttendancePercent(wins, losses int) float64 {
	totalGames := wins + losses
	if totalGames < 4 {
		return util.GenerateFloatFromRange(0.90, 1.00)
	}

	winRate := float64(wins) / float64(totalGames)

	switch {
	case winRate >= 0.75:
		return util.GenerateFloatFromRange(0.95, 1.05)
	case winRate >= 0.5:
		return util.GenerateFloatFromRange(0.85, 0.94)
	case winRate >= 0.35:
		return util.GenerateFloatFromRange(0.65, 0.84)
	default:
		return util.GenerateFloatFromRange(0.4, 0.64)
	}
}

// TopN returns the top n standings (or fewer if not enough teams).
func CreateTestResultsDirectory() error {
	// Create test_results directory if it doesn't exist
	if _, err := os.Stat("test_results"); os.IsNotExist(err) {
		err := os.Mkdir("test_results", 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateAndRunTestGames(ts structs.Timestamp, db *gorm.DB) {
	fmt.Println("Generating random test games...")

	// Get all college teams for random matchups
	allTeams := GetAllCollegeTeams()
	if len(allTeams) < 2 {
		fmt.Println("Not enough teams to generate test games")
		return
	}

	teamMap := GetCollegeTeamMap()

	// Generate random games (between 5-15 games)
	numGames := 5 + rand.Intn(11) // Random number between 5-15
	testGames := make([]structs.CollegeGame, 0, numGames)

	for i := 0; i < numGames; i++ {
		// Pick two random teams
		homeIdx := rand.Intn(len(allTeams))
		awayIdx := rand.Intn(len(allTeams))

		// Ensure home and away teams are different
		for awayIdx == homeIdx {
			awayIdx = rand.Intn(len(allTeams))
		}

		homeTeam := allTeams[homeIdx]
		awayTeam := allTeams[awayIdx]

		// Generate a test game
		game := generateCollegeGame(
			ts.SeasonID,
			ts.WeekID,
			ts.Week,
			homeTeam.ID,
			awayTeam.ID,
			ts.GetGameDay(),
			fmt.Sprintf("Test Game %d", i+1),
			teamMap,
			ts.IsPreseason,
		)

		testGames = append(testGames, game)
	}

	fmt.Printf("Generated %d test games\n", len(testGames))

	// Run the games without database operations
	runTestGamesOnly(testGames, ts, db)
}

// runTestGamesOnly runs games in testing mode with comprehensive error handling and debugging.
// Features enhanced crash debugging to identify nil pointer dereference locations:
//   - Detailed panic recovery in PrepareGames goroutines with step-by-step validation
//   - Comprehensive data validation for rosters, lineups, and arena data
//   - Stack trace capture for critical errors
//   - Safe game execution with isolated error handling
//   - CSV export only (no database operations)
//
// This function helps identify "where the invalid memory address occurred" by:
//  1. Adding panic recovery at every critical function call
//  2. Validating data availability before use
//  3. Providing detailed error messages with context
//  4. Capturing full stack traces for crashes
func runTestGamesOnly(collegeGames []structs.CollegeGame, ts structs.Timestamp, db *gorm.DB) {
	fmt.Println("Running test games (CSV export only)...")

	// Ensure test_results directory exists
	err := CreateTestResultsDirectory()
	if err != nil {
		fmt.Printf("Warning: Could not create test_results directory: %v\n", err)
	}

	// For testing mode, we need to check if we have the required data or create minimal fallbacks
	if !hasRequiredTestData() {
		fmt.Println("Warning: Missing required game data for testing. Creating minimal test games...")
		gameDTOs := createMinimalTestGameDTOs(collegeGames, ts)
		runMinimalTestGames(gameDTOs, ts)
		return
	}

	// Create empty standings maps for testing
	collegeStandingsMap := make(map[uint]structs.CollegeStandings)
	proStandingsMap := make(map[uint]structs.ProfessionalStandings)

	// Add default standings for all teams in the games to prevent nil pointer issues
	for _, game := range collegeGames {
		if _, exists := collegeStandingsMap[game.HomeTeamID]; !exists {
			collegeStandingsMap[game.HomeTeamID] = structs.CollegeStandings{
				BaseStandings: structs.BaseStandings{
					TeamID:      game.HomeTeamID,
					TotalWins:   5, // Default values for attendance calculation
					TotalOTWins: 2,
					TotalLosses: 3,
				},
			}
		}
		if _, exists := collegeStandingsMap[game.AwayTeamID]; !exists {
			collegeStandingsMap[game.AwayTeamID] = structs.CollegeStandings{
				BaseStandings: structs.BaseStandings{
					TeamID:      game.AwayTeamID,
					TotalWins:   5,
					TotalOTWins: 2,
					TotalLosses: 3,
				},
			}
		}
	}

	// Prepare games (no pro games for testing)
	proGames := []structs.ProfessionalGame{}

	// Try to prepare games with enhanced error handling
	var gameDTOs []structs.GameDTO

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("FATAL ERROR preparing test games: %v\n", r)
				fmt.Println("This indicates critical missing data. Check:")
				fmt.Println("  1. All teams in the games have players in the database")
				fmt.Println("  2. All teams have lineups configured (forward, defender, goalie lines)")
				fmt.Println("  3. All teams have shootout lineups set")
				fmt.Println("  4. All games reference valid arena IDs")
				fmt.Println("")
				fmt.Println("Game preparation failed. No games will be run.")
				gameDTOs = []structs.GameDTO{} // Return empty slice on error
			}
		}()

		fmt.Printf("Preparing %d games for testing...\n", len(collegeGames))
		gameDTOs = PrepareGames(collegeGames, proGames, collegeStandingsMap, proStandingsMap)
		fmt.Printf("Successfully prepared %d games.\n", len(gameDTOs))
	}()

	// Exit early if game preparation failed
	if len(gameDTOs) == 0 {
		fmt.Println("No games to run due to preparation errors.")
		return
	} // RUN THE GAMES with enhanced error handling!
	var results []engine.GameState

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("CRITICAL ERROR during game execution: %v\n", r)
				fmt.Println("This indicates a fundamental problem with the game engine or data.")
				fmt.Printf("Successfully prepared %d games, but execution failed.\n", len(gameDTOs))
				fmt.Println("\nStack trace:")
				fmt.Println(getStackTrace())

				// Show which games were prepared successfully
				fmt.Println("\nPrepared games:")
				for i, game := range gameDTOs {
					homeValid := game.HomeStrategy.CollegeRoster != nil && len(game.HomeStrategy.Forwards) > 0
					awayValid := game.AwayStrategy.CollegeRoster != nil && len(game.AwayStrategy.Forwards) > 0
					fmt.Printf("  %d. %s vs %s (Home valid: %t, Away valid: %t)\n",
						i+1,
						game.GameInfo.HomeTeam,
						game.GameInfo.AwayTeam,
						homeValid,
						awayValid)
				}
				results = []engine.GameState{} // Ensure empty results on crash
			}
		}()

		fmt.Printf("Executing %d prepared games...\n", len(gameDTOs))
		results = engine.RunGames(gameDTOs)
		fmt.Printf("Successfully completed %d games.\n", len(results))
	}()

	// Exit if no results (due to crash or other error)
	if len(results) == 0 {
		fmt.Println("No game results to export.")
		return
	}

	// Export results to CSV only (no database operations)
	allPlayers := GetAllCollegePlayers()
	collegePlayerMap := MakeCollegePlayerMap(allPlayers)
	collegeTeamMap := GetCollegeTeamMap()

	testNumString := strconv.Itoa(int(ts.TestRunNum))

	testFilePath := fmt.Sprintf("test_results/test_run_%s", testNumString)

	// Create the test run directory
	osErr := os.MkdirAll(testFilePath, 0755)
	if osErr != nil {
		fmt.Printf("Error creating test run directory: %v\n", osErr)
		return
	}

	fmt.Printf("Created test run directory: %s\n", testFilePath)

	for i, r := range results {
		// Generate unique filenames for each game within the test run directory
		filename := fmt.Sprintf("%s/game_%d_%s_vs_%s", testFilePath, i+1, r.HomeTeam, r.AwayTeam)

		// Write box score CSV
		boxScoreFilename := filename + "_box_score.csv"
		err := WriteBoxScoreFile(r, boxScoreFilename)
		if err != nil {
			fmt.Printf("Error writing box score file: %v\n", err)
		}

		// Write play-by-play CSV
		pbpFilename := filename + "_play_by_play.csv"
		err = WritePlayByPlayCSVFile(r.Collector.PlayByPlays, pbpFilename, collegePlayerMap, collegeTeamMap)
		if err != nil {
			fmt.Printf("Error writing play-by-play file: %v\n", err)
		}
	}

	ts.IncrementTestNumber()
	repository.SaveTimestamp(ts, db)

	fmt.Printf("Completed %d test games and exported to CSV files\n", len(results))
}

// hasRequiredTestData checks if the necessary data exists for full game testing
func hasRequiredTestData() bool {
	// Check if we have college players and lineups in the system
	allPlayers := GetAllCollegePlayers()
	if len(allPlayers) == 0 {
		return false
	}

	lineupMap := GetCollegeLineupsMap()
	if len(lineupMap) == 0 {
		return false
	}

	return true
}

// createMinimalTestGameDTOs provides a fallback when full data isn't available
func createMinimalTestGameDTOs(collegeGames []structs.CollegeGame, ts structs.Timestamp) []structs.GameDTO {
	fmt.Println("ERROR: Cannot run test games without database setup.")
	fmt.Println("Testing mode requires:")
	fmt.Println("  1. College teams with players in the database")
	fmt.Println("  2. Team lineups configured")
	fmt.Println("  3. Shootout lineups set up")
	fmt.Println("")
	fmt.Println("To set up testing environment:")
	fmt.Println("  1. Initialize the database with teams and players")
	fmt.Println("  2. Run lineup generation: RunLineupsForAICollegeTeams()")
	fmt.Println("  3. Ensure all teams have complete rosters and lineups")
	fmt.Println("")
	fmt.Println("Alternatively, use IsTesting=false to test with existing game data.")
	return []structs.GameDTO{}
}

// runMinimalTestGames provides error information when testing can't proceed
func runMinimalTestGames(gameDTOs []structs.GameDTO, ts structs.Timestamp) {
	fmt.Println("Test games cannot be run due to missing database setup.")
	fmt.Println("Please initialize the application with proper team and player data first.")
}

// getStackTrace captures and formats the current stack trace for debugging
func getStackTrace() string {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, false)
	return string(buf[:stackSize])
}

func TopN(ss []*structs.CollegeStandings, n int) []*structs.CollegeStandings {
	if len(ss) < n {
		n = len(ss)
	}
	return ss[:n]
}

// Quarter pairs: given a sorted slice (seed #1 at index 0),
// return pairs [ (0,last), (1,last-1), ... ] of length count.
func SeededPairs(ss []*structs.CollegeStandings, count int) [][2]*structs.CollegeStandings {
	pairs := make([][2]*structs.CollegeStandings, 0, count)
	top, bottom := 0, len(ss)-1
	for i := 0; i < count && top < bottom; i++ {
		pairs = append(pairs, [2]*structs.CollegeStandings{ss[top], ss[bottom]})
		top++
		bottom--
	}
	return pairs
}
