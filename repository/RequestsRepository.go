package repository

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type RequestQuery struct {
	ID       string
	TeamID   string
	Username string
	Role     string
}

func GetAllCHLTeamRequests(isApproved bool) []structs.CollegeTeamRequest {
	db := dbprovider.GetInstance().GetDB()
	var teamRequests []structs.CollegeTeamRequest

	query := db.Model(&teamRequests)
	if err := query.Where("is_approved = ?", isApproved).Find(&teamRequests).Error; err != nil {
		return []structs.CollegeTeamRequest{}
	}
	return teamRequests
}

func GetAllPHLTeamRequests(isApproved bool) []structs.ProTeamRequest {
	db := dbprovider.GetInstance().GetDB()
	var teamRequests []structs.ProTeamRequest

	query := db.Model(&teamRequests)
	if err := query.Where("is_approved = ?", isApproved).Find(&teamRequests).Error; err != nil {
		return []structs.ProTeamRequest{}
	}
	return teamRequests
}

func FindCHLRequestRecord(qu RequestQuery) structs.CollegeTeamRequest {
	db := dbprovider.GetInstance().GetDB()

	req := structs.CollegeTeamRequest{}

	query := db.Model(&req)

	if len(qu.ID) > 0 {
		query = query.Where("id = ?", qu.ID)
	}

	if len(qu.TeamID) > 0 {
		query = query.Where("team_id = ?", qu.TeamID)
	}

	if len(qu.Username) > 0 {
		query = query.Where("username = ?", qu.Username)
	}
	if err := query.Find(&req).Error; err != nil {
		return structs.CollegeTeamRequest{}
	}
	return req
}

func FindPHLRequestRecord(qu RequestQuery) structs.ProTeamRequest {
	db := dbprovider.GetInstance().GetDB()

	req := structs.ProTeamRequest{}

	query := db.Model(&req)

	if len(qu.ID) > 0 {
		query = query.Where("id = ?", qu.ID)
	}

	if len(qu.TeamID) > 0 {
		query = query.Where("team_id = ?", qu.TeamID)
	}

	if len(qu.Username) > 0 {
		query = query.Where("username = ?", qu.Username)
	}

	if len(qu.Role) > 0 {
		query = query.Where("role = ?", qu.Role)
	}

	if err := query.Find(&req).Error; err != nil {
		return structs.ProTeamRequest{}
	}
	return req
}

func CreateCHLTeamRequest(db *gorm.DB, record structs.CollegeTeamRequest) error {
	if err := db.Create(&record).Error; err != nil {
		return err
	}
	return nil
}

func CreatePHLTeamRequest(db *gorm.DB, record structs.ProTeamRequest) error {
	if err := db.Create(&record).Error; err != nil {
		return err
	}
	return nil
}

func SaveCHLTeamRequest(db *gorm.DB, record structs.CollegeTeamRequest) error {
	if err := db.Save(&record).Error; err != nil {
		return err
	}
	return nil
}

func SavePHLTeamRequest(db *gorm.DB, record structs.ProTeamRequest) error {
	if err := db.Save(&record).Error; err != nil {
		return err
	}
	return nil
}
