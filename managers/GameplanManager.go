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

func GetCollegeShootoutLineupByTeamID(TeamID string) structs.CollegeShootoutLineup {
	return repository.FindCollegeShootoutLineupByTeamID(TeamID)
}

func GetCollegeShootoutLineups() map[uint]structs.CollegeShootoutLineup {
	shootoutLineups := repository.FindAllCollegeShootoutLineups()
	return MakeCollegeShootoutLineupMap(shootoutLineups)
}

func GetCollegeLineupsMap() map[uint][]structs.CollegeLineup {
	lineups := repository.FindAllCollegeLineups()
	return MakeCollegeLineupMap(lineups)
}

func GetProLineupsByTeamID(TeamID string) []structs.ProfessionalLineup {
	return repository.FindProLineupsByTeamID(TeamID)
}

func GetProShootoutLineupByTeamID(TeamID string) structs.ProfessionalShootoutLineup {
	return repository.FindProShootoutLineupByTeamID(TeamID)
}

func GetProShootoutLineups() map[uint]structs.ProfessionalShootoutLineup {
	shootoutLineups := repository.FindAllProShootoutLineups()
	return MakeProfessionalShootoutLineupMap(shootoutLineups)
}

func GetProLineupsMap() map[uint][]structs.ProfessionalLineup {
	lineups := repository.FindAllProLineups()
	return MakeProfessionalLineupMap(lineups)
}

func SaveCHLLineup(dto structs.UpdateLineupsDTO) structs.UpdateLineupsDTO {
	db := dbprovider.GetInstance().GetDB()
	chlLineups := dto.CHLLineups
	chlSOLineup := dto.CHLShootoutLineup
	chlPlayers := dto.CollegePlayers
	teamID := strconv.Itoa(int(dto.CHLTeamID))
	// Make map of each lineup?
	chlLineupMap := MakeIndCollegeLineupMap(chlLineups)
	// Make map of each updated CHL player
	chlPlayerMap := MakeCollegePlayerMap(chlPlayers)
	// Get CHL Lineup Records
	playerIDs := []string{}
	chlLineupRecords := repository.FindCollegeLineupsByTeamID(teamID)
	for _, c := range chlLineupRecords {
		updatedLineup := chlLineupMap[c.ID]
		c.MapIDsAndAllocations(updatedLineup.LineupPlayerIDs, updatedLineup.Allocations)

		// Iterate by player
		if c.LineType == 1 {
			cID := strconv.Itoa(int(c.CenterID))
			f1ID := strconv.Itoa(int(c.Forward1ID))
			f2ID := strconv.Itoa(int(c.Forward2ID))
			playerIDs = append(playerIDs, cID, f1ID, f2ID)
		} else if c.LineType == 2 {
			d1ID := strconv.Itoa(int(c.Defender1ID))
			d2ID := strconv.Itoa(int(c.Defender2ID))
			playerIDs = append(playerIDs, d1ID, d2ID)
		} else {
			gID := strconv.Itoa(int(c.GoalieID))
			playerIDs = append(playerIDs, gID)
		}

		repository.SaveCollegeLineupRecord(c, db)
	}

	chlSORecord := repository.FindCollegeShootoutLineupByTeamID(teamID)
	chlSORecord.AssignIDs(chlSOLineup.Shooter1ID, chlSOLineup.Shooter2ID, chlSOLineup.Shooter3ID,
		chlSOLineup.Shooter4ID, chlSOLineup.Shooter5ID, chlSOLineup.Shooter6ID)
	chlSORecord.AssignShotTypes(chlSOLineup.Shooter1ShotType, chlSOLineup.Shooter2ShotType, chlSOLineup.Shooter3ShotType,
		chlSOLineup.Shooter4ShotType, chlSOLineup.Shooter5ShotType, chlSOLineup.Shooter6ShotType)

	repository.SaveCollegeShootoutLineupRecord(chlSORecord, db)

	// Get CHL Players based on updated
	collegePlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{PlayerIDs: playerIDs})

	for _, p := range collegePlayers {
		updatedPlayer := chlPlayerMap[p.ID]
		p.AssignAllocations(updatedPlayer.Allocations)
		repository.SaveCollegeHockeyPlayerRecord(p, db)
	}

	return dto
}

