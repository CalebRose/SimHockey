package structs

import "gorm.io/gorm"

type BaseTeam struct {
	TeamName         string
	Mascot           string
	Abbreviation     string
	ConferenceID     uint8
	Conference       string
	Coach            string
	City             string
	State            string
	Country          string
	ArenaID          uint16
	Arena            string
	ArenaCapacity    uint16
	RecordAttendance uint16
	FirstPlayed      uint16
	DiscordID        string
	ColorOne         string
	ColorTwo         string
	ColorThree       string
	OverallGrade     string
	OffenseGrade     string
	DefenseGrade     string
}

type CollegeTeam struct {
	gorm.Model
	BaseTeam
	IsUserCoached        bool
	IsClub               bool
	IsActive             bool
	PlayersProgressed    bool
	ProgramPrestige      uint8
	ProfessionalPrestige uint8
	Traditions           uint8
	Facilities           uint8
	Atmosphere           uint8
	AcademicPrestige     uint8
	ConferencePrestige   uint8
	CoachReputation      uint8
}

type ProfessionalTeam struct {
	gorm.Model
	BaseTeam
}

type Arena struct {
	gorm.Model
	Name             string
	TeamID           uint
	City             string
	State            string
	Country          string
	Capacity         uint16
	RecordAttendance uint16
}
