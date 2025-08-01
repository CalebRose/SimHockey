package managers

import (
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/engine"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func UpdateSeasonStats(ts structs.Timestamp, gameDay string) {
	db := dbprovider.GetInstance().GetDB()

	weekId := strconv.Itoa(int(ts.WeekID))
	seasonId := strconv.Itoa(int(ts.SeasonID))
	collegeGameIDs := []string{}
	proGameIDs := []string{}
	_, collegeGameType := ts.GetCurrentGameType(true)
	_, proGameType := ts.GetCurrentGameType(true)
	games := GetCollegeGamesForCurrentMatchup(weekId, seasonId, gameDay, ts.IsPreseason)
	collegePlayerSeasonStatMap := GetCollegePlayerSeasonStatMap(seasonId, collegeGameType)
	proPlayerSeasonStatMap := GetProPlayerSeasonStatMap(seasonId, proGameType)
	collegeTeamSeasonStatMap := GetCollegeTeamSeasonStatMap(seasonId, collegeGameType)
	proTeamSeasonStatMap := GetProTeamSeasonStatMap(seasonId, proGameType)

	for _, game := range games {
		if !game.GameComplete {
			continue
		}
		matchId := strconv.Itoa(int(game.ID))
		collegeGameIDs = append(collegeGameIDs, matchId)

		homeTeamStats := repository.FindCollegeTeamStatsRecordByGame(strconv.Itoa(int(game.ID)), strconv.Itoa(int(game.HomeTeamID)))
		homeSeasonStats := collegeTeamSeasonStatMap[game.HomeTeamID]

		homeSeasonStats.BaseTeamStats.AddStatsToSeasonRecord(homeTeamStats.BaseTeamStats)
		homeSeasonStats.TeamSeasonStats.AddStatsToSeasonRecord(homeTeamStats.BaseTeamStats, game.IsPlayoffGame, game.IsShootout)

		repository.SaveCollegeTeamSeasonStatsRecord(homeSeasonStats, db)

		awayTeamStats := repository.FindCollegeTeamStatsRecordByGame(strconv.Itoa(int(game.ID)), strconv.Itoa(int(game.AwayTeamID)))

		awaySeasonStats := collegeTeamSeasonStatMap[game.AwayTeamID]

		awaySeasonStats.BaseTeamStats.AddStatsToSeasonRecord(awayTeamStats.BaseTeamStats)
		awaySeasonStats.TeamSeasonStats.AddStatsToSeasonRecord(awayTeamStats.BaseTeamStats, game.IsPlayoffGame, game.IsShootout)

		repository.SaveCollegeTeamSeasonStatsRecord(awaySeasonStats, db)

		playerStats := repository.FindCollegePlayerStatsRecordByGame(strconv.Itoa(int(game.ID)))

		for _, stat := range playerStats {
			if stat.TimeOnIce <= 0 {
				continue
			}
			playerSeasonStats := collegePlayerSeasonStatMap[stat.PlayerID]
			playerSeasonStats.AddStatsToSeasonRecord(stat.BasePlayerStats)

			// if stat.IsInjured {
			// id := strconv.Itoa(int(stat.PlayerID))
			// 	player := GetCollegePlayerByPlayerID(id)
			// 	player.SetInjury(stat.InjuryName, stat.InjuryType, int(stat.WeeksOfRecovery))
			// 	message := player.Position + " " + player.FirstName + " " + player.LastName + " has been injured for " + strconv.Itoa(int(stat.WeeksOfRecovery)) + "."
			// 	CreateNotification("CBB", message, "Injury", player.TeamID)
			// 	repository.SaveCollegePlayerRecord(player, db)
			// }

			repository.SaveCollegePlayerSeasonStatsRecord(playerSeasonStats, db)
		}
	}

	db.Model(&structs.CollegePlayerGameStats{}).Where("game_id in (?)", collegeGameIDs).Update("reveal_results", true)
	db.Model(&structs.CollegeTeamGameStats{}).Where("game_id in (?)", collegeGameIDs).Update("reveal_results", true)

	proGames := GetProfessionalGamesForCurrentMatchup(weekId, seasonId, gameDay, ts.IsPreseason)

	for _, game := range proGames {
		if !game.GameComplete {
			continue
		}
		matchId := strconv.Itoa(int(game.ID))
		proGameIDs = append(proGameIDs, matchId)

		homeTeamStats := repository.FindProTeamStatsRecordByGame(strconv.Itoa(int(game.ID)), strconv.Itoa(int(game.HomeTeamID)))
		homeSeasonStats := proTeamSeasonStatMap[game.HomeTeamID]

		homeSeasonStats.BaseTeamStats.AddStatsToSeasonRecord(homeTeamStats.BaseTeamStats)
		homeSeasonStats.TeamSeasonStats.AddStatsToSeasonRecord(homeTeamStats.BaseTeamStats, game.IsPlayoffGame, game.IsShootout)

		repository.SaveProTeamSeasonStatsRecord(homeSeasonStats, db)

		awayTeamStats := repository.FindProTeamStatsRecordByGame(strconv.Itoa(int(game.ID)), strconv.Itoa(int(game.AwayTeamID)))

		awaySeasonStats := proTeamSeasonStatMap[game.AwayTeamID]

		awaySeasonStats.BaseTeamStats.AddStatsToSeasonRecord(awayTeamStats.BaseTeamStats)
		awaySeasonStats.TeamSeasonStats.AddStatsToSeasonRecord(awayTeamStats.BaseTeamStats, game.IsPlayoffGame, game.IsShootout)

		repository.SaveProTeamSeasonStatsRecord(awaySeasonStats, db)

		playerStats := repository.FindProPlayerStatsRecordByGame(strconv.Itoa(int(game.ID)))

		for _, stat := range playerStats {
			if stat.TimeOnIce <= 0 {
				continue
			}
			playerSeasonStats := proPlayerSeasonStatMap[stat.PlayerID]
			playerSeasonStats.AddStatsToSeasonRecord(stat.BasePlayerStats)

			// if stat.IsInjured {
			// id := strconv.Itoa(int(stat.PlayerID))
			// 	player := GetCollegePlayerByPlayerID(id)
			// 	player.SetInjury(stat.InjuryName, stat.InjuryType, int(stat.WeeksOfRecovery))
			// 	message := player.Position + " " + player.FirstName + " " + player.LastName + " has been injured for " + strconv.Itoa(int(stat.WeeksOfRecovery)) + "."
			// 	CreateNotification("CBB", message, "Injury", player.TeamID)
			// 	repository.SaveCollegePlayerRecord(player, db)
			// }

			repository.SaveProPlayerSeasonStatsRecord(playerSeasonStats, db)
		}

		db.Model(&structs.ProfessionalPlayerGameStats{}).Where("game_id in (?)", proGameIDs).Update("reveal_results", true)
		db.Model(&structs.ProfessionalTeamGameStats{}).Where("game_id in (?)", proGameIDs).Update("reveal_results", true)
	}
}

