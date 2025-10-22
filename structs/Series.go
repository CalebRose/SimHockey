package structs

import "gorm.io/gorm"

// CollegeSeries - struct for conference tournaments
type CollegeSeries struct {
	gorm.Model
	BaseSeries
	NextGameID   uint
	NextGameHOA  string
	ConferenceID uint8
}

type ProSeries struct {
	gorm.Model
	BaseSeries
	NextGameID uint
	LeagueID   uint8
}

type BaseSeries struct {
	SeriesName        string
	SeasonID          uint
	HomeTeamID        uint
	HomeTeam          string
	HomeTeamCoach     string
	HomeTeamWins      uint
	HomeTeamSeriesWin bool
	HomeTeamRank      uint
	AwayTeamID        uint
	AwayTeam          string
	AwayTeamCoach     string
	AwayTeamWins      uint
	AwayTeamSeriesWin bool
	AwayTeamRank      uint
	GameCount         uint // Number of games played in the series
	BestOfCount       uint // 3, 5, or 7
	NextSeriesID      uint
	NextSeriesHOA     string
	IsInternational   bool
	IsPlayoffGame     bool
	IsTheFinals       bool
	SeriesComplete    bool
}

func (s *BaseSeries) AddTeam(isHome bool, id, rank uint, team, coach string) {
	if isHome {
		s.HomeTeam = team
		s.HomeTeamID = id
		s.HomeTeamRank = rank
		s.HomeTeamCoach = coach
	} else {
		s.AwayTeam = team
		s.AwayTeamID = id
		s.AwayTeamRank = rank
		s.AwayTeamCoach = coach
	}
	if s.HomeTeamID > 0 && s.AwayTeamID > 0 && s.HomeTeamRank > s.AwayTeamRank {
		tempID := s.AwayTeamID
		temp := s.AwayTeam
		tempC := s.AwayTeamCoach
		tempR := s.AwayTeamRank
		s.AwayTeamID = s.HomeTeamID
		s.AwayTeam = s.HomeTeam
		s.AwayTeamCoach = s.HomeTeamCoach
		s.AwayTeamRank = s.HomeTeamRank
		s.HomeTeamID = tempID
		s.HomeTeam = temp
		s.HomeTeamCoach = tempC
		s.HomeTeamRank = tempR
	}
	s.GameCount = 1
}

func (s *BaseSeries) UpdateWinCount(id int) {
	if id == int(s.HomeTeamID) {
		s.HomeTeamWins += 1
	} else {
		s.AwayTeamWins += 1
	}
	if s.GameCount < 7 {
		s.GameCount += 1
	}
	if s.HomeTeamWins > s.BestOfCount/2 {
		s.HomeTeamSeriesWin = true
		s.SeriesComplete = true
	}
	if s.AwayTeamWins > s.BestOfCount/2 {
		s.AwayTeamSeriesWin = true
		s.SeriesComplete = true
	}

}
