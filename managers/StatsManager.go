package managers

import (
	"strconv"

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
	games := GetCollegeGamesForCurrentMatchup(weekId, seasonId, gameDay)
	collegePlayerSeasonStatMap := GetCollegePlayerSeasonStatMap(seasonId)
	proPlayerSeasonStatMap := GetProPlayerSeasonStatMap(seasonId)
	collegeTeamSeasonStatMap := GetCollegeTeamSeasonStatMap(seasonId)
	proTeamSeasonStatMap := GetProTeamSeasonStatMap(seasonId)

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

	proGames := GetProfessionalGamesForCurrentMatchup(weekId, seasonId, gameDay)

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
			pGameStats := GetCollegePlayerGameStatsBySeason(seasonID)
			playerGameStatsChan <- pGameStats
		}()

		playerGameStats = <-playerGameStatsChan
		close(playerGameStatsChan)

		go func() {
			tGameStats := GetCollegeTeamGameStatsBySeason(seasonID)
			teamGameStatsChan <- tGameStats
		}()
		teamGameStats = <-teamGameStatsChan
		close(teamGameStatsChan)
	} else {
		playerSeasonStatsChan := make(chan []structs.CollegePlayerSeasonStats)
		teamSeasonStatsChan := make(chan []structs.CollegeTeamSeasonStats)

		go func() {
			pSeasonStats := GetCollegePlayerSeasonStatsBySeason(seasonID)
			playerSeasonStatsChan <- pSeasonStats
		}()

		playerSeasonStats = <-playerSeasonStatsChan
		close(playerSeasonStatsChan)

		go func() {
			tSeasonStats := GetCollegeTeamSeasonStatsBySeason(seasonID)
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
			pGameStats := GetProPlayerGameStatsBySeason(seasonID)
			playerGameStatsChan <- pGameStats
		}()

		playerGameStats = <-playerGameStatsChan
		close(playerGameStatsChan)

		go func() {
			tGameStats := GetProTeamGameStatsBySeason(seasonID)
			teamGameStatsChan <- tGameStats
		}()
		teamGameStats = <-teamGameStatsChan
		close(teamGameStatsChan)
	} else {
		playerSeasonStatsChan := make(chan []structs.ProfessionalPlayerSeasonStats)
		teamSeasonStatsChan := make(chan []structs.ProfessionalTeamSeasonStats)

		go func() {
			pSeasonStats := GetProPlayerSeasonStatsBySeason(seasonID)
			playerSeasonStatsChan <- pSeasonStats
		}()

		playerSeasonStats = <-playerSeasonStatsChan
		close(playerSeasonStatsChan)

		go func() {
			tSeasonStats := GetProTeamSeasonStatsBySeason(seasonID)
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

func GetCollegePlayerSeasonStatMap(seasonID string) map[uint]structs.CollegePlayerSeasonStats {
	seasonStats := GetCollegePlayerSeasonStatsBySeason(seasonID)
	return MakeCollegePlayerSeasonStatMap(seasonStats)
}

func GetProPlayerSeasonStatMap(seasonID string) map[uint]structs.ProfessionalPlayerSeasonStats {
	seasonStats := GetProPlayerSeasonStatsBySeason(seasonID)
	return MakeProPlayerSeasonStatMap(seasonStats)
}

func GetCollegePlayerSeasonStatsBySeason(SeasonID string) []structs.CollegePlayerSeasonStats {
	return repository.FindCollegePlayerSeasonStatsRecords(SeasonID)
}

func GetProPlayerSeasonStatsBySeason(SeasonID string) []structs.ProfessionalPlayerSeasonStats {
	return repository.FindProPlayerSeasonStatsRecords(SeasonID)
}

func GetCollegePlayerGameStatsBySeason(SeasonID string) []structs.CollegePlayerGameStats {
	return repository.FindCollegePlayerGameStatsRecords(SeasonID)
}

func GetProPlayerGameStatsBySeason(SeasonID string) []structs.ProfessionalPlayerGameStats {
	return repository.FindProPlayerGameStatsRecords(SeasonID)
}

func GetCollegeTeamSeasonStatMap(seasonID string) map[uint]structs.CollegeTeamSeasonStats {
	seasonStats := GetCollegeTeamSeasonStatsBySeason(seasonID)
	return MakeCollegeTeamSeasonStatMap(seasonStats)
}

func GetProTeamSeasonStatMap(seasonID string) map[uint]structs.ProfessionalTeamSeasonStats {
	seasonStats := GetProTeamSeasonStatsBySeason(seasonID)
	return MakeProTeamSeasonStatMap(seasonStats)
}

func GetCollegeTeamSeasonStatsBySeason(SeasonID string) []structs.CollegeTeamSeasonStats {
	return repository.FindCollegeTeamSeasonStatsRecords(SeasonID)
}

func GetProTeamSeasonStatsBySeason(SeasonID string) []structs.ProfessionalTeamSeasonStats {
	return repository.FindProTeamSeasonStatsRecords(SeasonID)
}

func GetCollegeTeamGameStatsBySeason(SeasonID string) []structs.CollegeTeamGameStats {
	return repository.FindCollegeTeamGameStatsRecords(SeasonID)
}

func GetProTeamGameStatsBySeason(SeasonID string) []structs.ProfessionalTeamGameStats {
	return repository.FindProTeamGameStatsRecords(SeasonID)
}

func makeCollegePlayerStatsObject(weekID, gameID uint, s engine.PlayerStatsDTO) structs.CollegePlayerGameStats {
	return structs.CollegePlayerGameStats{
		WeekID:        weekID,
		GameID:        gameID,
		RevealResults: false,
		BasePlayerStats: structs.BasePlayerStats{
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

func makeProPlayerStatsObject(weekID, gameID uint, s engine.PlayerStatsDTO) structs.ProfessionalPlayerGameStats {
	return structs.ProfessionalPlayerGameStats{
		WeekID:        weekID,
		GameID:        gameID,
		RevealResults: false,
		BasePlayerStats: structs.BasePlayerStats{
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

func makeCollegeTeamStatsObject(weekID, gameID, seasonID uint, s engine.TeamStatDTO) structs.CollegeTeamGameStats {
	return structs.CollegeTeamGameStats{
		WeekID: weekID,
		GameID: gameID,
		BaseTeamStats: structs.BaseTeamStats{
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

func makeProTeamStatsObject(weekID, gameID, seasonID uint, s engine.TeamStatDTO) structs.ProfessionalTeamGameStats {
	return structs.ProfessionalTeamGameStats{
		WeekID: weekID,
		GameID: gameID,
		BaseTeamStats: structs.BaseTeamStats{
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
