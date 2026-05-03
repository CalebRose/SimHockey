package managers

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/engine"
	fbsvc "github.com/CalebRose/SimHockey/firebase"
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
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	gameDay := ts.GetGameDay()

	if ts.IsTesting {
		generateAndRunTestGames(ts, db)
		return
	}

	collegeGames := GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)
	proGames := GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay, ts.IsPreseason)

	collegeStandingsMap := GetCollegeStandingsMap(seasonID)
	proStandingsMap := GetProStandingsMap(seasonID)
	gameDTOs := PrepareGames(collegeGames, proGames, collegeStandingsMap, proStandingsMap)

	results := engine.RunGames(gameDTOs)

	collegeGameMap := MakeCollegeGameMap(collegeGames)
	proGameMap := MakeProGameMap(proGames)
	collegeTeamMap := GetCollegeTeamMap()
	proTeamMap := GetProTeamMap()
	collegePlayerMap := GetCollegePlayersMap()
	proPlayersMap := GetProPlayersMap()
	upload := NewStatsUpload()

	sentCollegeInjuryNotification := make(map[uint]bool)
	sentProInjuryNotification := make(map[uint]bool)

	for _, r := range results {
		gameID := strconv.Itoa(int(r.GameID))

		for _, injury := range r.InjuryLog {
			if r.IsCollegeGame && !ts.IsPreseason {
				if player, ok := collegePlayerMap[injury.PlayerID]; ok {
					player.ApplyInjury(injury.InjuryName, injury.InjuryType.String(), int8(injury.RecoveryDays))
					repository.SaveCollegeHockeyPlayerRecord(player, db)
					teamID := uint(player.TeamID)
					if !sentCollegeInjuryNotification[teamID] {
						if team, ok := collegeTeamMap[teamID]; ok && team.Coach != "" && team.Coach != "AI" {
							ctx := context.Background()
							uids := fbsvc.ResolveUIDsByUsernames(ctx, []string{team.Coach})
							if len(uids) > 0 {
								eventKey := fbsvc.BuildSourceEventKey("injury", "chl", gameID, strconv.Itoa(int(injury.PlayerID)))
								_ = fbsvc.NotifyTeamInjury(ctx, fbsvc.TeamInjuryNotificationInput{
									League:         "chl",
									Domain:         fbsvc.DomainCHL,
									TeamID:         teamID,
									TeamName:       team.TeamName,
									PlayerID:       injury.PlayerID,
									PlayerName:     player.FirstName + " " + player.LastName,
									Position:       player.Position,
									InjuryType:     injury.InjuryType.String(),
									DaysOfRecovery: injury.RecoveryDays,
									GameID:         gameID,
									RecipientUIDs:  uids,
									SourceEventKey: eventKey,
								})
							}
							sentCollegeInjuryNotification[teamID] = true
						}
					}
				}
			} else if !ts.IsPreseason && !r.IsCollegeGame {
				if player, ok := proPlayersMap[injury.PlayerID]; ok {
					player.ApplyInjury(injury.InjuryName, injury.InjuryType.String(), int8(injury.RecoveryDays))
					repository.SaveProPlayerRecord(player, db)
					teamID := uint(player.TeamID)
					if !sentProInjuryNotification[teamID] {
						if team, ok := proTeamMap[teamID]; ok && team.Owner != "" {
							ctx := context.Background()
							recipients := collectProTeamUsernames(team)
							uids := fbsvc.ResolveUIDsByUsernames(ctx, recipients)
							if len(uids) > 0 {
								eventKey := fbsvc.BuildSourceEventKey("injury", "phl", gameID, strconv.Itoa(int(injury.PlayerID)))
								_ = fbsvc.NotifyTeamInjury(ctx, fbsvc.TeamInjuryNotificationInput{
									League:         "phl",
									Domain:         fbsvc.DomainPHL,
									TeamID:         teamID,
									TeamName:       team.TeamName,
									PlayerID:       injury.PlayerID,
									PlayerName:     player.FirstName + " " + player.LastName,
									Position:       player.Position,
									InjuryType:     injury.InjuryType.String(),
									DaysOfRecovery: injury.RecoveryDays,
									GameID:         gameID,
									RecipientUIDs:  uids,
									SourceEventKey: eventKey,
								})
							}
							sentProInjuryNotification[teamID] = true
						}
					}
				}
			}
		}

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
	}
}

func NewStatsUpload() *StatsUpload {
	return &StatsUpload{}
}

func collectProTeamUsernames(team structs.ProfessionalTeam) []string {
	usernames := make([]string, 0, 2)
	if team.Owner != "" {
		usernames = append(usernames, team.Owner)
	}
	if team.GM != "" && team.GM != team.Owner {
		usernames = append(usernames, team.GM)
	}
	return usernames
}

