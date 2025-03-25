package repository

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

type FaceDataQuery struct {
	PlayerIDs []string
}

func CreateFaceRecordsBatch(db *gorm.DB, fds []structs.FaceData, batchSize int) error {
	total := len(fds)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(fds[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func FindFaceDataRecords(faceDataQuery FaceDataQuery) []structs.FaceData {
	db := dbprovider.GetInstance().GetDB()

	var facialData []structs.FaceData

	query := db.Model(&facialData)

	if len(faceDataQuery.PlayerIDs) > 0 {
		query = query.Where("player_id in (?)", faceDataQuery.PlayerIDs)
	}

	if err := query.Find(&facialData).Error; err != nil {
		return []structs.FaceData{}
	}
	return facialData
}
