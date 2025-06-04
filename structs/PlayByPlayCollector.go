package structs

import "gorm.io/gorm"

type PbPCollector struct {
	PlayByPlays []PbP
}

func (c *PbPCollector) AppendPlay(play PbP) {
	c.PlayByPlays = append(c.PlayByPlays, play)
}

type PbP struct {
	GameID                uint
	Period                uint8
	TimeOnClock           uint16
	SecondsConsumed       uint8
	EventID               uint8 // Enum
	ZoneID                uint8 // Enum
	NextZoneID            uint8 // Enum
	Outcome               uint8 // Enum
	HomeTeamScore         uint8
	AwayTeamScore         uint8
	HomeTeamShootoutScore uint8
	AwayTeamShootoutScore uint8
	TeamID                uint8
	PuckCarrierID         uint
	PassedPlayerID        uint
	AssistingPlayerID     uint
	DefenderID            uint
	GoalieID              uint
	InjuryID              uint8
	InjuryType            uint8
	InjuryDuration        uint8
	PenaltyID             uint8
	Severity              uint8
	IsFight               bool
	IsBreakaway           bool
	IsShootout            bool
}

type CollegePlayByPlay struct {
	gorm.Model
	PbP
}

type ProPlayByPlay struct {
	gorm.Model
	PbP
}

type PlayByPlayResponse struct {
	GameID            uint
	PlayNumber        uint
	HomeTeamID        uint
	HomeTeamScore     uint8
	AwayTeamID        uint
	AwayTeamScore     uint8
	Period            uint8
	TimeOnClock       string
	SecondsConsumed   uint8
	Event             string // Enum
	Zone              string // Enum
	NextZone          string // Enum
	Outcome           string // Enum
	TeamID            uint8
	PuckCarrierID     uint
	PassedPlayerID    uint
	AssistingPlayerID uint
	DefenderID        uint
	GoalieID          uint
	InjuryID          uint8
	InjuryType        uint8
	InjuryDuration    uint8
	Penalty           string
	Severity          string
	IsFight           string
	IsBreakaway       bool
	Result            string
	StreamResult      []string
}

func (p *PlayByPlayResponse) AddPlayInformation(toc, event, zone, nextZone, outcome string) {
	p.TimeOnClock = toc
	p.Event = event
	p.Zone = zone
	p.NextZone = nextZone
	p.Outcome = outcome
}

func (p *PlayByPlayResponse) AddResult(result []string, isStream bool) {
	if isStream {
		p.StreamResult = result
	} else {
		p.Result = result[0]
	}
}

type GameResultsPlayer struct {
	ID uint
	BasePlayer
	BaseInjuryData
	BasePlayerStats
}

type StreamResponse struct {
	GameID            uint
	HomeTeamID        uint
	GameLabel         string
	HomeLabel         string
	HomeTeam          string
	HomeTeamCoach     string
	HomeTeamDiscordID string
	HomeTeamRank      uint
	AwayTeamID        uint
	AwayLabel         string
	AwayTeam          string
	AwayTeamCoach     string
	AwayTeamDiscordID string
	AwayTeamRank      uint
	ArenaID           uint
	Arena             string
	Capacity          uint
	Attendance        uint
	City              string
	State             string
	Country           string
	HomeTeamWin       bool
	AwayTeamWin       bool
	HTScore           uint
	ATScore           uint
	HTSScore          uint
	ATSScore          uint
	Streams           []PlayByPlayResponse
}