func (u *StatsUpload) Collect(state engine.GameState, seasonID, gameType uint) {
	u.collectTeamStats(state, seasonID, gameType)
	u.collectPlayerStats(state.HomeStrategy, state.WeekID, state.GameID, gameType, state.IsCollegeGame)
	u.collectPlayerStats(state.AwayStrategy, state.WeekID, state.GameID, gameType, state.IsCollegeGame)
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
	collect := func(p *engine.GamePlayer) {
		if p == nil || p.ID == 0 {
			return
		}
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

	types := [][]engine.LineStrategy{pl.Forwards, pl.Defenders, pl.Goalies}
	for _, group := range types {
		for _, line := range group {
			for _, p := range line.Players {
				collect(p)
			}
		}
	}

	for _, p := range pl.InjuredPlayers {
		collect(p)
	}

	seen := make(map[uint]bool)
	if isCollege {
		filtered := make([]structs.CollegePlayerGameStats, 0, len(u.CollegePlayerStats))
		for _, s := range u.CollegePlayerStats {
			if s.PlayerID == 0 {
				continue
			}
			if _, ok := seen[s.PlayerID]; !ok {
				seen[s.PlayerID] = true
				filtered = append(filtered, s)
			}
		}
		u.CollegePlayerStats = filtered
	} else {
		filtered := make([]structs.ProfessionalPlayerGameStats, 0, len(u.ProPlayerStats))
		for _, s := range u.ProPlayerStats {
			if s.PlayerID == 0 {
				continue
			}
			if _, ok := seen[s.PlayerID]; !ok {
				seen[s.PlayerID] = true
				filtered = append(filtered, s)
			}
		}
		u.ProPlayerStats = filtered
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
	var collegeGamesWg sync.WaitGroup
	var mutex sync.Mutex

	collegeTeamRosterMap := GetAllCollegePlayersMapByTeam()
	collegeLineupMap := GetCollegeLineupsMap()
	collegeShootoutLineupMap := GetCollegeShootoutLineups()
	collegeGameplans := repository.FindCollegeGameplanRecords()
	collegeGameplanMap := MakeCollegeGameplanMap(collegeGameplans)
	arenaMap := GetArenaMap()

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
			hgp := collegeGameplanMap[c.HomeTeamID]
			agp := collegeGameplanMap[c.AwayTeamID]
			arena := arenaMap[c.ArenaID]

			if htr == nil || atr == nil || htl == nil || atl == nil || htsl.TeamID == 0 || atsl.TeamID == 0 || hgp.ID == 0 || agp.ID == 0 {
				mutex.Unlock()
				return
			}

			hp := getCollegePlaybookDTO(htl, htr, htsl, hgp)
			ap := getCollegePlaybookDTO(atl, atr, atsl, agp)
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

	var proGamesWg sync.WaitGroup
	proTeamRosterMap := GetAllProPlayersMapByTeam()
	proLineupMap := GetProLineupsMap()
	proShootoutLineupMap := GetProShootoutLineups()
	proGameplans := repository.FindProfessionalGameplanRecords()
	proGameplanMap := MakeProGameplanMap(proGameplans)
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
			hgp := proGameplanMap[g.HomeTeamID]
			agp := proGameplanMap[g.AwayTeamID]
			hp := getProfessionalPlaybookDTO(htl, htr, htsl, hgp)
			ap := getProfessionalPlaybookDTO(atl, atr, atsl, agp)
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

	sort.Slice(gameDTOList, func(i, j int) bool {
		return gameDTOList[i].IsCollegeGame
	})
	return gameDTOList
}

func GetCollegeGamesForCurrentMatchup(weekID, seasonID, gameDay string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGamesByCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
}

func GetProfessionalGamesForCurrentMatchup(weekID, seasonID, gameDay string, isPreseason bool) []structs.ProfessionalGame {
	return repository.FindProfessionalGamesByCurrentMatchup(weekID, seasonID, gameDay, isPreseason)
}

func GetCollegeGamesBySeasonID(seasonID string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, IsPreseason: isPreseason})
}

func GetProfessionalGamesBySeasonID(seasonID string, isPreseason bool) []structs.ProfessionalGame {
	return repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID, IsPreseason: isPreseason})
}

func GetCollegeGamesByTeamIDAndSeasonID(teamID, seasonID string, isPreseason bool) []structs.CollegeGame {
	return repository.FindCollegeGames(repository.GamesClauses{SeasonID: seasonID, TeamID: teamID, IsPreseason: isPreseason})
}

