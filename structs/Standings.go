package structs

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CollegeStandings struct {
	gorm.Model
	BaseStandings
	Rank                           uint
	IsConferenceTournamentChampion bool
	IsPostSeasonQualified          bool
	IsQuarterfinals                bool
	IsFrozenFour                   bool
	PreseasonRank                  uint8   // Preseason ranking: 1-255 is fine
	PairwiseRank                   uint8   // Rankings: 1-255 is fine
	RPIRank                        uint8   // Rankings: 1-255 is fine
	RPI                            float32 // RPI Value: 0.000-1.000+ (need decimals)
	SOS                            float32 // Strength of Schedule: decimal value
	SOR                            float32 // Strength of Record: decimal value
	Tier1Wins                      uint8   // Counts: integers are fine
	Tier2Wins                      uint8   // Counts: integers are fine
	BadLosses                      uint8   // Counts: integers are fine
	ConferenceStrengthAdj          float32 // Conference strength: decimal value
}

func (cs *CollegeStandings) AssignRank(rank int) {
	cs.Rank = uint(rank)
}

// GetWinPercentage returns the team's win percentage including OT losses as half wins
func (cs *CollegeStandings) GetWinPercentage() float32 {
	totalGames := cs.TotalWins + cs.TotalLosses + cs.TotalOTLosses
	if totalGames == 0 {
		return 0.0
	}
	adjustedWins := float32(cs.TotalWins) + (float32(cs.TotalOTLosses) * 0.5)
	return adjustedWins / float32(totalGames)
}

// GetRPIDisplay returns RPI as a decimal for display purposes
func (cs *CollegeStandings) GetRPIDisplay() string {
	return fmt.Sprintf("%.3f", cs.RPI)
}

// GetSOSDisplay returns SOS as a decimal for display purposes
func (cs *CollegeStandings) GetSOSDisplay() string {
	return fmt.Sprintf("%.3f", cs.SOS)
}

// GetSORDisplay returns SOR as a decimal for display purposes
func (cs *CollegeStandings) GetSORDisplay() string {
	return fmt.Sprintf("%.3f", cs.SOR)
}

// GetQualityRecord returns a formatted string showing quality wins
func (cs *CollegeStandings) GetQualityRecord() string {
	return fmt.Sprintf("T1: %d, T2: %d, Bad: %d", cs.Tier1Wins, cs.Tier2Wins, cs.BadLosses)
}

// IsRanked returns true if the team is in the top 25 Pairwise rankings
func (cs *CollegeStandings) IsRanked() bool {
	return cs.PairwiseRank <= 25 && cs.PairwiseRank > 0
}

func (cs *CollegeStandings) UpdateSeasonStatus(game CollegeGame) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		// if game.IsInvitational && strings.Contains(game.GameTitle, "Finals") && !strings.Contains(game.GameTitle, "Semifinals") {
		// 	cs.InvitationalChampion = true
		// }
		if game.IsConferenceTournament && strings.Contains(game.GameTitle, "Finals") && !strings.Contains(game.GameTitle, "Semifinals") {
			cs.PostSeasonStatus = "Conference Tournament Champion"
			cs.IsConferenceTournamentChampion = true
		}
		if game.IsPlayoffGame && strings.Contains(game.GameTitle, "Quarterfinals") {
			cs.PostSeasonStatus = game.GameTitle
			cs.IsQuarterfinals = true
		} else if game.IsPlayoffGame && strings.Contains(game.GameTitle, "Frozen Four") {
			cs.PostSeasonStatus = game.GameTitle
			cs.IsFrozenFour = true
		} else if game.IsPlayoffGame && game.IsNationalChampionship {
			cs.PostSeasonStatus = "National Champion"
			cs.IsNationalChampion = true
		} else if game.IsPlayoffGame {
			cs.PostSeasonStatus = "Playoffs"
			cs.IsPostSeasonQualified = true
		}
	} else {
		if game.IsPlayoffGame && game.IsNationalChampionship {
			cs.PostSeasonStatus = "National Champion Runner-Up"
			cs.IsRunnerUp = true
		}
	}
}

type ProfessionalStandings struct {
	gorm.Model
	BaseStandings
	DivisionID uint
}

type BaseStandings struct {
	TeamID             uint
	TeamName           string
	SeasonID           uint
	Season             uint
	LeagueID           uint
	ConferenceID       uint
	TotalWins          uint8
	TotalLosses        uint8
	TotalOTWins        uint8
	TotalOTLosses      uint8
	ShootoutWins       uint8
	ShootoutLosses     uint8
	ConferenceWins     uint8
	ConferenceLosses   uint8
	ConferenceOTWins   uint8
	ConferenceOTLosses uint8
	RankedWins         uint8
	RankedLosses       uint8
	Points             uint8
	GoalsFor           uint8
	GoalsAgainst       uint8
	Streak             uint8
	IsWinStreak        bool
	HomeWins           uint8
	AwayWins           uint8
	PostSeasonStatus   string
	Coach              string
	IsRunnerUp         bool
	IsNationalChampion bool
}

