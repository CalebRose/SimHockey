package structs

import "time"

type MissedRecruitingUsers struct {
	User        string
	WeeksMissed uint
	TeamID      uint
	TeamName    string
}

type BrokenLineup struct {
	User           string
	TeamID         uint
	TeamName       string
	MissingPlayers uint   // If roster does not have enough players
	PlayerName     string // If a specific player is missing from the lineup due to injury, or if they transferred.
	Reason         string
}

type AdminReport struct {
	League        string
	TeamID        uint
	ReportType    string
	Message       string
	ResponseType  uint
	Status        string // "pending", "acknowledged", "resolved", "dismissed"
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AdminResponse string
	ResolvedBy    string
	ResolvedAt    time.Time
	Phase         uint
	WeekID        uint
	ReportedBy    string
}
