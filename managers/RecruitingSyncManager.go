package managers

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func SyncCollegeRecruiting() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	if !ts.IsRecruitingLocked {
		ts.ToggleLockRecruiting()
		repository.SaveTimestamp(ts, db)
		return
	}
	// Constants
	var mod1 float32 = 75
	weeksOfRecruiting := 17
	eligibleThresholdPercentage := float32(0.66)
	pointLimit := 20

	// Load Data
	teamProfiles := repository.FindTeamRecruitingProfiles(false)
	allRecruitProfiles := repository.FindRecruitPlayerProfileRecords("", "", false, false, true)
	recruits := repository.FindAllRecruits(false, false, false, false, true, "")
	recruitProfileMap := MakeRecruitProfileMapByRecruitID(allRecruitProfiles)
	teamProfileMap := MakeTeamProfileMap(teamProfiles)
	teamPointsMap := getTeamPointsMap()
	teamMap := GetCollegeTeamMap()

	logs := []structs.NewsLog{}
	pointAllocations := []structs.RecruitPointAllocation{}

	for _, r := range recruits {
		previousRecruitStatus := r.RecruitingStatus
		recruitProfiles := recruitProfileMap[r.ID]

		if len(recruitProfiles) == 0 {
			fmt.Println("Skipping over " + r.FirstName + " " + r.LastName + " because no one is recruiting them.")
			continue
		}

		recruitProfilesWithScholarship := []structs.RecruitPlayerProfile{}
		eligibleTeams := 0
		var totalPointsOnRecruit float32 = 0

		var eligiblePointThreshold float32 = 0

		var signThreshold float32

		pointsPlaced := false
		spendingCountAdjusted := false

		allocations := allocatePointsToRecruit(r, &recruitProfiles, float32(pointLimit), &spendingCountAdjusted, &pointsPlaced, ts, &teamPointsMap)

		pointAllocations = append(pointAllocations, allocations...)

		if !pointsPlaced && !spendingCountAdjusted {
			fmt.Println("Skipping over " + r.FirstName + " " + r.LastName)
			continue
		}
		sort.Sort(structs.ByPoints(recruitProfiles))

		for i := 0; i < len(recruitProfiles) && pointsPlaced; i++ {
			recruitTeamProfile := teamProfileMap[recruitProfiles[i].ProfileID]
			if recruitTeamProfile.TotalCommitments >= recruitTeamProfile.RecruitClassSize || (recruitProfiles[i].RemovedFromBoard && recruitProfiles[i].ScholarshipRevoked) {
				continue
			}
			if eligiblePointThreshold == 0 && recruitProfiles[i].Scholarship {
				eligiblePointThreshold = float32(recruitProfiles[i].TotalPoints) * eligibleThresholdPercentage
			}

			if recruitProfiles[i].Scholarship && recruitProfiles[i].TotalPoints >= eligiblePointThreshold {
				totalPointsOnRecruit += recruitProfiles[i].TotalPoints
				eligibleTeams += 1
				recruitProfilesWithScholarship = append(recruitProfilesWithScholarship, recruitProfiles[i])
			}
		}

		firstMod := float32(uint(mod1) - ts.Week)
		secondMod := float32(eligibleTeams) / float32(r.RecruitingModifier)
		thirdMod := math.Log10(float64(uint(weeksOfRecruiting) - ts.Week))
		signThreshold = firstMod * secondMod * float32(thirdMod)
		passedSigningThreshold := totalPointsOnRecruit > signThreshold && eligibleTeams > 0 && pointsPlaced
		r.ApplySigningStatus(totalPointsOnRecruit, signThreshold, passedSigningThreshold)

		if passedSigningThreshold {
			var winningTeamID uint = 0
			var odds float32 = 0

			for winningTeamID == 0 {
				percentageOdds := rand.Float32() * (totalPointsOnRecruit)
				var currentProbability float32 = 0

				for i := 0; i < len(recruitProfilesWithScholarship); i++ {
					// If a team has no available scholarships or if a team has 25 commitments, continue
					currentProbability += recruitProfilesWithScholarship[i].TotalPoints
					if percentageOdds <= currentProbability {
						// WINNING TEAM
						winningTeamID = recruitProfilesWithScholarship[i].ProfileID
						odds = float32(recruitProfilesWithScholarship[i].TotalPoints) / float32(totalPointsOnRecruit) * 100
						break
					}
				}

				if winningTeamID > 0 {
					recruitTeamProfile := teamProfileMap[winningTeamID]
					team := teamMap[winningTeamID]
					if recruitTeamProfile.TotalCommitments < recruitTeamProfile.RecruitClassSize {
						recruitTeamProfile.IncreaseCommitCount()
						recruitTeamProfile.AddStarPlayer(r.Stars)
						teamAbbreviation := team.Abbreviation
						r.AssignCollege(teamAbbreviation)
						message := r.FirstName + " " + r.LastName + ", " + strconv.Itoa(int(r.Stars)) + " star " + r.Position + " from " + r.State + ", " + r.Country + " has signed with " + team.TeamName + " with " + strconv.Itoa(int(odds)) + " percent odds."
						news := CreateNewsLogObject("CHL", message, "Commitment", int(winningTeamID), ts, false)
						logs = append(logs, news)

						for i := 0; i < len(recruitProfiles); i++ {
							if recruitProfiles[i].ProfileID == winningTeamID {
								recruitProfiles[i].SignPlayer()
							} else {
								recruitProfiles[i].LockPlayer()
								if recruitProfiles[i].Scholarship {
									tp := teamProfileMap[recruitProfiles[i].ProfileID]
									t := teamMap[recruitProfiles[i].ProfileID]
									tp.ReallocateScholarship()
									repository.SaveTeamProfileRecord(db, *tp)
									fmt.Println("Reallocated Scholarship to " + t.TeamName)
								}
							}
						}
					} else {
						recruitProfilesWithScholarship = filterOutRecruitingProfile(recruitProfilesWithScholarship, int(winningTeamID))
						// If there are no longer any teams contending due to reaching the max class size, break the loop
						winningTeamID = 0
						if len(recruitProfilesWithScholarship) == 0 {
							break
						}

						totalPointsOnRecruit = 0
						for _, rp := range recruitProfilesWithScholarship {
							totalPointsOnRecruit += rp.TotalPoints
						}
					}
				}
			}
			r.UpdateTeamID(winningTeamID)
		}

		for _, rp := range recruitProfiles {
			team := teamMap[rp.ID]
			repository.SaveRecruitProfileRecord(db, rp)
			fmt.Println("Save recruit profile from " + team.TeamName + " towards " + r.FirstName + " " + r.LastName)
		}

		if r.TeamID > 0 || previousRecruitStatus != r.RecruitingStatus {
			repository.SaveCollegeHockeyRecruitRecord(r, db)
		}
	}

	repository.CreatePointAllocationsRecordsBatch(db, pointAllocations, 100)
	repository.CreateNewsLogRecordsBatch(db, logs, 50)

	updateTeamRankings(teamProfiles, teamProfileMap, teamPointsMap, db, int(ts.Week))

	if ts.IsRecruitingLocked {
		ts.ToggleLockRecruiting()
	}

	repository.SaveTimestamp(ts, db)
}

