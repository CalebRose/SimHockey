package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type TeamClauses struct {
	IDs           []string
	ID            string
	LeagueID      string
	IsUserCoached string
}

func FindCollegeTeamRecord(id string) structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var collegeTeam structs.CollegeTeam
	err := db.Where("id = ?", id).Find(&collegeTeam).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

	}

	return collegeTeam
}

func FindProTeamRecord(id string) structs.ProfessionalTeam {
	db := dbprovider.GetInstance().GetDB()

	var proTeam structs.ProfessionalTeam
	err := db.Where("id = ?", id).Find(&proTeam).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

	}

	return proTeam
}

func FindAllCollegeTeams(clauses TeamClauses) []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var CollegeTeams []structs.CollegeTeam

	query := db.Model(&CollegeTeams)

	if len(clauses.IDs) > 0 {
		query = query.Where("id in (?)", clauses.IDs)
	}

	if len(clauses.ID) > 0 {
		query = query.Where("id = ?", clauses.ID)
	}

	if len(clauses.LeagueID) > 0 {
		query = query.Where("league_id = ?", clauses.LeagueID)
	}

	if len(clauses.IsUserCoached) > 0 {
		query = query.Where("is_user_coached = ?", true)
	}

	if err := query.Find(&CollegeTeams).Error; err != nil {
		return []structs.CollegeTeam{}
	}

	return CollegeTeams
}

func FindAllProTeams(clauses TeamClauses) []structs.ProfessionalTeam {
	db := dbprovider.GetInstance().GetDB()

	var proTeams []structs.ProfessionalTeam
	query := db.Model(&proTeams)

	if len(clauses.IDs) > 0 {
		query = query.Where("id in (?)", clauses.IDs)
	}

	if len(clauses.ID) > 0 {
		query = query.Where("id = ?", clauses.ID)
	}

	if len(clauses.LeagueID) > 0 {
		query = query.Where("league_id = ?", clauses.LeagueID)
	}

	if len(clauses.IsUserCoached) > 0 {
		query = query.Where("is_user_coached = ?", true)
	}

	if err := query.Find(&proTeams).Error; err != nil {
		return []structs.ProfessionalTeam{}
	}

	return proTeams
}

func SaveCollegeTeamRecord(db *gorm.DB, teamRecord structs.CollegeTeam) {
	// Will need to potentially add stats parameters here
	err := db.Save(&teamRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(teamRecord.ID)))
	}
}

func SaveProTeamRecord(db *gorm.DB, teamRecord structs.ProfessionalTeam) {
	// Will need to potentially add stats parameters here
	err := db.Save(&teamRecord).Error
	if err != nil {
		log.Panicln("Could not save college player " + strconv.Itoa(int(teamRecord.ID)))
	}
}