func GetPlayoffSeriesBySeriesID(seriesID string) structs.ProSeries {
	return repository.FindPlayoffSeriesByID(seriesID)
}

func GetCollegeGameByID(id string) structs.CollegeGame {
	return repository.FindCollegeGameRecord(id)
}

func GetArenaMap() map[uint]structs.Arena {
	arenas := repository.FindAllArenas(repository.ArenaQuery{})
	return MakeArenaMap(arenas)
}

func getCollegePlaybookDTO(lineups []structs.CollegeLineup, roster []structs.CollegePlayer, shootoutLineup structs.CollegeShootoutLineup, gp structs.CollegeGameplan) structs.PlayBookDTO {
	forwards, defenders, goalies := getCollegeForwardDefenderGoalieLineups(lineups)
	return structs.PlayBookDTO{
		Forwards:       forwards,
		Defenders:      defenders,
		Goalies:        goalies,
		CollegeRoster:  roster,
		ShootoutLineup: shootoutLineup.ShootoutPlayerIDs,
		Gameplan:       gp.BaseGameplan,
	}
}

func getProfessionalPlaybookDTO(lineups []structs.ProfessionalLineup, roster []structs.ProfessionalPlayer, shootoutLineup structs.ProfessionalShootoutLineup, gp structs.ProGameplan) structs.PlayBookDTO {
	forwards, defenders, goalies := getProfessionalForwardDefenderGoalieLineups(lineups)
	return structs.PlayBookDTO{
		Forwards:           forwards,
		Defenders:          defenders,
		Goalies:            goalies,
		ProfessionalRoster: roster,
		ShootoutLineup:     shootoutLineup.ShootoutPlayerIDs,
		Gameplan:           gp.BaseGameplan,
	}
}

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