func allocatePointsToRecruit(recruit structs.Recruit, recruitProfiles *[]structs.RecruitPlayerProfile, pointLimit float32, spendingCountAdjusted *bool, pointsPlaced *bool, timestamp structs.Timestamp, recruitProfilePointsMap *map[uint]float32) []structs.RecruitPointAllocation {
	allocations := []structs.RecruitPointAllocation{}
	for i := range *recruitProfiles {
		if (*recruitProfiles)[i].CurrentWeeksPoints == 0 {
			if (*recruitProfiles)[i].SpendingCount > 0 {
				(*recruitProfiles)[i].ResetSpendingCount()
				*spendingCountAdjusted = true
				fmt.Println("Resetting spending count for " + recruit.FirstName + " " + recruit.LastName + " for Team " + strconv.Itoa(int((*recruitProfiles)[i].ProfileID)))
			}
		} else {
			*pointsPlaced = true
		}

		var curr float32 = 0

		var modifier float32 = 1

		if (*recruitProfiles)[i].IsHomeState {
			modifier += 0.25
		}
		if (*recruitProfiles)[i].IsPipelineState {
			modifier += 0.15
		}

		curr = float32((*recruitProfiles)[i].CurrentWeeksPoints) * modifier

		if (*recruitProfiles)[i].CurrentWeeksPoints < 0 || (*recruitProfiles)[i].CurrentWeeksPoints > pointLimit {
			curr = 0
		}

		rpa := structs.RecruitPointAllocation{
			RecruitID:          (*recruitProfiles)[i].RecruitID,
			TeamProfileID:      (*recruitProfiles)[i].ProfileID,
			RecruitProfileID:   (*recruitProfiles)[i].ID,
			WeekID:             timestamp.WeekID,
			IsHomeStateApplied: (*recruitProfiles)[i].IsHomeState,
			IsPipelineApplied:  (*recruitProfiles)[i].IsPipelineState,
			Points:             (*recruitProfiles)[i].CurrentWeeksPoints,
			ModAffectedPoints:  curr,
		}

		(*recruitProfiles)[i].AddCurrentWeekPointsToTotal(curr)
		(*recruitProfilePointsMap)[(*recruitProfiles)[i].ProfileID] += float32((*recruitProfiles)[i].CurrentWeeksPoints)

		allocations = append(allocations, rpa)
	}

	return allocations
}

