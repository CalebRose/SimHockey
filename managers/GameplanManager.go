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

// ----------------------------------------------------------------
//  1. Helper: compare two Allocations structs field‐by‐field.
//     Returns true if any field differs.
//
// ----------------------------------------------------------------
func allocationsDiffer(a1, a2 structs.Allocations) bool {
	return a1.AGZShot != a2.AGZShot ||
		a1.AGZPass != a2.AGZPass ||
		a1.AGZPassBack != a2.AGZPassBack ||
		a1.AGZAgility != a2.AGZAgility ||
		a1.AGZStickCheck != a2.AGZStickCheck ||
		a1.AGZBodyCheck != a2.AGZBodyCheck ||

		a1.AZShot != a2.AZShot ||
		a1.AZPass != a2.AZPass ||
		a1.AZLongPass != a2.AZLongPass ||
		a1.AZAgility != a2.AZAgility ||
		a1.AZStickCheck != a2.AZStickCheck ||
		a1.AZBodyCheck != a2.AZBodyCheck ||

		a1.NPass != a2.NPass ||
		a1.NAgility != a2.NAgility ||
		a1.NStickCheck != a2.NStickCheck ||
		a1.NBodyCheck != a2.NBodyCheck ||

		a1.DZPass != a2.DZPass ||
		a1.DZPassBack != a2.DZPassBack ||
		a1.DZAgility != a2.DZAgility ||
		a1.DZStickCheck != a2.DZStickCheck ||
		a1.DZBodyCheck != a2.DZBodyCheck ||

		a1.DGZPass != a2.DGZPass ||
		a1.DGZLongPass != a2.DGZLongPass ||
		a1.DGZAgility != a2.DGZAgility ||
		a1.DGZStickCheck != a2.DGZStickCheck ||
		a1.DGZBodyCheck != a2.DGZBodyCheck
}

// ----------------------------------------------------------------
//
//  2. Helper: compare two CollegeLineup records (old vs updated DTO).
//     Checks both the player‐ID slots and the Allocations fields.
//     Returns true if any slot or allocation value is different.
//
//     Assumes `updated` is a DTO‐style CollegeLineup that has the
//     desired (LineupPlayerIDs + Allocations) packed in.
//
// ----------------------------------------------------------------
func chlLineupHasChanged(old *structs.CollegeLineup, updated *structs.CollegeLineup) bool {
	// Compare "slots" (LineType tells you how many IDs to check)
	if old.LineType == 1 {
		// forward line: CenterID, Forward1ID, Forward2ID
		if old.CenterID != updated.CenterID ||
			old.Forward1ID != updated.Forward1ID ||
			old.Forward2ID != updated.Forward2ID {
			return true
		}
	} else if old.LineType == 2 {
		// defender line: Defender1ID, Defender2ID
		if old.Defender1ID != updated.Defender1ID ||
			old.Defender2ID != updated.Defender2ID {
			return true
		}
	} else {
		// goalie line: GoalieID
		if old.GoalieID != updated.GoalieID {
			return true
		}
	}

	// If the player‐ID slots match, check allocations:
	if allocationsDiffer(old.Allocations, updated.Allocations) {
		return true
	}

	// If we got here, nothing changed:
	return false
}

func phlLineupHasChanged(old *structs.ProfessionalLineup, updated *structs.ProfessionalLineup) bool {
	// Compare "slots" (LineType tells you how many IDs to check)
	if old.LineType == 1 {
		// forward line: CenterID, Forward1ID, Forward2ID
		if old.CenterID != updated.CenterID ||
			old.Forward1ID != updated.Forward1ID ||
			old.Forward2ID != updated.Forward2ID {
			return true
		}
	} else if old.LineType == 2 {
		// defender line: Defender1ID, Defender2ID
		if old.Defender1ID != updated.Defender1ID ||
			old.Defender2ID != updated.Defender2ID {
			return true
		}
	} else {
		// goalie line: GoalieID
		if old.GoalieID != updated.GoalieID {
			return true
		}
	}

	// If the player‐ID slots match, check allocations:
	if allocationsDiffer(old.Allocations, updated.Allocations) {
		return true
	}

	// If we got here, nothing changed:
	return false
}

// ----------------------------------------------------------------
//  3. Helper: compare two CollegePlayer records’ Allocations.
//     Returns true if any allocation field is different.
//
// ----------------------------------------------------------------
func playerAllocationsDiffer(old *structs.BasePlayer, newAlloc structs.Allocations) bool {
	return allocationsDiffer(old.Allocations, newAlloc)
}