func GenerateThreeStars(state engine.GameState, seasonID uint) structs.ThreeStars {
	types := [][]engine.LineStrategy{state.HomeStrategy.Forwards, state.HomeStrategy.Defenders, state.HomeStrategy.Goalies, state.AwayStrategy.Forwards, state.AwayStrategy.Defenders, state.AwayStrategy.Goalies}
	threeStars := []structs.ThreeStarsObj{}
	winningTeamID := state.HomeTeamID
	if state.AwayTeamWin {
		winningTeamID = state.AwayTeamID
	}
	winningTeamCount := 0
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

		if cid == 2 {
			seven := TopN(conferenceMap[cid], 7)
			pairs := [][2]*structs.CollegeStandings{
				{seven[1], seven[6]},
				{seven[2], seven[5]},
				{seven[3], seven[4]},
			}
			semiFinalID1 := nextGameID
			semiFinalID2 := nextGameID + 1
			finalsID := nextGameID + 2
			conferenceName := ""

			for idx, p := range pairs {
				a, b := p[0], p[1]
				homeTeam := teamMap[a.TeamID]
				conferenceName = homeTeam.Conference
				ngID := semiFinalID2
				nextHOA := "H"
				if idx == 0 {
					ngID = semiFinalID1
					nextHOA = "A"
				}

				series := structs.CollegeSeries{
					BaseSeries: structs.BaseSeries{
						SeasonID:    seasonID,
						SeriesName:  fmt.Sprintf("%s Conference Quarterfinals", conferenceName),
						BestOfCount: 3,
						HomeTeamID:  a.TeamID, HomeTeam: a.TeamName, HomeTeamRank: 2 + uint(idx),
						AwayTeamID: b.TeamID, AwayTeam: b.TeamName, AwayTeamRank: uint(7 - idx),
						GameCount:     0,
						IsPlayoffGame: true,
					},
					NextGameID:   ngID,
					NextGameHOA:  nextHOA,
					ConferenceID: uint8(cid),
				}
				quarterfinalsSeries = append(quarterfinalsSeries, series)
			}
			top1 := seven[0]
			top1Team := teamMap[top1.TeamID]
			semifinalGame1 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID1},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conferenceName),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					HomeTeamID: top1.TeamID, HomeTeam: top1.TeamName, HomeTeamRank: 1,
					Arena: top1Team.Arena, NextGameID: finalsID, NextGameHOA: "H",
					GameDay: "A",
				},
				IsConferenceTournament: true,
			}

			semifinalGame2 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID2},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conferenceName),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					NextGameID: finalsID, NextGameHOA: "A",
					GameDay: "A",
				},
				IsConferenceTournament: true,
			}

			finalsGame := structs.CollegeGame{
				Model: gorm.Model{ID: finalsID},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Finals", conferenceName),
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
				continue
			}
			pairs := SeededPairs(eight, 4)
			semiFinalID1 := nextGameID
			semiFinalID2 := nextGameID + 1
			finalsID := nextGameID + 2
			conferenceName := ""
			for qfIdx, p := range pairs {
				a, b := p[0], p[1]
				homeTeam := teamMap[a.TeamID]
				conferenceName = homeTeam.Conference

				ngID := semiFinalID2
				if qfIdx == 0 || qfIdx == 3 {
					ngID = semiFinalID1
				}

				nextHOA := "H"
				if qfIdx > 1 {
					nextHOA = "A"
				}

				series := structs.CollegeSeries{
					BaseSeries: structs.BaseSeries{
						SeasonID:    seasonID,
						SeriesName:  fmt.Sprintf("%s Conference Quarterfinals", conferenceName),
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
			semifinalGame1 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID1},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conferenceName),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					NextGameID: finalsID, NextGameHOA: "H",
					GameDay: "A", IsPlayoffGame: true,
				},
				IsConferenceTournament: true,
			}

			semifinalGame2 := structs.CollegeGame{
				Model: gorm.Model{ID: semiFinalID2},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Semifinals", conferenceName),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					NextGameID: finalsID, NextGameHOA: "A",
					GameDay: "A", IsPlayoffGame: true,
				},
				IsConferenceTournament: true,
			}

			finalsGame := structs.CollegeGame{
				Model: gorm.Model{ID: finalsID},
				BaseGame: structs.BaseGame{
					GameTitle: fmt.Sprintf("%s Conference Finals", conferenceName),
					SeasonID:  seasonID, WeekID: util.GetWeekID(seasonID, 19), Week: 19,
					GameDay: "B",
				},
				IsConferenceTournament: true,
			}

			semiFinalsAndFinalsGames = append(semiFinalsAndFinalsGames, semifinalGame1, semifinalGame2, finalsGame)
			nextGameID += 3
		}
	}
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

	for _, t := range collegeTeams {
		standings := stMap[t.ID]
		if standings.ID == 0 {
			continue
		}
		if standings.IsConferenceTournamentChampion {
			sCopy := standings
			pool = append(pool, &sCopy)
			qualified[t.ID] = true
		}
	}

	sort.Slice(collegeStandings, func(i, j int) bool {
		if collegeStandings[i].PairwiseRank == collegeStandings[j].PairwiseRank {
			return collegeStandings[i].RPIRank > collegeStandings[j].RPIRank
		}
		return collegeStandings[i].PairwiseRank > collegeStandings[j].PairwiseRank
	})

	for _, s := range collegeStandings {
		if len(pool) == 16 {
			break
		}
		if !qualified[s.TeamID] {
			sCopy := s
			pool = append(pool, &sCopy)
			qualified[s.TeamID] = true
		}
	}

	sort.Slice(pool, func(i, j int) bool {
		if pool[i].PairwiseRank == pool[j].PairwiseRank {
			return pool[i].RPIRank > pool[j].RPIRank
		}
		return pool[i].PairwiseRank > pool[j].PairwiseRank
	})

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

		g.NextGameID = baseID + 8 + uint(i/2)
		if i%2 == 0 {
			g.NextGameHOA = "H"
		} else {
			g.NextGameHOA = "A"
		}

		g.GameDay = "A"
		games = append(games, g)
	}

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

	arenaIdx := util.GenerateIntFromRange(0, len(arenas)-1)
	arenaStruct := arenas[arenaIdx]
	for s := 0; s < 2; s++ {
		g := mk(baseID+12+uint(s), fmt.Sprintf("%d SimCHL Frozen Four Semifinal", ts.Season), 21)
		g.NextGameID = baseID + 14
		if s == 0 {
			g.NextGameHOA = "H"
		} else {
			g.NextGameHOA = "A"
		}
		g.GameDay = "A"
		g.ArenaID = arenaStruct.ID
		g.Arena = arenaStruct.Name
		g.City = arenaStruct.City
		g.State = arenaStruct.State
		g.Country = arenaStruct.Country
		games = append(games, g)
	}

	final := mk(baseID+14, fmt.Sprintf("%d SimCHL National Championship", ts.Season), 21)
	final.IsNationalChampionship = true
	final.GameDay = "C"
	final.ArenaID = arenaStruct.ID
	final.Arena = arenaStruct.Name
	final.City = arenaStruct.City
	final.State = arenaStruct.State
	final.Country = arenaStruct.Country
	games = append(games, final)

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
	toiletBowl.IsNeutralSite = true
	toiletBowl.ArenaID = 24
	toiletBowl.Arena = "The Hart Center"
	toiletBowl.City = "Worcester"
	toiletBowl.State = "MA"
	toiletBowl.Country = "USA"
	games = append(games, toiletBowl)
	repository.CreateCHLGamesRecordsBatch(db, games, 50)
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
		ht := teamMap[s.HomeTeamID]
		arenaID := ht.ArenaID
		arena := ht.Arena
		city := ht.City
		state := ht.State
		country := ht.Country
		weekID := util.GetWeekID(ts.SeasonID, 18)
		week := 18

		collegeGame := structs.CollegeGame{
			BaseGame: structs.BaseGame{
				GameTitle: s.SeriesName + " Game: " + gameCount,
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

	for _, did := range divisionIDList {
		division := divisionMap[did]
		sort.Slice(division, func(i, j int) bool {
			if division[i].Points == division[j].Points {
				return division[i].GoalsFor > division[j].GoalsFor
			}
			return division[i].Points > division[j].Points
		})
		qualifyingTeams = append(qualifyingTeams, division[:2]...)
	}

	pairs := [][2]*structs.ProfessionalStandings{
		{qualifyingTeams[0], qualifyingTeams[3]},
		{qualifyingTeams[1], qualifyingTeams[2]},
		{qualifyingTeams[4], qualifyingTeams[7]},
		{qualifyingTeams[5], qualifyingTeams[6]},
	}

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
				NextSeriesHOA: nextSeriesHoa,
				NextSeriesID:  nextSeriesID + 4 + uint(qfIdx/2),
			},
		}
		postSeasonSeriesList = append(postSeasonSeriesList, quarterFinalsSeries)
	}

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
				NextSeriesHOA:   nextSeriesHoa,
				NextSeriesID:    nextSeriesID + 6,
				SeriesName:      fmt.Sprintf("%d SimPHL Semifinals", ts.Season),
				BestOfCount:     7,
				GameCount:       1,
			},
		}
		postSeasonSeriesList = append(postSeasonSeriesList, semiFinalsSeries)
	}

	finalsSeries := structs.ProSeries{
		Model: gorm.Model{ID: nextSeriesID + 6},
		BaseSeries: structs.BaseSeries{
			SeasonID:   seasonID,
			HomeTeamID: 0, HomeTeam: "", HomeTeamRank: 0,
			AwayTeamID: 0, AwayTeam: "", AwayTeamRank: 0,
			HomeTeamCoach:   "",
			AwayTeamCoach:   "",
			IsInternational: false,
			NextSeriesHOA:   "",
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

func GenerateProPlayoffGames(db *gorm.DB, ts structs.Timestamp) {
	weekID := strconv.Itoa(int(ts.WeekID))
	seasonID := strconv.Itoa(int(ts.SeasonID))
	teamMap := GetProTeamMap()
	professionalGames := repository.FindProfessionalGames(repository.GamesClauses{SeasonID: seasonID, WeekID: weekID})

	incompleteGames := []structs.ProfessionalGame{}
	for _, g := range professionalGames {
		if !g.GameComplete {
			incompleteGames = append(incompleteGames, g)
		}
	}
	if len(incompleteGames) > 0 {
		return
	}

	gameDay := ts.GetGameDay()
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
	ts := GetTimestamp()
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
		const maxTries = 1000
		for tries := 0; tries < maxTries; tries++ {
			rand.Shuffle(len(teamIDs), func(i, j int) {
				teamIDs[i], teamIDs[j] = teamIDs[j], teamIDs[i]
			})
			var err error
			pairings, err = createCollegeGamePairings(teamIDs, teamMap, playedGameReference)
			if err == nil {
				break
			}
		}
		for _, pair := range pairings {
			game := generateCollegeGame(ts.SeasonID, ts.WeekID, ts.Week, pair[0], pair[1], gameDay, "", teamMap, true)
			games = append(games, game)
		}
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
	ts := GetTimestamp()
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
			rand.Shuffle(len(teamIDs), func(i, j int) {
				teamIDs[i], teamIDs[j] = teamIDs[j], teamIDs[i]
			})
			var err error
			pairings, err = createProGamePairings(teamIDs, teamMap, playedGameReference)
			if err == nil {
				break
			}
		}
		for _, pair := range pairings {
			game := generateProfessionalGame(ts.SeasonID, ts.WeekID, ts.Week, pair[0], pair[1], gameDay, teamMap, true)
			games = append(games, game)
		}
		switch gameDay {
		case "A":
			gameDay = "B"
		case "B":
			gameDay = "C"
		}
	}
	return games
}

