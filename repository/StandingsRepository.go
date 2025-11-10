package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type StandingsQuery struct {
	SeasonID     string
	ConferenceID string
	TeamID       string
}

func FindAllCollegeStandings(clauses StandingsQuery) []structs.CollegeStandings {
	var standings []structs.CollegeStandings
	db := dbprovider.GetInstance().GetDB()

	query := db.Model(&standings)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if len(clauses.ConferenceID) > 0 {
		query = query.Where("conference_id = ?", clauses.ConferenceID)
	}
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}

	if err := query.Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings).Error; err != nil {
		return []structs.CollegeStandings{}
	}

	return standings
}

func FindAllProfessionalStandings(clauses StandingsQuery) []structs.ProfessionalStandings {
	var standings []structs.ProfessionalStandings
	db := dbprovider.GetInstance().GetDB()

	query := db.Model(&standings)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if len(clauses.ConferenceID) > 0 {
		query = query.Where("conference_id = ?", clauses.ConferenceID)
	}
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}

	if err := query.Order("points desc").Find(&standings).Error; err != nil {
		return []structs.ProfessionalStandings{}
	}

	return standings
}

func CreateCollegeStandingsRecord(standingsRecord structs.CollegeStandings, db *gorm.DB) {
	err := db.Create(&standingsRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(standingsRecord.ID)))
	}
}

func CreateProfessionalStandingsRecord(standingsRecord structs.ProfessionalStandings, db *gorm.DB) {
	err := db.Create(&standingsRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(standingsRecord.ID)))
	}
}

func CreateCollegeStandingsRecordsBatch(db *gorm.DB, players []structs.CollegeStandings, batchSize int) error {
	total := len(players)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(players[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateProStandingsRecordsBatch(db *gorm.DB, players []structs.ProfessionalStandings, batchSize int) error {
	total := len(players)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(players[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func SaveCollegeStandingsRecord(standingsRecord structs.CollegeStandings, db *gorm.DB) {
	err := db.Save(&standingsRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(standingsRecord.ID)))
	}
}

func SaveProfessionalStandingsRecord(standingsRecord structs.ProfessionalStandings, db *gorm.DB) {
	err := db.Save(&standingsRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(standingsRecord.ID)))
	}
}