func SearchCollegeStats(seasonID, weekID, viewType, gameType string) structs.SearchStatsResponse {
	var (
		playerGameStats   []structs.CollegePlayerGameStats
		playerSeasonStats []structs.CollegePlayerSeasonStats
		teamGameStats     []structs.CollegeTeamGameStats
		teamSeasonStats   []structs.CollegeTeamSeasonStats
	)

	// Fetch week stats by season... will save time for the player
	if viewType == "WEEK" {
		playerGameStatsChan := make(chan []structs.CollegePlayerGameStats)
		teamGameStatsChan := make(chan []structs.CollegeTeamGameStats)
		go func() {
			pGameStats := GetCollegePlayerGameStatsByWeek(weekID, gameType)
			playerGameStatsChan <- pGameStats
		}()

		playerGameStats = <-playerGameStatsChan
		close(playerGameStatsChan)

		go func() {
			tGameStats := GetCollegeTeamGameStatsByWeek(seasonID, gameType)
			teamGameStatsChan <- tGameStats
		}()
		teamGameStats = <-teamGameStatsChan
		close(teamGameStatsChan)
	} else {
		playerSeasonStatsChan := make(chan []structs.CollegePlayerSeasonStats)
		teamSeasonStatsChan := make(chan []structs.CollegeTeamSeasonStats)

		go func() {
			pSeasonStats := GetCollegePlayerSeasonStatsBySeason(seasonID, gameType)
			playerSeasonStatsChan <- pSeasonStats
		}()

		playerSeasonStats = <-playerSeasonStatsChan
		close(playerSeasonStatsChan)

		go func() {
			tSeasonStats := GetCollegeTeamSeasonStatsBySeason(seasonID, gameType)
			teamSeasonStatsChan <- tSeasonStats
		}()
		teamSeasonStats = <-teamSeasonStatsChan
		close(teamSeasonStatsChan)
	}

	return structs.SearchStatsResponse{
		CHLPlayerGameStats:   playerGameStats,
		CHLPlayerSeasonStats: playerSeasonStats,
		CHLTeamGameStats:     teamGameStats,
		CHLTeamSeasonStats:   teamSeasonStats,
	}
}

