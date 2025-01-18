package managers

import (
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCollegeStandingsByConferenceIDAndSeasonID(conferenceID string, seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(seasonID, conferenceID, "")
}

func GetProfessionalStandingsBySeasonID(seasonID string) []structs.ProfessionalStandings {
	return repository.FindAllProfessionalStandings(seasonID, "", "")
}

func GetAllCollegeStandingsBySeasonID(seasonID string) []structs.CollegeStandings {
	return repository.FindAllCollegeStandings(seasonID, "", "")
}

func GetAllProfessionalStandingsBySeasonID(seasonID string) []structs.ProfessionalStandings {
	return repository.FindAllProfessionalStandings(seasonID, "", "")
}

func GetCollegeStandingsMap(seasonID string) map[uint]structs.CollegeStandings {
	standings := repository.FindAllCollegeStandings(seasonID, "", "")
	return MakeCollegeStandingsMap(standings)
}