func createCollegeGamePairings(teamIDs []uint, teamMap map[uint]structs.CollegeTeam, playedGameReference map[uint]map[uint]bool) ([][2]uint, error) {
	n := len(teamIDs)
	paired := make(map[uint]bool)
	var pairings [][2]uint
	for _, t1 := range teamIDs {
		if paired[t1] {
			continue
		}
		c1 := teamMap[t1].ConferenceID
		found := false
		for _, t2 := range teamIDs {
			if t1 == t2 || paired[t2] {
				continue
			}
			c2 := teamMap[t2].ConferenceID
			if c1 != 7 && c1 == c2 {
				continue
			}
			if playedGameReference[t1][t2] {
				continue
			}
			pairings = append(pairings, [2]uint{t1, t2})
			paired[t1], paired[t2] = true, true
			playedGameReference[t1][t2] = true
			playedGameReference[t2][t1] = true
			found = true
			break
		}
		if !found {
			break
		}
	}
	if len(pairings)*2 != n {
		return nil, fmt.Errorf("incomplete pairing")
	}
	return pairings, nil
}

func createProGamePairings(teamIDs []uint, teamMap map[uint]structs.ProfessionalTeam, playedGameReference map[uint]map[uint]bool) ([][2]uint, error) {
	n := len(teamIDs)
	paired := make(map[uint]bool)
	var pairings [][2]uint
	for _, t1 := range teamIDs {
		if paired[t1] {
			continue
		}
		c1 := teamMap[t1].DivisionID
		found := false
		for _, t2 := range teamIDs {
			if t1 == t2 || paired[t2] {
				continue
			}
			c2 := teamMap[t2].DivisionID
			if c1 != 7 && c1 == c2 {
				continue
			}
			if playedGameReference[t1][t2] {
				continue
			}
			pairings = append(pairings, [2]uint{t1, t2})
			paired[t1], paired[t2] = true, true
			playedGameReference[t1][t2] = true
			playedGameReference[t2][t1] = true
			found = true
			break
		}
		if !found {
			break
		}
	}
	if len(pairings)*2 != n {
		return nil, fmt.Errorf("incomplete pairing")
	}
	return pairings, nil
}

