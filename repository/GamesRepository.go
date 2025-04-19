package repository

import (
	"log"
	"strconv"

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
