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
				AGZPass:       10,
				AGZPassBack:   10,
				AGZStickCheck: 15,
				AGZBodyCheck:  15,
				AZShot:        20,
				AZPass:        10,
				AZLongPass:    10,
				AZAgility:     15,
				AZStickCheck:  15,
				AZBodyCheck:   15,
				NPass:         10,
				NAgility:      15,
				NStickCheck:   15,
				NBodyCheck:    15,
				DZPass:        15,
				DZPassBack:    0,
				DZAgility:     15,
				DZStickCheck:  15,
				DZBodyCheck:   15,
				DGZLongPass:   15,
				DGZPass:       15,
				DGZAgility:    15,
				DGZStickCheck: 15,
				DGZBodyCheck:  15,
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

	for _, t := range teams {
		if len(t.Owner) > 0 {
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
				AGZShot:       15,
				AGZPass:       15,
				AGZPassBack:   15,
				AGZStickCheck: 15,
				AGZBodyCheck:  15,
				AZShot:        15,
				AZPass:        15,
				AZLongPass:    15,
				AZAgility:     15,
				AZStickCheck:  15,
				AZBodyCheck:   15,
				NPass:         15,
				NAgility:      15,
				NStickCheck:   15,
				NBodyCheck:    15,
				DZPass:        15,
				DZPassBack:    0,
				DZAgility:     15,
				DZStickCheck:  15,
				DZBodyCheck:   15,
				DGZLongPass:   15,
				DGZPass:       15,
				DGZAgility:    15,
				DGZStickCheck: 15,
				DGZBodyCheck:  15,
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