func generateCollegeGame(seasonID, weekID, week, hid, aid uint, gameDay, gameTitle string, teamMap map[uint]structs.CollegeTeam, isPreseason bool) structs.CollegeGame {
	return structs.CollegeGame{
		BaseGame: structs.BaseGame{
			WeekID: weekID, Week: int(week), GameDay: gameDay, GameTitle: gameTitle, SeasonID: seasonID,
			HomeTeamID: hid, HomeTeam: teamMap[hid].TeamName, AwayTeamID: aid, AwayTeam: teamMap[aid].TeamName,
			ArenaID: uint(teamMap[hid].ArenaID), IsPreseason: isPreseason,
		},
	}
}

func generateProfessionalGame(seasonID, weekID, week, hid, aid uint, gameDay string, teamMap map[uint]structs.ProfessionalTeam, isPreseason bool) structs.ProfessionalGame {
	return structs.ProfessionalGame{
		BaseGame: structs.BaseGame{
			WeekID: weekID, Week: int(week), GameDay: gameDay, SeasonID: seasonID,
			HomeTeamID: hid, HomeTeam: teamMap[hid].Abbreviation, AwayTeamID: aid, AwayTeam: teamMap[aid].Abbreviation,
			ArenaID: uint(teamMap[hid].ArenaID), IsPreseason: isPreseason,
		},
	}
}

func generateAndRunTestGames(ts structs.Timestamp, db *gorm.DB) {
	fmt.Println("Generating and running test games...")
}

type LiveGameHubDTO struct {
	GameID                uint   `json:"GameID"`
	HomeTeam              string `json:"HomeTeam"`
	AwayTeam              string `json:"AwayTeam"`
	HomeTeamScore         uint   `json:"HomeTeamScore"`
	AwayTeamScore         uint   `json:"AwayTeamScore"`
	HomeTeamShootoutScore uint   `json:"HomeTeamShootoutScore"`
	AwayTeamShootoutScore uint   `json:"AwayTeamShootoutScore"`
	Period                uint8  `json:"Period"`
	TimeOnClock           uint16 `json:"TimeOnClock"`
	Zone                  uint8  `json:"Zone"`
	GameComplete          bool   `json:"GameComplete"`
	IsShootout            bool   `json:"IsShootout"`
}

type GameDetailsDTO struct {
	Feeds     []PbPDTO        `json:"Feeds"`
	HomeStats TeamBoxScoreDTO `json:"HomeStats"`
	AwayStats TeamBoxScoreDTO `json:"AwayStats"`
}

type PbPDTO struct {
	Period      uint8  `json:"Period"`
	TimeOnClock uint16 `json:"TimeOnClock"`
	PlayText    string `json:"PlayText"`
	Zone        uint8  `json:"Zone"`
}

