package repository

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
)

func FindDraftPickRecord(id string) structs.DraftPick {
	db := dbprovider.GetInstance().GetDB()
	preferences := structs.DraftPick{}

	db.Where("id = ?", id).Find(&preferences)

	return preferences
}
