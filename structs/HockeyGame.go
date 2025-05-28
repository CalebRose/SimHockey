package structs

import "gorm.io/gorm"

type BaseGame struct {
	WeekID                uint
	Week                  int
	SeasonID              uint
	GameDay               string
	HomeTeamRank          uint
	HomeTeamID            uint
	HomeTeam              string
	HomeTeamCoach         string
	HomeTeamWin           bool
	AwayTeamRank          uint
	AwayTeamID            uint
	AwayTeam              string
	AwayTeamCoach         string
	AwayTeamWin           bool
	StarOne               uint
	StarTwo               uint
	StarThree             uint
	HomeTeamScore         uint
	AwayTeamScore         uint
	HomeTeamShootoutScore uint8
	AwayTeamShootoutScore uint8
	ArenaID               uint
	Arena                 string
	AttendanceCount       uint32
	City                  string
	State                 string
	Country               string
	IsNeutralSite         bool
	IsConference          bool
	IsPlayoffGame         bool
	IsRivalryGame         bool
	GameComplete          bool
	IsOvertime            bool
	IsShootout            bool
	GameTitle             string // For rivalry match-ups, bowl games, championships, and more
	NextGameID            uint
	NextGameHOA           string
	IsPreseason           bool
}

func (b *BaseGame) HideScore() {
	b.HomeTeamScore = 0
	b.AwayTeamScore = 0
	b.HomeTeamWin = false
	b.AwayTeamWin = false
}

func (b *BaseGame) UpdateScore(homeTeamScore, awayTeamScore, hometeamShootoutScore, awayTeamShootoutScore uint, isOvertime, isShootout bool) {
	b.HomeTeamScore = homeTeamScore
	b.AwayTeamScore = awayTeamScore
	b.HomeTeamShootoutScore = uint8(hometeamShootoutScore)
	b.AwayTeamShootoutScore = uint8(awayTeamShootoutScore)
	b.IsOvertime = isOvertime
	b.IsShootout = isShootout
	if !isShootout {
		if b.HomeTeamScore > b.AwayTeamScore {
			b.HomeTeamWin = true
		} else {
			b.AwayTeamWin = true
		}
	} else {
		if hometeamShootoutScore > awayTeamShootoutScore {
			b.HomeTeamWin = true
		} else {
			b.AwayTeamWin = true
		}
	}

	b.GameComplete = true
}

func (b *BaseGame) UpdateCoach(TeamID uint, Username string) {
	if b.HomeTeamID == TeamID {
		b.HomeTeamCoach = Username
	} else if b.AwayTeamID == TeamID {
		b.AwayTeamCoach = Username
	}
}

func (b *BaseGame) UpdateThreeStars(stars ThreeStars) {
	b.StarOne = stars.StarOne
	b.StarTwo = stars.StarTwo
	b.StarThree = stars.StarThree
}

func (m *BaseGame) AddTeam(isHome bool, id, rank uint, team, coach, arena, city, state string) {
	if isHome {
		m.HomeTeam = team
		m.HomeTeamID = id
		m.HomeTeamRank = rank
		m.HomeTeamCoach = coach
	} else {
		m.AwayTeam = team
		m.AwayTeamID = id
		m.AwayTeamRank = rank
		m.AwayTeamCoach = coach
	}
	if !m.IsNeutralSite && isHome {
		m.Arena = arena
		m.City = city
		m.State = state
	}
}

func (m *BaseGame) AssignRank(id, rank uint) {
	isHome := id == m.HomeTeamID
	if isHome {
		m.HomeTeamRank = rank
	} else {
		m.AwayTeamRank = rank
	}
}

func (m *BaseGame) AddWeekData(id, week uint, timeslot string) {
	m.WeekID = id
	m.Week = int(week)
	m.GameDay = timeslot
}

func (m *BaseGame) Reset() {
	m.GameComplete = false
	m.HomeTeamWin = false
	m.HomeTeamScore = 0
	m.AwayTeamScore = 0
	m.AwayTeamWin = false
}

type CollegeGame struct {
	gorm.Model
	BaseGame
	IsNationalChampionship bool
}

type ProfessionalGame struct {
	gorm.Model
	BaseGame
	SeriesID        uint
	IsDivisional    bool
	IsStanleyCup    bool
	IsInternational bool
}

type PlayoffSeries struct {
	gorm.Model
	SeriesName      string // For Post-Season matchups
	SeasonID        uint
	HomeTeamID      uint
	HomeTeam        string
	HomeTeamCoach   string
	HomeTeamWins    uint
	HomeTeamWin     bool
	HomeTeamRank    uint
	AwayTeamID      uint
	AwayTeam        string
	AwayTeamCoach   string
	AwayTeamWins    uint
	AwayTeamWin     bool
	AwayTeamRank    uint
	GameCount       uint
	NextSeriesID    uint
	NextSeriesHOA   string
	IsInternational bool
	IsPlayoffGame   bool
	IsTheFinals     bool
	SeriesComplete  bool
}

func (s *PlayoffSeries) AddTeam(isHome bool, id, rank uint, team, coach string) {
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

func (s *PlayoffSeries) UpdateWinCount(id int) {
	if id == int(s.HomeTeamID) {
		s.HomeTeamWins += 1
	} else {
		s.AwayTeamWins += 1
	}
	if s.GameCount < 7 {
		s.GameCount += 1
	}
	if s.HomeTeamWins > 3 {
		s.HomeTeamWin = true
		s.SeriesComplete = true
	}
	if s.AwayTeamWins > 3 {
		s.AwayTeamWin = true
		s.SeriesComplete = true
	}

}
