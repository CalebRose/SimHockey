package structs

import (
	"gorm.io/gorm"
)

type BasePlayerStats struct {
	gorm.Model
	StartedGame          bool
	GameDay              string
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
	ShotsBlocked         uint16
	BodyChecks           uint16
	StickChecks          uint16
	IsInjured            bool
	DaysOfRecovery       int8
	InjuryName           string
	InjuryType           string
	GameType             uint8
}

func (s *BasePlayerStats) AddStatsToSeasonRecord(stat BasePlayerStats) {
	// on first call, capture identity
	if s.PlayerID == 0 {
		s.PlayerID = stat.PlayerID
		s.TeamID = stat.TeamID
		s.Team = stat.Team
	}
	// always keep seasonID in sync
	s.SeasonID = stat.SeasonID
	s.GameType = stat.GameType

	// raw count accumulation
	s.Goals += stat.Goals
	s.Assists += stat.Assists
	s.Points += stat.Points
	s.PlusMinus += stat.PlusMinus
	s.PenaltyMinutes += stat.PenaltyMinutes

	s.EvenStrengthGoals += stat.EvenStrengthGoals
	s.EvenStrengthPoints += stat.EvenStrengthPoints
	s.PowerPlayGoals += stat.PowerPlayGoals
	s.PowerPlayPoints += stat.PowerPlayPoints
	s.ShorthandedGoals += stat.ShorthandedGoals
	s.ShorthandedPoints += stat.ShorthandedPoints

	s.OvertimeGoals += stat.OvertimeGoals
	s.GameWinningGoals += stat.GameWinningGoals

	s.Shots += stat.Shots
	s.TimeOnIce += stat.TimeOnIce

	if s.Shots > 0 {
		s.ShootingPercentage = float32(s.Goals) / float32(s.Shots)
	}

	s.FaceOffsWon += stat.FaceOffsWon
	s.FaceOffs += stat.FaceOffs

	if s.FaceOffs > 0 {
		s.FaceOffWinPercentage = float32(s.FaceOffsWon) / float32(s.FaceOffs)
	}

	s.ShotsAgainst += stat.ShotsAgainst
	s.Saves += stat.Saves
	s.GoalsAgainst += stat.GoalsAgainst
	if s.Saves > 0 && s.ShotsAgainst > 0 {
		s.SavePercentage = float32(s.Saves) / float32(s.ShotsAgainst)
	}

	s.Shutouts += stat.Shutouts
	s.ShotsBlocked += stat.ShotsBlocked
	s.BodyChecks += stat.BodyChecks
	s.StickChecks += stat.StickChecks
}

type BaseTeamStats struct {
	gorm.Model
	SeasonID             uint
	TeamID               uint
	Team                 string
	GameDay              string
	GoalsFor             uint16
	GoalsAgainst         uint16
	Assists              uint16
	Points               uint16
	Period1Score         uint8
	Period2Score         uint8
	Period3Score         uint8
	OTScore              uint8
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
	GameType             uint8
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
}

func (s *BaseTeamStats) AddStatsToSeasonRecord(stat BaseTeamStats) {
	if s.TeamID == 0 {
		s.TeamID = stat.TeamID
		s.Team = stat.Team
	}
	s.GameType = stat.GameType
	s.SeasonID = stat.SeasonID

	s.GoalsFor += stat.GoalsFor
	s.GoalsAgainst += stat.GoalsAgainst
	s.Assists += stat.Assists
	s.Points += stat.Points

	s.Period1Score += stat.Period1Score
	s.Period2Score += stat.Period2Score
	s.Period3Score += stat.Period3Score
	s.OTScore += stat.OTScore

	s.PlusMinus += stat.PlusMinus
	s.PenaltyMinutes += stat.PenaltyMinutes

	s.EvenStrengthGoals += stat.EvenStrengthGoals
	s.EvenStrengthPoints += stat.EvenStrengthPoints
	s.PowerPlayGoals += stat.PowerPlayGoals
	s.PowerPlayPoints += stat.PowerPlayPoints
	s.ShorthandedGoals += stat.ShorthandedGoals
	s.ShorthandedPoints += stat.ShorthandedPoints

	s.OvertimeGoals += stat.OvertimeGoals
	s.GameWinningGoals += stat.GameWinningGoals

	s.Shots += stat.Shots
	// note: we’ll recalc ShootingPercentage from GoalsFor/Shots at season end

	if s.Shots > 0 && s.GoalsFor > 0 {
		s.ShootingPercentage = float32(s.GoalsFor) / float32(s.Shots)
	}

	s.FaceOffsWon += stat.FaceOffsWon
	s.FaceOffs += stat.FaceOffs
	// likewise FaceOffWinPercentage = FaceOffsWon / FaceOffs
	if s.FaceOffs > 0 && s.FaceOffWinPercentage > 0 {
		s.FaceOffWinPercentage = float32(s.FaceOffsWon) / float32(s.FaceOffs)
	}

	s.ShotsAgainst += stat.ShotsAgainst
	s.Saves += stat.Saves
	// SavePercentage = float32(Saves)/float32(ShotsAgainst)
	if s.Saves > 0 && s.ShotsAgainst > 0 {
		s.SavePercentage = float32(s.Saves) / float32(s.ShotsAgainst)
	}

	s.Shutouts += stat.Shutouts
}

