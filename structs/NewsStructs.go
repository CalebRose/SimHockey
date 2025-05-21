package structs

import "gorm.io/gorm"

type NewsLog struct {
	gorm.Model
	WeekID      uint
	Week        uint
	SeasonID    uint
	Season      uint
	MessageType string
	Message     string
	League      string
	TeamID      uint
	ShowLog     bool
}

type InboxResponse struct {
	CHLNotifications []Notification
	PHLNotifications []Notification
}

type Notification struct {
	gorm.Model
	TeamID           uint
	League           string
	NotificationType string
	Message          string
	Subject          string
	IsRead           bool
}

func (n *Notification) ToggleIsRead() {
	n.IsRead = !n.IsRead
}