func SaveCHLLineup(dto structs.UpdateLineupsDTO) structs.UpdateLineupsDTO {
	db := dbprovider.GetInstance().GetDB()
	incomingLineups := dto.CHLLineups
	incomingSOLineup := dto.CHLShootoutLineup
	incomingPlayers := dto.CollegePlayers
	teamID := strconv.Itoa(int(dto.CHLTeamID))
	// Make map of each lineup?
	chlLineupMap := MakeIndCollegeLineupMap(incomingLineups)
	// Make map of each updated CHL player
	chlPlayerMap := MakeCollegePlayerMap(incomingPlayers)
	// Get CHL Lineup Records
	playerIDs := []string{}
	existingLineupRecords := repository.FindCollegeLineupsByTeamID(teamID)
	for _, rec := range existingLineupRecords {
		updatedLineup := chlLineupMap[rec.ID]
		// Iterate by player
		if rec.LineType == 1 {
			cID := strconv.Itoa(int(rec.CenterID))
			f1ID := strconv.Itoa(int(rec.Forward1ID))
			f2ID := strconv.Itoa(int(rec.Forward2ID))
			playerIDs = append(playerIDs, cID, f1ID, f2ID)
		} else if rec.LineType == 2 {
			d1ID := strconv.Itoa(int(rec.Defender1ID))
			d2ID := strconv.Itoa(int(rec.Defender2ID))
			playerIDs = append(playerIDs, d1ID, d2ID)
		} else {
			gID := strconv.Itoa(int(rec.GoalieID))
			playerIDs = append(playerIDs, gID)
		}
		changedLineupCheck := chlLineupHasChanged(&rec, &updatedLineup)
		if !changedLineupCheck {
			continue
		}
		rec.MapIDsAndAllocations(updatedLineup.LineupPlayerIDs, updatedLineup.Allocations)
		repository.SaveCollegeLineupRecord(rec, db)
	}

	soRec := repository.FindCollegeShootoutLineupByTeamID(teamID)
	soRec.AssignIDs(incomingSOLineup.Shooter1ID, incomingSOLineup.Shooter2ID, incomingSOLineup.Shooter3ID,
		incomingSOLineup.Shooter4ID, incomingSOLineup.Shooter5ID, incomingSOLineup.Shooter6ID)
	soRec.AssignShotTypes(incomingSOLineup.Shooter1ShotType, incomingSOLineup.Shooter2ShotType, incomingSOLineup.Shooter3ShotType,
		incomingSOLineup.Shooter4ShotType, incomingSOLineup.Shooter5ShotType, incomingSOLineup.Shooter6ShotType)

	repository.SaveCollegeShootoutLineupRecord(soRec, db)

	// Get CHL Players based on updated
	collegePlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{PlayerIDs: playerIDs})

	for _, p := range collegePlayers {
		if p.ID == 0 {
			continue
		}
		updatedPlayer, exists := chlPlayerMap[p.ID]
		if !exists {
			continue
		}
		allocationChangeCheck := playerAllocationsDiffer(&p.BasePlayer, updatedPlayer.Allocations)
		if !allocationChangeCheck {
			continue
		}
		p.AssignAllocations(updatedPlayer.Allocations)
		repository.SaveCollegeHockeyPlayerRecord(p, db)
	}

	return dto
}

