package structs

import "gorm.io/gorm"

type TradePreferencesDTO struct {
	TeamID        uint
	Centers       bool
	CentersType   uint8
	Forwards      bool
	ForwardsType  uint8
	Defenders     bool
	DefendersType uint8
	Goalies       bool
	GoaliesType   uint8
	DraftPicks    bool
	DraftPickType uint8
}

type TradePreferences struct {
	gorm.Model
	TeamID        uint
	Centers       bool
	CentersType   uint8
	Forwards      bool
	ForwardsType  uint8
	Defenders     bool
	DefendersType uint8
	Goalies       bool
	GoaliesType   uint8
	DraftPicks    bool
	DraftPickType uint8
}

func (tp *TradePreferences) UpdatePreferences(pref TradePreferencesDTO) {
	tp.Centers = pref.Centers
	if tp.Centers {
		tp.CentersType = pref.CentersType
	}
	tp.Forwards = pref.Forwards
	if tp.Forwards {
		tp.ForwardsType = pref.ForwardsType
	}
	tp.Defenders = pref.Defenders
	if tp.Defenders {
		tp.DefendersType = pref.DefendersType
	}
	tp.Goalies = pref.Goalies
	if tp.Goalies {
		tp.GoaliesType = pref.GoaliesType
	}
	tp.DraftPicks = pref.DraftPicks
	if tp.DraftPicks {
		tp.DraftPickType = pref.DraftPickType
	}
}

type TradeProposal struct {
	gorm.Model
	TeamID                    uint
	RecepientTeamID           uint
	IsTradeAccepted           bool
	IsTradeRejected           bool
	IsSynced                  bool
	TeamTradeOptions          []TradeOption `gorm:"foreignKey:TradeProposalID"`
	RecepientTeamTradeOptions []TradeOption `gorm:"foreignKey:TradeProposalID"`
}

func (p *TradeProposal) ToggleSyncStatus() {
	p.IsSynced = true
}

func (p *TradeProposal) AssignID(id uint) {
	p.ID = id
}

func (p *TradeProposal) AcceptTrade() {
	p.IsTradeAccepted = true
}

func (p *TradeProposal) RejectTrade() {
	p.IsTradeRejected = true
}

type TradeOption struct {
	gorm.Model
	TradeProposalID  uint
	TeamID           uint
	PlayerID         uint
	DraftPickID      uint
	OptionType       string
	SalaryPercentage float64 // Will be a percentage that the recepient team (TEAM B) will pay for Y1. Will be between 0 and 100.
	// Player           Player    // `gorm:"foreignKey:PlayerID"`       // If the PlayerID is greater than 0, it will return a player.
	// Draftpick        DraftPick // `gorm:"foreignKey:DraftPickID"` // If the DraftPickID is greater than 0, it will return a draft pick.
}

type TradeOptionObj struct {
	ID               uint
	TradeProposalID  uint
	TeamID           uint
	PlayerID         uint
	DraftPickID      uint
	OptionType       string
	SalaryPercentage float64            // Will be a percentage that the recepient team (TEAM B) will pay. Will be between 0 and 100.
	Player           ProfessionalPlayer // If the PlayerID is greater than 0, it will return a player.
	Draftpick        DraftPick          // If the DraftPickID is greater than 0, it will return a draft pick.
}

func (to *TradeOptionObj) AssignPlayer(player ProfessionalPlayer) {
	to.Player = player
	to.PlayerID = player.ID
}

func (to *TradeOptionObj) AssignPick(pick DraftPick) {
	to.Draftpick = pick
	to.DraftPickID = pick.ID
}

type TradeProposalDTO struct {
	ID                        uint
	TeamID                    uint
	Team                      string
	RecepientTeamID           uint
	RecepientTeam             string
	IsTradeAccepted           bool
	IsTradeRejected           bool
	TeamTradeOptions          []TradeOptionObj
	RecepientTeamTradeOptions []TradeOptionObj
}

type TeamProposals struct {
	SentTradeProposals     []TradeProposalDTO
	ReceivedTradeProposals []TradeProposalDTO
}
