package engine

import "github.com/CalebRose/SimHockey/structs"

func LoadAllLineStrategies(pb structs.PlayBookDTO, gameRoster []GamePlayer) ([]LineStrategy, []LineStrategy, []LineStrategy, []uint) {
	rosterMap := getGameRosterMap(gameRoster)

	pbfl := pb.Forwards
	pbdl := pb.Defenders
	pbgl := pb.Goalies
	forwardLines, fIDs := LoadLineStrategies(pbfl, rosterMap)
	defenderLines, dIDs := LoadLineStrategies(pbdl, rosterMap)
	goalieLines, gIDs := LoadLineStrategies(pbgl, rosterMap)
	activeIDs := []uint{}
	activeIDs = append(activeIDs, fIDs...)
	activeIDs = append(activeIDs, dIDs...)
	activeIDs = append(activeIDs, gIDs...)

	return forwardLines, defenderLines, goalieLines, activeIDs
}

func LoadLineStrategies(lines []structs.BaseLineup, rosterMap map[uint]GamePlayer) ([]LineStrategy, []uint) {
	lineStrategies := []LineStrategy{}
	activeIDs := []uint{}
	players := []GamePlayer{}
	for _, l := range lines {

		if l.LineType == 1 {
			players = append(players, rosterMap[l.CenterID], rosterMap[l.Forward1ID], rosterMap[l.Forward2ID])
			activeIDs = append(activeIDs, l.CenterID, l.Forward1ID, l.Forward2ID)
		} else if l.LineType == 2 {
			players = append(players, rosterMap[l.Defender1ID], rosterMap[l.Defender2ID])
			activeIDs = append(activeIDs, l.Defender1ID, l.Defender2ID)
		} else {
			players = append(players, rosterMap[l.GoalieID])
			activeIDs = append(activeIDs, l.GoalieID)
		}

		ls := LineStrategy{
			Players:     players,
			Allocations: l.Allocations,
		}
		lineStrategies = append(lineStrategies, ls)
	}

	return lineStrategies, activeIDs
}

func LoadGameRoster(isCollegeGame bool, collegePlayers []structs.CollegePlayer, professionalPlayers []structs.ProfessionalPlayer) []GamePlayer {
	if isCollegeGame {
		return LoadCollegeRoster(collegePlayers)
	}
	return LoadProfessionalRoster(professionalPlayers)
}

func LoadCollegeRoster(roster []structs.CollegePlayer) []GamePlayer {
	players := []GamePlayer{}
	for _, p := range roster {
		gp := LoadCollegePlayer(p)
		players = append(players, gp)
	}
	return players
}

func LoadCollegePlayer(p structs.CollegePlayer) GamePlayer {
	gamePlayer := GamePlayer{
		ID:             p.ID,
		BasePlayer:     p.BasePlayer,
		CurrentStamina: int(p.Stamina),
	}
	gamePlayer.CalculateModifiers()

	return gamePlayer
}

func LoadProfessionalRoster(roster []structs.ProfessionalPlayer) []GamePlayer {
	players := []GamePlayer{}
	for _, p := range roster {
		gp := LoadProfessionalPlayer(p)
		players = append(players, gp)
	}
	return players
}

func LoadProfessionalPlayer(p structs.ProfessionalPlayer) GamePlayer {
	gamePlayer := GamePlayer{
		ID:             p.ID,
		BasePlayer:     p.BasePlayer,
		CurrentStamina: int(p.Stamina),
	}
	gamePlayer.CalculateModifiers()

	return gamePlayer
}

func getGameRosterMap(roster []GamePlayer) map[uint]GamePlayer {
	rosterMap := make(map[uint]GamePlayer)

	for _, p := range roster {
		rosterMap[p.ID] = p
	}

	return rosterMap
}

func LoadBenchPlayers(activeIDs []uint, roster []GamePlayer) []GamePlayer {
	benchPlayers := []GamePlayer{}
	activeIDMap := make(map[uint]bool)

	for _, id := range activeIDs {
		activeIDMap[id] = true
	}

	for _, p := range roster {
		if activeIDMap[p.ID] {
			continue
		}
		benchPlayers = append(benchPlayers, p)
	}
	return benchPlayers
}
