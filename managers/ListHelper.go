package managers

import (
	"sort"

	util "github.com/CalebRose/SimHockey/_util"
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
		if p.TransferStatus > 0 || (p.LeagueID == 2 && p.Age > 17) {
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

func MakeDraftablePlayerList(players []structs.CollegePlayer) []structs.DraftablePlayer {
	draftableList := []structs.DraftablePlayer{}

	for _, p := range players {
		if p.DraftedTeamID > 0 {
			continue
		}
		if p.Age < 18 {
			continue
		}
		draftable := structs.DraftablePlayer{
			Model:          p.Model,
			BasePlayer:     p.BasePlayer,
			BasePotentials: p.BasePotentials,
			BaseLetterGrades: structs.BaseLetterGrades{
				AgilityGrade:           util.GetLetterGrade(int(p.Agility), p.Year),
				FaceoffsGrade:          util.GetLetterGrade(int(p.Faceoffs), p.Year),
				LongShotAccuracyGrade:  util.GetLetterGrade(int(p.LongShotAccuracy), p.Year),
				LongShotPowerGrade:     util.GetLetterGrade(int(p.LongShotPower), p.Year),
				CloseShotAccuracyGrade: util.GetLetterGrade(int(p.CloseShotAccuracy), p.Year),
				CloseShotPowerGrade:    util.GetLetterGrade(int(p.CloseShotPower), p.Year),
				OneTimerGrade:          util.GetLetterGrade(int(p.OneTimer), p.Year),
				PassingGrade:           util.GetLetterGrade(int(p.Passing), p.Year),
				PuckHandlingGrade:      util.GetLetterGrade(int(p.PuckHandling), p.Year),
				StrengthGrade:          util.GetLetterGrade(int(p.Strength), p.Year),
				BodyCheckingGrade:      util.GetLetterGrade(int(p.BodyChecking), p.Year),
				StickCheckingGrade:     util.GetLetterGrade(int(p.StickChecking), p.Year),
				ShotBlockingGrade:      util.GetLetterGrade(int(p.ShotBlocking), p.Year),
				GoalkeepingGrade:       util.GetLetterGrade(int(p.Goalkeeping), p.Year),
				GoalieVisionGrade:      util.GetLetterGrade(int(p.GoalieVision), p.Year),
			},
		}

		draftableList = append(draftableList, draftable)
	}
	return draftableList
}

func MakeDraftablePlayerListWithGrades(players []structs.DraftablePlayer) []structs.DraftablePlayer {
	draftableList := []structs.DraftablePlayer{}

	for _, p := range players {
		draftable := structs.DraftablePlayer{
			Model:          p.Model,
			BasePlayer:     p.BasePlayer,
			BasePotentials: p.BasePotentials,
			BaseLetterGrades: structs.BaseLetterGrades{
				AgilityGrade:           util.GetLetterGrade(int(p.Agility), 3),
				FaceoffsGrade:          util.GetLetterGrade(int(p.Faceoffs), 3),
				LongShotAccuracyGrade:  util.GetLetterGrade(int(p.LongShotAccuracy), 3),
				LongShotPowerGrade:     util.GetLetterGrade(int(p.LongShotPower), 3),
				CloseShotAccuracyGrade: util.GetLetterGrade(int(p.CloseShotAccuracy), 3),
				CloseShotPowerGrade:    util.GetLetterGrade(int(p.CloseShotPower), 3),
				OneTimerGrade:          util.GetLetterGrade(int(p.OneTimer), 3),
				PassingGrade:           util.GetLetterGrade(int(p.Passing), 3),
				PuckHandlingGrade:      util.GetLetterGrade(int(p.PuckHandling), 3),
				StrengthGrade:          util.GetLetterGrade(int(p.Strength), 3),
				BodyCheckingGrade:      util.GetLetterGrade(int(p.BodyChecking), 3),
				StickCheckingGrade:     util.GetLetterGrade(int(p.StickChecking), 3),
				ShotBlockingGrade:      util.GetLetterGrade(int(p.ShotBlocking), 3),
				GoalkeepingGrade:       util.GetLetterGrade(int(p.Goalkeeping), 3),
				GoalieVisionGrade:      util.GetLetterGrade(int(p.GoalieVision), 3),
			},
		}

		draftableList = append(draftableList, draftable)
	}
	return draftableList
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

func MakeCollegePlayerListFromHistorics(list []structs.HistoricCollegePlayer) []structs.CollegePlayer {
	players := []structs.CollegePlayer{}
	for _, player := range list {
		players = append(players, structs.CollegePlayer{
			Model:          player.Model,
			BasePlayer:     player.BasePlayer,
			BasePotentials: player.BasePotentials,
		})
	}
	return players
}

func MakeProfessionalPlayerListFromHistorics(list []structs.RetiredPlayer) []structs.ProfessionalPlayer {
	players := []structs.ProfessionalPlayer{}
	for _, player := range list {
		players = append(players, structs.ProfessionalPlayer{
			Model:          player.Model,
			BasePlayer:     player.BasePlayer,
			BasePotentials: player.BasePotentials,
		})
	}
	return players
}

func MakeBasePlayerList(collegePlayers []structs.CollegePlayer, proPlayers []structs.ProfessionalPlayer) []structs.BasePlayer {
	players := []structs.BasePlayer{}
	for _, player := range collegePlayers {
		players = append(players, player.BasePlayer)
	}
	for _, player := range proPlayers {
		players = append(players, player.BasePlayer)
	}
	return players
}