func SavePHLLineup(dto structs.UpdateLineupsDTO) structs.UpdateLineupsDTO {
	db := dbprovider.GetInstance().GetDB()

	phlLineups := dto.PHLLineups
	phlSOLineup := dto.PHLShootoutLineup
	phlPlayers := dto.ProPlayers
	teamID := strconv.Itoa(int(dto.CHLTeamID))
	// Make map of each lineup?
	phlLineupMap := MakeIndProLineupMap(phlLineups)
	// Make map of each updated CHL player
	phlPlayerMap := MakeProfessionalPlayerMap(phlPlayers)
	// Get CHL Lineup Records
	playerIDs := []string{}
	phlLineupRecords := repository.FindProLineupsByTeamID(teamID)
	for _, p := range phlLineupRecords {
		updatedLineup := phlLineupMap[p.ID]
		p.MapIDsAndAllocations(updatedLineup.LineupPlayerIDs, updatedLineup.Allocations)

		// Iterate by player
		if p.LineType == 1 {
			cID := strconv.Itoa(int(p.CenterID))
			f1ID := strconv.Itoa(int(p.Forward1ID))
			f2ID := strconv.Itoa(int(p.Forward2ID))
			playerIDs = append(playerIDs, cID, f1ID, f2ID)
		} else if p.LineType == 2 {
			d1ID := strconv.Itoa(int(p.Defender1ID))
			d2ID := strconv.Itoa(int(p.Defender2ID))
			playerIDs = append(playerIDs, d1ID, d2ID)
		} else {
			gID := strconv.Itoa(int(p.GoalieID))
			playerIDs = append(playerIDs, gID)
		}

		repository.SaveProfessionalLineupRecord(p, db)
	}

	phlSORecord := repository.FindProShootoutLineupByTeamID(teamID)
	phlSORecord.AssignIDs(phlSOLineup.Shooter1ID, phlSOLineup.Shooter2ID, phlSOLineup.Shooter3ID,
		phlSOLineup.Shooter4ID, phlSOLineup.Shooter5ID, phlSOLineup.Shooter6ID)
	phlSORecord.AssignShotTypes(phlSOLineup.Shooter1ShotType, phlSOLineup.Shooter2ShotType, phlSOLineup.Shooter3ShotType,
		phlSOLineup.Shooter4ShotType, phlSOLineup.Shooter5ShotType, phlSOLineup.Shooter6ShotType)

	repository.SaveProfessionalShootoutLineupRecord(phlSORecord, db)

	// Get CHL Players based on updated
	proPlayers := repository.FindAllProPlayers(repository.PlayerQuery{PlayerIDs: playerIDs})

	for _, p := range proPlayers {
		updatedPlayer := phlPlayerMap[p.ID]
		p.AssignAllocations(updatedPlayer.Allocations)
		repository.SaveProPlayerRecord(p, db)
	}

	return dto
}

