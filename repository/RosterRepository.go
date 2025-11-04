package repository

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type PlayerQuery struct {
	TeamID         string
	PlayerIDs      []string
	TransferStatus string
	LeagueID       string
	IsInjured      string
	IsFreeAgent    string
	OverallDesc    bool
}

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

func FindLatestGlobalPlayerRecord() structs.GlobalPlayer {
	db := dbprovider.GetInstance().GetDB()

	var lastPlayerRecord structs.GlobalPlayer
	err := db.Last(&lastPlayerRecord).Error
	if err != nil {
		return lastPlayerRecord
	}

	return lastPlayerRecord
}

// College Players
func FindAllCollegePlayers(clauses PlayerQuery) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.CollegePlayer

	query := db.Model(&CollegePlayers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
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

	if err := query.Find(&CollegePlayers).Error; err != nil {
		return []structs.CollegePlayer{}
	}

	return CollegePlayers
}

func FindCollegePlayersByTeamID(TeamID string) []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.CollegePlayer

	err := db.Order("overall desc").Where("team_id = ?", TeamID).Find(&CollegePlayers).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	return CollegePlayers
}

func FindAllHistoricCollegePlayers() []structs.HistoricCollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.HistoricCollegePlayer
	err := db.Find(&CollegePlayers).Error
	if err != nil {
		log.Printf("Error querying for college players: %v", err)

	}

	return CollegePlayers
}

// Professional Players
func FindAllProPlayers(clauses PlayerQuery) []structs.ProfessionalPlayer {
	db := dbprovider.GetInstance().GetDB()

	var proPlayers []structs.ProfessionalPlayer

	query := db.Model(&proPlayers)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.PlayerIDs) > 0 {
		query = query.Where("id in (?)", clauses.PlayerIDs)
	}

	if len(clauses.LeagueID) > 0 {
		query = query.Where("league_id = ?", clauses.LeagueID)
	}

	if len(clauses.IsInjured) > 0 {
		query = query.Where("is_injured = ?", true)
	}

	if len(clauses.IsFreeAgent) > 0 {
		query = query.Where("is_free_agent", true)
	}

	if clauses.OverallDesc {
		query = query.Order("overall desc")
	}

	if err := query.Find(&proPlayers).Error; err != nil {
		return []structs.ProfessionalPlayer{}
	}

	return proPlayers
}

func FindAllHistoricProPlayers() []structs.RetiredPlayer {
	db := dbprovider.GetInstance().GetDB()

	var retiredPlayers []structs.RetiredPlayer
	err := db.Find(&retiredPlayers).Error
	if err != nil {
		log.Printf("Error querying for college players: %v", err)

	}

	return retiredPlayers
}

func CreateRetiredPlayer(playerRecord structs.RetiredPlayer, db *gorm.DB) {
	err := db.Create(&playerRecord).Error
	if err != nil {
		log.Panicln("Could not create retired player " + strconv.Itoa(int(playerRecord.ID)))
	}
}

func DeleteProPlayerRecord(playerRecord structs.ProfessionalPlayer, db *gorm.DB) {
	err := db.Delete(&playerRecord).Error
	if err != nil {
		log.Panicln("Could not delete pro player " + strconv.Itoa(int(playerRecord.ID)))
	}
}
