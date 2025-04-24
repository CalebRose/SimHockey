package structs

import "gorm.io/gorm"

type PbPCollector struct {
	PlayByPlays []PbP
}

func (c *PbPCollector) AppendPlay(play PbP) {
	c.PlayByPlays = append(c.PlayByPlays, play)
}

type PbP struct {
	GameID            uint
	Period            uint8
	TimeOnClock       uint16
	SecondsConsumed   uint8
	EventID           uint8 // Enum
	ZoneID            uint8 // Enum
	NextZoneID        uint8 // Enum
	Outcome           uint8 // Enum
	HomeTeamScore     uint8
	AwayTeamScore     uint8
	TeamID            uint8
	PuckCarrierID     uint
	PassedPlayerID    uint
	AssistingPlayerID uint
	DefenderID        uint
	GoalieID          uint
	InjuryID          uint8
	InjuryType        uint8
	InjuryDuration    uint8
	PenaltyID         uint8
	Severity          uint8
	IsFight           bool
	IsBreakaway       bool
}

type CollegePlayByPlay struct {
	gorm.Model
	PbP
}

type ProPlayByPlay struct {
	gorm.Model
	PbP
}
