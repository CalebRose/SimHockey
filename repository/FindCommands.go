package repository

import (
	"fmt"
	"log"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
)

func FindTimestamp() structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()

	var timestamp structs.Timestamp

	err := db.First(&timestamp).Error
	if err != nil {
		log.Printf("Error querying for timestamp: %v", err)
	}

	return timestamp
}

// College Players
func FindAllCollegePlayers() []structs.CollegePlayer {
	db := dbprovider.GetInstance().GetDB()

	var CollegePlayers []structs.CollegePlayer

	err := db.Find(&CollegePlayers).Error
	if err != nil {
		log.Printf("Error querying for college players: %v", err)

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

func FindAllCollegeTeams() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var CollegeTeams []structs.CollegeTeam
	err := db.Find(&CollegeTeams).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

	}

	return CollegeTeams
}

func FindAllAvailableCollegeTeams() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var teams []structs.CollegeTeam

	db.Where("coach IN (?,?)", "", "AI").Find(&teams)

	return teams
}

func FindAllCollegeLineups() []structs.CollegeLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.CollegeLineup
	err := db.Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

	}

	return lineups
}

func FindCollegeLineupsByTeamID(TeamID string) []structs.CollegeLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.CollegeLineup
	err := db.Where("team_id = ?", TeamID).Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

	}

	return lineups
}

func FindCollegeGamesByCurrentMatchup(weekID, seasonID, gameDay string) []structs.CollegeGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.CollegeGame
	err := db.Where("week_id = ? AND season_id = ? AND game_day = ?", weekID, seasonID, gameDay).Find(&games).Error
	if err != nil {
		log.Printf("Error querying for college games: %v", err)

	}

	return games
}

func FindProfessionalGamesByCurrentMatchup(weekID, seasonID, gameDay string) []structs.ProfessionalGame {
	db := dbprovider.GetInstance().GetDB()

	var games []structs.ProfessionalGame
	err := db.Where("week_id = ? AND season_id = ? AND game_day = ?", weekID, seasonID, gameDay).Find(&games).Error
	if err != nil {
		log.Printf("Error querying for professional games: %v", err)

	}

	return games
}

func FindAllArenas() []structs.Arena {
	var arenas []structs.Arena
	db := dbprovider.GetInstance().GetDB()
	err := db.Find(&arenas).Error
	if err != nil {
		log.Fatal(err)
	}
	return arenas
}
