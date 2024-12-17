package managers

import (
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllCollegeTeams() []structs.CollegeTeam {
	return repository.FindAllCollegeTeams()
}

func GetCollegeTeamMap() map[uint]structs.CollegeTeam {
	teams := repository.FindAllCollegeTeams()
	return MakeCollegeTeamMap(teams)
}
