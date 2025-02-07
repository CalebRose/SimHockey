package repository

import (
	"log"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

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

func FindAllCollegeTeams() []structs.CollegeTeam {
	db := dbprovider.GetInstance().GetDB()

	var CollegeTeams []structs.CollegeTeam
	err := db.Find(&CollegeTeams).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

	}

	return CollegeTeams
}

func FindAllProTeams() []structs.ProfessionalTeam {
	db := dbprovider.GetInstance().GetDB()

	var proTeams []structs.ProfessionalTeam
	err := db.Find(&proTeams).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

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
