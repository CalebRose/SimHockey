package structs

import "gorm.io/gorm"

// Weights for archetypes
var archetypeWeights = map[string]map[string]map[string]float64{
	"C": {
		"Enforcer": {
			"PuckHandling":      1.1,
			"Strength":          1.3,
			"Agility":           1.1,
			"WristShotPower":    0.75,
			"SlapshotPower":     0.75,
			"Faceoffs":          1,
			"SlapshotAccuracy":  1,
			"WristShotAccuracy": 1,
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
			"WristShotAccuracy": 0.85,
			"SlapshotAccuracy":  0.85,
			"WristShotPower":    0.85,
			"SlapshotPower":     0.85,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Power": {
			"Strength":          1.2,
			"SlapshotPower":     1.2,
			"WristShotPower":    0.8,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"Agility":           1,
			"Faceoffs":          1,
			"SlapshotAccuracy":  1,
			"WristShotAccuracy": 1,
			"Passing":           1,
			"PuckHandling":      1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Sniper": {
			"WristShotPower":    1.15,
			"WristShotAccuracy": 1.2,
			"Passing":           1.15,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"SlapshotPower":     0.9,
			"Strength":          0.8,
			"Agility":           1,
			"Faceoffs":          1,
			"SlapshotAccuracy":  1,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
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
			"WristShotPower":    0.75,
			"SlapshotPower":     0.75,
			"Faceoffs":          1,
			"SlapshotAccuracy":  1,
			"WristShotAccuracy": 1,
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
			"WristShotAccuracy": 0.85,
			"SlapshotAccuracy":  0.85,
			"WristShotPower":    0.85,
			"SlapshotPower":     0.85,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
			"BodyChecking":      1,
			"StickChecking":     1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Power": {
			"Strength":          1.2,
			"SlapshotPower":     1.2,
			"WristShotPower":    0.8,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"Agility":           1,
			"Faceoffs":          1,
			"SlapshotAccuracy":  1,
			"WristShotAccuracy": 1,
			"Passing":           1,
			"PuckHandling":      1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Sniper": {
			"WristShotPower":    1.15,
			"WristShotAccuracy": 1.2,
			"Passing":           1.15,
			"StickChecking":     0.9,
			"BodyChecking":      0.9,
			"SlapshotPower":     0.9,
			"Strength":          0.8,
			"Agility":           1,
			"Faceoffs":          1,
			"SlapshotAccuracy":  1,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
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
			"WristShotAccuracy": 0.85,
			"SlapshotAccuracy":  0.85,
			"WristShotPower":    0.85,
			"SlapshotPower":     0.85,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
			"Passing":           1,
			"ShotBlocking":      1,
			"Goalkeeping":       1,
			"GoalieVision":      1,
			"GoalieRebound":     1,
		},
		"Offensive": {
			"Passing":           1.15,
			"PuckHandling":      1.15,
			"WristShotAccuracy": 1.15,
			"WristShotPower":    1.15,
			"StickChecking":     1.2,
			"BodyChecking":      0.8,
			"Strength":          0.85,
			"ShotBlocking":      0.85,
			"SlapshotAccuracy":  0.85,
			"SlapshotPower":     0.85,
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
			"SlapshotAccuracy":  0.9,
			"SlapshotPower":     0.9,
			"WristShotAccuracy": 0.9,
			"WristShotPower":    0.9,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
			"Passing":           1,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
			"Passing":           1,
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
			"SlapshotAccuracy":  1,
			"SlapshotPower":     1,
			"WristShotAccuracy": 1,
			"WristShotPower":    1,
			"Passing":           1,
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
	WristShotAccuracy    uint8 // Accuracy on non-close shots
	WristShotPower       uint8 // Power on non-close shots. High power means less shotblocking
	SlapshotAccuracy     uint8 // Accuracy on close shots. Great on pass plays
	SlapshotPower        uint8 // Power on Close shots
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
	HomeStatePreference  uint8 // How likely the player cares about being near home (where applicable)
	ProgramPref          uint8
	ProfDevPref          uint8
	TraditionsPref       uint8
	FacilitiesPref       uint8
	AtmospherePref       uint8
	AcademicsPref        uint8
	ConferencePrestige   uint8
	PlayerMorale         uint8 // Overall Morale of the player; used for transfer intention & FA
	Personality          string
	RelativeID           uint
	RelativeType         uint
	Notes                string
	HasProgressed        bool
	// Individual weight modifiers (-10 to 10)
	Allocations
}

type BasePlayerProgressions struct {
	Agility              int // How fast a player can go in a zone without a defense check
	Faceoffs             int // Ability to win faceoffs
	WristShotAccuracy    int // Accuracy on non-close shots
	WristShotPower       int // Power on non-close shots. High power means less shotblocking
	SlapshotAccuracy     int // Accuracy on close shots. Great on pass plays
	SlapshotPower        int // Power on Close shots
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
	b.SlapshotAccuracy = ProgressAttribute(b.SlapshotAccuracy, progressions.SlapshotAccuracy)
	b.SlapshotPower = ProgressAttribute(b.SlapshotPower, progressions.SlapshotPower)
	b.WristShotAccuracy = ProgressAttribute(b.WristShotAccuracy, progressions.WristShotAccuracy)
	b.WristShotPower = ProgressAttribute(b.WristShotPower, progressions.WristShotPower)
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
	sum := int(attr) + additive
	if sum > 50 {
		sum = 50
	}
	if sum < 1 {
		sum = 1
	}
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
			value = b.Faceoffs
		case "WristShotAccuracy":
			value = b.WristShotAccuracy
		case "WristShotPower":
			value = b.WristShotPower
		case "SlapshotAccuracy":
			value = b.SlapshotAccuracy
		case "SlapshotPower":
			value = b.SlapshotPower
		case "Passing":
			value = b.Passing
		case "PuckHandling":
			value = b.PuckHandling
		case "Strength":
			value = b.Strength
		case "BodyChecking":
			value = b.BodyChecking
		case "StickChecking":
			value = b.StickChecking
		case "ShotBlocking":
			value = b.ShotBlocking
		case "GoalieVision":
			value = b.GoalieVision
		case "Goalkeeping":
			value = b.Goalkeeping
			// Add other attributes as needed
		}
		weightedSum += float64(value) * weight
		totalWeight += weight
	}

	// Normalize to 1–50 range
	if totalWeight > 0 {
		b.Overall = uint8((weightedSum / totalWeight)) // * 50.0
	} else {
		b.Overall = 0
	}
}

