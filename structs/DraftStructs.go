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
	IsCalledUp             bool
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

func (p *DraftPick) MapDraftPickResults(drafteeID uint) {
	p.DrafteeID = drafteeID
}

type ScoutingProfile struct {
	gorm.Model
	PlayerID           uint
	TeamID             uint
	ShowAttribute1     bool
	ShowAttribute2     bool
	ShowAttribute3     bool
	ShowAttribute4     bool
	ShowAttribute5     bool
	ShowAttribute6     bool
	ShowAttribute7     bool
	ShowAttribute8     bool
	ShowAttribute9     bool
	ShowAttribute10    bool
	ShowAttribute11    bool
	ShowAttribute12    bool
	ShowAttribute13    bool
	ShowAttribute14    bool
	ShowAttribute15    bool
	ShowPotAttribute1  bool
	ShowPotAttribute2  bool
	ShowPotAttribute3  bool
	ShowPotAttribute4  bool
	ShowPotAttribute5  bool
	ShowPotAttribute6  bool
	ShowPotAttribute7  bool
	ShowPotAttribute8  bool
	ShowPotAttribute9  bool
	ShowPotAttribute10 bool
	ShowPotAttribute11 bool
	ShowPotAttribute12 bool
	ShowPotAttribute13 bool
	ShowPotAttribute14 bool
	ShowPotAttribute15 bool
	RemovedFromBoard   bool
	ShowCount          uint8
}

func (sp *ScoutingProfile) RevealAttribute(attr string) {
	switch attr {
	case "ShowAttribute1":
		sp.ShowAttribute1 = true
	case "ShowAttribute2":
		sp.ShowAttribute2 = true
	case "ShowAttribute3":
		sp.ShowAttribute3 = true
	case "ShowAttribute4":
		sp.ShowAttribute4 = true
	case "ShowAttribute5":
		sp.ShowAttribute5 = true
	case "ShowAttribute6":
		sp.ShowAttribute6 = true
	case "ShowAttribute7":
		sp.ShowAttribute7 = true
	case "ShowAttribute8":
		sp.ShowAttribute8 = true
	case "ShowAttribute9":
		sp.ShowAttribute9 = true
	case "ShowAttribute10":
		sp.ShowAttribute10 = true
	case "ShowAttribute11":
		sp.ShowAttribute11 = true
	case "ShowAttribute12":
		sp.ShowAttribute12 = true
	case "ShowAttribute13":
		sp.ShowAttribute13 = true
	case "ShowAttribute14":
		sp.ShowAttribute14 = true
	case "ShowAttribute15":
		sp.ShowAttribute15 = true
	case "ShowPotAttribute1":
		sp.ShowPotAttribute1 = true
	case "ShowPotAttribute2":
		sp.ShowPotAttribute2 = true
	case "ShowPotAttribute3":
		sp.ShowPotAttribute3 = true
	case "ShowPotAttribute4":
		sp.ShowPotAttribute4 = true
	case "ShowPotAttribute5":
		sp.ShowPotAttribute5 = true
	case "ShowPotAttribute6":
		sp.ShowPotAttribute6 = true
	case "ShowPotAttribute7":
		sp.ShowPotAttribute7 = true
	case "ShowPotAttribute8":
		sp.ShowPotAttribute8 = true
	case "ShowPotAttribute9":
		sp.ShowPotAttribute9 = true
	case "ShowPotAttribute10":
		sp.ShowPotAttribute10 = true
	case "ShowPotAttribute11":
		sp.ShowPotAttribute11 = true
	case "ShowPotAttribute12":
		sp.ShowPotAttribute12 = true
	case "ShowPotAttribute13":
		sp.ShowPotAttribute13 = true
	case "ShowPotAttribute14":
		sp.ShowPotAttribute14 = true
	case "ShowPotAttribute15":
		sp.ShowPotAttribute15 = true
	}
	sp.ShowCount++
}

func (sp *ScoutingProfile) RemoveFromBoard() {
	sp.RemovedFromBoard = true
}

func (sp *ScoutingProfile) ReplaceOnBoard() {
	sp.RemovedFromBoard = false
}

type ScoutingProfileDTO struct {
	PlayerID uint
	TeamID   uint
}

type RevealAttributeDTO struct {
	ScoutProfileID uint
	Attribute      string
	Points         uint16
	TeamID         uint
}

type ScoutingDataResponse struct {
	DrafteeSeasonStats CollegePlayerSeasonStats
	TeamStandings      CollegeStandings
}

type ExportDraftPicksDTO struct {
	DraftPicks []DraftPick
}

type ProWarRoom struct {
	gorm.Model
	TeamID         uint
	Team           string
	ScoutingPoints uint16
	SpentPoints    uint16
}

func (w *ProWarRoom) ResetSpentPoints() {
	w.SpentPoints = 0
}

func (w *ProWarRoom) AddToSpentPoints(points uint16) {
	w.SpentPoints += points
}

type ProDraftPageResponse struct {
	WarRoomMap       map[uint]ProWarRoom
	DraftablePlayers []DraftablePlayer
	ScoutingProfiles []ScoutingProfile
	DraftPicks       map[uint][]DraftPick
}
