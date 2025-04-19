package structs

import "gorm.io/gorm"

type BaseGame struct {
	WeekID          uint
	Week            int
	SeasonID        uint
	GameDay         string
	HomeTeamRank    uint
	HomeTeamID      uint
	HomeTeam        string
	HomeTeamCoach   string
	HomeTeamWin     bool
	AwayTeamRank    uint
	AwayTeamID      uint
	AwayTeam        string
	AwayTeamCoach   string
	AwayTeamWin     bool
	StarOne         uint
	StarTwo         uint
	StarThree       uint
	HomeTeamScore   uint
	AwayTeamScore   uint
	ArenaID         uint
	Arena           string
	AttendanceCount uint32
	City            string
	State           string
	Country         string
	IsNeutralSite   bool
	IsConference    bool
	IsPlayoffGame   bool
	IsRivalryGame   bool
	GameComplete    bool
	IsOvertime      bool
	IsShootout      bool
	GameTitle       string // For rivalry match-ups, bowl games, championships, and more
	NextGameID      uint
	NextGameHOA     string
}

func (b *BaseGame) HideScore() {
	b.HomeTeamScore = 0
	b.AwayTeamScore = 0
	b.HomeTeamWin = false
	b.AwayTeamWin = false
}

func (b *BaseGame) UpdateScore(homeTeamScore uint, awayTeamScore uint) {
	b.HomeTeamScore = homeTeamScore
	b.AwayTeamScore = awayTeamScore
	if b.HomeTeamScore > b.AwayTeamScore {
		b.HomeTeamWin = true
	} else {
		b.AwayTeamWin = true
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
	IsDivisional bool
	IsStanleyCup bool
}
