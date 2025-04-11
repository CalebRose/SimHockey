package managers

import (
	"fmt"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllCHLNewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "CHL").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetAllPHLNewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ?", "PHL").Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func CreateNewsLog(league, message, messageType string, teamID int, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	seasonID := ts.SeasonID
	weekID := ts.WeekID
	week := ts.Week

	news := structs.NewsLog{
		League:      league,
		Message:     message,
		MessageType: messageType,
		SeasonID:    uint(seasonID),
		WeekID:      uint(weekID),
		Week:        uint(week),
		TeamID:      uint(teamID),
	}

	repository.CreateNewsLog(news, db)
}

func GetNotificationByTeamIDAndLeague(league, teamID string) []structs.Notification {
	return repository.FindNotificationRecords("", league, teamID)
}

func CreateNotification(league, message, messageType string, teamID uint) {
	db := dbprovider.GetInstance().GetDB()

	notification := structs.Notification{
		League:           league,
		Message:          message,
		NotificationType: messageType,
		TeamID:           teamID,
	}

	repository.CreateNotification(notification, db)
}