func (s *TeamSeasonStats) AddStatsToSeasonRecord(stat BaseTeamStats, isPostSeason, isShootout bool) {
	if isPostSeason {
		s.StatType = 3
	} else {
		s.StatType = 2
	}
	s.GamesPlayed++

	// 1) regulation vs OT wins
	if stat.GoalsFor > stat.GoalsAgainst {
		if stat.OTScore == 0 {
			s.RegulationWins++
		} else {
			s.RegulationOvertimeWins++
			if isShootout {
				s.ShootoutsWon++
			}
			// if you track shootout as separate from OT,
			// you could detect that here and bump s.ShootoutsWon++
		}
	}

	// shorthand for rolling‐average updates
	gp := float32(s.GamesPlayed)
	prevFactor := (gp - 1) / gp
	newFactor := 1 / gp

	// 2) goals per game
	s.GFGP = s.GFGP*prevFactor + float32(stat.GoalsFor)*newFactor
	s.GAGP = s.GAGP*prevFactor + float32(stat.GoalsAgainst)*newFactor

	// 3) shots per game
	s.ShotsPerGame = s.ShotsPerGame*prevFactor + float32(stat.Shots)*newFactor
	s.ShotsAgainstPerGame = s.ShotsAgainstPerGame*prevFactor + float32(stat.ShotsAgainst)*newFactor

	// 4) power‐play %  (requires you define PP opportunities)
	//    here’s an example if you tracked raw PP chances in BaseTeamStats.PowerPlayPoints:
	if stat.PowerPlayGoals+stat.PowerPlayPoints > 0 {
		gamePP := float32(stat.PowerPlayGoals) /
			float32(stat.PowerPlayGoals+stat.PowerPlayPoints) * 100
		s.PowerPlayPercentage = s.PowerPlayPercentage*prevFactor + gamePP*newFactor
	}

	// 5) penalty‐kill % (if you track raw PK opportunities similarly)
	//    PK% = (PK saves / PK shots against) * 100
	//    assuming stat.Saves and stat.ShotsAgainst under shorthanded situations:
	if stat.ShotsAgainst > 0 {
		pkPct := (float32(stat.ShotsAgainst-stat.GoalsAgainst) /
			float32(stat.ShotsAgainst)) * 100
		s.PenaltyKillPercentage = s.PenaltyKillPercentage*prevFactor + pkPct*newFactor
	}

	// 6) “net” percentages
	s.PowerPlayNetPercentage = s.PowerPlayPercentage - s.PenaltyKillPercentage
	s.PenaltyKillNetPercentage = s.PenaltyKillPercentage - s.PowerPlayPercentage
}

type CollegePlayerSeasonStats struct {
	BasePlayerStats
	StatType            uint8
	GamesPlayed         uint8
	GamesStarted        uint8
	PointsPerGamePlayed float32
}

func (s *CollegePlayerSeasonStats) AddStatsToSeasonRecord(stat BasePlayerStats) {
	// accumulate raw counts & ids
	s.BasePlayerStats.AddStatsToSeasonRecord(stat)
	if stat.StartedGame {
		s.GamesStarted++
	}
	s.StatType = 1
	s.GameType = stat.GameType
	s.GamesPlayed++
	// If `stat` had a `Started` flag, you could do:
	// if stat.Started { s.GamesStarted++ }

	// rolling average: Points per game over season
	gp := float32(s.GamesPlayed)
	s.PointsPerGamePlayed = ((s.PointsPerGamePlayed * (gp - 1)) + float32(stat.Points)) / gp
}

