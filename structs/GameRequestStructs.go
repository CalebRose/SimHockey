package structs

import "gorm.io/gorm"

// GameRequest is the base struct for OOC game scheduling requests.
type GameRequest struct {
	gorm.Model
	HomeTeamID       uint
	AwayTeamID       uint
	SendingTeamID    uint
	RequestingTeamID uint
	IsAccepted       bool
	IsApproved       bool
	ArenaID          uint
	Arena            string
	IsNeutralSite    bool
	SeasonID         uint
	WeekID           uint
	Week             uint
	Timeslot         string
}

func (g *GameRequest) Accepted() {
	g.IsAccepted = true
}

func (g *GameRequest) Approved() {
	g.IsApproved = true
}

// CHLGameRequest represents an OOC game scheduling request for SimCHL.
type CHLGameRequest struct {
	GameRequest
}
