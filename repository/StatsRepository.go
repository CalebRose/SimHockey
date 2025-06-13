package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func CreatePHLPlayerGameStatsRecordBatch(db *gorm.DB, games []structs.ProfessionalPlayerGameStats, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCHLPlayerGameStatsRecordBatch(db *gorm.DB, games []structs.CollegePlayerGameStats, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreatePHLTeamGameStatsRecordBatch(db *gorm.DB, games []structs.ProfessionalTeamGameStats, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCHLTeamGameStatsRecordBatch(db *gorm.DB, games []structs.CollegeTeamGameStats, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreatePHLPlayByPlayRecordBatch(db *gorm.DB, games []structs.ProPlayByPlay, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCHLPlayByPlayRecordBatch(db *gorm.DB, games []structs.CollegePlayByPlay, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func MassDeleteProPlayByPlaysExceptShotsOnGoal(db *gorm.DB) error {
	err := db.Where("event_id != ?", "29").Delete(&structs.ProPlayByPlay{}).Error
	if err != nil {
		log.Panicln("Could not mass delete pro play by plays")
	}
	return err
}

func MassDeleteCollegePlayByPlaysExceptShotsOnGoal(db *gorm.DB) error {
	err := db.Where("event_id != ?", "29").Delete(&structs.CollegePlayByPlay{}).Error
	if err != nil {
		log.Panicln("Could not mass delete college play by plays")
	}
	return err
}

func FindCollegePlayerSeasonStatsRecords(SeasonID, gameType string) []structs.CollegePlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&playerStats)

	return playerStats
}

func FindProPlayerSeasonStatsRecords(SeasonID, gameType string) []structs.ProfessionalPlayerSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.ProfessionalPlayerSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&playerStats)

	return playerStats
}

func FindCollegePlayerGameStatsRecords(SeasonID, GameID string) []structs.CollegePlayerGameStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerGameStats

	query := db.Model(&playerStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Order("points desc").Find(&playerStats)

	return playerStats
}

func FindProPlayerGameStatsRecords(SeasonID, GameID string) []structs.ProfessionalPlayerGameStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.ProfessionalPlayerGameStats
	query := db.Model(&playerStats)
	if len(SeasonID) > 0 {
		query = query.Where("season_id = ?", SeasonID)
	}

	if len(GameID) > 0 {
		query = query.Where("game_id = ?", GameID)
	}

	query.Order("points desc").Find(&playerStats)

	return playerStats
}

func FindCollegeTeamSeasonStatsRecords(SeasonID, gameType string) []structs.CollegeTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindProTeamSeasonStatsRecords(SeasonID, gameType string) []structs.ProfessionalTeamSeasonStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.ProfessionalTeamSeasonStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindCollegeTeamGameStatsRecords(SeasonID, gameType string) []structs.CollegeTeamGameStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.CollegeTeamGameStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindProTeamGameStatsRecords(SeasonID, gameType string) []structs.ProfessionalTeamGameStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats []structs.ProfessionalTeamGameStats

	db.Order("points desc").Where("season_id = ? AND game_type = ?", SeasonID, gameType).Find(&teamStats)

	return teamStats
}

func FindCollegeTeamStatsRecordByGame(gameID, teamID string) structs.CollegeTeamGameStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.CollegeTeamGameStats

	db.Order("points desc").Where("game_id = ? AND team_id = ?", gameID, teamID).Find(&teamStats)

	return teamStats
}

func FindProTeamStatsRecordByGame(gameID, teamID string) structs.ProfessionalTeamGameStats {
	db := dbprovider.GetInstance().GetDB()

	var teamStats structs.ProfessionalTeamGameStats

	db.Order("points desc").Where("game_id = ? AND team_id = ?", gameID, teamID).Find(&teamStats)

	return teamStats
}

func FindCollegePlayerStatsRecordByGame(gameID string) []structs.CollegePlayerGameStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.CollegePlayerGameStats

	db.Order("points desc").Where("game_id = ?", gameID).Find(&playerStats)

	return playerStats
}

func FindProPlayerStatsRecordByGame(gameID string) []structs.ProfessionalPlayerGameStats {
	db := dbprovider.GetInstance().GetDB()

	var playerStats []structs.ProfessionalPlayerGameStats

	db.Order("points desc").Where("game_id = ?", gameID).Find(&playerStats)

	return playerStats
}

func SaveCollegePlayerSeasonStatsRecord(stats structs.CollegePlayerSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Fatalln("Could not save season stats for " + strconv.Itoa(int(stats.PlayerID)))
	}
}

func SaveProPlayerSeasonStatsRecord(stats structs.ProfessionalPlayerSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Fatalln("Could not save season stats for " + strconv.Itoa(int(stats.PlayerID)))
	}
}

func SaveCollegeTeamSeasonStatsRecord(stats structs.CollegeTeamSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Fatalln("Could not save season stats for " + strconv.Itoa(int(stats.TeamID)))
	}
}

func SaveProTeamSeasonStatsRecord(stats structs.ProfessionalTeamSeasonStats, db *gorm.DB) {
	err := db.Save(&stats).Error
	if err != nil {
		log.Fatalln("Could not save season stats for " + strconv.Itoa(int(stats.TeamID)))
	}
}

func FindCHLPlayByPlaysRecordsByGameID(id string) []structs.CollegePlayByPlay {
	db := dbprovider.GetInstance().GetDB()

	plays := []structs.CollegePlayByPlay{}

	db.Where("game_id = ?", id).Find(&plays)

	return plays
}

func FindPHLPlayByPlaysRecordsByGameID(id string) []structs.ProPlayByPlay {
	db := dbprovider.GetInstance().GetDB()

	plays := []structs.ProPlayByPlay{}

	db.Where("game_id = ?", id).Find(&plays)

	return plays
}
