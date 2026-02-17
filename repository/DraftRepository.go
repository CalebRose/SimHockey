package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type ScoutProfileQuery struct {
	ID                       string
	PlayerID                 string
	TeamID                   string
	FilterOutRemovedProfiles bool
}

func FindDraftablePlayerRecord(clauses ScoutProfileQuery) structs.DraftablePlayer {
	db := dbprovider.GetInstance().GetDB()

	var draftablePlayer structs.DraftablePlayer

	query := db.Model(&draftablePlayer)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.PlayerID) > 0 {
		query = query.Where("id = ?", clauses.PlayerID)
	}

	if err := query.Find(&draftablePlayer).Error; err != nil {
		return structs.DraftablePlayer{}
	}

	return draftablePlayer
}

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

func SaveDraftablePlayerRecord(player structs.DraftablePlayer, db *gorm.DB) {
	err := db.Save(&player).Error
	if err != nil {
		log.Panicln("Could not save draftable player " + strconv.Itoa(int(player.ID)))
	}
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

func FindProWarRoomRecord(clause ScoutProfileQuery) structs.ProWarRoom {
	db := dbprovider.GetInstance().GetDB()
	warRooms := structs.ProWarRoom{}

	query := db.Model(&warRooms)

	if len(clause.TeamID) > 0 {
		query = query.Where("team_id = ?", clause.TeamID)
	}

	if err := query.Find(&warRooms).Error; err != nil {
		return structs.ProWarRoom{}
	}
	return warRooms
}

func FindProWarRooms() []structs.ProWarRoom {
	db := dbprovider.GetInstance().GetDB()
	warRooms := []structs.ProWarRoom{}

	if err := db.Find(&warRooms).Error; err != nil {
		return []structs.ProWarRoom{}
	}
	return warRooms
}

func SaveProWarRoomRecord(warRoom structs.ProWarRoom, db *gorm.DB) {
	err := db.Save(&warRoom).Error
	if err != nil {
		log.Panicln("Could not save pro war room " + strconv.Itoa(int(warRoom.ID)))
	}
}

func FindScoutingProfiles(clauses ScoutProfileQuery) []structs.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()
	var scoutingProfiles []structs.ScoutingProfile

	query := db.Model(&scoutingProfiles)
	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}
	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if clauses.FilterOutRemovedProfiles {
		query = query.Where("removed_from_board = ?", false)
	}

	if err := query.Find(&scoutingProfiles).Error; err != nil {
		return []structs.ScoutingProfile{}
	}
	return scoutingProfiles
}

func FindScoutingProfile(clauses ScoutProfileQuery) structs.ScoutingProfile {
	db := dbprovider.GetInstance().GetDB()
	var scoutingProfiles structs.ScoutingProfile

	query := db.Model(&scoutingProfiles)
	if len(clauses.ID) > 0 {
		query = query.Where("id = ?", clauses.ID)
	}

	if len(clauses.PlayerID) > 0 {
		query = query.Where("player_id = ?", clauses.PlayerID)
	}
	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if err := query.Find(&scoutingProfiles).Error; err != nil {
		return structs.ScoutingProfile{}
	}
	return scoutingProfiles
}

func CreateScoutingProfileRecord(profile structs.ScoutingProfile, db *gorm.DB) {
	err := db.Create(&profile).Error
	if err != nil {
		log.Panicln("Could not create scouting profile record")
	}
}

func SaveScoutingProfileRecord(profile structs.ScoutingProfile, db *gorm.DB) {
	err := db.Save(&profile).Error
	if err != nil {
		log.Panicln("Could not save scouting profile record")
	}
}

func DeleteScoutingProfileRecord(id string, db *gorm.DB) {
	err := db.Delete(&structs.ScoutingProfile{}, id).Error
	if err != nil {
		log.Panicln("Could not delete scouting profile record with id " + id)
	}
}