func SearchProStats(seasonID, weekID, viewType, gameType string) structs.SearchStatsResponse {
	var (
		playerGameStats   []structs.ProfessionalPlayerGameStats
		playerSeasonStats []structs.ProfessionalPlayerSeasonStats
		teamGameStats     []structs.ProfessionalTeamGameStats
		teamSeasonStats   []structs.ProfessionalTeamSeasonStats
	)

	// Fetch week stats by season... will save time for the player
	if viewType == "WEEK" {
		playerGameStatsChan := make(chan []structs.ProfessionalPlayerGameStats)
		teamGameStatsChan := make(chan []structs.ProfessionalTeamGameStats)
		go func() {
			pGameStats := GetProPlayerGameStatsByWeek(weekID, gameType)
			playerGameStatsChan <- pGameStats
		}()

		playerGameStats = <-playerGameStatsChan
		close(playerGameStatsChan)

		go func() {
			tGameStats := GetProTeamGameStatsByWeek(weekID, gameType)
			teamGameStatsChan <- tGameStats
		}()
		teamGameStats = <-teamGameStatsChan
		close(teamGameStatsChan)
	} else {
		playerSeasonStatsChan := make(chan []structs.ProfessionalPlayerSeasonStats)
		teamSeasonStatsChan := make(chan []structs.ProfessionalTeamSeasonStats)

		go func() {
			pSeasonStats := GetProPlayerSeasonStatsBySeason(seasonID, gameType)
			playerSeasonStatsChan <- pSeasonStats
		}()

		playerSeasonStats = <-playerSeasonStatsChan
		close(playerSeasonStatsChan)

		go func() {
			tSeasonStats := GetProTeamSeasonStatsBySeason(seasonID, gameType)
			teamSeasonStatsChan <- tSeasonStats
		}()
		teamSeasonStats = <-teamSeasonStatsChan
		close(teamSeasonStatsChan)
	}

	return structs.SearchStatsResponse{
		PHLPlayerGameStats:   playerGameStats,
		PHLPlayerSeasonStats: playerSeasonStats,
		PHLTeamGameStats:     teamGameStats,
		PHLTeamSeasonStats:   teamSeasonStats,
	}
}

func GetCollegePlayerSeasonStatMap(seasonID, gameType string) map[uint]structs.CollegePlayerSeasonStats {
	seasonStats := GetCollegePlayerSeasonStatsBySeason(seasonID, gameType)
	return MakeCollegePlayerSeasonStatMap(seasonStats)
}

func GetProPlayerSeasonStatMap(seasonID, gameType string) map[uint]structs.ProfessionalPlayerSeasonStats {
	seasonStats := GetProPlayerSeasonStatsBySeason(seasonID, gameType)
	return MakeProPlayerSeasonStatMap(seasonStats)
}

func GetCollegePlayerSeasonStatsBySeason(SeasonID, gameType string) []structs.CollegePlayerSeasonStats {
	return repository.FindCollegePlayerSeasonStatsRecords(SeasonID, gameType)
}

