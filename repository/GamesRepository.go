package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func CreatePHLGamesRecordsBatch(db *gorm.DB, games []structs.ProfessionalGame, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCHLGamesRecordsBatch(db *gorm.DB, games []structs.CollegeGame, batchSize int) error {
	total := len(games)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(games[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCHLSeriesRecordsBatch(db *gorm.DB, series []structs.CollegeSeries, batchSize int) error {
	total := len(series)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(series[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreatePHLSeriesRecordsBatch(db *gorm.DB, series []structs.ProSeries, batchSize int) error {
	total := len(series)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(series[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func SaveCollegeGameRecord(gameRecord structs.CollegeGame, db *gorm.DB) {
	err := db.Save(&gameRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(gameRecord.ID)))
	}
}

func SaveProfessionalGameRecord(gameRecord structs.ProfessionalGame, db *gorm.DB) {
	err := db.Save(&gameRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(gameRecord.ID)))
	}
}

func FindCollegeGamesByCurrentMatchup(weekID, seasonID, gameDay string, isPreseason bool) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame
	err := db.Where("week_id = ? AND season_id = ? AND game_day = ? AND is_preseason = ?", weekID, seasonID, gameDay, isPreseason).Find(&games).Error
	if err != nil {
		log.Printf("Error querying for college games: %v", err)

	}

	return games
}

func FindProfessionalGamesByCurrentMatchup(weekID, seasonID, gameDay string, isPreseason bool) []structs.ProfessionalGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.ProfessionalGame
	err := db.Where("week_id = ? AND season_id = ? AND game_day = ? AND is_preseason = ?", weekID, seasonID, gameDay, isPreseason).Find(&games).Error
	if err != nil {
		log.Printf("Error querying for professional games: %v", err)

	}

	return games
}

type GamesClauses struct {
	SeasonID      string
	WeekID        string
	TeamID        string
	IsPreseason   bool
	Timeslot      string
	GameCompleted string
}

func FindCollegeGames(clauses GamesClauses) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame

	query := db.Model(&games)
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if len(clauses.TeamID) > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", clauses.TeamID, clauses.TeamID)
	}
	if len(clauses.WeekID) > 0 {
		query = query.Where("week_id = ?", clauses.WeekID)
	}

	if len(clauses.Timeslot) > 0 {
		query = query.Where("game_day = ?", clauses.Timeslot)
	}

	if len(clauses.GameCompleted) > 0 {
		gameCompleted := clauses.GameCompleted == "Y"
		query = query.Where("game_complete = ?", gameCompleted)
	}

	if err := query.Order("week_id asc").Where("is_preseason = ?", clauses.IsPreseason).Find(&games).Error; err != nil {
		return []structs.CollegeGame{}
	}

	return games
}

func FindProfessionalGames(clauses GamesClauses) []structs.ProfessionalGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.ProfessionalGame
	query := db.Model(&games)
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}
	if len(clauses.TeamID) > 0 {
		query = query.Where("home_team_id = ? OR away_team_id = ?", clauses.TeamID, clauses.TeamID)
	}
	if len(clauses.WeekID) > 0 {
		query = query.Where("week_id = ?", clauses.WeekID)
	}

	if len(clauses.Timeslot) > 0 {
		query = query.Where("game_day = ?", clauses.Timeslot)
	}
	if len(clauses.GameCompleted) > 0 {
		gameCompleted := clauses.GameCompleted == "Y"
		query = query.Where("game_complete = ?", gameCompleted)
	}

	if err := query.Order("week_id asc").Where("is_preseason = ?", clauses.IsPreseason).Find(&games).Error; err != nil {
		return []structs.ProfessionalGame{}
	}
	return games
}

func FindCollegeGameRecord(id string) structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games structs.CollegeGame

	query := db.Model(&games)
	if len(id) > 0 {
		query = query.Where("id = ?", id)
	}
	if err := query.Order("week_id asc").Find(&games).Error; err != nil {
		return structs.CollegeGame{}
	}

	return games
}

func FindProfessionalGameRecord(id string) structs.ProfessionalGame {
	db := dbprovider.GetInstance().GetDB()

	var games structs.ProfessionalGame
	query := db.Model(&games)
	if len(id) > 0 {
		query = query.Where("id = ?", id)
	}
	if err := query.Order("week_id asc").Find(&games).Error; err != nil {
		return structs.ProfessionalGame{}
	}
	return games
}

func FindPlayoffSeriesByID(seriesID string) structs.ProSeries {
	db := dbprovider.GetInstance().GetDB()

	var series structs.ProSeries

	db.Where("id = ?", seriesID).Find(&series)

	return series
}

func SavePlayoffSeriesRecord(seriesRecord structs.ProSeries, db *gorm.DB) {
	err := db.Save(&seriesRecord).Error
	if err != nil {
		log.Panicln("Could not save playoff series " + strconv.Itoa(int(seriesRecord.ID)))
	}
}

func FindLatestCHLGameID() uint {
	db := dbprovider.GetInstance().GetDB()
	var games structs.CollegeGame
	err := db.Last(&games).Error
	if err != nil {
		return games.ID
	}

	return games.ID
}

func FindLatestPHLGameID() uint {
	db := dbprovider.GetInstance().GetDB()
	var games structs.ProfessionalGame
	err := db.Last(&games).Error
	if err != nil {
		return games.ID
	}

	return games.ID
}

func FindLatestPHLSeriesID() uint {
	db := dbprovider.GetInstance().GetDB()
	var series structs.ProSeries
	err := db.Last(&series).Error
	if err != nil {
		return series.ID
	}

	return series.ID
}

func FindCollegeSeriesRecords(seasonID string) []structs.CollegeSeries {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeSeries

	query := db.Model(&games)
	if len(seasonID) > 0 {
		query = query.Where("season_id = ?", seasonID)
	}
	if err := query.Find(&games).Error; err != nil {
		return []structs.CollegeSeries{}
	}

	return games
}

func FindCollegeSeriesRecord(id string) structs.CollegeSeries {
	db := dbprovider.GetInstance().GetDB()

	var games structs.CollegeSeries

	query := db.Model(&games)
	if len(id) > 0 {
		query = query.Where("id = ?", id)
	}
	if err := query.Find(&games).Error; err != nil {
		return structs.CollegeSeries{}
	}

	return games
}

func FindProSeriesRecords(seasonID string) []structs.ProSeries {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.ProSeries
	query := db.Model(&games)
	if len(seasonID) > 0 {
		query = query.Where("season_id = ?", seasonID)
	}
	if err := query.Find(&games).Error; err != nil {
		return []structs.ProSeries{}
	}
	return games
}

func SaveCollegeSeriesRecord(gameRecord structs.CollegeSeries, db *gorm.DB) {
	err := db.Save(&gameRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(gameRecord.ID)))
	}
}

func SaveProSeriesRecord(gameRecord structs.ProSeries, db *gorm.DB) {
	err := db.Save(&gameRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(gameRecord.ID)))
	}
}
