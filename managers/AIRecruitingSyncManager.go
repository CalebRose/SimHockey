package managers

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func FillAIRecruitingBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := repository.FindTeamRecruitingProfiles(true)
	fmt.Println("Loading recruits...")
	UnsignedRecruits := GetAllUnsignedRecruits()
	allCollegePlayersMap := GetAllCollegePlayersMapByTeam()
	recruitProfiles := repository.FindRecruitPlayerProfileRecords("", "", true, true, false)
	recruitProfileMap := MakeRecruitProfileMapByRecruitID(recruitProfiles)
	existingBoardMap := MakeRecruitProfileMapByProfileID(recruitProfiles)
	fmt.Println("Loaded all unsigned recruits.")

	boardCount := 75

	for _, team := range AITeams {
		count := 0
		if !team.IsAI || team.TotalCommitments >= team.RecruitClassSize {
			continue
		}
		fmt.Println("Iterating through " + team.Team + ".")
		existingBoard := existingBoardMap[team.ID]
		teamRecruitProfileMap := MakeRecruitProfileMapByRecruitID(existingBoard)
		teamNeeds := getRecruitingNeeds(allCollegePlayersMap[team.ID])

		// Get Current Count of the existing board
		for _, r := range existingBoard {
			if r.RemovedFromBoard {
				continue
			}

			if r.IsSigned {
				teamNeeds[r.Recruit.Position] -= 1
			}

			count++
		}

		for k := range teamNeeds {
			if teamNeeds[k] > 0 {
				teamNeeds[k] *= 4
			}
		}

		for _, croot := range UnsignedRecruits {
			if count >= boardCount {
				break
			}
			starMin := team.AIStarMin
			starMax := team.AIStarMax
			if ts.Week >= 18 {
				starMin = 1
				starMax = 3
			}
			if (teamNeeds[croot.Position] < 1) ||
				(croot.Stars > starMax) || (croot.Stars < starMin) {
				continue
			}

			crootProfiles := recruitProfileMap[croot.ID]
			teamCount := 0

			for _, crootProfile := range crootProfiles {
				if crootProfile.RemovedFromBoard {
					continue
				}
				teamCount++
			}

			leadingVal := IsAITeamContendingForCroot(crootProfiles)
			if leadingVal > 15 {
				continue
			}

			// Check and see if the croot already exists on the player's board
			crootProfile := teamRecruitProfileMap[croot.ID]
			if crootProfile[0].ProfileID == team.ID || crootProfile[0].ID > 0 || crootProfile[0].RemovedFromBoard || crootProfile[0].IsLocked {
				fmt.Println(croot.FirstName + " " + croot.LastName + " is already on " + team.Team + "'s board.")
				continue
			}

			oddsObject := getRecruitingOdds(ts, croot, team, allCollegePlayersMap[team.ID])

			chance := util.GenerateIntFromRange(1, 100)
			highlyContestedMod := 1
			if oddsObject.IsCloseToHome {
				highlyContestedMod = 3
			} else if oddsObject.IsPipeline {
				highlyContestedMod = 2
			}

			willAddToBoard := isHighlyContestedCroot(highlyContestedMod, teamCount, int(ts.Week))

			addPlayer := chance <= oddsObject.Odds && willAddToBoard
			if !addPlayer {
				continue
			}
			playerProfile := structs.RecruitPlayerProfile{
				RecruitID:          croot.ID,
				ProfileID:          team.ID,
				SeasonID:           ts.SeasonID,
				TotalPoints:        0,
				CurrentWeeksPoints: 0,
				SpendingCount:      0,
				Scholarship:        false,
				ScholarshipRevoked: false,
				IsHomeState:        oddsObject.IsCloseToHome,
				IsPipelineState:    oddsObject.IsPipeline,
				IsSigned:           false,
				IsLocked:           false,
			}

			repository.CreateRecruitProfileRecord(db, playerProfile)
			teamNeeds[croot.Position] -= 1
			recruitProfileMap[croot.ID] = append(recruitProfileMap[croot.ID], playerProfile)
			sort.Slice(recruitProfileMap[croot.ID], func(i, j int) bool {
				return recruitProfileMap[croot.ID][i].TotalPoints > recruitProfileMap[croot.ID][j].TotalPoints
			})
			count++
		}
	}
}