func GetProPlayerSeasonStatsBySeason(SeasonID, gameType string) []structs.ProfessionalPlayerSeasonStats {
	return repository.FindProPlayerSeasonStatsRecords(SeasonID, gameType)
}

func GetCollegePlayerGameStatsByWeek(WeekID, gameType string) []structs.CollegePlayerGameStats {
	return repository.FindCollegePlayerGameStatsRecords("", WeekID, gameType, "")
}

func GetCollegePlayerGameStatsBySeason(SeasonID, gameType string) []structs.CollegePlayerGameStats {
	return repository.FindCollegePlayerGameStatsRecords(SeasonID, "", gameType, "")
}

func GetCollegePlayerGameStatsByGame(GameID string) []structs.CollegePlayerGameStats {
	return repository.FindCollegePlayerGameStatsRecords("", "", "", GameID)
}

func GetProPlayerGameStatsByWeek(WeekID, gameType string) []structs.ProfessionalPlayerGameStats {
	return repository.FindProPlayerGameStatsRecords("", WeekID, gameType, "")
}

func GetProPlayerGameStatsBySeason(SeasonID, gameType string) []structs.ProfessionalPlayerGameStats {
	return repository.FindProPlayerGameStatsRecords(SeasonID, "", gameType, "")
}

func GetCollegeTeamSeasonStatMap(seasonID, gameType string) map[uint]structs.CollegeTeamSeasonStats {
	seasonStats := GetCollegeTeamSeasonStatsBySeason(seasonID, gameType)
	return MakeCollegeTeamSeasonStatMap(seasonStats)
}

func GetProTeamSeasonStatMap(seasonID, gameType string) map[uint]structs.ProfessionalTeamSeasonStats {
	seasonStats := GetProTeamSeasonStatsBySeason(seasonID, gameType)
	return MakeProTeamSeasonStatMap(seasonStats)
}

func GetCollegeTeamSeasonStatsBySeason(SeasonID, gameType string) []structs.CollegeTeamSeasonStats {
	return repository.FindCollegeTeamSeasonStatsRecords(SeasonID, gameType)
}

func GetProTeamSeasonStatsBySeason(SeasonID, gameType string) []structs.ProfessionalTeamSeasonStats {
	return repository.FindProTeamSeasonStatsRecords(SeasonID, gameType)
}

func GetCollegeTeamGameStatsByWeek(WeekID, gameType string) []structs.CollegeTeamGameStats {
	return repository.FindCollegeTeamGameStatsRecords("", WeekID, gameType, "")
}

func GetProTeamGameStatsByWeek(WeekID, gameType string) []structs.ProfessionalTeamGameStats {
	return repository.FindProTeamGameStatsRecords("", WeekID, gameType, "")
}

func GetCollegeTeamGameStatsBySeason(SeasonID, gameType string) []structs.CollegeTeamGameStats {
	return repository.FindCollegeTeamGameStatsRecords(SeasonID, "", gameType, "")
}

func GetProTeamGameStatsBySeason(SeasonID, gameType string) []structs.ProfessionalTeamGameStats {
	return repository.FindProTeamGameStatsRecords(SeasonID, "", gameType, "")
}