// Will need to include logic for OT wins/losses
func (cs *BaseStandings) UpdateStandings(game BaseGame) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cs.TotalWins += 1
		if isAway {
			cs.AwayWins += 1
		} else {
			cs.HomeWins += 1
		}
		if game.IsOvertime {
			cs.TotalOTWins += 1
		}
		if game.IsShootout {
			cs.ShootoutWins += 1
		}
		if game.IsConference {
			cs.ConferenceWins += 1
			if game.IsOvertime {
				cs.ConferenceOTWins += 1
			}
		}
		if cs.IsWinStreak {
			cs.Streak += 1
		} else {
			cs.Streak = 1
			cs.IsWinStreak = true
		}
		cs.Points += 3
	} else {
		if cs.IsWinStreak {
			cs.Streak = 1
			cs.IsWinStreak = false
		} else {
			cs.Streak += 1
		}
		if game.IsOvertime {
			cs.TotalOTLosses += 1
			cs.Points += 1
		} else {
			cs.TotalLosses += 1
		}
		if game.IsShootout {
			cs.ShootoutLosses += 1
		}
		if game.IsConference {
			if game.IsOvertime {
				cs.ConferenceOTLosses += 1
			} else {
				cs.ConferenceLosses += 1
			}
		}
	}
	if isAway {
		cs.GoalsFor += uint8(game.AwayTeamScore)
		cs.GoalsAgainst += uint8(game.HomeTeamScore)
	} else {
		cs.GoalsFor += uint8(game.HomeTeamScore)
		cs.GoalsAgainst += uint8(game.AwayTeamScore)
	}
}

func (cs *BaseStandings) SubtractStandings(game BaseGame) {
	isAway := cs.TeamID == game.AwayTeamID
	winner := (!isAway && game.HomeTeamWin) || (isAway && game.AwayTeamWin)
	if winner {
		cs.TotalWins -= 1
		if isAway {
			cs.AwayWins -= 1
		} else {
			cs.HomeWins -= 1
		}
		if game.IsConference {
			cs.ConferenceWins -= 1
		}
		cs.Streak -= 1
		if game.IsOvertime {
			cs.TotalOTWins++
			if game.IsConference {
				cs.ConferenceOTWins++
			}
		}
		if game.IsShootout {
			cs.ShootoutWins++
		}
	} else {
		cs.TotalLosses -= 1
		cs.Streak = 0
		if game.IsConference {
			cs.ConferenceLosses -= 1
		}
		if game.IsOvertime {
			cs.TotalOTLosses++
			if game.IsConference {
				cs.ConferenceOTLosses++
			}
		}
		if game.IsShootout {
			cs.ShootoutLosses++
		}
	}
	if isAway {
		cs.GoalsFor -= uint8(game.AwayTeamScore)
		cs.GoalsAgainst -= uint8(game.HomeTeamScore)
	} else {
		cs.GoalsFor -= uint8(game.HomeTeamScore)
		cs.GoalsAgainst -= uint8(game.AwayTeamScore)
	}
}

func (cs *BaseStandings) ResetStandings() {
	cs.TotalLosses = 0
	cs.TotalWins = 0
	cs.ConferenceLosses = 0
	cs.ConferenceWins = 0
	cs.TotalOTWins = 0
	cs.TotalOTLosses = 0
	cs.ConferenceOTWins = 0
	cs.ConferenceOTLosses = 0
	cs.ShootoutLosses = 0
	cs.ShootoutWins = 0
	cs.PostSeasonStatus = ""
	cs.Streak = 0
	cs.Points = 0
	cs.GoalsFor = 0
	cs.GoalsAgainst = 0
	cs.HomeWins = 0
	cs.AwayWins = 0
	cs.RankedWins = 0
	cs.RankedLosses = 0
}

func (cs *BaseStandings) MaskGames(wins, losses, confWins, confLosses, otWins, otLosses, soWins, soLosses uint8) {
	cs.TotalWins = wins
	cs.TotalLosses = losses
	cs.TotalOTWins = otWins
	cs.TotalOTLosses = otLosses
	cs.ShootoutWins = soWins
	cs.ShootoutLosses = soLosses
	cs.ConferenceWins = confWins
	cs.ConferenceLosses = confLosses
}

func (cs *BaseStandings) CalculateConferencePoints() {
	conferenceWinScore := cs.ConferenceWins * 3
	conferenceOTLScore := cs.ConferenceOTLosses * 1
	cs.Points = conferenceWinScore + conferenceOTLScore
}