func RunLineupsForAICollegeTeams() {
	db := dbprovider.GetInstance().GetDB()
	teams := GetAllCollegeTeams()
	shootoutMap := GetCollegeShootoutLineups()
	collegeGameplans := repository.FindCollegeGameplanRecords()
	collegeGameplanMap := MakeCollegeGameplanMap(collegeGameplans)

	for _, t := range teams {
		gameplan := collegeGameplanMap[t.ID]
		if gameplan.ID == 0 || !gameplan.IsAI {
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
			return GetNonGoalieSortExpression(gameplan.CenterSortPreference1, gameplan.CenterSortPreference2, gameplan.CenterSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
		})

		sort.Slice(fPlayers, func(i, j int) bool {
			return GetNonGoalieSortExpression(gameplan.ForwardSortPreference1, gameplan.ForwardSortPreference2, gameplan.ForwardSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
		})

		sort.Slice(dPlayers, func(i, j int) bool {
			return GetNonGoalieSortExpression(gameplan.DefenderSortPreference1, gameplan.DefenderSortPreference2, gameplan.DefenderSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
		})

		sort.Slice(gPlayers, func(i, j int) bool {
			return GetGoalieSortExpression(gameplan.GoalieSortPreference, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
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
			pass := 20
			longPass := 0
			backPass := 0
			closeShot := 20
			longShot := 20
			bodyCheck := 15
			stickCheck := 15
			if gameplan.LongerPassesEnabled {
				pass = 10
				longPass = 10
				backPass = 10
			}
			if isForwardLine {
				if gameplan.ForwardShotPreference == 1 {
					closeShot = 25
					longShot = 15
				} else if gameplan.ForwardShotPreference == 3 {
					closeShot = 15
					longShot = 25
				}
				if gameplan.ForwardCheckPreference == 1 {
					bodyCheck = 20
					stickCheck = 10
				} else if gameplan.ForwardCheckPreference == 3 {
					bodyCheck = 10
					stickCheck = 20
				}
			}
			if isDefenderLine {
				if gameplan.DefenderShotPreference == 1 {
					closeShot = 25
					longShot = 15
				} else if gameplan.DefenderShotPreference == 3 {
					closeShot = 15
					longShot = 25
				}
				if gameplan.DefenderCheckPreference == 1 {
					bodyCheck = 20
					stickCheck = 10
				} else if gameplan.DefenderCheckPreference == 3 {
					bodyCheck = 10
					stickCheck = 20
				}
			}
			allocations := structs.Allocations{
				AGZShot:       int8(closeShot),
				AGZPass:       int8(pass),
				AGZPassBack:   int8(backPass),
				AGZStickCheck: int8(stickCheck),
				AGZBodyCheck:  int8(bodyCheck),
				AZShot:        int8(longShot),
				AZPass:        int8(pass),
				AZLongPass:    int8(longPass),
				AZAgility:     15,
				AZStickCheck:  int8(stickCheck),
				AZBodyCheck:   int8(bodyCheck),
				NPass:         10,
				NAgility:      15,
				NStickCheck:   int8(stickCheck),
				NBodyCheck:    int8(bodyCheck),
				DZPass:        15,
				DZPassBack:    0,
				DZAgility:     15,
				DZStickCheck:  int8(stickCheck),
				DZBodyCheck:   int8(bodyCheck),
				DGZLongPass:   int8(longPass),
				DGZPass:       int8(pass),
				DGZAgility:    15,
				DGZStickCheck: int8(stickCheck),
				DGZBodyCheck:  int8(bodyCheck),
			}

			lineup.MapIDsAndAllocations(lineUpIDs, allocations)
			repository.SaveCollegeLineupRecord(lineup, db)
		}

		// Shootout Lineup
		fIdx := 0
		dIdx := 0
		cIdx := 0
		count := 0
		shootOutLineup := shootoutMap[t.ID]
		if shootOutLineup.TeamID == 0 {
			shootOutLineup = structs.CollegeShootoutLineup{
				ShootoutPlayerIDs: structs.ShootoutPlayerIDs{
					TeamID: t.ID,
				},
			}
		}
		s1ID := 0
		s1ST := 1
		s2ID := 0
		s2ST := 1
		s3ID := 0
		s3ST := 1
		s4ID := 0
		s4ST := 1
		s5ID := 0
		s5ST := 1
		s6ID := 0
		s6ST := 1

		for count < 6 {
			f := fPlayers[fIdx]
			d := dPlayers[dIdx]
			c := cPlayers[cIdx]
			chosenPlayer := structs.CollegePlayer{}

			if (cIdx < len(cPlayers)) && c.Overall >= f.Overall && c.Overall >= d.Overall {
				chosenPlayer = c
				cIdx++
			} else if (fIdx < len(fPlayers)) && f.Overall >= c.Overall && f.Overall >= d.Overall {
				chosenPlayer = f
				fIdx++
			} else if (dIdx < len(dPlayers)) && d.Overall >= f.Overall && d.Overall >= c.Overall {
				chosenPlayer = d
				dIdx++
			}
			// If no one is chosen, it means that all lists have been accumulated
			if chosenPlayer.ID == 0 {
				break
			}
			wristShotOvr := chosenPlayer.LongShotAccuracy + chosenPlayer.LongShotPower
			slapShotOvr := chosenPlayer.CloseShotAccuracy + chosenPlayer.CloseShotPower
			if s1ID == 0 {
				s1ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s1ST = 2
				}
			} else if s2ID == 0 {
				s2ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s2ST = 2
				}
			} else if s3ID == 0 {
				s3ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s3ST = 2
				}
			} else if s4ID == 0 {
				s4ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s4ST = 2
				}
			} else if s5ID == 0 {
				s5ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s5ST = 2
				}
			} else if s6ID == 0 {
				s6ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s6ST = 2
				}
			}
			count++
		}
		shootOutLineup.AssignIDs(uint(s1ID), uint(s2ID), uint(s3ID), uint(s4ID), uint(s5ID), uint(s6ID))
		shootOutLineup.AssignShotTypes(uint8(s1ST), uint8(s2ST), uint8(s3ST), uint8(s4ST), uint8(s5ST), uint8(s6ST))
		repository.SaveCollegeShootoutLineupRecord(shootOutLineup, db)
	}
}