func makeCollegePlayerStatsObject(weekID, gameID, gameType uint, s engine.PlayerStatsDTO) structs.CollegePlayerGameStats {
	return structs.CollegePlayerGameStats{
		WeekID:        weekID,
		GameID:        gameID,
		RevealResults: false,
		BasePlayerStats: structs.BasePlayerStats{
			StartedGame:          s.StartedGame,
			GameType:             uint8(gameType),
			GameDay:              s.GameDay,
			PlayerID:             s.PlayerID,
			TeamID:               s.TeamID,
			SeasonID:             s.SeasonID,
			Goals:                s.Goals,
			Assists:              s.Assists,
			Points:               s.Points,
			PlusMinus:            s.PlusMinus,
			PenaltyMinutes:       s.PenaltyMinutes,
			EvenStrengthGoals:    s.EvenStrengthGoals,
			EvenStrengthPoints:   s.EvenStrengthPoints,
			PowerPlayGoals:       s.PowerPlayGoals,
			PowerPlayPoints:      s.PowerPlayPoints,
			ShorthandedGoals:     s.ShorthandedGoals,
			ShorthandedPoints:    s.ShorthandedPoints,
			OvertimeGoals:        s.OvertimeGoals,
			GameWinningGoals:     s.GameWinningGoals,
			Shots:                s.Shots,
			ShootingPercentage:   float32(s.Goals) / float32(s.Shots),
			TimeOnIce:            s.TimeOnIce,
			FaceOffWinPercentage: s.FaceOffWinPercentage,
			FaceOffsWon:          s.FaceOffsWon,
			FaceOffs:             s.FaceOffs,
			GoalieWins:           s.GoalieWins,
			GoalieLosses:         s.GoalieLosses,
			GoalieTies:           s.GoalieTies,
			OvertimeLosses:       s.OvertimeLosses,
			ShotsAgainst:         s.ShotsAgainst,
			Saves:                s.Saves,
			GoalsAgainst:         s.GoalsAgainst,
			SavePercentage:       s.SavePercentage,
			Shutouts:             s.Shutouts,
			ShotsBlocked:         s.ShotsBlocked,
			BodyChecks:           s.BodyChecks,
			StickChecks:          s.StickChecks,
		},
	}
}

func makeProPlayerStatsObject(weekID, gameID, gameType uint, s engine.PlayerStatsDTO) structs.ProfessionalPlayerGameStats {
	return structs.ProfessionalPlayerGameStats{
		WeekID:        weekID,
		GameID:        gameID,
		RevealResults: false,
		BasePlayerStats: structs.BasePlayerStats{
			GameType:             uint8(gameType),
			StartedGame:          s.StartedGame,
			GameDay:              s.GameDay,
			PlayerID:             s.PlayerID,
			TeamID:               s.TeamID,
			SeasonID:             s.SeasonID,
			Goals:                s.Goals,
			Assists:              s.Assists,
			Points:               s.Points,
			PlusMinus:            s.PlusMinus,
			PenaltyMinutes:       s.PenaltyMinutes,
			EvenStrengthGoals:    s.EvenStrengthGoals,
			EvenStrengthPoints:   s.EvenStrengthPoints,
			PowerPlayGoals:       s.PowerPlayGoals,
			PowerPlayPoints:      s.PowerPlayPoints,
			ShorthandedGoals:     s.ShorthandedGoals,
			ShorthandedPoints:    s.ShorthandedPoints,
			OvertimeGoals:        s.OvertimeGoals,
			GameWinningGoals:     s.GameWinningGoals,
			Shots:                s.Shots,
			ShootingPercentage:   s.ShootingPercentage,
			TimeOnIce:            s.TimeOnIce,
			FaceOffWinPercentage: s.FaceOffWinPercentage,
			FaceOffsWon:          s.FaceOffsWon,
			FaceOffs:             s.FaceOffs,
			GoalieWins:           s.GoalieWins,
			GoalieLosses:         s.GoalieLosses,
			GoalieTies:           s.GoalieTies,
			OvertimeLosses:       s.OvertimeLosses,
			ShotsAgainst:         s.ShotsAgainst,
			Saves:                s.Saves,
			GoalsAgainst:         s.GoalsAgainst,
			SavePercentage:       s.SavePercentage,
			Shutouts:             s.Shutouts,
			ShotsBlocked:         s.ShotsBlocked,
			BodyChecks:           s.BodyChecks,
			StickChecks:          s.StickChecks,
		},
	}
}