func AllocatePointsToAIBoards() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	ts := GetTimestamp()

	AITeams := repository.FindTeamRecruitingProfiles(true)
	fmt.Println("Loading recruits...")
	allCollegePlayersMap := GetAllCollegePlayersMapByTeam()
	recruitProfiles := repository.FindRecruitPlayerProfileRecords("", "", true, true, false)
	existingBoardMap := MakeRecruitProfileMapByProfileID(recruitProfiles)
	fmt.Println("Loaded all unsigned recruits.")
	for _, team := range AITeams {
		if team.SpentPoints >= team.WeeklyPoints || team.TotalCommitments >= team.RecruitClassSize {
			continue
		}

		teamRecruits := existingBoardMap[team.ID]

		teamNeedsMap := getRecruitingNeeds(allCollegePlayersMap[team.ID])

		// Safety check to make sure teams aren't recruiting too many in one position
		for _, croot := range teamRecruits {
			if croot.IsSigned && uint16(croot.ProfileID) == uint16(team.ID) && ts.Week < 17 {
				teamNeedsMap[croot.Recruit.Position] -= 1
			}
		}

		for _, croot := range teamRecruits {
			pointsRemaining := team.WeeklyPoints - team.SpentPoints
			if team.SpentPoints >= team.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			if croot.IsSigned || croot.CurrentWeeksPoints > 0 || croot.ScholarshipRevoked {
				continue
			}

			removeCrootFromBoard := false
			var num float32 = 0
			recruitID := strconv.Itoa(int(croot.RecruitID))

			if (croot.IsLocked && croot.ProfileID != uint(croot.Recruit.TeamID)) || teamNeedsMap[croot.Recruit.Position] <= 0 || croot.RemovedFromBoard {
				removeCrootFromBoard = true
			}

			if !removeCrootFromBoard {
				profiles := repository.FindRecruitPlayerProfileRecords("", recruitID, false, false, true)

				if croot.PreviousWeekPoints > 0 {
					leadingTeamVal := IsAITeamContendingForCroot(profiles)

					if croot.PreviousWeekPoints+croot.TotalPoints >= leadingTeamVal*0.66 || leadingTeamVal < 15 {
						num = croot.PreviousWeekPoints
						if num > pointsRemaining {
							num = pointsRemaining
						}
					} else {
						removeCrootFromBoard = true
					}
				} else {
					maxChance := 2
					if ts.Week > 3 {
						maxChance = 4
					}
					chance := util.GenerateIntFromRange(1, maxChance)
					if (chance < 2 && ts.Week <= 3) || (chance < 4 && ts.Week > 3) {
						continue
					}

					min := team.AIMinThreshold
					max := team.AIMaxThreshold

					num = float32(util.GenerateIntFromRange(int(min), int(max)))
					if num > pointsRemaining {
						num = pointsRemaining
					}

					leadingTeamVal := IsAITeamContendingForCroot(profiles)

					if float32(num)+croot.TotalPoints < leadingTeamVal*0.66 {
						removeCrootFromBoard = true
					}
					if leadingTeamVal < 15 {
						removeCrootFromBoard = false
					}
				}
			}

			if removeCrootFromBoard || (team.ScholarshipsAvailable == 0 && !croot.Scholarship) {
				if croot.Scholarship {
					croot.ToggleScholarship(false, true)
					team.ReallocateScholarship()
				}
				croot.ToggleRemoveFromBoard()
				fmt.Println("Because " + croot.Recruit.FirstName + " " + croot.Recruit.LastName + " is heavily considering other teams, they are being removed from " + team.Team + "'s Recruiting Board.")
				db.Save(&croot)
				continue
			}

			if ts.Week == 20 {
				num = 2
			}

			croot.AllocateCurrentWeekPoints(num)
			if !croot.Scholarship && team.ScholarshipsAvailable > 0 {
				croot.ToggleScholarship(true, false)
				team.SubtractScholarshipsAvailable()
			}

			team.AIAllocateSpentPoints(num)
			repository.SaveRecruitProfileRecord(db, croot)
			fmt.Println(team.Team + " allocating " + strconv.Itoa(int(num)) + " points to " + croot.Recruit.FirstName + " " + croot.Recruit.LastName)

		}
		// Save Team Profile after iterating through recruits
		fmt.Println("Saved " + team.Team + " Recruiting Board!")
		repository.SaveTeamProfileRecord(db, team)
	}

	ts.ToggleAIrecruitingSync()
	repository.SaveTimestamp(ts, db)
}