func RunLineupsForAIProTeams() {
	db := dbprovider.GetInstance().GetDB()
	teams := repository.FindAllProTeams()
	shootoutMap := GetProShootoutLineups()
	gameplans := repository.FindProfessionalGameplanRecords()
	proGameplanMap := MakeProGameplanMap(gameplans)

	for _, t := range teams {
		gameplan := proGameplanMap[t.ID]
		if gameplan.ID == 0 || !gameplan.IsAI {
			continue
		}
		fmt.Println("Iterating over Team: " + t.Abbreviation)

		teamID := strconv.Itoa(int(t.ID))

		roster := GetProPlayersByTeamID(teamID)
		lineups := GetProLineupsByTeamID(teamID)
		// rosterMap := MakeCollegePlayerMap(roster)
		lineUpCheckMap := make(map[uint]bool)

		// Prepare roster lineups

		cPlayers := []structs.ProfessionalPlayer{}
		fPlayers := []structs.ProfessionalPlayer{}
		dPlayers := []structs.ProfessionalPlayer{}
		gPlayers := []structs.ProfessionalPlayer{}

		// Allocate players to proper list
		for _, p := range roster {
			if p.IsInjured {
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
			return GetNonGoalieSortExpression(gameplan.CenterSortPreference1, gameplan.CenterSortPreference2, gameplan.CenterSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
		})

		sort.Slice(fPlayers, func(i, j int) bool {
			return GetNonGoalieSortExpression(gameplan.ForwardSortPreference1, gameplan.ForwardSortPreference2, gameplan.ForwardSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
		})

		sort.Slice(dPlayers, func(i, j int) bool {
			return GetNonGoalieSortExpression(gameplan.DefenderSortPreference1, gameplan.DefenderSortPreference2, gameplan.DefenderSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
		})

		sort.Slice(gPlayers, func(i, j int) bool {
			return GetGoalieSortExpression(gameplan.GoalieSortPreference, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer)
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
			pass := 20
			longPass := 0
			backPass := 0
			closeShot := 20
			longShot := 20
			bodyCheck := 15
			stickCheck := 15
			if gameplan.LongerPassesEnabled {
				pass = 10
				longPass = 10
				backPass = 10
			}
			if isForwardLine {
				if gameplan.ForwardShotPreference == 1 {
					closeShot = 25
					longShot = 15
				} else if gameplan.ForwardShotPreference == 3 {
					closeShot = 15
					longShot = 25
				}
				if gameplan.ForwardCheckPreference == 1 {
					bodyCheck = 20
					stickCheck = 10
				} else if gameplan.ForwardCheckPreference == 3 {
					bodyCheck = 10
					stickCheck = 20
				}
			}
			if isDefenderLine {
				if gameplan.DefenderShotPreference == 1 {
					closeShot = 25
					longShot = 15
				} else if gameplan.DefenderShotPreference == 3 {
					closeShot = 15
					longShot = 25
				}
				if gameplan.DefenderCheckPreference == 1 {
					bodyCheck = 20
					stickCheck = 10
				} else if gameplan.DefenderCheckPreference == 3 {
					bodyCheck = 10
					stickCheck = 20
				}
			}
			allocations := structs.Allocations{
				AGZShot:       int8(closeShot),
				AGZPass:       int8(pass),
				AGZPassBack:   int8(backPass),
				AGZStickCheck: int8(stickCheck),
				AGZBodyCheck:  int8(bodyCheck),
				AZShot:        int8(longShot),
				AZPass:        int8(pass),
				AZLongPass:    int8(longPass),
				AZAgility:     15,
				AZStickCheck:  int8(stickCheck),
				AZBodyCheck:   int8(bodyCheck),
				NPass:         10,
				NAgility:      15,
				NStickCheck:   int8(stickCheck),
				NBodyCheck:    int8(bodyCheck),
				DZPass:        15,
				DZPassBack:    0,
				DZAgility:     15,
				DZStickCheck:  int8(stickCheck),
				DZBodyCheck:   int8(bodyCheck),
				DGZLongPass:   int8(longPass),
				DGZPass:       int8(pass),
				DGZAgility:    15,
				DGZStickCheck: int8(stickCheck),
				DGZBodyCheck:  int8(bodyCheck),
			}

			lineup.MapIDsAndAllocations(lineUpIDs, allocations)
			repository.SaveProfessionalLineupRecord(lineup, db)
		}

		// Shootout Lineup
		fIdx := 0
		dIdx := 0
		cIdx := 0
		count := 0
		shootOutLineup := shootoutMap[t.ID]
		if shootOutLineup.TeamID == 0 {
			shootOutLineup = structs.ProfessionalShootoutLineup{
				ShootoutPlayerIDs: structs.ShootoutPlayerIDs{
					TeamID: t.ID,
				},
			}
		}
		s1ID := 0
		s1ST := 1
		s2ID := 0
		s2ST := 1
		s3ID := 0
		s3ST := 1
		s4ID := 0
		s4ST := 1
		s5ID := 0
		s5ST := 1
		s6ID := 0
		s6ST := 1

		for count < 6 {
			f := structs.ProfessionalPlayer{}
			if fIdx < len(fPlayers) {
				f = fPlayers[fIdx]
			}
			d := structs.ProfessionalPlayer{}
			if dIdx < len(dPlayers) {
				d = dPlayers[dIdx]
			}
			c := structs.ProfessionalPlayer{}
			if cIdx < len(cPlayers) {
				c = cPlayers[cIdx]
			}
			chosenPlayer := structs.ProfessionalPlayer{}

			if (cIdx < len(cPlayers)) && c.Overall >= f.Overall && c.Overall >= d.Overall {
				chosenPlayer = c
				cIdx++
			} else if (fIdx < len(fPlayers)) && f.Overall >= c.Overall && f.Overall >= d.Overall {
				chosenPlayer = f
				fIdx++
			} else if (dIdx < len(dPlayers)) && d.Overall >= f.Overall && d.Overall >= c.Overall {
				chosenPlayer = d
				dIdx++
			}
			// If no one is chosen, it means that all lists have been accumulated
			if chosenPlayer.ID == 0 {
				break
			}
			wristShotOvr := chosenPlayer.LongShotAccuracy + chosenPlayer.LongShotPower
			slapShotOvr := chosenPlayer.CloseShotAccuracy + chosenPlayer.CloseShotPower
			if s1ID == 0 {
				s1ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s1ST = 2
				}
			} else if s2ID == 0 {
				s2ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s2ST = 2
				}
			} else if s3ID == 0 {
				s3ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s3ST = 2
				}
			} else if s4ID == 0 {
				s4ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s4ST = 2
				}
			} else if s5ID == 0 {
				s5ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s5ST = 2
				}
			} else if s6ID == 0 {
				s6ID = int(chosenPlayer.ID)
				if wristShotOvr > slapShotOvr {
					s6ST = 2
				}
			}
			count++
		}
		shootOutLineup.AssignIDs(uint(s1ID), uint(s2ID), uint(s3ID), uint(s4ID), uint(s5ID), uint(s6ID))
		shootOutLineup.AssignShotTypes(uint8(s1ST), uint8(s2ST), uint8(s3ST), uint8(s4ST), uint8(s5ST), uint8(s6ST))
		repository.SaveProfessionalShootoutLineupRecord(shootOutLineup, db)
	}
}