type CollegePlayerGameStats struct {
	WeekID        uint
	GameID        uint
	RevealResults bool
	BasePlayerStats
}

type CollegeTeamSeasonStats struct {
	BaseTeamStats
	TeamSeasonStats
}

type CollegeTeamGameStats struct {
	WeekID        uint
	GameID        uint
	RevealResults bool
	BaseTeamStats
}

type ProfessionalPlayerSeasonStats struct {
	StatType     uint8
	GamesPlayed  uint8
	GamesStarted uint8
	BasePlayerStats
}

func (s *ProfessionalPlayerSeasonStats) AddStatsToSeasonRecord(stat BasePlayerStats) {
	s.StatType = 1
	if stat.StartedGame {
		s.GamesStarted++
	}
	// accumulate player counts
	s.BasePlayerStats.AddStatsToSeasonRecord(stat)
}

type ProfessionalPlayerGameStats struct {
	WeekID        uint
	GameID        uint
	RevealResults bool
	BasePlayerStats
}

type ProfessionalTeamSeasonStats struct {
	BaseTeamStats
	TeamSeasonStats
}

type ProfessionalTeamGameStats struct {
	WeekID        uint
	GameID        uint
	RevealResults bool
	BaseTeamStats
}

type SearchStatsResponse struct {
	CHLPlayerGameStats   []CollegePlayerGameStats
	CHLPlayerSeasonStats []CollegePlayerSeasonStats
	CHLTeamGameStats     []CollegeTeamGameStats
	CHLTeamSeasonStats   []CollegeTeamSeasonStats
	PHLPlayerGameStats   []ProfessionalPlayerGameStats
	PHLPlayerSeasonStats []ProfessionalPlayerSeasonStats
	PHLTeamGameStats     []ProfessionalTeamGameStats
	PHLTeamSeasonStats   []ProfessionalTeamSeasonStats
}

type GameResultsResponse struct {
	CHLHomeStats   []CollegePlayerGameStats
	CHLAwayStats   []CollegePlayerGameStats
	CHLPlayByPlays []PlayByPlayResponse
	PHLHomeStats   []ProfessionalPlayerGameStats
	PHLAwayStats   []ProfessionalPlayerGameStats
	PHLPlayByPlays []PlayByPlayResponse
	Score          ScoreBoard
}

type ScoreBoard struct {
	P1Home              int
	P2Home              int
	P3Home              int
	OTHome              int
	P1Away              int
	P2Away              int
	P3Away              int
	OTAway              int
	HomeShootoutScore   int
	AwayShootoutScore   int
	HomeOffensiveScheme string
	HomeDefensiveScheme string
	AwayOffensiveScheme string
	AwayDefensiveScheme string
}

type ThreeStars struct {
	StarOne   uint
	StarTwo   uint
	StarThree uint
}

type ThreeStarsObj struct {
	GameID   uint
	PlayerID uint
	TeamID   uint
	Points   float32
}

func (t *ThreeStarsObj) MapPoints(stats BasePlayerStats, wonGame bool) {
	if stats.Goals > 0 {
		t.Points += float32(stats.Goals) * 1.15
	}
	if stats.GameWinningGoals > 0 {
		t.Points += float32(stats.GameWinningGoals) * 1.25
	}
	if stats.ShorthandedGoals > 0 {
		t.Points += float32(stats.ShorthandedGoals) * 0.75
	}
	if stats.OvertimeGoals > 0 {
		t.Points += float32(stats.OvertimeGoals) * 0.8
	}
	if stats.Assists > 0 {
		t.Points += float32(stats.Assists)
	}
	if stats.BodyChecks > 0 {
		t.Points += float32(stats.BodyChecks) * 0.05
	}
	if stats.StickChecks > 0 {
		t.Points += float32(stats.StickChecks) * 0.05
	}
	if stats.Saves > 0 {
		t.Points += float32(stats.Saves) * stats.SavePercentage
	}
	if stats.GoalsAgainst > 0 {
		t.Points -= float32(stats.GoalsAgainst)
	}

	if !wonGame {
		t.Points = t.Points * 0.9
	}
}