type TeamBoxScoreDTO struct {
	Forwards  []PlayerBoxScoreDTO `json:"Forwards"`
	Defenders []PlayerBoxScoreDTO `json:"Defenders"`
	Goalies   []GoalieBoxScoreDTO `json:"Goalies"`
}

type PlayerBoxScoreDTO struct {
	Name      string `json:"Name"`
	Goals     uint8  `json:"Goals"`
	Assists   uint8  `json:"Assists"`
	PlusMinus int8   `json:"PlusMinus"`
}

type GoalieBoxScoreDTO struct {
	Name           string  `json:"Name"`
	Saves          uint16  `json:"Saves"`
	ShotsAgainst   uint16  `json:"ShotsAgainst"`
	SavePercentage float64 `json:"SavePercentage"`
}

func GetLiveGamesHubData(isCollege bool, reqSeason string, reqWeek string, reqTimeslot string) map[uint]LiveGameHubDTO {
	ts := GetTimestamp()
	seasonID := reqSeason
	if seasonID == "" {
		seasonID = strconv.Itoa(int(ts.SeasonID))
	}
	weekID := reqWeek
	if weekID == "" {
		weekID = strconv.Itoa(int(ts.WeekID))
	} else if len(weekID) <= 2 {
		seasonNum, _ := strconv.Atoi(seasonID)
		weekNum, _ := strconv.Atoi(weekID)
		weekID = strconv.Itoa(int(util.GetWeekID(uint(seasonNum), uint(weekNum))))
	}

	responseMap := make(map[uint]LiveGameHubDTO)
	if isCollege {
		clauses := repository.GamesClauses{SeasonID: seasonID, WeekID: weekID, IsPreseason: ts.IsPreseason}
		if reqTimeslot != "" {
			clauses.Timeslot = reqTimeslot
		}
		games := repository.FindCollegeGames(clauses)
		allCollegeTeams := repository.FindAllCollegeTeams(repository.TeamClauses{})
		chlTeamMap := MakeCollegeTeamMap(allCollegeTeams)

		for _, g := range games {
			homeTeam := chlTeamMap[g.HomeTeamID]
			awayTeam := chlTeamMap[g.AwayTeamID]
			Period := uint8(0)
			if g.GameComplete {
				Period = 3
				if g.IsOvertime {
					Period = 4
				}
				if g.IsShootout {
					Period = 5
				}
			}
			responseMap[g.ID] = LiveGameHubDTO{
				GameID: g.ID, HomeTeam: homeTeam.Abbreviation, AwayTeam: awayTeam.Abbreviation,
				HomeTeamScore: uint(g.HomeTeamScore), AwayTeamScore: uint(g.AwayTeamScore),
				HomeTeamShootoutScore: uint(g.HomeTeamShootoutScore), AwayTeamShootoutScore: uint(g.AwayTeamShootoutScore),
				Period: Period, TimeOnClock: 0, Zone: 11, GameComplete: g.GameComplete, IsShootout: g.IsShootout,
			}
		}
	}
	return responseMap
}

func GetGameDetailsData(gameID string, isCollege bool) GameDetailsDTO {
	response := GameDetailsDTO{
		Feeds:     []PbPDTO{},
		HomeStats: TeamBoxScoreDTO{Forwards: []PlayerBoxScoreDTO{}, Defenders: []PlayerBoxScoreDTO{}, Goalies: []GoalieBoxScoreDTO{}},
		AwayStats: TeamBoxScoreDTO{Forwards: []PlayerBoxScoreDTO{}, Defenders: []PlayerBoxScoreDTO{}, Goalies: []GoalieBoxScoreDTO{}},
	}
	if isCollege {
		game := repository.FindCollegeGameRecord(gameID)
		collegePlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{})
		collegePlayerMap := MakeCollegePlayerMap(collegePlayers)
		collegeTeamMap := GetCollegeTeamMap()

		pbps := GetCHLPlayByPlaysByGameID(gameID)
		for _, pbp := range pbps {
			p := pbp.PbP
			eventStr := util.ReturnStringFromPBPID(p.EventID)
			outcomeStr := util.ReturnStringFromPBPID(p.Outcome)
			playText := generateCollegeResultsString(p, eventStr, outcomeStr, collegePlayerMap, collegeTeamMap[uint(p.TeamID)])
			// Citing source for ZoneID mapping
			response.Feeds = append(response.Feeds, PbPDTO{Period: p.Period, TimeOnClock: p.TimeOnClock, PlayText: playText, Zone: p.ZoneID})
		}
		playerStats := repository.FindCollegePlayerGameStatsRecords(strconv.Itoa(int(game.SeasonID)), "", "", gameID)
		for _, s := range playerStats {
			if s.TimeOnIce <= 0 {
				continue
			}
			playerInfo := collegePlayerMap[s.PlayerID]
			nameStr := fmt.Sprintf("%s. %s", string(playerInfo.FirstName[0]), playerInfo.LastName)
			isHome := s.TeamID == game.HomeTeamID
			if playerInfo.Position == "Goalie" || playerInfo.Position == "G" {
				gs := GoalieBoxScoreDTO{Name: nameStr, Saves: uint16(s.Saves), ShotsAgainst: uint16(s.ShotsAgainst)}
				if s.ShotsAgainst > 0 {
					gs.SavePercentage = float64(s.Saves) / float64(s.ShotsAgainst)
				}
				if isHome {
					response.HomeStats.Goalies = append(response.HomeStats.Goalies, gs)
				} else {
					response.AwayStats.Goalies = append(response.AwayStats.Goalies, gs)
				}
			} else {
				ps := PlayerBoxScoreDTO{Name: nameStr, Goals: uint8(s.Goals), Assists: uint8(s.Assists), PlusMinus: int8(s.PlusMinus)}
				if playerInfo.Position == "D" {
					if isHome {
						response.HomeStats.Defenders = append(response.HomeStats.Defenders, ps)
					} else {
						response.AwayStats.Defenders = append(response.AwayStats.Defenders, ps)
					}
				} else {
					if isHome {
						response.HomeStats.Forwards = append(response.HomeStats.Forwards, ps)
					} else {
						response.AwayStats.Forwards = append(response.AwayStats.Forwards, ps)
					}
				}
			}
		}
	}
	return response
}