func makeCollegeTeamStatsObject(weekID, gameID, seasonID, gameType uint, s engine.TeamStatDTO) structs.CollegeTeamGameStats {
	return structs.CollegeTeamGameStats{
		WeekID: weekID,
		GameID: gameID,
		BaseTeamStats: structs.BaseTeamStats{
			GameType:             uint8(gameType),
			GameDay:              s.GameDay,
			SeasonID:             seasonID,
			TeamID:               s.TeamID,
			Team:                 s.Team,
			GoalsFor:             s.GoalsFor,
			GoalsAgainst:         s.GoalsAgainst,
			Assists:              s.Assists,
			Points:               s.Points,
			Period1Score:         s.Period1Score,
			Period2Score:         s.Period2Score,
			Period3Score:         s.Period3Score,
			OTScore:              s.OTScore,
			PlusMinus:            s.PlusMinus,
			PenaltyMinutes:       s.PenaltyMinutes,
			EvenStrengthGoals:    s.EvenStrengthGoals,
			EvenStrengthPoints:   s.EvenStrengthPoints,
			PowerPlayGoals:       s.PowerPlayGoals,
			PowerPlayPoints:      s.PowerPlayPoints,
			ShorthandedGoals:     s.ShorthandedGoals,
			ShorthandedPoints:    s.ShorthandedPoints,
			OvertimeGoals:        s.OvertimeGoals,
			Shots:                s.Shots,
			ShootingPercentage:   s.ShootingPercentage,
			FaceOffWinPercentage: s.FaceOffWinPercentage,
			FaceOffsWon:          s.FaceOffsWon,
			FaceOffs:             s.FaceOffs,
			ShotsAgainst:         s.ShotsAgainst,
			Saves:                s.Saves,
			SavePercentage:       s.SavePercentage,
			Shutouts:             s.Shutouts,
		},
		RevealResults: false,
	}
}

func makeProTeamStatsObject(weekID, gameID, seasonID, gameType uint, s engine.TeamStatDTO) structs.ProfessionalTeamGameStats {
	return structs.ProfessionalTeamGameStats{
		WeekID: weekID,
		GameID: gameID,
		BaseTeamStats: structs.BaseTeamStats{
			GameType:             uint8(gameType),
			GameDay:              s.GameDay,
			SeasonID:             seasonID,
			TeamID:               s.TeamID,
			Team:                 s.Team,
			GoalsFor:             s.GoalsFor,
			GoalsAgainst:         s.GoalsAgainst,
			Assists:              s.Assists,
			Points:               s.Points,
			Period1Score:         s.Period1Score,
			Period2Score:         s.Period2Score,
			Period3Score:         s.Period3Score,
			OTScore:              s.OTScore,
			PlusMinus:            s.PlusMinus,
			PenaltyMinutes:       s.PenaltyMinutes,
			EvenStrengthGoals:    s.EvenStrengthGoals,
			EvenStrengthPoints:   s.EvenStrengthPoints,
			PowerPlayGoals:       s.PowerPlayGoals,
			PowerPlayPoints:      s.PowerPlayPoints,
			ShorthandedGoals:     s.ShorthandedGoals,
			ShorthandedPoints:    s.ShorthandedPoints,
			OvertimeGoals:        s.OvertimeGoals,
			Shots:                s.Shots,
			ShootingPercentage:   s.ShootingPercentage,
			FaceOffWinPercentage: s.FaceOffWinPercentage,
			FaceOffsWon:          s.FaceOffsWon,
			FaceOffs:             s.FaceOffs,
			ShotsAgainst:         s.ShotsAgainst,
			Saves:                s.Saves,
			SavePercentage:       s.SavePercentage,
			Shutouts:             s.Shutouts,
		},
		RevealResults: false,
	}
}

