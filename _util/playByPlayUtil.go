package util

// EventID Map
var eventIDMap = map[uint8]string{
	FaceoffID:          "Faceoff",
	PhysDefenseCheckID: "Physical Defense Check",
	DexDefenseCheckID:  "Dexterity Defense Check",
	PassCheckID:        "Pass Check",
	AgilityCheckID:     "Agility Check",
	WristshotCheckID:   "Wrist Shot Check",
	SlapshotCheckID:    "Slap Shot Check",
	PenaltyCheckID:     "Penalty Check",
	PuckBattleID:       "Puck Battle",
	PuckBattleWinID:    "Puck Battle Win",
	PuckBattleLoseID:   "Puck Battle Lose",
	PuckScrambleID:     "Puck Scramble",
	PuckCoveredID:      "Puck Covered",
	LongPassCheckID:    "Long Pass Check",
	PassBackCheckID:    "Pass Back Check",
	InjuryCheckID:      "Injury Check",
	MinorInjuryID:      "Minor Injury",
	ModerateInjuryID:   "Moderate Injury",
	SevereInjuryID:     "Severe Injury",
	CriticalInjuryID:   "Critical Injury",

	// Zones
	HomeGoalZoneID: "Home Goal Zone",
	HomeZoneID:     "Home Zone",
	NeutralZoneID:  "Neutral Zone",
	AwayZoneID:     "Away Zone",
	AwayGoalZoneID: "Away Goal Zone",

	// Outcomes
	DefenseTakesPuckID:   "Defense Takes Puck",
	CarrierKeepsPuckID:   "Carrier Keeps Puck",
	DefenseStopAgilityID: "Defense Stops Agility",
	OffenseMovesUpID:     "Offense Moves Up",
	GeneralPenaltyID:     "General Penalty",
	FightPenaltyID:       "Fight Penalty",
	InterceptedPassID:    "Intercepted Pass",
	ReceivedPassID:       "Received Pass",
	ReceivedLongPassID:   "Received Long Pass",
	ReceivedBackPassID:   "Received Back Pass",
	NoOneOpenID:          "No One Open",
	HomeFaceoffWinID:     "Home Faceoff Win",
	AwayFaceoffWinID:     "Away Faceoff Win",
	InAccurateShotID:     "Inaccurate Shot",
	ShotBlockedID:        "Shot Blocked",
	GoalieSaveID:         "Goalie Save",
	GoalieReboundID:      "Goalie Rebound",
	ShotOnGoalID:         "Shot on Goal",
	GoalieHoldID:         "Goalie Hold",
	EnteringShootout:     "Entering Shootout",
	WSShootoutID:         "Shootout",
	CSShootoutID:         "Shootout",
}

func ReturnStringFromPBPID(id uint8) string {
	if val, exists := eventIDMap[id]; exists {
		return val
	}
	return "Unknown Event" // Return default for invalid IDs
}
