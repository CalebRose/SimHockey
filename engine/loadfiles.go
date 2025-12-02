package engine

import (
	"fmt"
	"sort"

	"github.com/CalebRose/SimHockey/structs"
)

func LoadAllLineStrategies(pb structs.PlayBookDTO, gameRoster []*GamePlayer) ([]LineStrategy, []LineStrategy, []LineStrategy, []uint) {
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

func LoadLineStrategies(lines []structs.BaseLineup, rosterMap map[uint]*GamePlayer) ([]LineStrategy, []uint) {
	lineStrategies := []LineStrategy{}
	activeIDs := []uint{}
	triggerGoalieReSort := false
	goalieLines := []LineStrategy{}

	for _, l := range lines {
		players := []*GamePlayer{}
		switch l.LineType {
		case 1:
			// Forward line - check for nil players before adding
			center := rosterMap[l.CenterID]
			forward1 := rosterMap[l.Forward1ID]
			forward2 := rosterMap[l.Forward2ID]

			if center != nil {
				players = append(players, center)
				activeIDs = append(activeIDs, l.CenterID)
			}
			if forward1 != nil {
				players = append(players, forward1)
				activeIDs = append(activeIDs, l.Forward1ID)
			}
			if forward2 != nil {
				players = append(players, forward2)
				activeIDs = append(activeIDs, l.Forward2ID)
			}

			// Log warning if lineup is incomplete
			if len(players) < 3 {
				fmt.Printf("WARNING: Forward line incomplete - only %d/3 players available (C:%t F1:%t F2:%t)\n",
					len(players), center != nil, forward1 != nil, forward2 != nil)
			}
		case 2:
			// Defender line - check for nil players before adding
			defender1 := rosterMap[l.Defender1ID]
			defender2 := rosterMap[l.Defender2ID]

			if defender1 != nil {
				players = append(players, defender1)
				activeIDs = append(activeIDs, l.Defender1ID)
			}
			if defender2 != nil {
				players = append(players, defender2)
				activeIDs = append(activeIDs, l.Defender2ID)
			}

			// Log warning if lineup is incomplete
			if len(players) < 2 {
				fmt.Printf("WARNING: Defender line incomplete - only %d/2 players available (D1:%t D2:%t)\n",
					len(players), defender1 != nil, defender2 != nil)
			}
		default:
			// Goalie line
			goalie := rosterMap[l.GoalieID]
			if goalie != nil {
				if goalie.GoalieStamina < GoalieStaminaThreshold {
					triggerGoalieReSort = true
				}
				players = append(players, goalie)
				activeIDs = append(activeIDs, l.GoalieID)
			} else {
				fmt.Printf("WARNING: Goalie position empty - no goalie available for lineup\n")
			}
		}

		ls := LineStrategy{
			Players:     players,
			Allocations: l.Allocations,
			CenterID:    l.CenterID,
			Forward1ID:  l.Forward1ID,
			Forward2ID:  l.Forward2ID,
			Defender1ID: l.Defender1ID,
			Defender2ID: l.Defender2ID,
		}
		if l.LineType == 3 && triggerGoalieReSort {
			goalieLines = append(goalieLines, ls)
		} else {
			lineStrategies = append(lineStrategies, ls)
		}
	}

	// Trigger if the first goalie line is below the threshold
	if triggerGoalieReSort {
		// Only two goalie lines, so really just swap the two
		sort.Slice(goalieLines, func(i, j int) bool {
			iPlayer := goalieLines[i].Players[0]
			jPlayer := goalieLines[j].Players[0]
			return iPlayer.GoalieStamina > jPlayer.GoalieStamina
		})

		lineStrategies = append(lineStrategies, goalieLines...)
	}

	return lineStrategies, activeIDs
}

func LoadGameRoster(isCollegeGame bool, collegePlayers []structs.CollegePlayer, professionalPlayers []structs.ProfessionalPlayer, seasonID uint, gameDay string, isHome bool, hra float64) []*GamePlayer {
	if isCollegeGame {
		return LoadCollegeRoster(collegePlayers, seasonID, gameDay, isHome, hra)
	}
	return LoadProfessionalRoster(professionalPlayers, seasonID, gameDay, isHome, hra)
}

func LoadCollegeRoster(roster []structs.CollegePlayer, seasonID uint, gameDay string, isHome bool, hra float64) []*GamePlayer {
	players := []*GamePlayer{}
	for _, p := range roster {
		gp := LoadCollegePlayer(p, seasonID, gameDay, isHome, hra)
		players = append(players, &gp)
	}
	return players
}

func LoadCollegePlayer(p structs.CollegePlayer, seasonID uint, gameDay string, isHome bool, hra float64) GamePlayer {
	gamePlayer := GamePlayer{
		ID:             p.ID,
		BasePlayer:     p.BasePlayer,
		CurrentStamina: int(p.Stamina),
		Stats: PlayerStatsDTO{
			PlayerID: p.ID,
			TeamID:   uint(p.TeamID),
			SeasonID: seasonID,
			GameDay:  gameDay,
		},
	}
	gamePlayer.CalculateModifiers(isHome, hra)

	return gamePlayer
}

func LoadProfessionalRoster(roster []structs.ProfessionalPlayer, seasonID uint, gameDay string, isHome bool, hra float64) []*GamePlayer {
	players := []*GamePlayer{}
	for _, p := range roster {
		gp := LoadProfessionalPlayer(p, seasonID, gameDay, isHome, hra)
		players = append(players, &gp)
	}
	return players
}

func LoadProfessionalPlayer(p structs.ProfessionalPlayer, seasonID uint, gameDay string, isHome bool, hra float64) GamePlayer {
	gamePlayer := GamePlayer{
		ID:             p.ID,
		BasePlayer:     p.BasePlayer,
		CurrentStamina: int(p.Stamina),
		Stats: PlayerStatsDTO{
			PlayerID: p.ID,
			TeamID:   uint(p.TeamID),
			SeasonID: seasonID,
			GameDay:  gameDay,
		},
	}
	gamePlayer.CalculateModifiers(isHome, hra)

	return gamePlayer
}

func getGameRosterMap(roster []*GamePlayer) map[uint]*GamePlayer {
	rosterMap := make(map[uint]*GamePlayer)

	for _, p := range roster {
		rosterMap[p.ID] = p
	}

	return rosterMap
}

func LoadBenchPlayers(activeIDs []uint, roster []*GamePlayer) []*GamePlayer {
	benchPlayers := []*GamePlayer{}
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
