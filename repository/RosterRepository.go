package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func FindCollegePlayer(id string) structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayer structs.CollegePlayer

	db.Where("id = ?", id).Find(&CollegePlayer)

	return CollegePlayer
}

func SaveCollegeHockeyPlayerRecord(playerRecord structs.CollegePlayer, db *gorm.DB) {
	err := db.Save(&playerRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(playerRecord.ID)))
	}
}

func FindProPlayer(id string) structs.ProfessionalPlayer {
	db := dbprovider.GetInstance().GetDB()

	var proPlayer structs.ProfessionalPlayer

	db.Where("id = ?", id).Find(&proPlayer)

	return proPlayer
}

func SaveProPlayerRecord(playerRecord structs.ProfessionalPlayer, db *gorm.DB) {
	err := db.Save(&playerRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(playerRecord.ID)))
	}
}

func FindProContract(playerID string) structs.ProContract {
	db := dbprovider.GetInstance().GetDB()

	var proPlayer structs.ProContract

	db.Where("player_id = ?", playerID).Find(&proPlayer)

	return proPlayer
}

func SaveProContractRecord(record structs.ProContract, db *gorm.DB) {
	err := db.Save(&record).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(record.ID)))
	}
}

func SaveProCapsheetRecord(record structs.ProCapsheet, db *gorm.DB) {
	err := db.Save(&record).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(record.ID)))
	}
}
