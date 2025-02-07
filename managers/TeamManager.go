package managers

import (
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCollegeTeamByTeamID(teamID string) structs.CollegeTeam {
	return repository.FindCollegeTeamRecord(teamID)
}

func GetAllCollegeTeams() []structs.CollegeTeam {
	return repository.FindAllCollegeTeams()
}

func GetCollegeTeamMap() map[uint]structs.CollegeTeam {
	teams := repository.FindAllCollegeTeams()
	return MakeCollegeTeamMap(teams)
}

func GetProTeamByTeamID(teamID string) structs.ProfessionalTeam {
	return repository.FindProTeamRecord(teamID)
}

func GetAllProfessionalTeams() []structs.ProfessionalTeam {
	return repository.FindAllProTeams()
}

func GetProTeamMap() map[uint]structs.ProfessionalTeam {
	teams := repository.FindAllProTeams()
	return MakeProTeamMap(teams)
}
