package repository

import (
	"log"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func FindNotificationRecords(id, league, teamID string) []structs.Notification {
	db := dbprovider.GetInstance().GetDB()

	noti := []structs.Notification{}

	query := db.Model(&noti)

	if len(id) > 0 {
		query = query.Where("id = ?", id)
	}

	if len(league) > 0 {
		query = query.Where("league = ?", league)
	}

	if len(teamID) > 0 {
		query = query.Where("team_id = ?", teamID)
	}

	if err := query.Find(&noti).Error; err != nil {
		return []structs.Notification{}
	}
	return noti
}

func CreateNewsLog(news structs.NewsLog, db *gorm.DB) {
	err := db.Create(&news).Error
	if err != nil {
		log.Panicln("Could not create notification record!")
	}
}

func CreateNewsLogRecordsBatch(db *gorm.DB, newsLogs []structs.NewsLog, batchSize int) error {
	total := len(newsLogs)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)

		if err := db.CreateInBatches(newsLogs[i:end], batchSize).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNotification(noti structs.Notification, db *gorm.DB) {
	err := db.Create(&noti).Error
	if err != nil {
		log.Panicln("Could not create notification record!")
	}
}

func SaveNotification(noti structs.Notification, db *gorm.DB) {
	err := db.Save(&noti).Error
	if err != nil {
		log.Panicln("Could not save notification record!")
	}
}

func DeleteNotificationRecord(noti structs.Notification, db *gorm.DB) {
	err := db.Delete(&noti).Error
	if err != nil {
		log.Panicln("Could not delete old notification record.")
	}
}