func ResetAIBoardsForCompletedTeams() {
	db := dbprovider.GetInstance().GetDB()

	AITeams := repository.FindTeamRecruitingProfiles(true)

	for _, team := range AITeams {
		// If a team already has the maximum allowed for their recruiting class, take all Recruit Profiles for that team where the recruit hasn't signed, and reset their total points.
		// This is so that these unsigned recruits can be recruited for and will allow the AI to put points onto those recruits.

		if team.TotalCommitments >= team.RecruitClassSize {
			teamProfiles := repository.FindRecruitPlayerProfileRecords(strconv.Itoa(int(team.ID)), "", false, false, true)

			for _, profile := range teamProfiles {
				if profile.IsSigned || profile.IsLocked || profile.TotalPoints == 0 {
					continue
				}
				profile.ResetTotalPoints()
				repository.SaveRecruitProfileRecord(db, profile)
				db.Save(&profile)
			}
			team.ResetSpentPoints()
			repository.SaveTeamProfileRecord(db, team)
		}
	}
}

func isHighlyContestedCroot(mod int, teams int, CollegeWeek int) bool {
	if CollegeWeek == 20 && teams > 1 {
		return false
	}
	chance := util.GenerateIntFromRange(1, 5)
	chance += mod

	return chance > teams
}

func getRecruitingNeeds(playerRoster []structs.CollegePlayer) map[string]int {
	needsMap := make(map[string]int)

	// The roster is in dire shape, just open up the team needs
	if len(playerRoster) <= 20 {
		needsMap["C"] = 2
		needsMap["F"] = 3
		needsMap["D"] = 2
		needsMap["G"] = 2
		return needsMap
	}

	for _, player := range playerRoster {
		if player.IsRedshirting {
			continue
		}
		if (player.Year == 4 && !player.IsRedshirt) || (player.Year == 5 && player.IsRedshirt) {
			needsMap[player.Position] += 1
		}
	}

	return needsMap
}

func IsAITeamContendingForCroot(profiles []structs.RecruitPlayerProfile) float32 {
	if len(profiles) == 0 {
		return 0
	}
	var leadingVal float32 = 0
	for _, profile := range profiles {
		if profile.TotalPoints != 0 && profile.TotalPoints > float32(leadingVal) {
			leadingVal = profile.TotalPoints
		}
	}

	return leadingVal
}

func getRecruitingOdds(ts structs.Timestamp, croot structs.Recruit, team structs.RecruitingTeamProfile, roster []structs.CollegePlayer) structs.RecruitingOdds {
	odds := 5
	if ts.Week > 5 {
		odds = 10
	} else if ts.Week > 14 && odds < 20 {
		odds = 20
	}

	if croot.State == team.State {
		odds = 25
	}

	country := croot.Country
	isCloseToHome := country == util.USA && croot.State == team.State
	isPipeline := false
	pipelineMap := make(map[string]int)
	if !isCloseToHome {
		// Check Pipeline
		for _, p := range roster {
			if p.Country == util.USA || p.Country == util.Canada {
				pipelineMap[p.State] += 1
			} else {
				pipelineMap[p.Country] += 1
			}
		}

		isPipeline = ((country == util.USA || country == util.Canada) && pipelineMap[croot.State] > 6) || pipelineMap[croot.Country] > 6
	}

	if isCloseToHome {
		odds += 25
	} else if isPipeline {
		odds += 15
	}

	return structs.RecruitingOdds{
		Odds:          odds,
		IsCloseToHome: isCloseToHome,
		IsPipeline:    isPipeline,
	}
}