type BasePotentials struct {
	// Potential Attributes
	// Each attribute has a chance to grow at a different rate. These are all small modifiers
	AgilityPotential           uint8 // Ability to switch between zones
	FaceoffsPotential          uint8 // Ability to win faceoffs
	SlapshotAccuracyPotential  uint8 // Accuracy on close shots
	SlapshotPowerPotential     uint8 // Power on close shots. High power means less shotblocking
	WristShotAccuracyPotential uint8 // Accuracy on far shots. Great on pass plays
	WristShotPowerPotential    uint8 // Accuracy on far shots
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
}

func (cp *CollegePlayer) AssignTeam(teamID uint, team string) {
	cp.TeamID = uint16(teamID)
	cp.Team = team
}

func (cp *CollegePlayer) ProgressPlayer(progressions BasePlayerProgressions) {
	cp.Progress(progressions)
	cp.Year++
	cp.GetOverall()
}

type HistoricCollegePlayer struct {
	CollegePlayer
}

type ProfessionalPlayer struct {
	gorm.Model
	BasePlayer
	BasePotentials
	BaseInjuryData
	Year int
}

type RetiredPlayer struct {
	ProfessionalPlayer
}

type Recruit struct {
	gorm.Model
	BasePlayer
	BasePotentials
	BaseInjuryData
	IsSigned        bool
	College         string
	IsCustomCroot   bool
	CustomCrootFor  string
	RecruitModifier float64
	RivalsRank      float64
	ESPNRank        float64
	Rank247         float64
	TopRankModifier float64
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

/*
    // New attributes for in-depth simulation
    OneTimers     uint8 // Ability to execute one-timer shots effectively
    Deception     uint8 // Ability to fake out opponents and goalies
    Endurance      uint8 // Stamina for prolonged play without fatigue
    OffensiveAwareness uint8 // Ability to read plays and position offensively
    DefensiveAwareness uint8 // Ability to read plays and position defensively
    Leadership    uint8 // Influence on team morale and performance
    PenaltyKilling uint8 // Effectiveness in killing penalties
    PowerPlay     uint8 // Effectiveness during power plays
    ShotCreativity uint8 // Ability to create unique shot opportunities
    Breakaway     uint8 // Skill in one-on-one situations with the goalie
    AgilityWithPuck uint8 // Agility while handling the puck
    FaceoffDefense uint8 // Ability to defend against faceoff
	// Additional attributes for gameplay depth
    ShotRelease   uint8 // Speed of shot release
    DefensivePositioning uint8 // Ability to position oneself defensively
    OffensivePositioning uint8 // Ability to position oneself offensively
    CheckingPower  uint8 // Strength of body checks
    StickHandling  uint8 // Skill in maneuvering the puck with the stick
    GameSense      uint8 // Overall understanding of the game and decision-making
    ClutchFactor    uint8 // Performance under pressure in critical moments
    AgilityInTraffic uint8 // Ability to maneuver in crowded areas
    TransitionPlay uint8 // Skill in transitioning from defense to offense
    ShotSelection  uint8 // Ability to choose the right shot in various situations
    Communication  uint8 // Effectiveness in communicating with teammates
    Recovery       uint8 // Ability to recover quickly from falls or hits
    FaceoffOffense uint8 // Ability to win offensive faceoffs
    FaceoffDefense uint8 // Ability to win defensive faceoffs
    GoalieVision   uint8 // Goalkeeper's ability to track the puck
    GoalieReboundControl uint8 // Goalkeeper's ability to control rebounds
    PenaltyDrawing uint8 // Ability to draw penalties from opponents
    EnduranceRecovery uint8 // Speed of recovery after exertion
*/