func SavePHLLineup(dto structs.UpdateLineupsDTO) structs.UpdateLineupsDTO {
	db := dbprovider.GetInstance().GetDB()

	incomingLineups := dto.PHLLineups
	incomingSOLineup := dto.PHLShootoutLineup
	incomingPlayers := dto.ProPlayers
	teamID := strconv.Itoa(int(dto.PHLTeamID))
	// Make map of each lineup?
	phlLineupMap := MakeIndProLineupMap(incomingLineups)
	// Make map of each updated CHL player
	phlPlayerMap := MakeProfessionalPlayerMap(incomingPlayers)
	// Get CHL Lineup Records
	playerIDs := []string{}
	existingRecords := repository.FindProLineupsByTeamID(teamID)
	for _, rec := range existingRecords {
		updatedLineup := phlLineupMap[rec.ID]

		// Iterate by player
		if rec.LineType == 1 {
			cID := strconv.Itoa(int(rec.CenterID))
			f1ID := strconv.Itoa(int(rec.Forward1ID))
			f2ID := strconv.Itoa(int(rec.Forward2ID))
			playerIDs = append(playerIDs, cID, f1ID, f2ID)
		} else if rec.LineType == 2 {
			d1ID := strconv.Itoa(int(rec.Defender1ID))
			d2ID := strconv.Itoa(int(rec.Defender2ID))
			playerIDs = append(playerIDs, d1ID, d2ID)
		} else {
			gID := strconv.Itoa(int(rec.GoalieID))
			playerIDs = append(playerIDs, gID)
		}
		changedLineupCheck := phlLineupHasChanged(&rec, &updatedLineup)
		if !changedLineupCheck {
			continue
		}
		rec.MapIDsAndAllocations(updatedLineup.LineupPlayerIDs, updatedLineup.Allocations)

		repository.SaveProfessionalLineupRecord(rec, db)
	}

	phlSORecord := repository.FindProShootoutLineupByTeamID(teamID)
	phlSORecord.AssignIDs(incomingSOLineup.Shooter1ID, incomingSOLineup.Shooter2ID, incomingSOLineup.Shooter3ID,
		incomingSOLineup.Shooter4ID, incomingSOLineup.Shooter5ID, incomingSOLineup.Shooter6ID)
	phlSORecord.AssignShotTypes(incomingSOLineup.Shooter1ShotType, incomingSOLineup.Shooter2ShotType, incomingSOLineup.Shooter3ShotType,
		incomingSOLineup.Shooter4ShotType, incomingSOLineup.Shooter5ShotType, incomingSOLineup.Shooter6ShotType)

	repository.SaveProfessionalShootoutLineupRecord(phlSORecord, db)

	// Get CHL Players based on updated
	proPlayers := repository.FindAllProPlayers(repository.PlayerQuery{PlayerIDs: playerIDs})

	for _, p := range proPlayers {
		if p.ID == 0 {
			continue
		}
		updatedPlayer, exists := phlPlayerMap[p.ID]
		if !exists {
			continue
		}
		allocationChangeCheck := playerAllocationsDiffer(&p.BasePlayer, updatedPlayer.Allocations)
		if !allocationChangeCheck {
			continue
		}
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
			return GetNonGoalieSortExpression(gameplan.ForwardSortPreference1, gameplan.ForwardSortPreference2, gameplan.ForwardSortPreference3, fPlayers[i].BasePlayer, fPlayers[j].BasePlayer)
		})

		sort.Slice(dPlayers, func(i, j int) bool {
			return GetNonGoalieSortExpression(gameplan.DefenderSortPreference1, gameplan.DefenderSortPreference2, gameplan.DefenderSortPreference3, dPlayers[i].BasePlayer, dPlayers[j].BasePlayer)
		})

		sort.Slice(gPlayers, func(i, j int) bool {
			return GetGoalieSortExpression(gameplan.GoalieSortPreference, gPlayers[i].BasePlayer, gPlayers[j].BasePlayer)
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
			bodyCheck := 10
			stickCheck := 10
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
					bodyCheck = 15
					stickCheck = 5
				} else if gameplan.ForwardCheckPreference == 3 {
					bodyCheck = 5
					stickCheck = 15
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
					bodyCheck = 15
					stickCheck = 5
				} else if gameplan.DefenderCheckPreference == 3 {
					bodyCheck = 5
					stickCheck = 15
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
			if p.IsInjured || p.IsAffiliatePlayer {
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
			return GetNonGoalieSortExpression(gameplan.ForwardSortPreference1, gameplan.ForwardSortPreference2, gameplan.ForwardSortPreference3, fPlayers[i].BasePlayer, fPlayers[j].BasePlayer)
		})

		sort.Slice(dPlayers, func(i, j int) bool {
			return GetNonGoalieSortExpression(gameplan.DefenderSortPreference1, gameplan.DefenderSortPreference2, gameplan.DefenderSortPreference3, dPlayers[i].BasePlayer, dPlayers[j].BasePlayer)
		})

		sort.Slice(gPlayers, func(i, j int) bool {
			return GetGoalieSortExpression(gameplan.GoalieSortPreference, gPlayers[i].BasePlayer, gPlayers[j].BasePlayer)
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
			bodyCheck := 10
			stickCheck := 10
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
					bodyCheck = 15
					stickCheck = 5
				} else if gameplan.ForwardCheckPreference == 3 {
					bodyCheck = 5
					stickCheck = 15
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
					bodyCheck = 15
					stickCheck = 5
				} else if gameplan.DefenderCheckPreference == 3 {
					bodyCheck = 5
					stickCheck = 15
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
	} else if pref1 == 10 {
		iVal1 = i.Passing
		jVal1 = j.Passing
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
	} else if pref2 == 10 {
		iVal2 = i.Passing
		jVal2 = j.Passing
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
	} else if pref3 == 10 {
		iVal3 = i.Passing
		jVal3 = j.Passing
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
