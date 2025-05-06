package structs

import (
	util "github.com/CalebRose/SimHockey/_util"
	"gorm.io/gorm"
)

// Weights for archetypes
var archetypeWeights = map[string]map[string]map[string]float64{
	"C": {
		"Enforcer": {
			"PuckHandling":      1.1,
			"Strength":          1.3,
			"Agility":           1.1,
			"LongShotPower":     0.75,
			"CloseShotPower":    0.75,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"LongShotAccuracy":  1,
			"Passing":           1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Grinder": {
			"BodyChecking":      1.15,
			"StickChecking":     1.15,
			"Passing":           1.3,
			"Strength":          1.2,
			"PuckHandling":      0.8,
			"LongShotAccuracy":  0.85,
			"CloseShotAccuracy": 0.85,
			"LongShotPower":     0.85,
			"CloseShotPower":    0.85,
			"Agility":           1,
			"Faceoffs":          1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Playmaker": {
			"Faceoffs":          1.1,
			"Passing":           1.1,
			"PuckHandling":      1.1,
			"Strength":          0.7,
			"Agility":           1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Power": {
			"Strength":          1.2,
			"CloseShotPower":    1.2,
			"LongShotPower":     0.8,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"Agility":           1,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"LongShotAccuracy":  1,
			"Passing":           1,
			"PuckHandling":      1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Sniper": {
			"LongShotPower":     1.15,
			"LongShotAccuracy":  1.2,
			"Passing":           1.15,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"CloseShotPower":    0.9,
			"Strength":          0.8,
			"Agility":           1,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"PuckHandling":      1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Two-Way": {
			"BodyChecking":      1.05,
			"StickChecking":     1.05,
			"Passing":           1.1,
			"PuckHandling":      0.9,
			"Strength":          0.9,
			"Agility":           1,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
	},
	"F": {
		"Enforcer": {
			"PuckHandling":      1.1,
			"Strength":          1.3,
			"Agility":           1.1,
			"LongShotPower":     0.75,
			"CloseShotPower":    0.75,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"LongShotAccuracy":  1,
			"Passing":           1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Grinder": {
			"BodyChecking":      1.15,
			"StickChecking":     1.15,
			"Passing":           1.3,
			"Strength":          1.2,
			"PuckHandling":      0.8,
			"LongShotAccuracy":  0.85,
			"CloseShotAccuracy": 0.85,
			"LongShotPower":     0.85,
			"CloseShotPower":    0.85,
			"Agility":           1,
			"Faceoffs":          1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Playmaker": {
			"Passing":           1.15,
			"PuckHandling":      1.15,
			"Strength":          0.7,
			"Agility":           1,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Power": {
			"Strength":          1.2,
			"CloseShotPower":    1.2,
			"LongShotPower":     0.8,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"Agility":           1,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"LongShotAccuracy":  1,
			"Passing":           1,
			"PuckHandling":      1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Sniper": {
			"LongShotPower":     1.15,
			"LongShotAccuracy":  1.2,
			"Passing":           1.15,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"CloseShotPower":    0.9,
			"Strength":          0.8,
			"Agility":           1,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"PuckHandling":      1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Two-Way": {
			"BodyChecking":      1.05,
			"StickChecking":     1.05,
			"Passing":           1.1,
			"PuckHandling":      0.9,
			"Agility":           0.9,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"Strength":          1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
	},
	"D": {
		"Defensive": {
			"BodyChecking":      1.2,
			"StickChecking":     1.2,
			"Strength":          1.15,
			"ShotBlocking":      1.15,
			"PuckHandling":      0.9,
			"LongShotAccuracy":  0.85,
			"CloseShotAccuracy": 0.85,
			"LongShotPower":     0.85,
			"CloseShotPower":    0.85,
			"Agility":           1,
			"Faceoffs":          1,
			"Passing":           1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Enforcer": {
			"Strength":          1.15,
			"BodyChecking":      1.25,
			"Agility":           1.15,
			"PuckHandling":      0.7,
			"StickChecking":     0.75,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"Passing":           1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Offensive": {
			"Passing":           1.15,
			"PuckHandling":      1.15,
			"LongShotAccuracy":  1.15,
			"LongShotPower":     1.15,
			"StickChecking":     1.2,
			"BodyChecking":      0.8,
			"Strength":          0.85,
			"ShotBlocking":      0.85,
			"CloseShotAccuracy": 0.85,
			"CloseShotPower":    0.85,
			"Agility":           1,
			"Faceoffs":          1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Two-Way": {
			"BodyChecking":      1.2,
			"StickChecking":     1.2,
			"Passing":           1.2,
			"PuckHandling":      0.9,
			"Agility":           0.9,
			"CloseShotAccuracy": 0.9,
			"CloseShotPower":    0.9,
			"LongShotAccuracy":  0.9,
			"LongShotPower":     0.9,
			"Faceoffs":          1,
			"Strength":          1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
	},
	"G": {
		"Stand-Up": {
			"GoalieVision":      1.25,
			"Strength":          1.25,
			"Goalkeeping":       0.75,
			"Agility":           0.75,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"Passing":           0.85,
			"PuckHandling":      1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"GoalieRebound":     1,
		},
		"Hybrid": {
			"GoalieVision":      1.1,
			"Strength":          1.1,
			"Goalkeeping":       1.1,
			"Agility":           1.1,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"Passing":           0.85,
			"PuckHandling":      1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"GoalieRebound":     1,
		},
		"Butterfly": {
			"GoalieVision":      1.25,
			"Strength":          1.25,
			"Goalkeeping":       0.75,
			"Agility":           0.75,
			"Faceoffs":          1,
			"CloseShotAccuracy": 1,
			"CloseShotPower":    1,
			"LongShotAccuracy":  1,
			"LongShotPower":     1,
			"Passing":           0.85,
			"PuckHandling":      1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"GoalieRebound":     1,
		},
	},

	// Add more archetypes as needed
}

type GlobalPlayer struct {
	gorm.Model
	RecruitID            uint
	CollegePlayerID      uint
	ProfessionalPlayerID uint
}

func (p *GlobalPlayer) AssignID(newID uint) {
	p.ID = newID
}

type BasePlayer struct {
	FirstName            string
	LastName             string
	Position             string
	Archetype            string
	TeamID               uint16
	Team                 string
	Height               uint8
	Weight               uint16
	Stars                uint8
	Age                  uint8
	Overall              uint8
	Agility              uint8 // How fast a player can go in a zone without a defense check
	Faceoffs             uint8 // Ability to win faceoffs
	LongShotAccuracy     uint8 // Accuracy on non-close shots
	LongShotPower        uint8 // Power on non-close shots. High power means less shotblocking
	CloseShotAccuracy    uint8 // Accuracy on close shots. Great on pass plays
	CloseShotPower       uint8 // Power on Close shots
	OneTimer             uint8 // Shots bassed on passing. Essentially a modifier that gets greater with each pass in a zone
	Passing              uint8 // Passing ability
	PuckHandling         uint8 // Ability to handle the puck when going between zones.
	Strength             uint8 // General modifier on all physical attributes. Also used in fights
	BodyChecking         uint8 // Physical defense check.
	StickChecking        uint8 // Non-phyisical defense check
	ShotBlocking         uint8 // Ability for defensemen to block a shot being made
	Goalkeeping          uint8 // Goalkeepers' ability to block a shot
	GoalieVision         uint8 // Goalkeepers' vision
	GoalieReboundControl uint8 // Ability to control a rebound
	Discipline           uint8 // Ability to not foul out
	Aggression           uint8 // Ability to avoid fights. Lower the number == more likely to cause fights
	Stamina              uint8 // Endurance in a game
	InjuryRating         uint8 // How likely the player doesn't get injured in a game
	DisciplineDeviation  uint8 // Modifier that adjusts the player's behavior in games
	InjuryDeviation      uint8 // Modifier that adjusts the player's injury rating in games
	PrimeAge             uint8 // The peak age of the player
	Clutch               int8  // Determines how well the player performs in big games. Modifier on all attributes applied prior to start of a game
	HighSchool           string
	City                 string
	State                string
	Country              string
	OriginalTeamID       uint
	OriginalTeam         string
	PreviousTeamID       uint8
	PreviousTeam         string
	Competitiveness      uint8 // How competitive the player is in games. Modifier on morale
	TeamLoyalty          uint8 // How likely a player will stay at their current program
	PlaytimePreference   uint8 // How likely the player wants to be on the rink
	PlayerMorale         uint8
	Personality          string
	RelativeID           uint
	RelativeType         uint
	Notes                string
	HasProgressed        bool
	PlayerPreferences
	// Individual weight modifiers (-10 to 10)
	Allocations
}

type PlayerPreferences struct {
	ProgramPref        uint8
	ProfDevPref        uint8
	TraditionsPref     uint8
	FacilitiesPref     uint8
	AtmospherePref     uint8
	AcademicsPref      uint8
	ConferencePref     uint8
	CoachPref          uint8
	SeasonMomentumPref uint8
}

type BasePlayerProgressions struct {
	Agility              int // How fast a player can go in a zone without a defense check
	Faceoffs             int // Ability to win faceoffs
	LongShotAccuracy     int // Accuracy on non-close shots
	LongShotPower        int // Power on non-close shots. High power means less shotblocking
	CloseShotAccuracy    int // Accuracy on close shots. Great on pass plays
	CloseShotPower       int // Power on Close shots
	OneTimer             int // Shots bassed on passing. Essentially a modifier that gets greater with each pass in a zone
	Passing              int // Passing ability
	PuckHandling         int // Ability to handle the puck when going between zones.
	Strength             int // General modifier on all physical attributes. Also used in fights
	BodyChecking         int // Physical defense check.
	StickChecking        int // Non-phyisical defense check
	ShotBlocking         int // Ability for defensemen to block a shot being made
	Goalkeeping          int // Goalkeepers' ability to block a shot
	GoalieVision         int // Goalkeepers' vision
	GoalieReboundControl int // Ability to control a rebound
}

func (b *BasePlayer) Progress(progressions BasePlayerProgressions) {
	b.Age++
	b.HasProgressed = true
	b.Agility = ProgressAttribute(b.Agility, progressions.Agility)
	b.Faceoffs = ProgressAttribute(b.Faceoffs, progressions.Faceoffs)
	b.CloseShotAccuracy = ProgressAttribute(b.CloseShotAccuracy, progressions.CloseShotAccuracy)
	b.CloseShotPower = ProgressAttribute(b.CloseShotPower, progressions.CloseShotPower)
	b.LongShotAccuracy = ProgressAttribute(b.LongShotAccuracy, progressions.LongShotAccuracy)
	b.LongShotPower = ProgressAttribute(b.LongShotPower, progressions.LongShotPower)
	b.Passing = ProgressAttribute(b.Passing, progressions.Passing)
	b.PuckHandling = ProgressAttribute(b.PuckHandling, progressions.PuckHandling)
	b.Strength = ProgressAttribute(b.Strength, progressions.Strength)
	b.BodyChecking = ProgressAttribute(b.BodyChecking, progressions.BodyChecking)
	b.StickChecking = ProgressAttribute(b.StickChecking, progressions.StickChecking)
	b.ShotBlocking = ProgressAttribute(b.ShotBlocking, progressions.ShotBlocking)
	b.Goalkeeping = ProgressAttribute(b.Goalkeeping, progressions.Goalkeeping)
	b.GoalieVision = ProgressAttribute(b.GoalieVision, progressions.GoalieVision)
	b.GoalieReboundControl = ProgressAttribute(b.GoalieReboundControl, progressions.GoalieReboundControl)
}

func ProgressAttribute(attr uint8, additive int) uint8 {
	sum := max(min(int(attr)+additive, 50), 1)
	return uint8(sum)
}

func (b *BasePlayer) GetOverall() {
	weights := archetypeWeights[b.Position][b.Archetype]
	totalWeight := 0.0
	weightedSum := 0.0

	for attr, weight := range weights {
		var value uint8
		switch attr {
		case "Agility":
			value = b.Agility

		case "Faceoffs":
			if b.Position != util.Goalie {
				value = b.Faceoffs
			}
		case "LongShotAccuracy":
			if b.Position != util.Goalie {
				value = b.LongShotAccuracy
			}
		case "LongShotPower":
			if b.Position != util.Goalie {
				value = b.LongShotPower
			}
		case "CloseShotAccuracy":
			if b.Position != util.Goalie {
				value = b.CloseShotAccuracy
			}
		case "CloseShotPower":
			if b.Position != util.Goalie {
				value = b.CloseShotPower
			}
		case "Passing":
			value = b.Passing
		case "PuckHandling":
			if b.Position != util.Goalie {
				value = b.PuckHandling
			}

		case "Strength":
			value = b.Strength
		case "BodyChecking":
			if b.Position != util.Goalie {
				value = b.BodyChecking
			}
		case "StickChecking":
			if b.Position != util.Goalie {
				value = b.StickChecking
			}

		case "ShotBlocking":
			if b.Position != util.Goalie {
				value = b.ShotBlocking
			}
		case "GoalieVision":
			if b.Position == util.Goalie {
				value = b.GoalieVision
			}

		case "Goalkeeping":
			if b.Position == util.Goalie {
				value = b.Goalkeeping
			}

			// Add other attributes as needed
		}
		weightedSum += float64(value) * weight
		if value > 0 {
			totalWeight += weight
		}
	}

	// Normalize to 1–50 range
	if totalWeight > 0 {
		b.Overall = uint8((weightedSum / totalWeight)) // * 50.0
	} else {
		b.Overall = 0
	}
}

func (cp *BasePlayer) AssignTeam(teamID uint, team string) {
	cp.TeamID = uint16(teamID)
	cp.Team = team
}

func (cp *BasePlayer) AssignAllocations(updatedAllocations Allocations) {
	cp.Allocations = updatedAllocations
}

type BasePotentials struct {
	// Potential Attributes
	// Each attribute has a chance to grow at a different rate. These are all small modifiers
	AgilityPotential           uint8 // Ability to switch between zones
	FaceoffsPotential          uint8 // Ability to win faceoffs
	CloseShotAccuracyPotential uint8 // Accuracy on close shots
	CloseShotPowerPotential    uint8 // Power on close shots. High power means less shotblocking
	LongShotAccuracyPotential  uint8 // Accuracy on far shots. Great on pass plays
	LongShotPowerPotential     uint8 // Accuracy on far shots
	PassingPotential           uint8 // Power on close shots. Great on pass plays
	PuckHandlingPotential      uint8 // Ability to handle the puck when going between zones.
	StrengthPotential          uint8 // General modifier on all physical attributes. Also used in fights
	BodyCheckingPotential      uint8 // Physical defense check.
	StickCheckingPotential     uint8 // Non-phyisical defense check
	ShotBlockingPotential      uint8 // Ability for defensemen to block a shot being made
	GoalkeepingPotential       uint8 // Goalkeepers' ability to block a shot
	GoalieVisionPotential      uint8 // Goalkeepers' ability to block a shot
	GoalieReboundPotential     uint8 // Goalkeepers' ability to block a shot
}

type BaseInjuryData struct {
	IsInjured      bool
	DaysOfRecovery int8
	InjuryName     string
	InjuryType     string
	InjuryCount    uint8
	Regression     uint8
	DecayRate      float32
}

type CollegePlayer struct {
	gorm.Model
	BasePlayer
	BasePotentials
	BaseInjuryData
	Year               int
	IsRedshirt         bool
	IsRedshirting      bool
	TransferStatus     uint8
	TransferLikeliness string
	DraftedTeamID      uint
	DraftedTeam        string
	DraftedRound       uint
	DraftPickID        uint
	DraftedPick        uint
	Stats              []CollegePlayerGameStats `gorm:"foreignKey:PlayerID;references:ID"`
	SeasonStats        CollegePlayerSeasonStats `gorm:"foreignKey:PlayerID;references:ID"`
	Profiles           []TransferPortalProfile  `gorm:"foreignKey:CollegePlayerID"`
}

func (cp *CollegePlayer) AddSeasonStats(seasonStats CollegePlayerSeasonStats) {
	cp.SeasonStats = seasonStats
}

func (cp *CollegePlayer) ProgressPlayer(progressions BasePlayerProgressions) {
	cp.Progress(progressions)
	cp.Year++
	if cp.IsRedshirting {
		cp.CompleteRedshirt()
	}
	cp.GetOverall()
}

func (cp *CollegePlayer) AssignID(id uint) {
	cp.ID = id
}

func (cp *CollegePlayer) RedshirtPlayer() {
	cp.IsRedshirting = true
}

func (cp *CollegePlayer) CompleteRedshirt() {
	cp.IsRedshirting = false
	cp.IsRedshirt = true
}

func (cp *CollegePlayer) WillTransfer() {
	cp.TransferStatus = 2
	cp.PreviousTeam = cp.Team
	cp.PreviousTeamID = uint8(cp.TeamID)
	cp.Team = ""
	cp.TeamID = 0
}

func (cp *CollegePlayer) WillReturn() {
	cp.TransferStatus = 0
	cp.Team = cp.PreviousTeam
	cp.TeamID = uint16(cp.PreviousTeamID)
	cp.PreviousTeam = ""
	cp.PreviousTeamID = 0
}

func (cp *CollegePlayer) SignWithNewTeam(teamID int, teamAbbr string) {
	cp.TransferStatus = 0
	cp.Team = teamAbbr
	cp.TeamID = uint16(teamID)
	cp.TransferLikeliness = ""
}

type HistoricCollegePlayer struct {
	CollegePlayer
}

type ProfessionalPlayer struct {
	gorm.Model
	BasePlayer
	BasePotentials
	BaseInjuryData
	CollegeID             uint
	Year                  int
	IsAffiliatePlayer     bool
	IsWaived              bool
	IsFreeAgent           bool
	IsOnTradeBlock        bool
	IsAcceptingOffers     bool
	IsNegotiating         bool
	DraftedTeamID         uint8
	DraftedTeam           string
	DraftedRound          uint8
	DraftPickID           uint
	DraftedPick           uint16
	MinimumValue          float32
	HasProgressed         bool
	Rejections            int8
	AffiliateTeamID       uint16
	Marketability         uint8   // How marketable / in demand a player's jersey will be
	JerseyPrice           float32 // Price of jersey, can be set by user
	MarketPreference      uint8
	CompetitivePreference uint8
	FinancialPreference   uint8
	IsEligibleForPlay     bool
	Stats                 []ProfessionalPlayerGameStats `gorm:"foreignKey:PlayerID;references:ID"`
	SeasonStats           ProfessionalPlayerSeasonStats `gorm:"foreignKey:PlayerID;references:ID"`
	Contract              ProContract                   `gorm:"foreignKey:PlayerID"`
	Offers                []FreeAgencyOffer             `gorm:"foreignKey:PlayerID"`
	WaiverOffer           []WaiverOffer                 `gorm:"foreignKey:PlayerID"`
	Extensions            []ExtensionOffer              `gorm:"foreignKey:PlayerID"`
}

func (cp *ProfessionalPlayer) AddSeasonStats(seasonStats ProfessionalPlayerSeasonStats) {
	cp.SeasonStats = seasonStats
}

func (cp *ProfessionalPlayer) ProgressPlayer(progressions BasePlayerProgressions) {
	cp.Progress(progressions)
	cp.Year++
	cp.GetOverall()
}

func (np *ProfessionalPlayer) ToggleIsFreeAgent() {
	np.PreviousTeamID = uint8(np.TeamID)
	np.PreviousTeam = np.Team
	np.IsFreeAgent = true
	np.TeamID = 0
	np.Team = ""
	np.IsAcceptingOffers = true
	np.IsNegotiating = false
	np.IsOnTradeBlock = false
	np.IsAffiliatePlayer = false
	np.Rejections = 0
	np.IsWaived = false
}

func (np *ProfessionalPlayer) SignPlayer(TeamID uint, Abbr string, isEligible, ToAffiliate bool) {
	np.IsFreeAgent = false
	np.IsWaived = false
	np.TeamID = uint16(TeamID)
	np.Team = Abbr
	np.IsAcceptingOffers = false
	np.IsNegotiating = false
	np.IsAffiliatePlayer = false
	np.IsEligibleForPlay = isEligible
	np.IsAffiliatePlayer = ToAffiliate
}

func (np *ProfessionalPlayer) ToggleAffiliation() {
	np.IsAffiliatePlayer = !np.IsAffiliatePlayer
	np.IsNegotiating = false
	if np.IsAffiliatePlayer {
		np.IsAcceptingOffers = true
	}
}

func (np *ProfessionalPlayer) ToggleTradeBlock() {
	np.IsOnTradeBlock = !np.IsOnTradeBlock
}

func (np *ProfessionalPlayer) RemoveFromTradeBlock() {
	np.IsOnTradeBlock = false
}

func (cp *ProfessionalPlayer) WaivePlayer() {
	cp.PreviousTeamID = uint8(cp.TeamID)
	cp.PreviousTeam = cp.Team
	cp.TeamID = 0
	cp.Team = ""
	cp.RemoveFromTradeBlock()
	cp.IsWaived = true
}

func (np *ProfessionalPlayer) ConvertWaivedPlayerToFA() {
	np.IsWaived = false
	np.IsFreeAgent = true
	np.IsAcceptingOffers = true
}

func (np *ProfessionalPlayer) ToggleIsNegotiating() {
	np.IsNegotiating = true
	np.IsAcceptingOffers = false
}

func (np *ProfessionalPlayer) WaitUntilAfterDraft() {
	np.IsNegotiating = false
	np.IsAcceptingOffers = false
}

func (np *ProfessionalPlayer) AssignPreferences(m, c, f uint8) {
	np.MarketPreference = m
	np.CompetitivePreference = c
	np.FinancialPreference = f
}

func (cp *ProfessionalPlayer) AssignID(id uint) {
	cp.ID = id
}

type RetiredPlayer struct {
	ProfessionalPlayer
}

type Recruit struct {
	gorm.Model
	BasePlayer
	BasePotentials
	BaseInjuryData
	IsSigned              bool
	College               string
	IsCustomCroot         bool
	CustomCrootFor        string
	RecruitModifier       float32
	CompositeRank         float32
	RivalsRank            float32
	ESPNRank              float32
	Rank247               float32
	TopRankModifier       float32
	RecruitingModifier    float32                // For signing threshold
	RecruitingStatus      string                 // For signing progress
	RecruitPlayerProfiles []RecruitPlayerProfile `gorm:"foreignKey:RecruitID"`
}

func (r *Recruit) AssignID(id uint) {
	r.ID = id
}

func (r *Recruit) AssignCollege(abbr string) {
	r.College = abbr
}

func (r *Recruit) UpdateTeamID(id uint) {
	r.TeamID = uint16(id)
	if id > 0 {
		r.IsSigned = true
	}
}

func (r *Recruit) AssignRelativeData(rID, rType, teamID uint, team, notes string) {
	r.RelativeID = rID
	r.RelativeType = rType
	r.Notes = notes
	if teamID > 0 {
		r.UpdateTeamID(teamID)
		r.AssignCollege(team)
	}
}

func (r *Recruit) AssignTwinData(lastName, city, state, highschool string) {
	r.LastName = lastName
	r.City = city
	r.State = state
	r.HighSchool = highschool
}

func (r *Recruit) ApplySigningStatus(num, threshold float32, signing bool) {
	percentage := num / threshold

	if threshold == 0 || num == 0 || percentage < 0.26 {
		r.RecruitingStatus = "Not Ready"
	} else if percentage < 0.51 {
		r.RecruitingStatus = "Hearing Offers"
	} else if percentage < 0.76 {
		r.RecruitingStatus = "Narrowing Down Offers"
	} else if percentage < 0.96 {
		r.RecruitingStatus = "Finalizing Decisions"
	} else if percentage < 1 {
		r.RecruitingStatus = "Ready to Sign"
	} else {
		r.RecruitingStatus = "Signed"
	}

	if signing {
		r.RecruitingStatus = "Signed"
	}
}

func (r *Recruit) AssignRankValues(rank247 float32, espnRank float32, rivalsRank float32, modifier float32) {
	r.Rank247 = rank247
	r.ESPNRank = espnRank
	r.RivalsRank = rivalsRank
	r.CompositeRank = (rank247 + espnRank + rivalsRank) / 3
	r.TopRankModifier = modifier
}

func (r *Recruit) AssignRecruitingModifier(recruitingMod float32) {
	r.RecruitingModifier = recruitingMod
}
