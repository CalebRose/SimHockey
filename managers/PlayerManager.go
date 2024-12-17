package managers

import (
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

// College Functions
func GetAllCollegePlayers() []structs.CollegePlayer {
	return repository.FindAllCollegePlayers()
}

func GetAllHistoricCollegePlayers() []structs.HistoricCollegePlayer {
	return repository.FindAllHistoricCollegePlayers()
}

func GetCollegePlayersByTeamID(TeamID string) []structs.CollegePlayer {
	return repository.FindCollegePlayersByTeamID(TeamID)
}

func GetCollegePlayerMapByTeamID(TeamID string) map[uint]structs.CollegePlayer {
	players := repository.FindCollegePlayersByTeamID(TeamID)
	return MakeCollegePlayerMap(players)
}

func GetAllCollegePlayersMapByTeam() map[uint][]structs.CollegePlayer {
	players := repository.FindAllCollegePlayers()
	return MakeCollegePlayerMapByTeamID(players)
}

// Professional Functions