func GetCHLGameResultsByGameID(gameID string) structs.GameResultsResponse {
	game := repository.FindCollegeGameRecord(gameID)
	htID := strconv.Itoa(int(game.HomeTeamID))
	atID := strconv.Itoa(int(game.AwayTeamID))
	allStats := repository.FindCollegePlayerStatsRecordByGame(gameID)
	collegePlayerMap := GetCollegePlayersMap()
	teamMap := GetCollegeTeamMap()
	homeTeamStats := repository.FindCollegeTeamStatsRecordByGame(gameID, htID)
	awayTeamStats := repository.FindCollegeTeamStatsRecordByGame(gameID, atID)
	homePlayerStats := MakeCollegePlayerGameStatsListByTeamID(allStats, game.HomeTeamID)
	awayPlayerStats := MakeCollegePlayerGameStatsListByTeamID(allStats, game.AwayTeamID)
	score := structs.ScoreBoard{
		P1Home:            int(homeTeamStats.Period1Score),
		P2Home:            int(homeTeamStats.Period2Score),
		P3Home:            int(homeTeamStats.Period3Score),
		OTHome:            int(homeTeamStats.OTScore),
		HomeShootoutScore: int(game.HomeTeamShootoutScore),
		P1Away:            int(awayTeamStats.Period1Score),
		P2Away:            int(awayTeamStats.Period2Score),
		P3Away:            int(awayTeamStats.Period3Score),
		OTAway:            int(awayTeamStats.OTScore),
		AwayShootoutScore: int(game.AwayTeamShootoutScore),
	}

	playByPlays := GetCHLPlayByPlaysByGameID(gameID)
	// Generate the Play By Play Response
	playbyPlayResponseList := GenerateCHLPlayByPlayResponse(playByPlays, teamMap, collegePlayerMap, false, game.HomeTeamID, game.AwayTeamID)

	return structs.GameResultsResponse{
		CHLHomeStats:   homePlayerStats,
		CHLAwayStats:   awayPlayerStats,
		CHLPlayByPlays: playbyPlayResponseList,
		Score:          score,
	}
}

func GetPHLGameResultsByGameID(gameID string) structs.GameResultsResponse {
	game := repository.FindProfessionalGameRecord(gameID)
	htID := strconv.Itoa(int(game.HomeTeamID))
	atID := strconv.Itoa(int(game.AwayTeamID))
	allStats := repository.FindProPlayerStatsRecordByGame(gameID)
	playerMap := GetProPlayersMap()
	teamMap := GetProTeamMap()
	homeTeamStats := repository.FindProTeamStatsRecordByGame(gameID, htID)
	awayTeamStats := repository.FindProTeamStatsRecordByGame(gameID, atID)
	homePlayerStats := MakeProPlayerGameStatsListByTeamID(allStats, game.HomeTeamID)
	awayPlayerStats := MakeProPlayerGameStatsListByTeamID(allStats, game.AwayTeamID)
	score := structs.ScoreBoard{
		P1Home:            int(homeTeamStats.Period1Score),
		P2Home:            int(homeTeamStats.Period2Score),
		P3Home:            int(homeTeamStats.Period3Score),
		OTHome:            int(homeTeamStats.OTScore),
		HomeShootoutScore: int(game.HomeTeamShootoutScore),
		P1Away:            int(awayTeamStats.Period1Score),
		P2Away:            int(awayTeamStats.Period2Score),
		P3Away:            int(awayTeamStats.Period3Score),
		OTAway:            int(awayTeamStats.OTScore),
		AwayShootoutScore: int(game.AwayTeamShootoutScore),
	}

	playByPlays := GetPHLPlayByPlaysByGameID(gameID)
	// Generate the Play By Play Response
	playbyPlayResponseList := GeneratePHLPlayByPlayResponse(playByPlays, teamMap, playerMap, false, game.HomeTeamID, game.AwayTeamID)
	return structs.GameResultsResponse{
		PHLHomeStats:   homePlayerStats,
		PHLAwayStats:   awayPlayerStats,
		PHLPlayByPlays: playbyPlayResponseList,
		Score:          score,
	}
}

