package structs

import "gorm.io/gorm"

type DraftPick struct {
	gorm.Model
	SeasonID               uint
	Season                 uint
	DrafteeID              uint
	DraftRound             uint
	DraftNumber            uint
	TeamID                 uint
	Team                   string
	OriginalTeamID         uint
	OriginalTeam           string
	PreviousTeamID         uint
	PreviousTeam           string
	DraftValue             float64
	Notes                  string
	SelectedPlayerID       uint
	SelectedPlayerName     string
	SelectedPlayerPosition string
	SelectedPlayerType     uint8
	PickupStatus           uint8
	IsCompensation         bool
	IsVoid                 bool
}

func (p *DraftPick) TradePick(id uint, team string) {
	p.PreviousTeamID = p.TeamID
	p.PreviousTeam = p.Team
	p.TeamID = id
	p.Team = team
	if p.PreviousTeamID == p.OriginalTeamID {
		p.Notes = "From " + p.OriginalTeam
	} else {
		p.Notes = "From " + p.PreviousTeam + " via " + p.OriginalTeam
	}
}

func (p *DraftPick) MapValuesToDraftPick(id, draftRound, draftNumber, teamID uint, team string, draftValue float64, isComp, isVoid bool) {
	if p.ID == 0 {
		p.ID = id
	}
	p.DraftRound = draftRound
	p.DraftNumber = draftNumber
	p.TeamID = teamID
	p.Team = team
	p.DraftValue = draftValue
	p.IsCompensation = isComp
	p.IsVoid = isVoid
}
