package repository

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
)

type ArenaQuery struct {
	Country       string
	GreaterThanID string
	LessThanID    string
	Capacity      string
	TeamID        string
}

func FindAllArenas(clauses ArenaQuery) []structs.Arena {
	db := dbprovider.GetInstance().GetDB()

	var arenas []structs.Arena

	query := db.Model(&arenas)

	if len(clauses.TeamID) > 0 {
		query = query.Where("team_id = ?", clauses.TeamID)
	}

	if len(clauses.GreaterThanID) > 0 {
		query = query.Where("id > ?", clauses.GreaterThanID)
	}
	if len(clauses.LessThanID) > 0 {
		query = query.Where("id < ?", clauses.LessThanID)
	}

	if len(clauses.Capacity) > 0 {
		query = query.Where("capacity > ?", clauses.Capacity)
	}

	if len(clauses.Country) > 0 {
		query = query.Where("country = ?", clauses.Country)
	}

	if err := query.Find(&arenas).Error; err != nil {
		return []structs.Arena{}
	}

	return arenas
}
