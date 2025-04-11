package managers

import "github.com/CalebRose/SimHockey/structs"

func MakeCollegeInjuryList(players []structs.CollegePlayer) []structs.CollegePlayer {
	injuryList := []structs.CollegePlayer{}

	for _, p := range players {
		if p.IsInjured {
			injuryList = append(injuryList, p)
		}
	}
	return injuryList
}

func MakeCollegePortalList(players []structs.CollegePlayer) []structs.CollegePlayer {
	portalList := []structs.CollegePlayer{}

	for _, p := range players {
		if p.TransferStatus > 0 {
			portalList = append(portalList, p)
		}
	}
	return portalList
}

func MakeProInjuryList(players []structs.ProfessionalPlayer) []structs.ProfessionalPlayer {
	injuryList := []structs.ProfessionalPlayer{}

	for _, p := range players {
		if p.IsInjured {
			injuryList = append(injuryList, p)
		}
	}
	return injuryList
}

func MakeProAffiliateList(players []structs.ProfessionalPlayer) []structs.ProfessionalPlayer {
	playerList := []structs.ProfessionalPlayer{}

	for _, p := range players {
		if p.IsAffiliatePlayer {
			playerList = append(playerList, p)
		}
	}
	return playerList
}