func CreateGameplans() {
	db := dbprovider.GetInstance().GetDB()

	proGamePlans := []structs.ProGameplan{}
	collegeGameplans := []structs.CollegeGameplan{}

	chlTeams := repository.FindAllCollegeTeams()
	proTeams := repository.FindAllProTeams()

	for _, team := range chlTeams {
		gameplan := structs.CollegeGameplan{
			BaseGameplan: structs.BaseGameplan{
				TeamID:                  team.ID,
				IsAI:                    !team.IsUserCoached,
				ForwardShotPreference:   2,
				DefenderShotPreference:  2,
				ForwardCheckPreference:  2,
				DefenderCheckPreference: 2,
				LongerPassesEnabled:     false,
				CenterSortPreference1:   1,
				ForwardSortPreference1:  1,
				DefenderSortPreference1: 1,
				GoalieSortPreference:    1,
			},
		}

		collegeGameplans = append(collegeGameplans, gameplan)
	}

	for _, team := range proTeams {
		gameplan := structs.ProGameplan{
			BaseGameplan: structs.BaseGameplan{
				TeamID:                  team.ID,
				IsAI:                    len(team.Owner) == 0,
				ForwardShotPreference:   2,
				DefenderShotPreference:  2,
				ForwardCheckPreference:  2,
				DefenderCheckPreference: 2,
				CenterSortPreference1:   1,
				ForwardSortPreference1:  1,
				DefenderSortPreference1: 1,
				LongerPassesEnabled:     false,
			},
		}

		proGamePlans = append(proGamePlans, gameplan)
	}

	repository.CreateCollegeGameplanRecordsBatch(db, collegeGameplans, 50)
	repository.CreateProfessionalGameplanRecordsBatch(db, proGamePlans, 20)
}

