package repository

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
)

type StandingsQuery struct {
	SeasonID     string
	ConferenceID string
	TeamID       string
}

func FindAllCollegeStandings(clauses StandingsQuery) []structs.CollegeStandings {
	var standings []structs.CollegeStandings
	db := dbprovider.GetInstance().GetDB()

	query := db.Model(&standings)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if len(clauses.ConferenceID) > 0 {
		query = query.Where("conference_id = ?", clauses.ConferenceID)
	}
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}

	if err := query.Order("conference_losses asc").Order("conference_wins desc").
		Order("total_losses asc").Order("total_wins desc").Find(&standings).Error; err != nil {
		return []structs.CollegeStandings{}
	}

	return standings
}

func FindAllProfessionalStandings(clauses StandingsQuery) []structs.ProfessionalStandings {
	var standings []structs.ProfessionalStandings
	db := dbprovider.GetInstance().GetDB()

	query := db.Model(&standings)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}
	if len(clauses.ConferenceID) > 0 {
		query = query.Where("conference_id = ?", clauses.ConferenceID)
	}
	if len(clauses.SeasonID) > 0 {
		query = query.Where("season_id = ?", clauses.SeasonID)
	}

	if err := query.Order("points desc").Find(&standings).Error; err != nil {
		return []structs.ProfessionalStandings{}
	}

	return standings
}
