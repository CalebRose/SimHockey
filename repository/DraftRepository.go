package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func FindAllDraftablePlayers(clauses PlayerQuery) []structs.DraftablePlayer {
	db := dbprovider.GetInstance().GetDB()

	var DraftablePlayers []structs.DraftablePlayer

	query := db.Model(&DraftablePlayers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.CollegeID) > 0 {
		query = query.Where("college_id = ?", clauses.CollegeID)
	}

	if len(clauses.PlayerIDs) > 0 {
		query = query.Where("id in (?)", clauses.PlayerIDs)
	}

	if len(clauses.TransferStatus) > 0 {
		query = query.Where("transfer_status = ?", clauses.TransferStatus)
	}

	if len(clauses.LeagueID) > 0 {
		query = query.Where("league_id = ?", clauses.LeagueID)
	}

	if len(clauses.IsInjured) > 0 {
		query = query.Where("is_injured = ?", true)
	}

	if clauses.OverallDesc {
		query = query.Order("overall desc")
	}

	if err := query.Find(&DraftablePlayers).Error; err != nil {
		return []structs.DraftablePlayer{}
	}

	return DraftablePlayers
}

func FindDraftPickRecord(id string) structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()
	preferences := structs.DraftPick{}

	db.Where("id = ?", id).Find(&preferences)

	return preferences
}

func FindDraftPicks(seasonID string) []structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()
	preferences := []structs.DraftPick{}

	query := db.Model(&preferences)

	if len(seasonID) > 0 {
		query = query.Where("season_id = ?", seasonID)
	}

	if err := query.Find(&preferences).Error; err != nil {
		return []structs.DraftPick{}
	}

	return preferences
}

func CreateDraftPickRecordsBatch(db *gorm.DB, picks []structs.DraftPick, batchSize int) error {
	total := len(picks)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(picks[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func SaveDraftPickRecord(pickRecord structs.DraftPick, db *gorm.DB) {
	err := db.Save(&pickRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(pickRecord.ID)))
	}
}
