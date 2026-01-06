package engine

import (
	"math"
	"math/rand"

	util "github.com/CalebRose/SimHockey/_util"
)

// InjuryEventType represents the type of event that caused an injury
type InjuryEventType uint8

const (
	BodyCheckEvent InjuryEventType = iota
	StickCheckEvent
	MissedShotBlocked
	MissedPassInterception
	WildPuck
	Fighting
	PuckContact
	Fall
	PenaltyEvent
)

// InjurySeverity represents how severe an injury is
type InjurySeverity uint8

const (
	Minor InjurySeverity = iota
	Moderate
	Severe
	Critical
)

// InjuryType represents the specific type of injury
type InjuryType uint8

const (
	// Upper Body
	Concussion InjuryType = iota
	ShoulderSeparation
	BrokenWrist
	BrokenHand
	ElbowInjury
	RibInjury
	BackStrain
	Cut
	Bruise

	// Lower Body
	GroinStrain
	KneeSprain
	AnkleSprain
	HipPointer
	HamstringStrain

	// Other
	GeneralSoreness
)

// Injury represents a player injury
type Injury struct {
	PlayerID      uint
	InjuryType    InjuryType
	InjuryName    string
	Severity      InjurySeverity
	RecoveryDays  int
	CausedByEvent InjuryEventType
	WeekID        uint
	GameID        uint
}

// InjuryData holds the injury information for each injury type
type InjuryData struct {
	Name       string
	MinDays    int
	MaxDays    int
	Severity   InjurySeverity
	BodyRegion string
}

// injuryDatabase contains all possible injuries and their data
var injuryDatabase = map[InjuryType]InjuryData{
	Concussion:         {"Concussion", 7, 28, Severe, "Head"},
	ShoulderSeparation: {"Shoulder Separation", 14, 42, Moderate, "Upper Body"},
	BrokenWrist:        {"Broken Wrist", 28, 56, Severe, "Upper Body"},
	BrokenHand:         {"Broken Hand", 21, 42, Severe, "Upper Body"},
	ElbowInjury:        {"Elbow Injury", 7, 21, Minor, "Upper Body"},
	RibInjury:          {"Rib Injury", 14, 35, Moderate, "Upper Body"},
	BackStrain:         {"Back Strain", 7, 21, Minor, "Upper Body"},
	Cut:                {"Cut/Laceration", 1, 7, Minor, "Various"},
	Bruise:             {"Bruising", 1, 5, Minor, "Various"},
	GroinStrain:        {"Groin Strain", 7, 21, Minor, "Lower Body"},
	KneeSprain:         {"Knee Sprain", 14, 42, Moderate, "Lower Body"},
	AnkleSprain:        {"Ankle Sprain", 7, 28, Minor, "Lower Body"},
	HipPointer:         {"Hip Pointer", 7, 14, Minor, "Lower Body"},
	HamstringStrain:    {"Hamstring Strain", 7, 21, Minor, "Lower Body"},
	GeneralSoreness:    {"General Soreness", 1, 3, Minor, "Various"},
}

// Event-based injury probabilities (injuries that can occur from each event type)
var eventInjuryMap = map[InjuryEventType][]InjuryType{
	BodyCheckEvent: {
		Concussion, ShoulderSeparation, KneeSprain, RibInjury,
		BackStrain, Bruise, AnkleSprain, HipPointer, GeneralSoreness,
	},
	StickCheckEvent: {
		BrokenWrist, BrokenHand, Cut, Bruise, ElbowInjury, GeneralSoreness,
	},
	MissedShotBlocked: {
		BrokenWrist, BrokenHand, Cut, Bruise, RibInjury, GeneralSoreness,
	},
	MissedPassInterception: {
		Cut, Bruise, ElbowInjury, GeneralSoreness,
	},
	WildPuck: {
		Cut, Bruise, Concussion, BrokenWrist, BrokenHand, GeneralSoreness,
	},
	Fighting: {
		Cut, Bruise, BrokenHand, Concussion, ElbowInjury, GeneralSoreness,
	},
	PuckContact: {
		Cut, Bruise, BrokenWrist, BrokenHand, RibInjury, GeneralSoreness,
	},
	Fall: {
		Concussion, ShoulderSeparation, KneeSprain, AnkleSprain,
		BackStrain, Bruise, HipPointer, HamstringStrain, GeneralSoreness,
	},
	PenaltyEvent: {
		Concussion, ShoulderSeparation, KneeSprain, AnkleSprain,
		BackStrain, Bruise, HipPointer, HamstringStrain, GeneralSoreness,
	},
}

// GetPossibleInjuries returns a slice of possible injuries for a given event type
func GetPossibleInjuries(eventType InjuryEventType) []Injury {
	possibleInjuryTypes, exists := eventInjuryMap[eventType]
	if !exists {
		return []Injury{}
	}

	var injuries []Injury
	for _, injuryType := range possibleInjuryTypes {
		injuryData := injuryDatabase[injuryType]

		// Calculate recovery days within the range
		recoveryDays := injuryData.MinDays + rand.Intn(injuryData.MaxDays-injuryData.MinDays+1)

		// Reduce the number of recovery days by the number of games simulated per week (4 per week, so right now do 0.5)
		estimatedGamesMissed := EstimatedGamesMissed(uint8(recoveryDays))

		injury := Injury{
			InjuryType:    injuryType,
			InjuryName:    injuryData.Name,
			Severity:      injuryData.Severity,
			RecoveryDays:  estimatedGamesMissed,
			CausedByEvent: eventType,
		}

		injuries = append(injuries, injury)
	}

	return injuries
}