func GenerateCHLPlayByPlayResponse(playByPlays []structs.CollegePlayByPlay, teamMap map[uint]structs.CollegeTeam, playerMap map[uint]structs.CollegePlayer, isStream bool, ht, at uint) []structs.PlayByPlayResponse {
	results := []structs.PlayByPlayResponse{}
	for idx, play := range playByPlays {
		if play.Outcome == util.CarrierKeepsPuckID {
			continue
		}
		timeOnClock := FormatTimeToClock(play.TimeOnClock)
		event := util.ReturnStringFromPBPID(play.EventID)
		outcome := util.ReturnStringFromPBPID(play.Outcome)
		possessingTeam := teamMap[uint(play.TeamID)]
		zone := getZoneLabel(play.ZoneID)
		nextZone := getZoneLabel(play.NextZoneID)
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		isFight := "No"
		if play.IsFight {
			isFight = "Yes"
		}
		result := generateCollegeResultsString(play.PbP, event, outcome, playerMap, possessingTeam)

		res := structs.PlayByPlayResponse{
			GameID:            play.GameID,
			PlayNumber:        uint(idx) + 1,
			HomeTeamID:        ht,
			HomeTeamScore:     play.HomeTeamScore,
			AwayTeamID:        at,
			AwayTeamScore:     play.AwayTeamScore,
			Period:            play.Period,
			TimeOnClock:       timeOnClock,
			SecondsConsumed:   play.SecondsConsumed,
			Event:             event,
			Outcome:           outcome,
			Zone:              zone,
			NextZone:          nextZone,
			TeamID:            play.TeamID,
			PuckCarrierID:     play.PuckCarrierID,
			PassedPlayerID:    play.PassedPlayerID,
			AssistingPlayerID: play.AssistingPlayerID,
			DefenderID:        play.DefenderID,
			GoalieID:          play.GoalieID,
			InjuryID:          play.InjuryID,
			InjuryType:        play.InjuryType,
			InjuryDuration:    play.InjuryDuration,
			Penalty:           penalty,
			Severity:          severity,
			IsFight:           isFight,
			IsBreakaway:       play.IsBreakaway,
			Result:            result,
		}

		results = append(results, res)
	}
	return results
}

func GeneratePHLPlayByPlayResponse(playByPlays []structs.ProPlayByPlay, teamMap map[uint]structs.ProfessionalTeam, playerMap map[uint]structs.ProfessionalPlayer, isStream bool, ht, at uint) []structs.PlayByPlayResponse {
	results := []structs.PlayByPlayResponse{}
	for idx, play := range playByPlays {
		if play.Outcome == util.CarrierKeepsPuckID {
			continue
		}
		timeOnClock := FormatTimeToClock(play.TimeOnClock)
		event := util.ReturnStringFromPBPID(play.EventID)
		outcome := util.ReturnStringFromPBPID(play.Outcome)
		possessingTeam := teamMap[uint(play.TeamID)]
		zone := getZoneLabel(play.ZoneID)
		nextZone := getZoneLabel(play.NextZoneID)
		penalty := getPenaltyByID(uint(play.PenaltyID))
		severity := getSeverityByID(play.Severity)
		isFight := "No"
		if play.IsFight {
			isFight = "Yes"
		}
		result := generateProResultsString(play.PbP, event, outcome, playerMap, possessingTeam)

		res := structs.PlayByPlayResponse{
			GameID:            play.GameID,
			PlayNumber:        uint(idx) + 1,
			HomeTeamID:        ht,
			HomeTeamScore:     play.HomeTeamScore,
			AwayTeamID:        at,
			AwayTeamScore:     play.AwayTeamScore,
			Period:            play.Period,
			TimeOnClock:       timeOnClock,
			SecondsConsumed:   play.SecondsConsumed,
			Event:             event,
			Outcome:           outcome,
			Zone:              zone,
			NextZone:          nextZone,
			TeamID:            play.TeamID,
			PuckCarrierID:     play.PuckCarrierID,
			PassedPlayerID:    play.PassedPlayerID,
			AssistingPlayerID: play.AssistingPlayerID,
			DefenderID:        play.DefenderID,
			GoalieID:          play.GoalieID,
			InjuryID:          play.InjuryID,
			InjuryType:        play.InjuryType,
			InjuryDuration:    play.InjuryDuration,
			Penalty:           penalty,
			Severity:          severity,
			IsFight:           isFight,
			IsBreakaway:       play.IsBreakaway,
			Result:            result,
		}

		results = append(results, res)
	}
	return results
}

func GetCHLPlayByPlaysByGameID(id string) []structs.CollegePlayByPlay {
	return repository.FindCHLPlayByPlaysRecordsByGameID(id)
}

func GetPHLPlayByPlaysByGameID(id string) []structs.ProPlayByPlay {
	return repository.FindPHLPlayByPlaysRecordsByGameID(id)
}