func updateTeamRankings(teamRecruitingProfiles []structs.RecruitingTeamProfile, teamMap map[uint]*structs.RecruitingTeamProfile, recruitProfilePointsMap map[uint]float32, db *gorm.DB, week int) {
	// Update rank system for all teams
	var maxESPNScore float32 = 0
	var minESPNScore float32 = 10000
	var maxRivalsScore float32 = 0
	var minRivalsScore float32 = 10000
	var max247Score float32 = 0
	var min247Score float32 = 10000

	for i := 0; i < len(teamRecruitingProfiles); i++ {
		tp := teamMap[teamRecruitingProfiles[i].ID]

		signedRecruits := repository.FindAllRecruits(false, true, true, true, true, strconv.Itoa(int(teamRecruitingProfiles[i].TeamID)))

		teamRecruitingProfiles[i].UpdateTotalSignedRecruits(uint8(len(signedRecruits)))

		team247Rank := get247TeamRanking(signedRecruits)
		teamESPNRank := getESPNTeamRanking(signedRecruits)
		teamRivalsRank := getRivalsTeamRanking(signedRecruits)
		if teamESPNRank > maxESPNScore {
			maxESPNScore = teamESPNRank
		}
		if teamESPNRank < minESPNScore {
			minESPNScore = teamESPNRank
		}
		if teamRivalsRank > maxRivalsScore {
			maxRivalsScore = teamRivalsRank
		}
		if teamRivalsRank < minRivalsScore {
			minRivalsScore = teamRivalsRank
		}
		if team247Rank > max247Score {
			max247Score = team247Rank
		}
		if team247Rank < min247Score {
			min247Score = team247Rank
		}

		tp.Assign247Rank(team247Rank)
		tp.AssignESPNRank(teamESPNRank)
		tp.AssignRivalsRank(teamRivalsRank)
	}

	espnDivisor := (maxESPNScore - minESPNScore)
	divisor247 := (max247Score - min247Score)
	rivalsDivisor := (maxRivalsScore - minRivalsScore)
	for _, rp := range teamRecruitingProfiles {
		tp := teamMap[rp.ID]
		if recruitProfilePointsMap[rp.ID] > float32(rp.WeeklyPoints) {
			tp.ApplyCaughtCheating()
		}

		var avg float32 = 0
		if espnDivisor > 0 && divisor247 > 0 && rivalsDivisor > 0 {
			distributionESPN := (tp.ESPNScore - minESPNScore) / espnDivisor
			distribution247 := (tp.Rank247Score - min247Score) / divisor247
			distributionRivals := (tp.RivalsScore - minRivalsScore) / rivalsDivisor

			avg = (distributionESPN + distribution247 + distributionRivals)

			tp.AssignCompositeRank(avg)
		}
		tp.ResetSpentPoints()
		tp.ResetScoutingPoints(week)

		// Save TEAM Recruiting Profile
		repository.SaveTeamProfileRecord(db, *tp)
		fmt.Println("Saved Rank Scores for Team " + rp.Team)
	}
}

func getTeamPointsMap() map[uint]float32 {
	return map[uint]float32{
		1:  0,
		2:  0,
		3:  0,
		4:  0,
		5:  0,
		6:  0,
		7:  0,
		8:  0,
		9:  0,
		10: 0,
		11: 0,
		12: 0,
		13: 0,
		14: 0,
		15: 0,
		16: 0,
		17: 0,
		18: 0,
		19: 0,
		20: 0,
		21: 0,
		22: 0,
		23: 0,
		24: 0,
		25: 0,
		26: 0,
		27: 0,
		28: 0,
		29: 0,
		30: 0,
		31: 0,
		32: 0,
		33: 0,
		34: 0,
		35: 0,
		36: 0,
		37: 0,
		38: 0,
		39: 0,
		40: 0,
		41: 0,
		42: 0,
		43: 0,
		44: 0,
		45: 0,
		46: 0,
		47: 0,
		48: 0,
		49: 0,
		50: 0,
		51: 0,
		52: 0,
		53: 0,
		54: 0,
		55: 0,
		56: 0,
		57: 0,
		58: 0,
		59: 0,
		60: 0,
		61: 0,
		62: 0,
		63: 0,
		64: 0,
		65: 0,
		66: 0,
	}
}

func filterOutRecruitingProfile(profiles []structs.RecruitPlayerProfile, ID int) []structs.RecruitPlayerProfile {
	var rp []structs.RecruitPlayerProfile

	for _, profile := range profiles {
		if int(profile.ID) != ID {
			rp = append(rp, profile)
		}
	}

	return rp
}

func get247TeamRanking(signedCroots []structs.Recruit) float32 {
	stddev := 10

	var Rank247 float32 = 0

	for idx, croot := range signedCroots {

		rank := float64((idx - 1) / stddev)

		expo := (-0.5 * (math.Pow(rank, 2)))

		weightedScore := (float64(croot.Rank247) - 20) * math.Pow(math.E, expo)

		Rank247 += float32(weightedScore)
	}

	return Rank247
}

func getESPNTeamRanking(signedCroots []structs.Recruit) float32 {

	var espnRank float32 = 0

	for _, croot := range signedCroots {
		espnRank += croot.ESPNRank
	}

	return espnRank
}

func getRivalsTeamRanking(signedCroots []structs.Recruit) float32 {

	var rivalsRank float32 = 0

	for _, croot := range signedCroots {
		rivalsRank += croot.RivalsRank
	}

	return rivalsRank
}
