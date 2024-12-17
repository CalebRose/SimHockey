package managers

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetCollegeLineupsByTeamID(TeamID string) []structs.CollegeLineup {
	return repository.FindCollegeLineupsByTeamID(TeamID)
}

func GetCollegeLineupsMap() map[uint][]structs.CollegeLineup {
	lineups := repository.FindAllCollegeLineups()
	return MakeCollegeLineupMap(lineups)
}

func RunLineupsForAICollegeTeams() {
	db := dbprovider.GetInstance().GetDB()
	teams := GetAllCollegeTeams()

	for _, t := range teams {
		if t.IsUserCoached {
			continue
		}
		fmt.Println("Iterating over Team: " + t.Abbreviation)

		teamID := strconv.Itoa(int(t.ID))

		roster := GetCollegePlayersByTeamID(teamID)
		lineups := GetCollegeLineupsByTeamID(teamID)
		// rosterMap := MakeCollegePlayerMap(roster)
		lineUpCheckMap := make(map[uint]bool)

		// Prepare roster lineups

		cPlayers := []structs.CollegePlayer{}
		fPlayers := []structs.CollegePlayer{}
		dPlayers := []structs.CollegePlayer{}
		gPlayers := []structs.CollegePlayer{}

		// Allocate players to proper list
		for _, p := range roster {
			if p.IsInjured || p.IsRedshirt {
				continue
			}

			if p.Position == Center {
				cPlayers = append(cPlayers, p)
			} else if p.Position == Forward {
				fPlayers = append(fPlayers, p)
			} else if p.Position == Defender {
				dPlayers = append(dPlayers, p)
			} else if p.Position == Goalie {
				gPlayers = append(gPlayers, p)
			}
		}

		// Sort
		sort.Slice(cPlayers, func(i, j int) bool {
			return cPlayers[i].Overall > cPlayers[j].Overall
		})

		sort.Slice(fPlayers, func(i, j int) bool {
			return fPlayers[i].Overall > fPlayers[j].Overall
		})

		sort.Slice(dPlayers, func(i, j int) bool {
			return dPlayers[i].Overall > dPlayers[j].Overall
		})

		sort.Slice(gPlayers, func(i, j int) bool {
			return gPlayers[i].Overall > gPlayers[j].Overall
		})

		for _, lineup := range lineups {
			// Get Lineup and then line
			isForwardLine := lineup.LineType == 1
			isDefenderLine := lineup.LineType == 2
			centerID := 0
			forward1ID := 0
			forward2ID := 0
			defender1ID := 0
			defender2ID := 0
			goalieID := 0
			if isForwardLine {
				for _, p := range cPlayers {
					if lineUpCheckMap[p.ID] {
						continue
					}
					lineUpCheckMap[p.ID] = true
					centerID = int(p.ID)
					break
				}
				for _, p := range fPlayers {
					if lineUpCheckMap[p.ID] {
						continue
					}
					lineUpCheckMap[p.ID] = true
					if forward1ID == 0 {
						forward1ID = int(p.ID)
					} else {
						forward2ID = int(p.ID)
					}
					if forward1ID > 0 && forward2ID > 0 {
						break
					}
				}
			} else if isDefenderLine {
				for _, p := range dPlayers {
					if lineUpCheckMap[p.ID] {
						continue
					}
					lineUpCheckMap[p.ID] = true
					if defender1ID == 0 {
						defender1ID = int(p.ID)
					} else {
						defender2ID = int(p.ID)
					}
					if defender1ID > 0 && defender2ID > 0 {
						break
					}
				}
			} else {
				for _, p := range gPlayers {
					if lineUpCheckMap[p.ID] {
						continue
					}
					lineUpCheckMap[p.ID] = true
					goalieID = int(p.ID)
					break
				}
			}
			lineUpIDs := structs.LineupPlayerIDs{
				CenterID:    uint(centerID),
				Forward1ID:  uint(forward1ID),
				Forward2ID:  uint(forward2ID),
				Defender1ID: uint(defender1ID),
				Defender2ID: uint(defender2ID),
				GoalieID:    uint(goalieID),
			}
			allocations := structs.Allocations{
				AGZShot:       20,
				AGZPass:       20,
				AGZStickCheck: 20,
				AGZBodyCheck:  20,
				AZSlapshot:    20,
				AZWristshot:   20,
				AZPass:        20,
				AZAgility:     20,
				AZStickCheck:  20,
				AZBodyCheck:   20,
				NPass:         20,
				NAgility:      20,
				NStickCheck:   20,
				NBodyCheck:    20,
				DZPass:        20,
				DZAgility:     20,
				DZStickCheck:  20,
				DZBodyCheck:   20,
				DGZPass:       20,
				DGZAgility:    20,
				DGZStickCheck: 20,
				DGZBodyCheck:  20,
			}

			lineup.MapIDsAndAllocations(lineUpIDs, allocations)
			repository.SaveCollegeLineupRecord(lineup, db)
		}
	}
}
