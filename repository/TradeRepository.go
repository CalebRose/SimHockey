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

func FindAllTradeProposalsRecords(clauses TradeClauses) []structs.TradeProposal {
	db := dbprovider.GetInstance().GetDB()
	proposal := []structs.TradeProposal{}

	query := db.Model(&proposal)

	if clauses.PreloadTradeOptions {
		query = query.Preload("TeamTradeOptions").Preload("RecepientTeamTradeOptions")
	}

	if clauses.IsAccepted {
		query = query.Where("is_trade_accepted = ?", true)
	}

	if clauses.IsRejected {
		query = query.Where("is_trade_rejected = ?", true)
	}

	if err := query.Find(&proposal).Error; err != nil {
		return []structs.TradeProposal{}
	}

	return proposal
}

func FindTradePreferencesByTeamID(id string) structs.TradePreferences {
	db := dbprovider.GetInstance().GetDB()
	preferences := structs.TradePreferences{}

	db.Where("id = ?", id).Find(&preferences)

	return preferences
}

func FindTradeProposalRecord(clauses TradeClauses, id string) structs.TradeProposal {
	db := dbprovider.GetInstance().GetDB()
	proposal := structs.TradeProposal{}

	query := db.Model(&proposal)

	if clauses.PreloadTradeOptions {
		query = query.Preload("TeamTradeOptions")
	}

	if clauses.IsAccepted {
		query = query.Where("is_trade_accepted = ?", true)
	}

	if clauses.IsRejected {
		query = query.Where("is_trade_rejected = ?", true)
	}

	if err := query.Where("id = ?", id).Find(&proposal).Error; err != nil {
		return structs.TradeProposal{}
	}

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
	if err := db.Create(&tp).Error; err != nil {
		return err
	}

	return nil
}

func CreateTradeOptionRecord(db *gorm.DB, tp structs.TradeProposal) error {
	if err := db.Create(&tp).Error; err != nil {
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

	err := db.Unscoped().Order("id DESC").First(&latestProposal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1
		}
		log.Fatalln("ERROR! Could not find latest record" + err.Error())
	}

	return latestProposal.ID + 1
}
