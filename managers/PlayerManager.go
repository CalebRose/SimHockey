package managers

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

// College Functions
func GetAllCollegePlayers() []structs.CollegePlayer {
	return repository.FindAllCollegePlayers(repository.PlayerQuery{})
}

func GetAllHistoricCollegePlayers() []structs.HistoricCollegePlayer {
	return repository.FindAllHistoricCollegePlayers()
}

func GetCollegePlayersByTeamID(TeamID string) []structs.CollegePlayer {
	return repository.FindCollegePlayersByTeamID(TeamID)
}

func GetCollegePlayersMap() map[uint]structs.CollegePlayer {
	players := repository.FindAllCollegePlayers(repository.PlayerQuery{})
	return MakeCollegePlayerMap(players)
}

func GetCollegePlayerMapByTeamID(TeamID string) map[uint]structs.CollegePlayer {
	players := repository.FindCollegePlayersByTeamID(TeamID)
	return MakeCollegePlayerMap(players)
}

func GetCollegePlayerByID(id string) structs.CollegePlayer {
	return repository.FindCollegePlayer(id)
}

func GetAllCollegePlayersMapByTeam() map[uint][]structs.CollegePlayer {
	players := repository.FindAllCollegePlayers(repository.PlayerQuery{})
	return MakeCollegePlayerMapByTeamID(players)
}

// Professional Functions
func GetAllProPlayers() []structs.ProfessionalPlayer {
	return repository.FindAllProPlayers(repository.PlayerQuery{})
}

func GetAllRetiredPlayers() []structs.RetiredPlayer {
	return repository.FindAllHistoricProPlayers()
}

func GetProPlayersByTeamID(TeamID string) []structs.ProfessionalPlayer {
	return repository.FindAllProPlayers(repository.PlayerQuery{TeamID: TeamID})
}

func GetProPlayersMap() map[uint]structs.ProfessionalPlayer {
	players := repository.FindAllProPlayers(repository.PlayerQuery{})
	return MakeProfessionalPlayerMap(players)
}

func GetProPlayerMapByTeamID(TeamID string) map[uint]structs.ProfessionalPlayer {
	players := repository.FindAllProPlayers(repository.PlayerQuery{TeamID: TeamID})
	return MakeProfessionalPlayerMap(players)
}

func GetAllProPlayersMapByTeam() map[uint][]structs.ProfessionalPlayer {
	players := repository.FindAllProPlayers(repository.PlayerQuery{})
	return MakeProfessionalPlayerMapByTeamID(players)
}

func RecoverPlayers() {
	db := dbprovider.GetInstance().GetDB()

	collegePlayers := GetAllCollegePlayers()

	for _, p := range collegePlayers {
		if !p.IsInjured {
			continue
		}
		p.RecoveryCheck()
		repository.SaveCollegeHockeyPlayerRecord(p, db)
	}

	for _, p := range collegePlayers {
		if p.Position != Goalie || p.GoalieStamina == 100 {
			continue
		}
		p.RecoverGoalieStamina()
		repository.SaveCollegeHockeyPlayerRecord(p, db)
	}

	proPlayers := GetAllProPlayers()

	for _, p := range proPlayers {
		if !p.IsInjured {
			continue
		}

		p.RecoveryCheck()
		repository.SaveProPlayerRecord(p, db)
	}

	for _, p := range proPlayers {
		if p.Position != Goalie || p.GoalieStamina == 100 {
			continue
		}
		p.RecoverGoalieStamina()
		repository.SaveProPlayerRecord(p, db)
	}

}
