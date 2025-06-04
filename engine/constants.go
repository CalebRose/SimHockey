package engine

const (
	HomeGoal                string  = "Home Goal"
	HomeZone                string  = "Home Zone"
	NeutralZone             string  = "Neutral Zone"
	AwayZone                string  = "Away Zone"
	AwayGoal                string  = "Away Goal"
	Defender                string  = "D"
	Forward                 string  = "F"
	Center                  string  = "C"
	Goalie                  string  = "G"
	Rebound                 string  = "Rebound"
	Defense                 string  = "Defense"
	ShotBlock               string  = "ShotBlock"
	Faceoff                 string  = "Faceoff"
	Pass                    string  = "Pass"
	VeryEasyReq             float64 = 6
	EasyReq                 float64 = 8
	BaseReq                 float64 = 10
	SlightlyDiffReq         float64 = 12
	DiffReq                 float64 = 14
	ToughReq                float64 = 17
	CritSuccess             int     = 20
	CritFail                int     = 1
	Heads                   int     = 1
	Tails                   int     = 1
	ModifierFactor          float64 = 1.3 // Adjust as needed for your testing
	ScaleFactor             float64 = 1.7 // Adjust as needed for your testing
	MinorPenalty            string  = "Minor Penalty"
	MajorPenalty            string  = "Major Penalty"
	Misconduct              string  = "Misconduct"
	GameMisconduct          string  = "Game Misconduct"
	MatchPenalty            string  = "Match Penalty"
	BodyCheck               string  = "BodyCheck"
	StickCheck              string  = "StickCheck"
	General                 string  = "General"
	Fight                   string  = "Fight"
	ShootoutMomenumModifier float64 = 0.375
	RegularPeriodTime       uint16  = 1200
	OvertimePeriodTime      uint16  = 300
	MaxTimeOnClock          uint16  = 65000
	GoalieStaminaThreshold  uint8   = 30
)

// Event Constants
const (
	// Event IDs
	FaceoffID          uint8 = 1
	PhysDefenseCheckID uint8 = 2
	DexDefenseCheckID  uint8 = 3
	PassCheckID        uint8 = 4
	AgilityCheckID     uint8 = 5
	WristshotCheckID   uint8 = 6
	SlapshotCheckID    uint8 = 7
	PenaltyCheckID     uint8 = 8
	EnteringShootout   uint8 = 9
	WSShootoutID       uint8 = 10
	CSShootoutID       uint8 = 11

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
)

// PenaltyIDs
