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
		draftablePlayerType := uint8(0)
		if p.LeagueID > 1 {
			draftablePlayerType = 3
		}
		draftable := structs.DraftablePlayer{
			Model:          p.Model,
			BasePlayer:     p.BasePlayer,
			BasePotentials: p.BasePotentials,
			BaseLetterGrades: structs.BaseLetterGrades{
				OverallGrade:           util.GetDrafteeGrade(p.ID, "Overall", int(p.Overall), true),
				AgilityGrade:           util.GetDrafteeGrade(p.ID, "Agility", int(p.Agility), false),
				FaceoffsGrade:          util.GetDrafteeGrade(p.ID, "Faceoffs", int(p.Faceoffs), false),
				LongShotAccuracyGrade:  util.GetDrafteeGrade(p.ID, "LongShotAccuracy", int(p.LongShotAccuracy), false),
				LongShotPowerGrade:     util.GetDrafteeGrade(p.ID, "LongShotPower", int(p.LongShotPower), false),
				CloseShotAccuracyGrade: util.GetDrafteeGrade(p.ID, "CloseShotAccuracy", int(p.CloseShotAccuracy), false),
				CloseShotPowerGrade:    util.GetDrafteeGrade(p.ID, "CloseShotPower", int(p.CloseShotPower), false),
				OneTimerGrade:          util.GetDrafteeGrade(p.ID, "OneTimer", int(p.OneTimer), false),
				PassingGrade:           util.GetDrafteeGrade(p.ID, "Passing", int(p.Passing), false),
				PuckHandlingGrade:      util.GetDrafteeGrade(p.ID, "PuckHandling", int(p.PuckHandling), false),
				StrengthGrade:          util.GetDrafteeGrade(p.ID, "Strength", int(p.Strength), false),
				BodyCheckingGrade:      util.GetDrafteeGrade(p.ID, "BodyChecking", int(p.BodyChecking), false),
				StickCheckingGrade:     util.GetDrafteeGrade(p.ID, "StickChecking", int(p.StickChecking), false),
				ShotBlockingGrade:      util.GetDrafteeGrade(p.ID, "ShotBlocking", int(p.ShotBlocking), false),
				GoalkeepingGrade:       util.GetDrafteeGrade(p.ID, "Goalkeeping", int(p.Goalkeeping), false),
				GoalieVisionGrade:      util.GetDrafteeGrade(p.ID, "GoalieVision", int(p.GoalieVision), false),
			},
			DraftablePlayerType: draftablePlayerType,
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
				OverallGrade:           util.GetDrafteeGrade(p.ID, "Overall", int(p.Overall), true),
				AgilityGrade:           util.GetDrafteeGrade(p.ID, "Agility", int(p.Agility), false),
				FaceoffsGrade:          util.GetDrafteeGrade(p.ID, "Faceoffs", int(p.Faceoffs), false),
				LongShotAccuracyGrade:  util.GetDrafteeGrade(p.ID, "LongShotAccuracy", int(p.LongShotAccuracy), false),
				LongShotPowerGrade:     util.GetDrafteeGrade(p.ID, "LongShotPower", int(p.LongShotPower), false),
				CloseShotAccuracyGrade: util.GetDrafteeGrade(p.ID, "CloseShotAccuracy", int(p.CloseShotAccuracy), false),
				CloseShotPowerGrade:    util.GetDrafteeGrade(p.ID, "CloseShotPower", int(p.CloseShotPower), false),
				OneTimerGrade:          util.GetDrafteeGrade(p.ID, "OneTimer", int(p.OneTimer), false),
				PassingGrade:           util.GetDrafteeGrade(p.ID, "Passing", int(p.Passing), false),
				PuckHandlingGrade:      util.GetDrafteeGrade(p.ID, "PuckHandling", int(p.PuckHandling), false),
				StrengthGrade:          util.GetDrafteeGrade(p.ID, "Strength", int(p.Strength), false),
				BodyCheckingGrade:      util.GetDrafteeGrade(p.ID, "BodyChecking", int(p.BodyChecking), false),
				StickCheckingGrade:     util.GetDrafteeGrade(p.ID, "StickChecking", int(p.StickChecking), false),
				ShotBlockingGrade:      util.GetDrafteeGrade(p.ID, "ShotBlocking", int(p.ShotBlocking), false),
				GoalkeepingGrade:       util.GetDrafteeGrade(p.ID, "Goalkeeping", int(p.Goalkeeping), false),
				GoalieVisionGrade:      util.GetDrafteeGrade(p.ID, "GoalieVision", int(p.GoalieVision), false),
			},
			DraftablePlayerType: p.DraftablePlayerType,
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