// GetRandomInjury returns a random injury for a given event type
func GetRandomInjury(eventType InjuryEventType) *Injury {
	possibleInjuries := GetPossibleInjuries(eventType)
	if len(possibleInjuries) == 0 {
		return nil
	}

	randomIndex := rand.Intn(len(possibleInjuries))
	injury := possibleInjuries[randomIndex]

	return &injury
}

// GetInjuryProbabilityByEvent returns injury probability weights for each event type
func GetInjuryProbabilityByEvent() map[InjuryEventType]float64 {
	return map[InjuryEventType]float64{
		BodyCheckEvent:         0.0006,  // 0.06% chance per body check (most dangerous)
		StickCheckEvent:        0.0002,  // 0.02% chance per stick check
		PuckContact:            0.0001,  // 0.01% chance per puck contact
		Fall:                   0.0001,  // 0.01% chance per fall
		Fighting:               0.0010,  // 0.10% chance per fight (higher risk)
		PenaltyEvent:           0.0004,  // 0.04% chance per penalty (physical play)
		MissedShotBlocked:      0.00005, // 0.005% chance per blocked shot
		MissedPassInterception: 0.00005, // 0.005% chance per interception
		WildPuck:               0.0001,  // 0.01% chance per wild puck
	}
}

// CalculateInjuryRisk calculates the risk of injury based on various factors
func CalculateInjuryRisk(eventType InjuryEventType, playerInjuryRating int, gameIntensity float64) float64 {
	baseProbability := GetInjuryProbabilityByEvent()[eventType]

	// Adjust based on player toughness (0-100 scale)
	toughnessModifier := 1.0 - (float64(playerInjuryRating) / 100.0 * 0.3) // Up to 30% reduction

	// Adjust based on game intensity (0.5-2.0 scale)
	intensityModifier := gameIntensity

	finalRisk := baseProbability * toughnessModifier * intensityModifier

	return finalRisk
}

// IsPlayerInjured determines if a player gets injured based on calculated risk
func IsPlayerInjured(injuryRisk float64) bool {
	chance := rand.Float64()
	if chance < injuryRisk {
		return true
	}
	return false
}

// GetInjuryMapByEvent returns a map of injuries organized by event type
func GetInjuryMapByEvent(eventType InjuryEventType) map[string]interface{} {
	possibleInjuries := GetPossibleInjuries(eventType)

	result := map[string]interface{}{
		"event_type":        eventType,
		"event_name":        getEventTypeName(eventType),
		"possible_injuries": possibleInjuries,
		"injury_count":      len(possibleInjuries),
	}

	return result
}

// Helper function to get event type name
func getEventTypeName(eventType InjuryEventType) string {
	switch eventType {
	case BodyCheckEvent:
		return "Body Check"
	case StickCheckEvent:
		return "Stick Check"
	case MissedShotBlocked:
		return "Missed Shot/Blocked Shot"
	case MissedPassInterception:
		return "Missed Pass/Puck Interception"
	case WildPuck:
		return "Wild Puck"
	case Fighting:
		return "Fighting"
	case PuckContact:
		return "Puck Contact"
	case Fall:
		return "Fall/Collision"
	default:
		return "Unknown Event"
	}
}

// GetSeverityString returns a string representation of injury severity
func (s InjurySeverity) String() string {
	switch s {
	case Minor:
		return "Minor"
	case Moderate:
		return "Moderate"
	case Severe:
		return "Severe"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// GetInjuryTypeString returns a string representation of injury type
func (i InjuryType) String() string {
	if data, exists := injuryDatabase[i]; exists {
		return data.Name
	}
	return "Unknown Injury"
}

func GetInjurySeverity(severity InjurySeverity) string {
	switch severity {
	case 0:
		return "Minor"
	case 1:
		return "Moderate"
	case 2:
		return "Severe"
	default:
		return "Critical"
	}
}

func HandleInjuryEvent(gs *GameState, eventType InjuryEventType, player *GamePlayer) {
	injury := GetRandomInjury(eventType)
	if injury != nil {
		severity := GetInjurySeverity(injury.Severity)
		injuryNameID := util.GetInjuryIDByName(injury.InjuryName)
		player.RecordInjury(injuryNameID, uint8(injury.Severity), uint8(injury.RecoveryDays))
		player.ApplyInjury(injury.InjuryName, severity, int8(injury.RecoveryDays))

		// Remove player from lineup due to injury
		RemovePlayerFromGame(gs, player)

		// Set the PlayerID before logging the injury
		injury.PlayerID = player.ID

		// Log Injury to Game State
		gs.LogInjury(*injury)
		injurySeverityID := injury.Severity + 45 // Offset to avoid conflict with other event IDs
		RecordPlay(gs, InjuryCheckID, uint8(injurySeverityID), 0, uint8(injury.InjuryType), uint8(injury.Severity), uint8(injury.RecoveryDays), 0, 0, false, player.ID, 0, 0, 0, 0, false)
	}
}

// Convert injury days to estimated games missed based on typical NHL schedule
func EstimatedGamesMissed(recoveryDays uint8) int {
	// NHL averages ~3.5 games per week during season
	gamesPerDay := 0.5
	return int(math.Ceil(float64(recoveryDays) * gamesPerDay))
}