func SaveCollegeGameplanSettings(updatedGameplan structs.CollegeGameplan) structs.CollegeGameplan {
	db := dbprovider.GetInstance().GetDB()
	gameplanRecords := repository.FindCollegeGameplanRecords()
	gameplanMap := MakeCollegeGameplanMap(gameplanRecords)

	gameplan := gameplanMap[updatedGameplan.TeamID]
	gameplan.UpdateGameplan(updatedGameplan.BaseGameplan)

	repository.SaveCollegeGameplanRecord(gameplan, db)
	return gameplan
}

func SaveProGameplanSettings(updatedGameplan structs.ProGameplan) structs.ProGameplan {
	db := dbprovider.GetInstance().GetDB()
	gameplanRecords := repository.FindProfessionalGameplanRecords()
	gameplanMap := MakeProGameplanMap(gameplanRecords)

	gameplan := gameplanMap[updatedGameplan.TeamID]
	gameplan.UpdateGameplan(updatedGameplan.BaseGameplan)

	repository.SaveProfessionalGameplanRecord(gameplan, db)
	return gameplan
}

func GetNonGoalieSortExpression(pref1, pref2, pref3 uint8, i structs.BasePlayer, j structs.BasePlayer) bool {
	iVal := i.Overall
	jVal := j.Overall
	// Skip over checks if the preference is only overall
	if pref1 == 1 {
		return iVal > jVal
	}
	iVal1 := uint8(0)
	iVal2 := uint8(0)
	iVal3 := uint8(0)
	jVal1 := uint8(0)
	jVal2 := uint8(0)
	jVal3 := uint8(0)
	if pref1 == 2 {
		iVal1 = i.CloseShotAccuracy
		jVal1 = j.CloseShotAccuracy
	} else if pref1 == 3 {
		iVal1 = i.LongShotAccuracy
		jVal1 = j.LongShotAccuracy
	} else if pref1 == 4 {
		iVal1 = i.Agility
		jVal1 = j.Agility
	} else if pref1 == 5 {
		iVal1 = i.PuckHandling
		jVal1 = j.PuckHandling
	} else if pref1 == 6 {
		iVal1 = i.Strength
		jVal1 = j.Strength
	} else if pref1 == 7 {
		iVal1 = i.BodyChecking
		jVal1 = j.BodyChecking
	} else if pref1 == 8 {
		iVal1 = i.StickChecking
		jVal1 = j.StickChecking
	} else if pref1 == 9 {
		iVal1 = i.Faceoffs
		jVal1 = j.Faceoffs
	}
	if pref2 == 2 {
		iVal2 = i.CloseShotAccuracy
		jVal2 = j.CloseShotAccuracy
	} else if pref2 == 3 {
		iVal2 = i.LongShotAccuracy
		jVal2 = j.LongShotAccuracy
	} else if pref2 == 4 {
		iVal2 = i.Agility
		jVal2 = j.Agility
	} else if pref2 == 5 {
		iVal2 = i.PuckHandling
		jVal2 = j.PuckHandling
	} else if pref2 == 6 {
		iVal2 = i.Strength
		jVal2 = j.Strength
	} else if pref2 == 7 {
		iVal2 = i.BodyChecking
		jVal2 = j.BodyChecking
	} else if pref2 == 8 {
		iVal2 = i.StickChecking
		jVal2 = j.StickChecking
	} else if pref2 == 9 {
		iVal2 = i.Faceoffs
		jVal2 = j.Faceoffs
	}
	if pref3 == 2 {
		iVal3 = i.CloseShotAccuracy
		jVal3 = j.CloseShotAccuracy
	} else if pref3 == 3 {
		iVal3 = i.LongShotAccuracy
		jVal3 = j.LongShotAccuracy
	} else if pref3 == 4 {
		iVal3 = i.Agility
		jVal3 = j.Agility
	} else if pref3 == 5 {
		iVal3 = i.PuckHandling
		jVal3 = j.PuckHandling
	} else if pref3 == 6 {
		iVal3 = i.Strength
		jVal3 = j.Strength
	} else if pref3 == 7 {
		iVal3 = i.BodyChecking
		jVal3 = j.BodyChecking
	} else if pref3 == 8 {
		iVal3 = i.StickChecking
		jVal3 = j.StickChecking
	} else if pref3 == 9 {
		iVal3 = i.Faceoffs
		jVal3 = j.Faceoffs
	}
	finalIVal := iVal1 + iVal2 + iVal3
	finalJVal := jVal1 + jVal2 + jVal3
	return finalIVal > finalJVal
}

func GetGoalieSortExpression(preference uint8, i structs.BasePlayer, j structs.BasePlayer) bool {
	iVal := i.Overall
	jVal := j.Overall
	if preference == 2 {
		iVal = i.Goalkeeping
		jVal = j.Goalkeeping
	} else if preference == 3 {
		iVal = i.GoalieVision
		jVal = j.GoalieVision
	}
	return iVal > jVal
}
