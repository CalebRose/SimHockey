package repository

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
)

type TransferPortalQuery struct {
	ProfileID        string
	CollegePlayerID  string
	RemovedFromBoard string
	OrderByPoints    bool
}

func FindTransferPortalProfileRecords(clauses TransferPortalQuery) []structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	var profiles []structs.TransferPortalProfile

	query := db.Model(&profiles)

	if len(clauses.ProfileID) > 0 {
		query = query.Where("profile_id = ?", clauses.ProfileID)
	}

	if len(clauses.CollegePlayerID) > 0 {
		query = query.Where("college_player_id = ?", clauses.CollegePlayerID)
	}

	if len(clauses.RemovedFromBoard) > 0 {
		isRemoved := clauses.RemovedFromBoard == "Y"
		query = query.Where("removed_from_board = ?", isRemoved)
	}

	if clauses.OrderByPoints {
		query = query.Order("total_points desc")
	}

	if err := query.Find(&profiles).Error; err != nil {
		return []structs.TransferPortalProfile{}
	}

	return profiles
}
