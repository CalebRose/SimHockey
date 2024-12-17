package managers

import "github.com/CalebRose/SimHockey/structs"

func MakeCollegePlayerMap(players []structs.CollegePlayer) map[uint]structs.CollegePlayer {
	playerMap := make(map[uint]structs.CollegePlayer)

	for _, p := range players {
		playerMap[p.ID] = p
	}

	return playerMap
}

func MakeCollegePlayerMapByTeamID(players []structs.CollegePlayer) map[uint][]structs.CollegePlayer {
	playerMap := make(map[uint][]structs.CollegePlayer)

	for _, p := range players {
		if len(playerMap[uint(p.TeamID)]) > 0 {
			playerMap[uint(p.TeamID)] = append(playerMap[uint(p.TeamID)], p)
		} else {
			playerMap[uint(p.TeamID)] = []structs.CollegePlayer{p}
		}
	}

	return playerMap
}

func MakeArenaMap(arenas []structs.Arena) map[uint]structs.Arena {
	arenaMap := make(map[uint]structs.Arena)

	for _, a := range arenas {
		arenaMap[a.ID] = a
	}

	return arenaMap
}

func MakeCollegeTeamMap(collegeTeams []structs.CollegeTeam) map[uint]structs.CollegeTeam {
	teamMap := make(map[uint]structs.CollegeTeam)
	for _, t := range collegeTeams {
		teamMap[t.ID] = t
	}
	return teamMap
}

func MakeCollegeLineupMap(lineups []structs.CollegeLineup) map[uint][]structs.CollegeLineup {
	lineupMap := make(map[uint][]structs.CollegeLineup)

	for _, l := range lineups {
		if len(lineupMap[uint(l.TeamID)]) > 0 {
			lineupMap[uint(l.TeamID)] = append(lineupMap[uint(l.TeamID)], l)
		} else {
			lineupMap[uint(l.TeamID)] = []structs.CollegeLineup{l}
		}
	}

	return lineupMap
}