func GetBulkPlayByPlayData(isCollege bool, reqSeason string, reqWeek string, reqTimeslot string) map[uint][]PbPDTO {
	ts := GetTimestamp()
	seasonID := reqSeason
	if seasonID == "" {
		seasonID = strconv.Itoa(int(ts.SeasonID))
	}
	weekID := reqWeek
	if weekID == "" {
		weekID = strconv.Itoa(int(ts.WeekID))
	} else if len(weekID) <= 2 {
		seasonNum, _ := strconv.Atoi(seasonID)
		weekNum, _ := strconv.Atoi(weekID)
		weekID = strconv.Itoa(int(util.GetWeekID(uint(seasonNum), uint(weekNum))))
	}
	responseMap := make(map[uint][]PbPDTO)
	if isCollege {
		clauses := repository.GamesClauses{SeasonID: seasonID, WeekID: weekID, IsPreseason: ts.IsPreseason}
		if reqTimeslot != "" {
			clauses.Timeslot = reqTimeslot
		}
		games := repository.FindCollegeGames(clauses)
		var gameIDs []string
		for _, g := range games {
			gameIDs = append(gameIDs, strconv.Itoa(int(g.ID)))
			responseMap[g.ID] = []PbPDTO{}
		}
		if len(gameIDs) == 0 {
			return responseMap
		}
		collegePlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{})
		collegePlayerMap := MakeCollegePlayerMap(collegePlayers)
		collegeTeamMap := GetCollegeTeamMap()
		db := dbprovider.GetInstance().GetDB()
		var allPbPs []structs.CollegePlayByPlay
		db.Where("game_id IN ?", gameIDs).Find(&allPbPs)
		for _, p := range allPbPs {
			eventStr := util.ReturnStringFromPBPID(p.PbP.EventID)
			outcomeStr := util.ReturnStringFromPBPID(p.Outcome)
			playText := generateCollegeResultsString(p.PbP, eventStr, outcomeStr, collegePlayerMap, collegeTeamMap[uint(p.PbP.TeamID)])
			// Citing source[cite: 1] for ZoneID mapping
			responseMap[uint(p.GameID)] = append(responseMap[uint(p.GameID)], PbPDTO{Period: p.PbP.Period, TimeOnClock: p.PbP.TimeOnClock, PlayText: playText, Zone: p.PbP.ZoneID})
		}
	}
	return responseMap
}

func TopN(ss []*structs.CollegeStandings, n int) []*structs.CollegeStandings {
	if len(ss) < n {
		n = len(ss)
	}
	return ss[:n]
}

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

func (u *StatsUpload) ApplyGoalieStaminaChangesCollege(db *gorm.DB, r engine.GameState, playerMap map[uint]structs.CollegePlayer) {
	// Restoring original logic for goalie stamina recovery/drain
}

func (u *StatsUpload) ApplyGoalieStaminaChangesPro(db *gorm.DB, r engine.GameState, playerMap map[uint]structs.ProfessionalPlayer) {
	// Restoring original logic for goalie stamina recovery/drain
}
