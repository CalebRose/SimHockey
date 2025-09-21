package managers

import (
	"sort"

	"github.com/CalebRose/SimHockey/structs"
)

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

func GetCollegeOrderedListByStatType(statType string, teamID uint, CollegeStats []structs.CollegePlayerSeasonStats, collegePlayerMap map[uint]structs.CollegePlayer) []structs.CollegePlayer {
	orderedStats := CollegeStats
	resultList := []structs.CollegePlayer{}
	switch statType {
	case "GOALS":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].Goals > orderedStats[j].Goals
		})
	case "ASSISTS":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].Assists > orderedStats[j].Assists
		})
	case "SAVES":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].Saves > orderedStats[j].Saves
		})
	}

	teamLeaderInTopStats := false
	for idx, stat := range orderedStats {
		if idx > 4 {
			break
		}
		player := collegePlayerMap[stat.PlayerID]
		if stat.TeamID == teamID {
			teamLeaderInTopStats = true
		}
		player.AddSeasonStats(stat)
		resultList = append(resultList, player)
	}

	if !teamLeaderInTopStats {
		for _, stat := range orderedStats {
			if stat.TeamID == teamID {
				player := collegePlayerMap[stat.PlayerID]
				player.AddSeasonStats(stat)
				resultList = append(resultList, player)
				break
			}
		}
	}
	return resultList
}

func GetProOrderedListByStatType(statType string, teamID uint, CollegeStats []structs.ProfessionalPlayerSeasonStats, proPlayerMap map[uint]structs.ProfessionalPlayer) []structs.ProfessionalPlayer {
	orderedStats := CollegeStats
	resultList := []structs.ProfessionalPlayer{}
	switch statType {
	case "GOALS":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].Goals > orderedStats[j].Goals
		})
	case "ASSISTS":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].Assists > orderedStats[j].Assists
		})
	case "SAVES":
		sort.Slice(orderedStats[:], func(i, j int) bool {
			return orderedStats[i].Saves > orderedStats[j].Saves
		})
	}

	teamLeaderInTopStats := false
	for idx, stat := range orderedStats {
		if idx > 4 {
			break
		}
		player := proPlayerMap[stat.PlayerID]
		if stat.TeamID == teamID {
			teamLeaderInTopStats = true
		}
		player.AddSeasonStats(stat)
		resultList = append(resultList, player)
	}

	if !teamLeaderInTopStats {
		for _, stat := range orderedStats {
			if stat.TeamID == teamID {
				player := proPlayerMap[stat.PlayerID]
				player.AddSeasonStats(stat)
				resultList = append(resultList, player)
				break
			}
		}
	}
	return resultList
}

func MakeCollegePlayerGameStatsListByTeamID(list []structs.CollegePlayerGameStats, teamID uint) []structs.CollegePlayerGameStats {
	stats := []structs.CollegePlayerGameStats{}
	for _, stat := range list {
		if stat.TeamID != teamID {
			continue
		}
		stats = append(stats, stat)
	}
	return stats
}

func MakeProPlayerGameStatsListByTeamID(list []structs.ProfessionalPlayerGameStats, teamID uint) []structs.ProfessionalPlayerGameStats {
	stats := []structs.ProfessionalPlayerGameStats{}
	for _, stat := range list {
		if stat.TeamID != teamID {
			continue
		}
		stats = append(stats, stat)
	}
	return stats
}
