package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func CreateCollegeLineupRecordsBatch(db *gorm.DB, lineups []structs.CollegeLineup, batchSize int) error {
	total := len(lineups)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(lineups[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateProfessionalLineupRecordsBatch(db *gorm.DB, lineups []structs.ProfessionalLineup, batchSize int) error {
	total := len(lineups)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		if err := db.CreateInBatches(lineups[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateCollegeGameplanRecordsBatch(db *gorm.DB, lineups []structs.CollegeGameplan, batchSize int) error {
	total := len(lineups)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(lineups[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateProfessionalGameplanRecordsBatch(db *gorm.DB, lineups []structs.ProGameplan, batchSize int) error {
	total := len(lineups)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(lineups[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func SaveCollegeLineupRecord(lineupRecord structs.CollegeLineup, db *gorm.DB) {
	err := db.Save(&lineupRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(lineupRecord.ID)))
	}
}

func SaveCollegeShootoutLineupRecord(lineupRecord structs.CollegeShootoutLineup, db *gorm.DB) {
	err := db.Save(&lineupRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(lineupRecord.TeamID)))
	}
}

func SaveProfessionalLineupRecord(lineupRecord structs.ProfessionalLineup, db *gorm.DB) {
	err := db.Save(&lineupRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(lineupRecord.ID)))
	}
}

func SaveProfessionalShootoutLineupRecord(lineupRecord structs.ProfessionalShootoutLineup, db *gorm.DB) {
	err := db.Save(&lineupRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(lineupRecord.TeamID)))
	}
}

func SaveCollegeGameplanRecord(lineupRecord structs.CollegeGameplan, db *gorm.DB) {
	err := db.Save(&lineupRecord).Error
	if err != nil {
		log.Panicln("Could not save college gameplan " + strconv.Itoa(int(lineupRecord.ID)))
	}
}

func SaveProfessionalGameplanRecord(lineupRecord structs.ProGameplan, db *gorm.DB) {
	err := db.Save(&lineupRecord).Error
	if err != nil {
		log.Panicln("Could not save pro gameplan " + strconv.Itoa(int(lineupRecord.ID)))
	}
}

func FindCollegeGameplanRecord(id string) structs.CollegeGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans structs.CollegeGameplan

	query := db.Model(&gameplans)

	if err := query.Where("team_id = ?", id).Find(&gameplans).Error; err != nil {
		return structs.CollegeGameplan{}
	}

	return gameplans
}

func FindProGameplanRecord(id string) structs.ProGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans structs.ProGameplan

	query := db.Model(&gameplans)

	if err := query.Where("team_id = ?", id).Find(&gameplans).Error; err != nil {
		return structs.ProGameplan{}
	}

	return gameplans
}

func FindCollegeGameplanRecords() []structs.CollegeGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans []structs.CollegeGameplan

	query := db.Model(&gameplans)

	if err := query.Find(&gameplans).Error; err != nil {
		return []structs.CollegeGameplan{}
	}

	return gameplans
}

func FindProfessionalGameplanRecords() []structs.ProGameplan {
	db := dbprovider.GetInstance().GetDB()

	var gameplans []structs.ProGameplan

	query := db.Model(&gameplans)

	if err := query.Find(&gameplans).Error; err != nil {
		return []structs.ProGameplan{}
	}

	return gameplans
}
