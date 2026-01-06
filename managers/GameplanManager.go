package managers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
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
	switch old.LineType {
	case 1:
		// forward line: CenterID, Forward1ID, Forward2ID
		if old.CenterID != updated.CenterID ||
			old.Forward1ID != updated.Forward1ID ||
			old.Forward2ID != updated.Forward2ID {
			return true
		}
	case 2:
		// defender line: Defender1ID, Defender2ID
		if old.Defender1ID != updated.Defender1ID ||
			old.Defender2ID != updated.Defender2ID {
			return true
		}
	default:
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
	switch old.LineType {
	case 1:
		// forward line: CenterID, Forward1ID, Forward2ID
		if old.CenterID != updated.CenterID ||
			old.Forward1ID != updated.Forward1ID ||
			old.Forward2ID != updated.Forward2ID {
			return true
		}
	case 2:
		// defender line: Defender1ID, Defender2ID
		if old.Defender1ID != updated.Defender1ID ||
			old.Defender2ID != updated.Defender2ID {
			return true
		}
	default:
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
		switch rec.LineType {
		case 1:
			cID := strconv.Itoa(int(rec.CenterID))
			f1ID := strconv.Itoa(int(rec.Forward1ID))
			f2ID := strconv.Itoa(int(rec.Forward2ID))
			playerIDs = append(playerIDs, cID, f1ID, f2ID)
		case 2:
			d1ID := strconv.Itoa(int(rec.Defender1ID))
			d2ID := strconv.Itoa(int(rec.Defender2ID))
			playerIDs = append(playerIDs, d1ID, d2ID)
		default:
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
		switch rec.LineType {
		case 1:
			cID := strconv.Itoa(int(rec.CenterID))
			f1ID := strconv.Itoa(int(rec.Forward1ID))
			f2ID := strconv.Itoa(int(rec.Forward2ID))
			playerIDs = append(playerIDs, cID, f1ID, f2ID)
		case 2:
			d1ID := strconv.Itoa(int(rec.Defender1ID))
			d2ID := strconv.Itoa(int(rec.Defender2ID))
			playerIDs = append(playerIDs, d1ID, d2ID)
		default:
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
		// if gameplan.ID == 0 || !gameplan.IsAI {
		// 	continue
		// }
		fmt.Println("Iterating over Team: " + t.Abbreviation)

		teamID := strconv.Itoa(int(t.ID))

		roster := GetCollegePlayersByTeamID(teamID)

		// VALIDATE ROSTER BEFORE LINEUP GENERATION (warnings only)
		rosterComplete := validateCollegeRoster(t, roster)

		// VALIDATE AND FIX LINEUPS (only if roster is reasonably complete)
		if rosterComplete {
			if !validateAndFixCollegeLineups(t, roster, db) {
				fmt.Printf("WARNING: Team %s still has incomplete lineups after attempted fix\n", t.Abbreviation)
			}
		} else {
			fmt.Printf("SKIPPING lineup generation for %s due to incomplete roster\n", t.Abbreviation)
		}

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

			switch p.Position {
			case Center:
				cPlayers = append(cPlayers, p)
			case Forward:
				fPlayers = append(fPlayers, p)
			case Defender:
				dPlayers = append(dPlayers, p)
			case Goalie:
				gPlayers = append(gPlayers, p)
			}
		}

		// Sort with system compatibility weighting
		sort.Slice(cPlayers, func(i, j int) bool {
			return GetSystemWeightedSortExpression(gameplan.CenterSortPreference1, gameplan.CenterSortPreference2, gameplan.CenterSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer, gameplan.OffensiveSystem, gameplan.DefensiveSystem, gameplan.OffensiveIntensity, gameplan.DefensiveIntensity, false)
		})

		sort.Slice(fPlayers, func(i, j int) bool {
			return GetSystemWeightedSortExpression(gameplan.ForwardSortPreference1, gameplan.ForwardSortPreference2, gameplan.ForwardSortPreference3, fPlayers[i].BasePlayer, fPlayers[j].BasePlayer, gameplan.OffensiveSystem, gameplan.DefensiveSystem, gameplan.OffensiveIntensity, gameplan.DefensiveIntensity, false)
		})

		sort.Slice(dPlayers, func(i, j int) bool {
			return GetSystemWeightedSortExpression(gameplan.DefenderSortPreference1, gameplan.DefenderSortPreference2, gameplan.DefenderSortPreference3, dPlayers[i].BasePlayer, dPlayers[j].BasePlayer, gameplan.OffensiveSystem, gameplan.DefensiveSystem, gameplan.OffensiveIntensity, gameplan.DefensiveIntensity, false)
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
			longPass := 5
			backPass := 5
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
				switch gameplan.ForwardShotPreference {
				case 1:
					closeShot = 25
					longShot = 15
				case 3:
					closeShot = 15
					longShot = 25
				}
				switch gameplan.ForwardCheckPreference {
				case 1:
					bodyCheck = 15
					stickCheck = 5
				case 3:
					bodyCheck = 5
					stickCheck = 15
				}
			}
			if isDefenderLine {
				switch gameplan.DefenderShotPreference {
				case 1:
					closeShot = 25
					longShot = 15
				case 3:
					closeShot = 15
					longShot = 25
				}
				switch gameplan.DefenderCheckPreference {
				case 1:
					bodyCheck = 15
					stickCheck = 5
				case 3:
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
	teams := repository.FindAllProTeams(repository.TeamClauses{LeagueID: "1"})
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

		// VALIDATE ROSTER BEFORE LINEUP GENERATION (warnings only)
		rosterComplete := validateProRoster(t, roster)

		// Note: Pro lineup validation would be similar to college - implement if needed
		if !rosterComplete {
			fmt.Printf("SKIPPING lineup generation for %s due to incomplete roster\n", t.Abbreviation)
			continue // Skip this team if roster is incomplete
		}

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

			switch p.Position {
			case Center:
				cPlayers = append(cPlayers, p)
			case Forward:
				fPlayers = append(fPlayers, p)
			case Defender:
				dPlayers = append(dPlayers, p)
			case Goalie:
				gPlayers = append(gPlayers, p)
			}
		}

		// Sort with system compatibility weighting
		sort.Slice(cPlayers, func(i, j int) bool {
			return GetSystemWeightedSortExpression(gameplan.CenterSortPreference1, gameplan.CenterSortPreference2, gameplan.CenterSortPreference3, cPlayers[i].BasePlayer, cPlayers[j].BasePlayer, gameplan.OffensiveSystem, gameplan.DefensiveSystem, gameplan.OffensiveIntensity, gameplan.DefensiveIntensity, false)
		})

		sort.Slice(fPlayers, func(i, j int) bool {
			return GetSystemWeightedSortExpression(gameplan.ForwardSortPreference1, gameplan.ForwardSortPreference2, gameplan.ForwardSortPreference3, fPlayers[i].BasePlayer, fPlayers[j].BasePlayer, gameplan.OffensiveSystem, gameplan.DefensiveSystem, gameplan.OffensiveIntensity, gameplan.DefensiveIntensity, false)
		})

		sort.Slice(dPlayers, func(i, j int) bool {
			return GetSystemWeightedSortExpression(gameplan.DefenderSortPreference1, gameplan.DefenderSortPreference2, gameplan.DefenderSortPreference3, dPlayers[i].BasePlayer, dPlayers[j].BasePlayer, gameplan.OffensiveSystem, gameplan.DefensiveSystem, gameplan.OffensiveIntensity, gameplan.DefensiveIntensity, false)
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
				switch gameplan.ForwardShotPreference {
				case 1:
					closeShot = 25
					longShot = 15
				case 3:
					closeShot = 15
					longShot = 25
				}
				switch gameplan.ForwardCheckPreference {
				case 1:
					bodyCheck = 15
					stickCheck = 5
				case 3:
					bodyCheck = 5
					stickCheck = 15
				}
			}
			if isDefenderLine {
				switch gameplan.DefenderShotPreference {
				case 1:
					closeShot = 25
					longShot = 15
				case 3:
					closeShot = 15
					longShot = 25
				}
				switch gameplan.DefenderCheckPreference {
				case 1:
					bodyCheck = 15
					stickCheck = 5
				case 3:
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

	chlTeams := repository.FindAllCollegeTeams(repository.TeamClauses{})
	proTeams := repository.FindAllProTeams(repository.TeamClauses{})

	for _, team := range chlTeams {
		if team.ID < 73 {
			continue
		}
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
	// repository.CreateProfessionalGameplanRecordsBatch(db, proGamePlans, 20)
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
	switch pref1 {
	case 2:
		iVal1 = i.CloseShotAccuracy
		jVal1 = j.CloseShotAccuracy
	case 3:
		iVal1 = i.LongShotAccuracy
		jVal1 = j.LongShotAccuracy
	case 4:
		iVal1 = i.Agility
		jVal1 = j.Agility
	case 5:
		iVal1 = i.PuckHandling
		jVal1 = j.PuckHandling
	case 6:
		iVal1 = i.Strength
		jVal1 = j.Strength
	case 7:
		iVal1 = i.BodyChecking
		jVal1 = j.BodyChecking
	case 8:
		iVal1 = i.StickChecking
		jVal1 = j.StickChecking
	case 9:
		iVal1 = i.Faceoffs
		jVal1 = j.Faceoffs
	case 10:
		iVal1 = i.Passing
		jVal1 = j.Passing
	}
	switch pref2 {
	case 2:
		iVal2 = i.CloseShotAccuracy
		jVal2 = j.CloseShotAccuracy
	case 3:
		iVal2 = i.LongShotAccuracy
		jVal2 = j.LongShotAccuracy
	case 4:
		iVal2 = i.Agility
		jVal2 = j.Agility
	case 5:
		iVal2 = i.PuckHandling
		jVal2 = j.PuckHandling
	case 6:
		iVal2 = i.Strength
		jVal2 = j.Strength
	case 7:
		iVal2 = i.BodyChecking
		jVal2 = j.BodyChecking
	case 8:
		iVal2 = i.StickChecking
		jVal2 = j.StickChecking
	case 9:
		iVal2 = i.Faceoffs
		jVal2 = j.Faceoffs
	case 10:
		iVal2 = i.Passing
		jVal2 = j.Passing
	}
	switch pref3 {
	case 2:
		iVal3 = i.CloseShotAccuracy
		jVal3 = j.CloseShotAccuracy
	case 3:
		iVal3 = i.LongShotAccuracy
		jVal3 = j.LongShotAccuracy
	case 4:
		iVal3 = i.Agility
		jVal3 = j.Agility
	case 5:
		iVal3 = i.PuckHandling
		jVal3 = j.PuckHandling
	case 6:
		iVal3 = i.Strength
		jVal3 = j.Strength
	case 7:
		iVal3 = i.BodyChecking
		jVal3 = j.BodyChecking
	case 8:
		iVal3 = i.StickChecking
		jVal3 = j.StickChecking
	case 9:
		iVal3 = i.Faceoffs
		jVal3 = j.Faceoffs
	case 10:
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
	switch preference {
	case 2:
		iVal = i.Goalkeeping
		jVal = j.Goalkeeping
	case 3:
		iVal = i.GoalieVision
		jVal = j.GoalieVision
	}
	return iVal > jVal
}

// selectOptimalSystemsForRoster analyzes a college roster and returns the best offensive/defensive systems plus intensity
func selectOptimalSystemsForRoster(roster []structs.CollegePlayer) (structs.OffensiveSystemType, structs.DefensiveSystemType, uint8) {
	// Convert to BasePlayer for analysis
	baseRoster := make([]structs.BasePlayer, len(roster))
	for i, p := range roster {
		baseRoster[i] = p.BasePlayer
	}
	return analyzeRosterForSystems(baseRoster)
}

// selectOptimalSystemsForProRoster analyzes a pro roster and returns the best offensive/defensive systems plus intensity
func selectOptimalSystemsForProRoster(roster []structs.ProfessionalPlayer) (structs.OffensiveSystemType, structs.DefensiveSystemType, uint8) {
	// Convert to BasePlayer for analysis
	baseRoster := make([]structs.BasePlayer, len(roster))
	for i, p := range roster {
		baseRoster[i] = p.BasePlayer
	}
	return analyzeRosterForSystems(baseRoster)
}

// analyzeRosterForSystems performs the core system selection logic
func analyzeRosterForSystems(roster []structs.BasePlayer) (structs.OffensiveSystemType, structs.DefensiveSystemType, uint8) {
	// Analyze roster composition
	archCounts := analyzeArchetypeComposition(roster)
	avgOverall := calculateAverageOverall(roster)

	// Determine intensity based on roster quality and cohesion
	intensity := determineOptimalIntensity(archCounts, avgOverall)

	// Test all system combinations and find the best
	bestOffensive, bestOffensiveScore := findBestOffensiveSystem(archCounts, intensity)
	bestDefensive, bestDefensiveScore := findBestDefensiveSystem(archCounts, intensity)

	// Adjust intensity based on system compatibility
	intensity = adjustIntensityForSystemCompatibility(intensity, bestOffensiveScore, bestDefensiveScore)

	return bestOffensive, bestDefensive, intensity
}

// analyzeArchetypeComposition counts archetypes by position
func analyzeArchetypeComposition(roster []structs.BasePlayer) map[string]map[string]int {
	archCounts := map[string]map[string]int{
		"F": make(map[string]int), // Forwards (C + F)
		"D": make(map[string]int), // Defensemen
		"G": make(map[string]int), // Goalies
	}

	for _, player := range roster {
		var posCategory string
		switch player.Position {
		case "C", "F":
			posCategory = "F"
		case "D":
			posCategory = "D"
		case "G":
			posCategory = "G"
		default:
			continue
		}

		archCounts[posCategory][player.Archetype]++
	}

	return archCounts
}

// calculateAverageOverall determines roster quality
func calculateAverageOverall(roster []structs.BasePlayer) float64 {
	if len(roster) == 0 {
		return 50.0
	}

	total := 0
	count := 0
	for _, player := range roster {
		if player.Position != "G" { // Exclude goalies from overall calculation
			total += int(player.Overall)
			count++
		}
	}

	if count == 0 {
		return 50.0
	}

	return float64(total) / float64(count)
}

// determineOptimalIntensity calculates ideal intensity based on roster factors
func determineOptimalIntensity(archCounts map[string]map[string]int, avgOverall float64) uint8 {
	baseIntensity := 5 // Start at medium intensity

	// Higher overall = can handle higher intensity
	if avgOverall >= 80 {
		baseIntensity = 7
	} else if avgOverall >= 70 {
		baseIntensity = 6
	} else if avgOverall <= 55 {
		baseIntensity = 4
	} else if avgOverall <= 50 {
		baseIntensity = 3
	}

	// Check archetype diversity - more diverse = lower intensity to avoid penalties
	totalPlayers := 0
	uniqueArchetypes := 0

	for _, posArchs := range archCounts {
		for _, count := range posArchs {
			if count > 0 {
				totalPlayers += count
				uniqueArchetypes++
			}
		}
	}

	if totalPlayers > 0 {
		diversity := float64(uniqueArchetypes) / float64(totalPlayers)
		if diversity > 0.8 { // Very diverse roster
			baseIntensity-- // Lower intensity for mixed archetypes
		} else if diversity < 0.4 { // Specialized roster
			baseIntensity++ // Can handle higher intensity
		}
	}

	// Clamp to valid range
	if baseIntensity < 1 {
		baseIntensity = 1
	} else if baseIntensity > 10 {
		baseIntensity = 10
	}

	return uint8(baseIntensity)
}

// findBestOffensiveSystem tests all offensive systems and returns the best match
func findBestOffensiveSystem(archCounts map[string]map[string]int, intensity uint8) (structs.OffensiveSystemType, int) {
	systems := []structs.OffensiveSystemType{
		structs.Offensive122Forecheck,
		structs.Offensive212Forecheck,
		structs.Offensive113Forecheck,
		structs.OffensiveCycleGame,
		structs.OffensiveQuickTransition,
		structs.OffensiveUmbrella,
		structs.OffensiveEastWestMotion,
		structs.OffensiveCrashNet,
	}

	bestSystem := structs.Offensive122Forecheck
	bestScore := -1000

	for _, system := range systems {
		score := calculateSystemCompatibilityScore(system, 0, archCounts, intensity, true)
		if score > bestScore {
			bestScore = score
			bestSystem = system
		}
	}

	return bestSystem, bestScore
}

// findBestDefensiveSystem tests all defensive systems and returns the best match
func findBestDefensiveSystem(archCounts map[string]map[string]int, intensity uint8) (structs.DefensiveSystemType, int) {
	systems := []structs.DefensiveSystemType{
		structs.DefensiveBalanced,
		structs.DefensiveManToMan,
		structs.DefensiveZone,
		structs.DefensiveNeutralTrap,
		structs.DefensiveLeftWingLock,
		structs.DefensiveAggressiveForecheck,
		structs.DefensiveCollapsing,
		structs.DefensiveBox,
	}

	bestSystem := structs.DefensiveBalanced
	bestScore := -1000

	for _, system := range systems {
		score := calculateSystemCompatibilityScore(0, system, archCounts, intensity, false)
		if score > bestScore {
			bestScore = score
			bestSystem = system
		}
	}

	return bestSystem, bestScore
}

// calculateSystemCompatibilityScore evaluates how well a system matches the roster
func calculateSystemCompatibilityScore(offSystem structs.OffensiveSystemType, defSystem structs.DefensiveSystemType, archCounts map[string]map[string]int, intensity uint8, isOffensive bool) int {
	var mods structs.SystemModifiers

	if isOffensive {
		mods = structs.GetOffensiveSystemModifiers(offSystem, intensity)
	} else {
		mods = structs.GetDefensiveSystemModifiers(defSystem, intensity)
	}

	totalScore := 0

	// Calculate weighted score based on archetype fit and roster composition
	for archetype, weight := range mods.ArchetypeWeights {
		// Count players with this archetype across relevant positions
		playerCount := 0

		// Forwards and centers
		if fCount, exists := archCounts["F"][archetype]; exists {
			playerCount += fCount
		}

		// Defensemen
		if dCount, exists := archCounts["D"][archetype]; exists {
			playerCount += dCount
		}

		// Score = (archetype weight) * (number of players) * (intensity factor)
		score := int(weight) * playerCount
		totalScore += score
	}

	return totalScore
}

// adjustIntensityForSystemCompatibility fine-tunes intensity based on system synergy
func adjustIntensityForSystemCompatibility(baseIntensity uint8, offensiveScore, defensiveScore int) uint8 {
	// If both systems have high compatibility, can increase intensity
	if offensiveScore > 50 && defensiveScore > 50 {
		if baseIntensity < 10 {
			baseIntensity++
		}
	}

	// If either system has poor compatibility, decrease intensity
	if offensiveScore < -20 || defensiveScore < -20 {
		if baseIntensity > 1 {
			baseIntensity--
		}
	}

	return baseIntensity
}

// getDetailedSystemAnalysis provides comprehensive analysis for a team's optimal systems
func getDetailedSystemAnalysis(roster []structs.BasePlayer) string {
	archCounts := analyzeArchetypeComposition(roster)
	avgOverall := calculateAverageOverall(roster)
	intensity := determineOptimalIntensity(archCounts, avgOverall)

	analysis := "Roster Analysis:\n"
	analysis += fmt.Sprintf("  Average Overall: %.1f\n", avgOverall)
	analysis += fmt.Sprintf("  Recommended Intensity: %d\n", intensity)
	analysis += "  Archetype Composition:\n"

	for pos, archs := range archCounts {
		analysis += fmt.Sprintf("    %s: ", pos)
		for arch, count := range archs {
			if count > 0 {
				analysis += fmt.Sprintf("%s(%d) ", arch, count)
			}
		}
		analysis += "\n"
	}

	// Test all system combinations
	bestOffensive, bestOffensiveScore := findBestOffensiveSystem(archCounts, intensity)
	bestDefensive, bestDefensiveScore := findBestDefensiveSystem(archCounts, intensity)

	analysis += fmt.Sprintf("  Best Offensive System: %s (Score: %d)\n",
		structs.GetOffensiveSystemName(bestOffensive), bestOffensiveScore)
	analysis += fmt.Sprintf("  Best Defensive System: %s (Score: %d)\n",
		structs.GetDefensiveSystemName(bestDefensive), bestDefensiveScore)

	return analysis
}

func SelectOffensiveAndDefensiveSystemsForAllTeams_Offseason() {
	db := dbprovider.GetInstance().GetDB()

	// Get all teams that need system selection
	collegeTeams := repository.FindAllCollegeTeams(repository.TeamClauses{})
	proTeams := repository.FindAllProTeams(repository.TeamClauses{})

	// Process college teams
	collegeGameplans := repository.FindCollegeGameplanRecords()
	collegeGameplanMap := MakeCollegeGameplanMap(collegeGameplans)

	for _, team := range collegeTeams {
		gameplan := collegeGameplanMap[team.ID]
		if gameplan.ID == 0 || !gameplan.IsAI {
			continue // Skip user-controlled teams
		}

		fmt.Printf("Selecting systems for college team: %s\n", team.Abbreviation)

		// Get roster and analyze
		teamID := strconv.Itoa(int(team.ID))
		roster := GetCollegePlayersByTeamID(teamID)

		// Select optimal systems
		bestOffensive, bestDefensive, intensity := selectOptimalSystemsForRoster(roster)

		// Update gameplan
		gameplan.OffensiveSystem = uint8(bestOffensive)
		gameplan.DefensiveSystem = uint8(bestDefensive)
		gameplan.OffensiveIntensity = intensity
		gameplan.DefensiveIntensity = intensity

		repository.SaveCollegeGameplanRecord(gameplan, db)

		fmt.Printf("%s  Selected: %s (O) + %s (D) at intensity %d\n",
			team.TeamName,
			structs.GetOffensiveSystemName(structs.OffensiveSystemType(bestOffensive)),
			structs.GetDefensiveSystemName(structs.DefensiveSystemType(bestDefensive)),
			intensity)
	}

	// Process professional teams
	proGameplans := repository.FindProfessionalGameplanRecords()
	proGameplanMap := MakeProGameplanMap(proGameplans)

	for _, team := range proTeams {
		gameplan := proGameplanMap[team.ID]
		if gameplan.ID == 0 || !gameplan.IsAI {
			continue // Skip user-controlled teams
		}

		fmt.Printf("Selecting systems for pro team: %s\n", team.Abbreviation)

		// Get roster and analyze
		teamID := strconv.Itoa(int(team.ID))
		roster := GetProPlayersByTeamID(teamID)

		// Select optimal systems
		bestOffensive, bestDefensive, intensity := selectOptimalSystemsForProRoster(roster)

		// Update gameplan
		gameplan.OffensiveSystem = uint8(bestOffensive)
		gameplan.DefensiveSystem = uint8(bestDefensive)
		gameplan.OffensiveIntensity = intensity
		gameplan.DefensiveIntensity = intensity

		repository.SaveProfessionalGameplanRecord(gameplan, db)

		fmt.Printf("  Selected: %s (O) + %s (D) at intensity %d\n",
			structs.GetOffensiveSystemName(structs.OffensiveSystemType(bestOffensive)),
			structs.GetDefensiveSystemName(structs.DefensiveSystemType(bestDefensive)),
			intensity)
	}
}

// GetSystemWeightedSortExpression combines traditional sorting with system compatibility weighting
func GetSystemWeightedSortExpression(pref1, pref2, pref3 uint8, i structs.BasePlayer, j structs.BasePlayer, offensiveSystem, defensiveSystem, offensiveIntensity, defensiveIntensity uint8, isGoalie bool) bool {
	// Get base sorting scores
	iBaseScore := calculatePlayerSortScore(pref1, pref2, pref3, i)
	jBaseScore := calculatePlayerSortScore(pref1, pref2, pref3, j)

	// Calculate system compatibility scores (10% weight)
	iSystemScore := calculatePlayerSystemCompatibility(i, offensiveSystem, defensiveSystem, offensiveIntensity, defensiveIntensity)
	jSystemScore := calculatePlayerSystemCompatibility(j, offensiveSystem, defensiveSystem, offensiveIntensity, defensiveIntensity)

	// Apply 10% system weighting to base scores
	iTotalScore := float64(iBaseScore)*0.9 + float64(iSystemScore)*0.1
	jTotalScore := float64(jBaseScore)*0.9 + float64(jSystemScore)*0.1

	return iTotalScore > jTotalScore
}

// calculatePlayerSortScore returns a numeric score based on sorting preferences
func calculatePlayerSortScore(pref1, pref2, pref3 uint8, player structs.BasePlayer) int {
	if pref1 == 1 {
		return int(player.Overall)
	}

	score1 := getPlayerAttributeByPreference(pref1, player)
	score2 := getPlayerAttributeByPreference(pref2, player)
	score3 := getPlayerAttributeByPreference(pref3, player)

	return int(score1 + score2 + score3)
}

// getPlayerAttributeByPreference returns the attribute value based on preference number
func getPlayerAttributeByPreference(pref uint8, player structs.BasePlayer) uint8 {
	switch pref {
	case 2:
		return player.CloseShotAccuracy
	case 3:
		return player.LongShotAccuracy
	case 4:
		return player.Agility
	case 5:
		return player.PuckHandling
	case 6:
		return player.Strength
	case 7:
		return player.BodyChecking
	case 8:
		return player.StickChecking
	case 9:
		return player.Faceoffs
	case 10:
		return player.Passing
	default:
		return player.Overall
	}
}

// validateCollegeRoster checks if a team has the minimum required players by position
// Required: 4 centers, 8 forwards, 6 defenders, 2 goalies
// Returns true if roster is complete, false if missing players
func validateCollegeRoster(team structs.CollegeTeam, roster []structs.CollegePlayer) bool {
	// Count current players by position
	centers := 0
	forwards := 0
	defenders := 0
	goalies := 0

	for _, p := range roster {
		if p.IsRedshirt || p.IsInjured {
			continue // Skip unavailable players
		}
		switch p.Position {
		case "C":
			centers++
		case "F":
			forwards++
		case "D":
			defenders++
		case "G":
			goalies++
		}
	}

	fmt.Printf("Team %s current roster: C:%d F:%d D:%d G:%d\n", team.Abbreviation, centers, forwards, defenders, goalies)

	// Check if roster is complete
	isComplete := centers >= 4 && forwards >= 8 && defenders >= 6 && goalies >= 2

	if isComplete {
		fmt.Printf("Team %s has complete roster ✓\n", team.Abbreviation)
		return true
	}

	// Report what's missing without generating players
	fmt.Printf("⚠️  Team %s INCOMPLETE ROSTER - Missing: ", team.Abbreviation)
	missingPieces := []string{}
	if centers < 4 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d centers", 4-centers))
	}
	if forwards < 8 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d forwards", 8-forwards))
	}
	if defenders < 6 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d defenders", 6-defenders))
	}
	if goalies < 2 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d goalies", 2-goalies))
	}
	fmt.Printf("%s\n", strings.Join(missingPieces, ", "))
	fmt.Printf("   → Use recruiting/transfer portal to complete roster or generate walk-ons\n")
	notificationMessage := fmt.Sprintf("⚠️  Team %s INCOMPLETE ROSTER - Missing: ", team.Abbreviation)
	notificationMessage += strings.Join(missingPieces, ", ")
	notificationMessage += " → Use recruiting/transfer portal to complete roster or generate walk-ons"
	CreateNotification("CHL", notificationMessage, "Roster", team.ID)
	return false
} // validateProRoster checks if a pro team has the minimum required players by position
// Required: 4 centers, 8 forwards, 6 defenders, 2 goalies
// Returns true if roster is complete, false if missing players
func validateProRoster(team structs.ProfessionalTeam, roster []structs.ProfessionalPlayer) bool {
	// Count current players by position
	centers := 0
	forwards := 0
	defenders := 0
	goalies := 0

	for _, p := range roster {
		if p.IsInjured {
			continue // Skip unavailable players
		}
		switch p.Position {
		case "C":
			centers++
		case "F":
			forwards++
		case "D":
			defenders++
		case "G":
			goalies++
		}
	}

	fmt.Printf("Team %s current roster: C:%d F:%d D:%d G:%d\n", team.Abbreviation, centers, forwards, defenders, goalies)

	// Check if roster is complete
	isComplete := centers >= 4 && forwards >= 8 && defenders >= 6 && goalies >= 2

	if isComplete {
		fmt.Printf("Team %s has complete roster ✓\n", team.Abbreviation)
		return true
	}

	// Report what's missing without generating players
	fmt.Printf("⚠️  Team %s INCOMPLETE ROSTER - Missing: ", team.Abbreviation)
	missingPieces := []string{}
	if centers < 4 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d centers", 4-centers))
	}
	if forwards < 8 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d forwards", 8-forwards))
	}
	if defenders < 6 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d defenders", 6-defenders))
	}
	if goalies < 2 {
		missingPieces = append(missingPieces, fmt.Sprintf("%d goalies", 2-goalies))
	}
	fmt.Printf("%s\n", strings.Join(missingPieces, ", "))
	fmt.Printf("   → Use free agency/trades to complete roster or generate walk-ons\n")
	notificationMessage := fmt.Sprintf("⚠️  Team %s INCOMPLETE ROSTER - Missing: ", team.Abbreviation)
	notificationMessage += strings.Join(missingPieces, ", ")
	notificationMessage += " → Use recruiting/transfer portal to complete roster."
	CreateNotification("PHL", notificationMessage, "Roster", team.ID)
	return false
}

// validateAndFixLineups ensures a team has complete lineups (4 forward lines, 3 defender lines, 2 goalie lines)
func validateAndFixCollegeLineups(team structs.CollegeTeam, roster []structs.CollegePlayer, db *gorm.DB) bool {
	teamID := strconv.Itoa(int(team.ID))
	existingLineups := GetCollegeLineupsByTeamID(teamID)

	// Count existing lineups by type
	forwardLines := 0
	defenderLines := 0
	goalieLines := 0

	for _, lineup := range existingLineups {
		switch lineup.LineType {
		case 1: // Forward lines
			forwardLines++
		case 2: // Defender lines
			defenderLines++
		case 3: // Goalie lines
			goalieLines++
		}
	}

	fmt.Printf("Team %s existing lineups: Forward:%d/4 Defender:%d/3 Goalie:%d/2\n",
		team.Abbreviation, forwardLines, defenderLines, goalieLines)

	// Check if lineups are complete
	lineupsComplete := forwardLines >= 4 && defenderLines >= 3 && goalieLines >= 2

	if lineupsComplete {
		fmt.Printf("Team %s has complete lineups ✓\n", team.Abbreviation)
		return true
	}

	// Regenerating complete lineups (existing ones will be overwritten)
	fmt.Printf("Regenerating complete lineups for team %s\n", team.Abbreviation)

	// Sort players by position for lineup creation
	centers := []structs.CollegePlayer{}
	forwards := []structs.CollegePlayer{}
	defenders := []structs.CollegePlayer{}
	goalies := []structs.CollegePlayer{}

	for _, p := range roster {
		if p.IsRedshirt || p.IsInjured {
			continue
		}
		switch p.Position {
		case "C":
			centers = append(centers, p)
		case "F":
			forwards = append(forwards, p)
		case "D":
			defenders = append(defenders, p)
		case "G":
			goalies = append(goalies, p)
		}
	}

	// Sort by overall rating (best players first)
	sort.Slice(centers, func(i, j int) bool { return centers[i].Overall > centers[j].Overall })
	sort.Slice(forwards, func(i, j int) bool { return forwards[i].Overall > forwards[j].Overall })
	sort.Slice(defenders, func(i, j int) bool { return defenders[i].Overall > defenders[j].Overall })
	sort.Slice(goalies, func(i, j int) bool { return goalies[i].Overall > goalies[j].Overall })

	// Create 4 forward lines (C + 2F each)
	for line := 0; line < 4 && line < len(centers); line++ {
		f1Idx := min(line*2, len(forwards)-1)
		f2Idx := min(line*2+1, len(forwards)-1)

		if f1Idx < len(forwards) && f2Idx < len(forwards) {
			lineup := createCollegeForwardLineup(team.ID, line+1, centers[line], forwards[f1Idx], forwards[f2Idx])
			repository.SaveCollegeLineupRecord(lineup, db)
		}
	}

	// Create 3 defender lines (2D each)
	for line := 0; line < 3 && line*2+1 < len(defenders); line++ {
		d1Idx := line * 2
		d2Idx := line*2 + 1

		if d1Idx < len(defenders) && d2Idx < len(defenders) {
			lineup := createCollegeDefenderLineup(team.ID, line+1, defenders[d1Idx], defenders[d2Idx])
			repository.SaveCollegeLineupRecord(lineup, db)
		}
	}

	// Create 2 goalie lines
	for line := 0; line < 2 && line < len(goalies); line++ {
		lineup := createCollegeGoalieLineup(team.ID, line+1, goalies[line])
		repository.SaveCollegeLineupRecord(lineup, db)
	}

	return true
}

// Helper functions to create lineup structs
func createCollegeForwardLineup(teamID uint, lineNum int, center, f1, f2 structs.CollegePlayer) structs.CollegeLineup {
	return structs.CollegeLineup{
		BaseLineup: structs.BaseLineup{
			TeamID:   teamID,
			Line:     uint8(lineNum),
			LineType: 1, // Forward line
			LineupPlayerIDs: structs.LineupPlayerIDs{
				CenterID:   center.ID,
				Forward1ID: f1.ID,
				Forward2ID: f2.ID,
			},
			Allocations: getDefaultForwardAllocations(),
		},
	}
}

func createCollegeDefenderLineup(teamID uint, lineNum int, d1, d2 structs.CollegePlayer) structs.CollegeLineup {
	return structs.CollegeLineup{
		BaseLineup: structs.BaseLineup{
			TeamID:   teamID,
			Line:     uint8(lineNum),
			LineType: 2, // Defender line
			LineupPlayerIDs: structs.LineupPlayerIDs{
				Defender1ID: d1.ID,
				Defender2ID: d2.ID,
			},
			Allocations: getDefaultDefenderAllocations(),
		},
	}
}

func createCollegeGoalieLineup(teamID uint, lineNum int, goalie structs.CollegePlayer) structs.CollegeLineup {
	return structs.CollegeLineup{
		BaseLineup: structs.BaseLineup{
			TeamID:   teamID,
			Line:     uint8(lineNum),
			LineType: 3, // Goalie line
			LineupPlayerIDs: structs.LineupPlayerIDs{
				GoalieID: goalie.ID,
			},
			Allocations: getDefaultGoalieAllocations(),
		},
	}
}

// Default allocation functions
func getDefaultForwardAllocations() structs.Allocations {
	return structs.Allocations{
		AGZShot:       15,
		AGZPass:       15,
		AGZPassBack:   5,
		AGZAgility:    15,
		AGZStickCheck: 10,
		AGZBodyCheck:  10,
		AZShot:        15,
		AZPass:        15,
		AZLongPass:    10,
		AZAgility:     15,
		AZStickCheck:  10,
		AZBodyCheck:   10,
		NPass:         15,
		NAgility:      15,
		NStickCheck:   10,
		NBodyCheck:    10,
		DZPass:        15,
		DZPassBack:    5,
		DZAgility:     15,
		DZStickCheck:  10,
		DZBodyCheck:   10,
		DGZPass:       15,
		DGZLongPass:   10,
		DGZAgility:    15,
		DGZStickCheck: 10,
		DGZBodyCheck:  10,
	}
}

func getDefaultDefenderAllocations() structs.Allocations {
	return structs.Allocations{
		AGZShot:       5,
		AGZPass:       15,
		AGZPassBack:   10,
		AGZAgility:    15,
		AGZStickCheck: 15,
		AGZBodyCheck:  15,
		AZShot:        10,
		AZPass:        15,
		AZLongPass:    15,
		AZAgility:     10,
		AZStickCheck:  15,
		AZBodyCheck:   15,
		NPass:         15,
		NAgility:      15,
		NStickCheck:   15,
		NBodyCheck:    15,
		DZPass:        15,
		DZPassBack:    10,
		DZAgility:     15,
		DZStickCheck:  15,
		DZBodyCheck:   15,
		DGZPass:       15,
		DGZLongPass:   15,
		DGZAgility:    15,
		DGZStickCheck: 15,
		DGZBodyCheck:  15,
	}
}

func getDefaultGoalieAllocations() structs.Allocations {
	// Goalies use different allocations focused on positioning and saves
	return structs.Allocations{
		AGZShot:       0,
		AGZPass:       5,
		AGZPassBack:   10,
		AGZAgility:    20,
		AGZStickCheck: 5,
		AGZBodyCheck:  0,
		AZShot:        0,
		AZPass:        10,
		AZLongPass:    15,
		AZAgility:     15,
		AZStickCheck:  5,
		AZBodyCheck:   0,
		NPass:         10,
		NAgility:      20,
		NStickCheck:   5,
		NBodyCheck:    0,
		DZPass:        15,
		DZPassBack:    15,
		DZAgility:     20,
		DZStickCheck:  10,
		DZBodyCheck:   5,
		DGZPass:       20,
		DGZLongPass:   20,
		DGZAgility:    25,
		DGZStickCheck: 15,
		DGZBodyCheck:  10,
	}
}

// min helper function for integer comparison
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// calculatePlayerSystemCompatibility evaluates how well a player fits the team's systems
func calculatePlayerSystemCompatibility(player structs.BasePlayer, offensiveSystem, defensiveSystem, offensiveIntensity, defensiveIntensity uint8) int {
	// Get system modifiers for both offensive and defensive systems
	offensiveMods := structs.GetOffensiveSystemModifiers(structs.OffensiveSystemType(offensiveSystem), offensiveIntensity)
	defensiveMods := structs.GetDefensiveSystemModifiers(structs.DefensiveSystemType(defensiveSystem), defensiveIntensity)

	offensiveScore := 0
	defensiveScore := 0

	// Calculate offensive system fit
	if weight, exists := offensiveMods.ArchetypeWeights[player.Archetype]; exists {
		offensiveScore = int(weight)
	}

	// Calculate defensive system fit
	if weight, exists := defensiveMods.ArchetypeWeights[player.Archetype]; exists {
		defensiveScore = int(weight)
	}

	// Combine scores (average of offensive and defensive fit)
	totalScore := (offensiveScore + defensiveScore) / 2

	// Scale the score to be comparable to player attributes (0-100 range)
	// System weights are typically in range of -20 to +20, so we normalize and scale
	normalizedScore := (totalScore + 20) * 2 // Converts -20,+20 range to 0,80 range

	if normalizedScore < 0 {
		normalizedScore = 0
	} else if normalizedScore > 100 {
		normalizedScore = 100
	}

	return normalizedScore
}
