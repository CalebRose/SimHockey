package repository

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type SchedulerQuery struct {
	ID       string
	TeamID   string
	SeasonID string
	WeekID   string
}

func FindCHLGameRequestRecord(q SchedulerQuery) structs.CHLGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var request structs.CHLGameRequest
	query := db.Model(&structs.CHLGameRequest{})
	if q.ID != "" {
		query = query.Where("id = ?", q.ID)
	}
	query.First(&request)
	return request
}

func FindCHLGameRequestRecords(q SchedulerQuery) []structs.CHLGameRequest {
	db := dbprovider.GetInstance().GetDB()
	var requests []structs.CHLGameRequest
	query := db.Model(&structs.CHLGameRequest{})
	if q.TeamID != "" {
		query = query.Where("home_team_id = ? OR away_team_id = ?", q.TeamID, q.TeamID)
	}
	if q.SeasonID != "" {
		query = query.Where("season_id = ?", q.SeasonID)
	}
	if q.WeekID != "" {
		query = query.Where("week_id = ?", q.WeekID)
	}
	query.Find(&requests)
	return requests
}

func CreateCHLGameRequest(request structs.CHLGameRequest, db *gorm.DB) {
	db.Create(&request)
}

func SaveCHLGameRequest(request structs.CHLGameRequest, db *gorm.DB) {
	db.Save(&request)
}

func DeleteCHLGameRequest(request structs.CHLGameRequest, db *gorm.DB) {
	db.Delete(&request)
}
