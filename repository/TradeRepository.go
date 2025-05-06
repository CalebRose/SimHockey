package repository

import (
	"errors"
	"log"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type TradeClauses struct {
	IsAccepted          bool
	IsRejected          bool
	PreloadTradeOptions bool
}

func FindAllTradePreferenceRecords() []structs.TradePreferences {
	db := dbprovider.GetInstance().GetDB()
	preferences := []structs.TradePreferences{}

	db.Find(&preferences)

	return preferences
}

func FindAllTradeProposalsRecords() []structs.TradeProposal {
	db := dbprovider.GetInstance().GetDB()
	preferences := []structs.TradeProposal{}

	db.Preload("TeamTradeOptions").Find(&preferences)

	return preferences
}

func FindTradePreferencesByTeamID(id string) structs.TradePreferences {
	db := dbprovider.GetInstance().GetDB()
	preferences := structs.TradePreferences{}

	db.Where("id = ?", id).Find(&preferences)

	return preferences
}

func FindTradeProposalRecord(id string) structs.TradeProposal {
	db := dbprovider.GetInstance().GetDB()
	proposal := structs.TradeProposal{}

	db.Where("id = ?", id).Find(&proposal)

	return proposal
}

func SaveTradePreferencesRecord(tp structs.TradePreferences, db *gorm.DB) {
	err := db.Save(&tp).Error
	if err != nil {
		log.Panicln("Could not save timestamp")
	}
}

func SaveTradeProposalRecord(tp structs.TradeProposal, db *gorm.DB) {
	err := db.Save(&tp).Error
	if err != nil {
		log.Panicln("Could not save timestamp")
	}
}

func CreateTradeProposalRecord(db *gorm.DB, tp structs.TradeProposal) error {
	if err := db.Create(tp).Error; err != nil {
		return err
	}

	return nil
}

func CreateTradeOptionRecord(db *gorm.DB, tp structs.TradeProposal) error {
	if err := db.Create(tp).Error; err != nil {
		return err
	}

	return nil
}

func CreateTradeOptionRecordsBatch(db *gorm.DB, options []structs.TradeOption, batchSize int) error {
	total := len(options)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(options[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func FindLatestProposalInDB(db *gorm.DB) uint {
	var latestProposal structs.TradeProposal

	err := db.Last(&latestProposal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestProposal.ID + 1
}
