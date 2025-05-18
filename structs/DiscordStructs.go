package structs

type CollegeTeamResponseData struct {
	TeamData        CollegeTeam
	TeamStandings   CollegeStandings
	UpcomingMatches []CollegeGame
}

type ProTeamResponseData struct {
	TeamData        ProfessionalTeam
	TeamStandings   ProfessionalStandings
	UpcomingMatches []ProfessionalGame
}

type DiscordPlayer struct {
	PlayerID           uint
	FirstName          string
	LastName           string
	Position           string
	Archetype          string
	TeamID             uint16
	Team               string
	Height             uint8
	Weight             uint16
	Stars              uint8
	Age                uint8
	Year               uint8
	IsRedshirt         bool
	IsRedshirting      bool
	Overall            string
	Agility            string // How fast a player can go in a zone without a defense check
	Faceoffs           string // Ability to win faceoffs
	LongShotAccuracy   string // Accuracy on non-close shots
	LongShotPower      string // Power on non-close shots. High power means less shotblocking
	CloseShotAccuracy  string // Accuracy on close shots. Great on pass plays
	CloseShotPower     string // Power on Close shots
	Passing            string // Passing ability
	PuckHandling       string // Ability to handle the puck when going between zones.
	Strength           string // General modifier on all physical attributes. Also used in fights
	BodyChecking       string // Physical defense check.
	StickChecking      string // Non-phyisical defense check
	ShotBlocking       string // Ability for defensemen to block a shot being made
	Goalkeeping        string // Goalkeepers' ability to block a shot
	GoalieVision       string // Goalkeepers' vision
	Stamina            string
	InjuryRating       string
	HighSchool         string
	City               string
	State              string
	Country            string
	OriginalTeamID     uint
	OriginalTeam       string
	PreviousTeamID     uint8
	PreviousTeam       string
	Competitiveness    string // How competitive the player is in games. Modifier on morale
	TeamLoyalty        string // How likely a player will stay at their current program
	PlaytimePreference string // How likely the player wants to be on the rink
	PlayerMorale       uint8
	Personality        string
	RelativeID         uint
	RelativeType       uint
	Notes              string
	HasProgressed      bool
	PlayerPreferences
	CollegeStats CollegePlayerSeasonStats
	ProStats     ProfessionalPlayerSeasonStats
}

type ProDiscordPlayer struct {
	PlayerID              uint
	FirstName             string
	LastName              string
	Position              string
	Archetype             string
	TeamID                uint16
	Team                  string
	Height                uint8
	Weight                uint16
	Stars                 uint8
	Age                   uint8
	Overall               uint8
	Year                  uint8
	Agility               uint8 // How fast a player can go in a zone without a defense check
	Faceoffs              uint8 // Ability to win faceoffs
	LongShotAccuracy      uint8 // Accuracy on non-close shots
	LongShotPower         uint8 // Power on non-close shots. High power means less shotblocking
	CloseShotAccuracy     uint8 // Accuracy on close shots. Great on pass plays
	CloseShotPower        uint8 // Power on Close shots
	Passing               uint8 // Passing ability
	PuckHandling          uint8 // Ability to handle the puck when going between zones.
	Strength              uint8 // General modifier on all physical attributes. Also used in fights
	BodyChecking          uint8 // Physical defense check.
	StickChecking         uint8 // Non-phyisical defense check
	ShotBlocking          uint8 // Ability for defensemen to block a shot being made
	Goalkeeping           uint8 // Goalkeepers' ability to block a shot
	GoalieVision          uint8 // Goalkeepers' vision
	Stamina               string
	InjuryRating          string
	HighSchool            string
	City                  string
	State                 string
	Country               string
	OriginalTeamID        uint
	OriginalTeam          string
	PreviousTeamID        uint8
	PreviousTeam          string
	IsFreeAgent           bool
	MarketPreference      string // How competitive the player is in games. Modifier on morale
	CompetitivePreference string // How likely a player will stay at their current program
	FinancialPreference   string // How likely the player wants to be on the rink
	PlayerMorale          uint8
	Personality           string
	RelativeID            uint
	RelativeType          uint
	Notes                 string
	HasProgressed         bool
	PlayerPreferences
	CollegeStats CollegePlayerSeasonStats
	ProStats     ProfessionalPlayerSeasonStats
}

type TeamComparisonModel struct {
	TeamOneID       uint
	TeamOne         string
	TeamOneWins     uint
	TeamOneLosses   uint
	TeamOneOTWins   uint
	TeamOneOTLosses uint
	TeamOneSOWins   uint
	TeamOneSOLosses uint
	TeamOneStreak   uint
	TeamOneMSeason  int
	TeamOneMScore   string
	TeamTwoID       uint
	TeamTwo         string
	TeamTwoWins     uint
	TeamTwoLosses   uint
	TeamTwoOTWins   uint
	TeamTwoOTLosses uint
	TeamTwoSOWins   uint
	TeamTwoSOLosses uint
	TeamTwoStreak   uint
	TeamTwoMSeason  int
	TeamTwoMScore   string
	CurrentStreak   uint
	LatestWin       string
}
