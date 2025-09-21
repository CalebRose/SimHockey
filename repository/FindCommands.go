package repository

import (
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

func FindAllCollegeShootoutLineups() []structs.CollegeShootoutLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.CollegeShootoutLineup
	err := db.Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college lineups: %v", err)
	}

	return lineups
}

func FindCollegeShootoutLineupByTeamID(TeamID string) structs.CollegeShootoutLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups structs.CollegeShootoutLineup
	err := db.Where("team_id = ?", TeamID).Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college shootout lineups: %v", err)
	}

	return lineups
}

func FindAllProLineups() []structs.ProfessionalLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.ProfessionalLineup
	err := db.Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)

	}

	return lineups
}

func FindProLineupsByTeamID(TeamID string) []structs.ProfessionalLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.ProfessionalLineup
	err := db.Where("team_id = ?", TeamID).Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college teams: %v", err)
	}

	return lineups
}

func FindAllProShootoutLineups() []structs.ProfessionalShootoutLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups []structs.ProfessionalShootoutLineup
	err := db.Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college lineups: %v", err)
	}

	return lineups
}

func FindProShootoutLineupByTeamID(TeamID string) structs.ProfessionalShootoutLineup {
	db := dbprovider.GetInstance().GetDB()

	var lineups structs.ProfessionalShootoutLineup
	err := db.Where("team_id = ?", TeamID).Find(&lineups).Error
	if err != nil {
		log.Printf("Error querying for college shootout lineups: %v", err)
	}

	return lineups
}

func FindAllCollegePolls(weekID, seasonID, username string) []structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()
	submissions := []structs.CollegePollSubmission{}

	query := db.Model(&submissions)

	// Add conditional filtering based on provided parameters
	if len(weekID) > 0 && len(seasonID) > 0 {
		query = query.Where("week_id = ? AND season_id = ?", weekID, seasonID)
	}
	if len(username) > 0 {
		query = query.Where("username = ?", username)
	}

	// Execute the query and handle errors
	if err := query.Find(&submissions).Error; err != nil {
		return []structs.CollegePollSubmission{}
	}

	return submissions
}

func FindAllCollegeStandings(seasonID, conferenceID, teamID string) []structs.CollegeStandings {
	var standings []structs.CollegeStandings
	db := dbprovider.GetInstance().GetDB()

	query := db.Model(&standings)

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}
	if len(conferenceID) > 0 {
		query = query.Where("conference_id = ?", conferenceID)
	}
	if len(seasonID) > 0 {
		query = query.Where("season_id = ?", seasonID)
	}

	if err := query.Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings).Error; err != nil {
		return []structs.CollegeStandings{}
	}

	return standings
}

func FindAllProfessionalStandings(seasonID, conferenceID, teamID string) []structs.ProfessionalStandings {
	var standings []structs.ProfessionalStandings
	db := dbprovider.GetInstance().GetDB()

	query := db.Model(&standings)

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}
	if len(conferenceID) > 0 {
		query = query.Where("conference_id = ?", conferenceID)
	}
	if len(seasonID) > 0 {
		query = query.Where("season_id = ?", seasonID)
	}

	if err := query.Order("points desc").Find(&standings).Error; err != nil {
		return []structs.ProfessionalStandings{}
	}

	return standings
}

func FindCollegePollSubmission(id, weekID, seasonID, username string) structs.CollegePollSubmission {
	db := dbprovider.GetInstance().GetDB()

	submission := structs.CollegePollSubmission{}
	// Add conditional filtering based on provided parameters
	query := db.Model(&submission)
	if len(weekID) > 0 && len(seasonID) > 0 {
		query = query.Where("week_id = ? AND season_id = ?", weekID, seasonID)
	}
	if len(username) > 0 {
		query = query.Where("username = ?", username)
	}
	if len(id) > 0 {
		query = query.Where("id = ?", id)
	}

	if err := query.Find(&submission).Error; err != nil {
		return structs.CollegePollSubmission{}
	}

	return submission
}

func FindCollegePollOfficial(seasonID string) []structs.CollegePollOfficial {
	db := dbprovider.GetInstance().GetDB()
	officialPoll := []structs.CollegePollOfficial{}
	query := db.Model(&officialPoll)
	if len(seasonID) > 0 {
		query = query.Where("season_id = ?", seasonID)
	}

	if err := query.Find(&officialPoll).Error; err != nil {
		return []structs.CollegePollOfficial{}
	}

	return officialPoll
}

func FindCapsheetRecords() []structs.ProCapsheet {
	db := dbprovider.GetInstance().GetDB()

	capsheets := []structs.ProCapsheet{}

	err := db.Find(&capsheets).Error
	if err != nil {
		log.Fatal(err)
	}
	return capsheets
}
