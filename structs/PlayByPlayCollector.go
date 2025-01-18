package structs

type PbPCollector struct {
	PlayByPlays []PlayByPlay
}

func (c *PbPCollector) AppendPlay(play PlayByPlay) {
	c.PlayByPlays = append(c.PlayByPlays, play)
}

type PlayByPlay struct {
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
