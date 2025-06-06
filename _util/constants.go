package util

const (
	Center            string = "C"
	Forward           string = "F"
	Defender          string = "D"
	Goalie            string = "G"
	Enforcer          string = "Enforcer"
	Grinder           string = "Grinder"
	Playmaker         string = "Playmaker"
	Power             string = "Power"
	Sniper            string = "Sniper"
	TwoWay            string = "Two-Way"
	Defensive         string = "Defensive"
	Offensive         string = "Offensive"
	StandUp           string = "Stand-Up"
	Hybrid            string = "Hybrid"
	Butterfly         string = "Butterfly"
	Agility           string = "Agility"
	Faceoffs          string = "Faceoffs"
	CloseShotAccuracy string = "CloseShotAccuracy"
	CloseShotPower    string = "CloseShotPower"
	LongShotAccuracy  string = "LongShotAccuracy"
	LongShotPower     string = "LongShotPower"
	Passing           string = "Passing"
	PuckHandling      string = "PuckHandling"
	Strength          string = "Strength"
	BodyChecking      string = "BodyChecking"
	StickChecking     string = "StickChecking"
	ShotBlocking      string = "ShotBlocking"
	Goalkeeping       string = "Goalkeeping"
	GoalieVision      string = "GoalieVision"
	GoalieRebound     string = "GoalieRebound"
	GoalieHold        string = "Goalie Hold"

	// Event IDs
	FaceoffID          uint8 = 1
	PhysDefenseCheckID uint8 = 2
	DexDefenseCheckID  uint8 = 3
	PassCheckID        uint8 = 4
	AgilityCheckID     uint8 = 5
	WristshotCheckID   uint8 = 6
	SlapshotCheckID    uint8 = 7
	PenaltyCheckID     uint8 = 8

	// Zone IDs
	HomeGoalZoneID uint8 = 9
	HomeZoneID     uint8 = 10
	NeutralZoneID  uint8 = 11
	AwayZoneID     uint8 = 12
	AwayGoalZoneID uint8 = 13

	// Outcome IDs
	DefenseTakesPuckID   uint8 = 14
	CarrierKeepsPuckID   uint8 = 15
	DefenseStopAgilityID uint8 = 16
	OffenseMovesUpID     uint8 = 17
	GeneralPenaltyID     uint8 = 18
	FightPenaltyID       uint8 = 20
	InterceptedPassID    uint8 = 21
	ReceivedPassID       uint8 = 22
	HomeFaceoffWinID     uint8 = 23
	AwayFaceoffWinID     uint8 = 24
	InAccurateShotID     uint8 = 25
	ShotBlockedID        uint8 = 26
	GoalieSaveID         uint8 = 27
	GoalieReboundID      uint8 = 28
	ShotOnGoalID         uint8 = 29
	GoalieHoldID         uint8 = 30
	NoOneOpenID          uint8 = 31
	ReceivedLongPassID   uint8 = 32
	ReceivedBackPassID   uint8 = 33
	EnteringShootout     uint8 = 34
	WSShootoutID         uint8 = 35
	CSShootoutID         uint8 = 36

	// Free Agency Preferences
	Average                  uint8 = 1
	MarketCTH                uint8 = 2
	MarketCountrymen         uint8 = 3
	MarketLarge              uint8 = 4
	MarketNotLarge           uint8 = 5
	MarketSmall              uint8 = 6
	MarketNotSmall           uint8 = 7
	MarketLoyal              uint8 = 8
	MarketAvoidPrevTeam      uint8 = 9
	CompetitiveMentorship    uint8 = 2
	CompetitiveVeteranMentor uint8 = 3
	CompetitiveFirstLine     uint8 = 4
	CompetitiveSecondLine    uint8 = 5
	CompetitiveCompete       uint8 = 6
	FinancialShort           uint8 = 2
	FinancialLong            uint8 = 3
	FinancialLargeAAV        uint8 = 4

	// Recruiting
	USA    string = "USA"
	Canada string = "Canada"
	Sweden string = "Sweden"
	Russia string = "Russia"

	MaxGoalieStamina   uint8 = 100
	MinGoalieStamina   uint8 = 1
	ForwardShiftLimit  int   = 40
	DefenderShiftLimit int   = 25
)
