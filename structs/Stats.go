package structs

type BasePlayerStats struct {
	PlayerID             uint
	TeamID               uint
	SeasonID             uint
	Team                 string
	PreviousTeamID       uint
	PreviousTeam         string
	Goals                uint16
	Assists              uint16
	Points               uint16
	PlusMinus            int8
	PenaltyMinutes       uint16
	EvenStrengthGoals    uint8
	EvenStrengthPoints   uint8
	PowerPlayGoals       uint8
	PowerPlayPoints      uint8
	ShorthandedGoals     uint8
	ShorthandedPoints    uint8
	OvertimeGoals        uint8
	GameWinningGoals     uint8
	Shots                uint16
	ShootingPercentage   float32
	TimeOnIce            uint
	FaceOffWinPercentage float32
	FaceOffsWon          uint
	FaceOffs             uint
	GoalieWins           uint8
	GoalieLosses         uint8
	GoalieTies           uint8
	OvertimeLosses       uint8
	ShotsAgainst         uint16
	Saves                uint16
	GoalsAgainst         uint16
	SavePercentage       float32
	Shutouts             uint16
}

type BaseTeamStats struct {
	TeamID               uint
	Team                 string
	GoalsFor             uint16
	GoalsAgainst         uint16
	Assists              uint16
	Points               uint16
	PlusMinus            int8
	PenaltyMinutes       uint16
	EvenStrengthGoals    uint8
	EvenStrengthPoints   uint8
	PowerPlayGoals       uint8
	PowerPlayPoints      uint8
	ShorthandedGoals     uint8
	ShorthandedPoints    uint8
	OvertimeGoals        uint8
	GameWinningGoals     uint8
	Shots                uint16
	ShootingPercentage   float32
	FaceOffWinPercentage float32
	FaceOffsWon          uint
	FaceOffs             uint
	ShotsAgainst         uint16
	Saves                uint16
	SavePercentage       float32
	Shutouts             uint16
}

type TeamSeasonStats struct {
	StatType                 uint8
	GamesPlayed              uint8
	PointPercentage          float32
	RegulationWins           uint8
	RegulationOvertimeWins   uint8
	ShootoutsWon             uint8
	GFGP                     float32
	GAGP                     float32
	PowerPlayPercentage      float32
	PenaltyKillPercentage    float32
	PowerPlayNetPercentage   float32
	PenaltyKillNetPercentage float32
	ShotsPerGame             float32
	ShotsAgainstPerGame      float32
	FaceOffWinPercentage     float32
}

type CollegePlayerSeasonStats struct {
	BasePlayerStats
	StatType            uint8
	GamesPlayed         uint8
	GamesStarted        uint8
	PointsPerGamePlayed float32
}

type CollegePlayerGameStats struct {
	WeekID uint
	GameID uint
	BasePlayerStats
}

type CollegeTeamSeasonStats struct {
	BaseTeamStats
	TeamSeasonStats
}

type CollegeTeamGameStats struct {
	WeekID uint
	GameID uint
}

type ProfessionalPlayerSeasonStats struct {
	StatType uint8
	BasePlayerStats
}

type ProfessionalPlayerGameStats struct {
	WeekID uint
	GameID uint
	BasePlayerStats
}

type ProfessionalTeamSeasonStats struct {
	BaseTeamStats
	TeamSeasonStats
}

type ProfessionalTeamGameStats struct {
}
