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

	err := db.Where("league = ? AND show_log = ?", "CHL", true).Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func GetAllPHLNewsLogs() []structs.NewsLog {
	db := dbprovider.GetInstance().GetDB()

	var logs []structs.NewsLog

	err := db.Where("league = ? AND show_log = ?", "PHL", true).Find(&logs).Error
	if err != nil {
		fmt.Println(err)
	}

	return logs
}

func CreateNewsLog(league, message, messageType string, teamID int, ts structs.Timestamp, showLog bool) {
	db := dbprovider.GetInstance().GetDB()

	news := CreateNewsLogObject(league, message, messageType, teamID, ts, showLog)

	repository.CreateNewsLog(news, db)
}

func CreateNewsLogObject(league, message, messageType string, teamID int, ts structs.Timestamp, showLog bool) structs.NewsLog {
	seasonID := ts.SeasonID
	weekID := ts.WeekID
	week := ts.Week
	return structs.NewsLog{
		League:      league,
		Message:     message,
		MessageType: messageType,
		SeasonID:    uint(seasonID),
		WeekID:      uint(weekID),
		Week:        uint(week),
		TeamID:      uint(teamID),
		ShowLog:     showLog,
	}
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

func GetNotificationById(id string) structs.Notification {
	notifications := repository.FindNotificationRecords(id, "", "")

	if len(notifications) == 0 {
		return structs.Notification{}
	}

	return notifications[0]
}

func ToggleNotification(id string) {
	db := dbprovider.GetInstance().GetDB()

	noti := GetNotificationById(id)

	if noti.ID == 0 {
		return
	}

	noti.ToggleIsRead()
	repository.SaveNotification(noti, db)
}

func DeleteNotification(id string) {
	db := dbprovider.GetInstance().GetDB()

	noti := GetNotificationById(id)

	if noti.ID == 0 {
		return
	}

	repository.DeleteNotificationRecord(noti, db)
}
