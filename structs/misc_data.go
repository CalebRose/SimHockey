package structs

import util "github.com/CalebRose/SimHockey/_util"

type CrootLocation struct {
	City       string
	HighSchool string
}

type TeamRecordResponse struct {
	OverallWins             int
	OverallLosses           int
	CurrentSeasonWins       int
	CurrentSeasonLosses     int
	PostSeasonWins          int
	PostSeasonLosses        int
	ConferenceChampionships []string
	DivisionTitles          []string
	NationalChampionships   []string
	TopPlayers              []TopPlayer
}

func (t *TeamRecordResponse) AddTopPlayers(players []TopPlayer) {
	t.TopPlayers = players
}

type TopPlayer struct {
	PlayerID     uint
	FirstName    string
	LastName     string
	Position     string
	Archetype    string
	OverallGrade string
	Overall      int
}

func (t *TopPlayer) MapCollegePlayer(player CollegePlayer) {
	t.PlayerID = player.ID
	t.FirstName = player.FirstName
	t.LastName = player.LastName
	t.Position = player.Position
	t.Archetype = player.Archetype
	t.Overall = int(player.Overall)
	t.OverallGrade = util.GetLetterGrade(int(player.Overall), player.Year)
}

func (t *TopPlayer) MapProPlayer(player ProfessionalPlayer) {
	t.PlayerID = player.ID
	t.FirstName = player.FirstName
	t.LastName = player.LastName
	t.Position = player.Position
	t.Archetype = player.Archetype
	t.Overall = int(player.Overall)
	t.OverallGrade = util.GetLetterGrade(int(player.Overall), player.Year)
}

type CollegeGenObj struct {
	Year int
	Pos  string
}
