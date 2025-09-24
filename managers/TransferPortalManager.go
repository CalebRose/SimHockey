package managers

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func GetCollegePromiseByCollegePlayerID(id, teamID string) structs.CollegePromise {
	return repository.FindCollegePromiseRecord(repository.TransferPortalQuery{
		CollegePlayerID: id,
		TeamID:          teamID,
	})
}

func GetAllCollegePromises() []structs.CollegePromise {
	return repository.FindCollegePromiseRecords(repository.TransferPortalQuery{IsActive: "Y"})
}

func GetCollegePromiseByID(id string) structs.CollegePromise {
	return repository.FindCollegePromiseRecord(repository.TransferPortalQuery{ID: id})
}

func ProcessTransferIntention(w http.ResponseWriter) {
	db := dbprovider.GetInstance().GetDB()
	// w.Header().Set("Content-Disposition", "attachment;filename=transferStats.csv")
	// w.Header().Set("Transfer-Encoding", "chunked")
	// Initialize writer
	// writer := csv.NewWriter(w)
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	allCollegePlayers := GetAllCollegePlayers()
	collegePlayerGameStats := repository.FindCollegePlayerGameStatsRecords(seasonID, "", "", "")
	gameStatMap := MakeCollegePlayerGameStatsMap(collegePlayerGameStats)
	fullRosterMap := MakeCollegePlayerMapByTeamID(allCollegePlayers)
	teamProfiles := repository.FindTeamRecruitingProfiles(false)
	teamProfileMap := MakeTeamProfileMap(teamProfiles)
	transferCount := 0
	freshmanCount := 0
	redshirtFreshmanCount := 0
	sophomoreCount := 0
	redshirtSophomoreCount := 0
	juniorCount := 0
	redshirtJuniorCount := 0
	seniorCount := 0
	redshirtSeniorCount := 0
	lowCount := 0
	mediumCount := 0
	highCount := 0
	bigDrop := -25.0
	smallDrop := -10.0
	// tinyDrop := -5.0
	// tinyGain := 5.0
	smallGain := 10.0
	mediumGain := 15.0
	mediumDrop := -15.0
	bigGain := 25.0

	// HeaderRow := []string{
	// 	"Team", "First Name", "Last Name", "Stars",
	// 	"Archetype", "Position", "Year", "Age", "Redshirt Status",
	// 	"Overall", "Transfer Bias", "Transfer Status", "Transfer Weight", "Dice Roll",
	// 	"Age Mod", "Snap Mod", "star Mod", "DC Comp Mod", "Scheme Mod", "FCS Mod",
	// }

	// err := writer.Write(HeaderRow)
	// if err != nil {
	// 	log.Fatal("Cannot write header row", err)
	// }

	for _, p := range allCollegePlayers {
		// Do not include redshirts and all graduating players
		if p.IsRedshirting || p.TeamID > 194 || p.TeamID == 0 {
			continue
		}
		// Weight will be the initial barrier required for a player to consider transferring.
		// The lower the number gets, the more likely the player will transfer
		transferWeight := 0.0

		// Modifiers on reasons why they would transfer
		snapMod := 0.0
		ageMod := 1.125
		starMod := 0.0
		depthChartCompetitionMod := 0.0
		schemeMod := 0.0
		// closeToHomeMod := 0.0

		// Check Snaps
		seasonStats := gameStatMap[p.ID]
		minutes := 0
		for _, s := range seasonStats {
			minutes += int(s.TimeOnIce)
		}
		minutesAccurate := minutes / 60
		regularGamesPerSeason := 34
		minutesPerGame := minutesAccurate / regularGamesPerSeason

		if minutesPerGame > 15 {
			snapMod = bigDrop
		} else if minutesPerGame > 12 {
			snapMod = smallDrop
		} else if minutesPerGame > 10 {
			snapMod = smallGain
		} else if minutesPerGame > 5 {
			snapMod = mediumGain
		} else {
			snapMod = bigGain
		}

		// Check Age
		// The more experienced the player is in the league,
		// the more likely they will transfer.
		/// Have this be a multiplicative factor to odds
		if p.Year == 1 {
			ageMod = .01
		} else if p.Year == 2 && p.IsRedshirt {
			ageMod = .1
		} else if p.Year == 2 && !p.IsRedshirt {
			ageMod = .4
		} else if p.Year == 3 && p.IsRedshirt {
			ageMod = .7
		} else if p.Year == 3 && !p.IsRedshirt {
			ageMod = 1
		} else if p.Year == 4 {
			ageMod = 1.25
		} else if p.Year == 5 {
			ageMod = 1.45
		}

		/// Higher star players are more likely to transfer
		switch p.Stars {
		case 0:
			starMod = 1
		case 1:
			starMod = .66
		case 2:
			starMod = .75
		case 3:
			starMod = util.GenerateFloatFromRange(0.9, 1.1)
		case 4:
			starMod = util.GenerateFloatFromRange(1.11, 1.3)
		case 5:
			starMod = util.GenerateFloatFromRange(1.31, 1.75)
		}

		// Check Team Position Rank
		teamRoster := fullRosterMap[uint(p.TeamID)]
		filteredRosterByPosition := filterRosterByPosition(teamRoster, p.Position)
		youngerPlayerAhead := false
		idFound := false
		for idx, pl := range filteredRosterByPosition {
			if pl.Age < p.Age && !idFound {
				youngerPlayerAhead = true
			}
			if pl.ID == p.ID {
				idFound = true
				// Check the index of the player.
				// If they're at the top of the list, they're considered to be starting caliber.
				if (p.Position == util.Goalie) && idx > 1 {
					depthChartCompetitionMod += 33
				}

				if (p.Position == util.Center) && idx > 2 {
					depthChartCompetitionMod += 33
				}

				if (p.Position == util.Forward ||
					p.Position == util.Defender) && idx > 3 {
					depthChartCompetitionMod += 33
				}
			}
		}

		// If there's a modifier applied and there's a younger player ahead on the roster, double the amount on the modifier
		if depthChartCompetitionMod > 0 {
			if youngerPlayerAhead {
				depthChartCompetitionMod += 33
			} else {
				depthChartCompetitionMod = .63 * depthChartCompetitionMod
			}
		}

		// Check for scheme based on Team Recruiting Profile.
		// If it is not a good fit for the player, they will want to transfer
		// Will Need to Lock Scheme Dropdown by halfway through the season or by end of season

		teamID := p.TeamID
		if teamID == 0 {
			teamID = uint16(p.PreviousTeamID)
		}
		teamProfile := teamProfileMap[uint(teamID)]
		schemeMod = getSchemeMod(teamProfile, p, mediumDrop, mediumGain)

		fcsMod := 1.0
		if p.TeamID > 134 && p.TeamID != 138 && p.TeamID != 206 {
			if p.Year > 2 && p.Overall > 39 {
				fcsMod += (0.1 * float64(p.Year))
			}
			if p.Personality == "Loyal" {
				fcsMod = 0.0
			}
		}

		/// Not playing = 25, low depth chart = 16 or 33, scheme = 10, if you're all 3, that's a ~60% chance of transferring pre- modifiers
		transferWeight = starMod * ageMod * (snapMod + depthChartCompetitionMod + schemeMod) * fcsMod
		diceRoll := util.GenerateIntFromRange(1, 100)

		// NOT INTENDING TO TRANSFER
		transferInt := int(transferWeight)
		if diceRoll > transferInt {
			continue
		}

		if p.Year == 1 {
			fmt.Println("Dice Roll: ", diceRoll)
		}

		// Is Intending to transfer
		p.DeclareTransferIntention(int(transferWeight))
		transferCount++
		if p.Year == 1 && !p.IsRedshirt {
			freshmanCount++
		} else if p.Year == 2 && p.IsRedshirt {
			redshirtFreshmanCount++
		} else if p.Year == 2 && !p.IsRedshirt {
			sophomoreCount++
		} else if p.Year == 3 && p.IsRedshirt {
			redshirtSophomoreCount++
		} else if p.Year == 3 && !p.IsRedshirt {
			juniorCount++
		} else if p.Year == 4 && p.IsRedshirt {
			redshirtJuniorCount++
		} else if p.Year == 4 && !p.IsRedshirt {
			seniorCount++
		} else if p.Year == 5 && p.IsRedshirt {
			redshirtSeniorCount++
		}

		if transferWeight < 30 {
			lowCount++
		} else if transferWeight < 70 {
			mediumCount++
		} else {
			highCount++
		}

		repository.SaveCollegeHockeyPlayerRecord(p, db)
		if p.Stars > 2 {
			message := "Breaking News! " + strconv.Itoa(int(p.Stars)) + " star " + p.Position + " " + p.FirstName + " " + p.LastName + " has announced their intention to transfer from " + p.Team + "!"
			CreateNewsLog("CHL", message, "Transfer Portal", int(p.TeamID), ts, true)
		}
		notificationMessage := strconv.Itoa(int(p.Stars)) + " star " + p.Position + " " + p.FirstName + " " + p.LastName + " has a " + p.TransferLikeliness + " likeliness of entering the transfer portal. Please navigate to the Roster page to submit a promise."
		CreateNotification("CHL", notificationMessage, "Transfer Intention", uint(p.TeamID))
		// fmt.Println(strconv.Itoa(p.Year)+" YEAR "+p.TeamAbbr+" "+p.Position+" "+p.FirstName+" "+p.LastName+" HAS ANNOUNCED THEIR INTENTION TO TRANSFER | Weight: ", int(transferWeight))
		// // db.Save(&p)
		// csvModel := structs.MapPlayerToCSVModel(p)
		// playerRow := []string{
		// 	p.TeamAbbr, csvModel.FirstName, csvModel.LastName, strconv.Itoa(p.Stars),
		// 	csvModel.Archetype, csvModel.Position,
		// 	csvModel.Year, strconv.Itoa(p.Age), csvModel.RedshirtStatus,
		// 	csvModel.OverallGrade, p.RecruitingBias, p.TransferLikeliness, strconv.Itoa(transferInt), strconv.Itoa(diceRoll),
		// 	fmt.Sprintf("%.3f", ageMod), fmt.Sprintf("%.3f", snapMod), fmt.Sprintf("%.3f", starMod), fmt.Sprintf("%.3f", depthChartCompetitionMod), fmt.Sprintf("%.3f", schemeMod), fmt.Sprintf("%.3f", fcsMod),
		// }

		// err = writer.Write(playerRow)
		// if err != nil {
		// 	log.Fatal("Cannot write player row to CSV", err)
		// }

		// writer.Flush()
		// err = writer.Error()
		// if err != nil {
		// 	log.Fatal("Error while writing to file ::", err)
		// }
	}
	transferPortalMessage := "Breaking News! About " + strconv.Itoa(transferCount) + " players intend to transfer from their current schools. Teams have one week to commit promises to retain players."
	CreateNewsLog("CHL", transferPortalMessage, "Transfer Portal", 0, ts, true)
	ts.EnactPromisePhase()
	repository.SaveTimestamp(ts, db)
	fmt.Println("Total number of players entering the transfer portal: ", transferCount)
	fmt.Println("Total number of freshmen entering the transfer portal: ", freshmanCount)
	fmt.Println("Total number of redshirt freshmen entering the transfer portal: ", redshirtFreshmanCount)
	fmt.Println("Total number of sophomores entering the transfer portal: ", sophomoreCount)
	fmt.Println("Total number of redshirt sophomores entering the transfer portal: ", redshirtSophomoreCount)
	fmt.Println("Total number of juniors entering the transfer portal: ", juniorCount)
	fmt.Println("Total number of redshirt juniors entering the transfer portal: ", redshirtJuniorCount)
	fmt.Println("Total number of seniors entering the transfer portal: ", seniorCount)
	fmt.Println("Total number of redshirt seniors entering the transfer portal: ", redshirtSeniorCount)
	fmt.Println("Total number of players with low likeliness to enter transfer portal: ", lowCount)
	fmt.Println("Total number of players with medium likeliness to enter transfer portal: ", mediumCount)
	fmt.Println("Total number of players with high likeliness to enter transfer portal: ", highCount)
}

/* Currently only utilizes season momentum for making promises. Which is fine. Keeps successful teams competitive. */
func AICoachPromisePhase() {
	db := dbprovider.GetInstance().GetDB()

	aiTeamProfiles := repository.FindTeamRecruitingProfiles(true)
	collegeTeamMap := GetCollegeTeamMap()

	for _, team := range aiTeamProfiles {
		if !team.IsAI || team.ID == 0 || team.IsUserTeam {
			continue
		}

		teamID := strconv.Itoa(int(team.ID))
		roster := repository.FindCollegePlayersByTeamID(teamID)
		for _, p := range roster {
			if p.TransferStatus > 1 || p.TransferStatus == 0 {
				continue
			}
			collegePlayerID := strconv.Itoa(int(p.ID))
			promise := GetCollegePromiseByCollegePlayerID(collegePlayerID, teamID)
			if promise.ID != 0 {
				continue
			}

			team := collegeTeamMap[uint(team.ID)]

			promiseOdds := getBasePromiseOdds(team, p.PlayerPreferences)
			diceRoll := util.GenerateIntFromRange(1, 100)

			if diceRoll < promiseOdds {
				// Commit Promise
				promiseWeight := "Medium"
				promiseType := ""
				benchmarkStr := ""
				promiseBenchmark := 0

				promiseType = "Minutes"
				promiseBenchmark = 10

				if p.Overall > 20 {
					promiseBenchmark += 6
				} else if p.Overall > 16 {
					promiseBenchmark += 3
				} else if p.Overall < 10 {
					promiseBenchmark -= 3
				} else if p.Overall < 7 {
					promiseBenchmark -= 6
				}
				promiseWeight = getPromiseWeightByTimeOrWins("Time on Ice", promiseBenchmark)
				if p.SeasonMomentumPref > 6 || p.ProgramPref > 6 {
					// Promise based on wins
					promiseBenchmark = 15
					promiseType = "Wins"
					if team.SeasonMomentum > 8 {
						promiseBenchmark += 10
					} else if team.SeasonMomentum > 5 {
						promiseBenchmark += 6
					} else if team.SeasonMomentum < 5 {
						promiseBenchmark -= 6
					} else if team.SeasonMomentum < 3 {
						promiseBenchmark -= 10
					}
					promiseWeight = getPromiseWeightByTimeOrWins("Wins", promiseBenchmark)
				}

				if promiseType == "" {
					continue
				}

				collegePromise := structs.CollegePromise{
					TeamID:          team.ID,
					CollegePlayerID: p.ID,
					PromiseType:     promiseType,
					PromiseWeight:   promiseWeight,
					Benchmark:       promiseBenchmark,
					BenchmarkStr:    benchmarkStr,
					IsActive:        true,
				}
				repository.CreateCollegePromiseRecord(collegePromise, db)
			}
		}
	}
}

func CreatePromise(promise structs.CollegePromise) structs.CollegePromise {
	db := dbprovider.GetInstance().GetDB()
	collegePlayerID := strconv.Itoa(int(promise.CollegePlayerID))
	profileID := strconv.Itoa(int(promise.TeamID))

	existingPromise := GetCollegePromiseByCollegePlayerID(collegePlayerID, profileID)
	if existingPromise.ID != 0 && existingPromise.ID > 0 {
		existingPromise.Reactivate(promise.PromiseType, promise.PromiseWeight, promise.Benchmark)
		repository.SaveCollegePromiseRecord(promise, db)
		assignPromiseToProfile(db, collegePlayerID, profileID, existingPromise.ID)
		return existingPromise
	}

	repository.CreateCollegePromiseRecord(promise, db)

	assignPromiseToProfile(db, collegePlayerID, profileID, promise.ID)

	return promise
}

func assignPromiseToProfile(db *gorm.DB, collegePlayerID, profileID string, id uint) {
	tpProfile := repository.FindTransferPortalProfileRecord(repository.TransferPortalQuery{
		CollegePlayerID: collegePlayerID,
		ProfileID:       profileID,
	})
	if tpProfile.ID > 0 {
		tpProfile.AssignPromise(id)
		repository.SaveTransferPortalProfileRecord(tpProfile, db)
	}
}

func UpdatePromise(promise structs.CollegePromise) {
	db := dbprovider.GetInstance().GetDB()
	id := strconv.Itoa(int(promise.ID))
	existingPromise := GetCollegePromiseByID(id)
	existingPromise.UpdatePromise(promise.PromiseType, promise.PromiseWeight, promise.Benchmark)
	repository.SaveCollegePromiseRecord(existingPromise, db)
}

func CancelPromise(id string) {
	db := dbprovider.GetInstance().GetDB()
	promise := GetCollegePromiseByID(id)
	promise.Deactivate()
	repository.SaveCollegePromiseRecord(promise, db)
}

func EnterTheTransferPortal() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	// Get All Teams
	teams := GetAllCollegeTeams()

	for _, t := range teams {
		teamID := strconv.Itoa(int(t.ID))
		roster := repository.FindCollegePlayersByTeamID(teamID)

		for _, p := range roster {
			if p.TransferStatus != 1 {
				continue
			}

			playerID := strconv.Itoa(int(p.ID))

			promise := GetCollegePromiseByCollegePlayerID(playerID, teamID)
			if promise.ID == 0 {
				p.WillTransfer()
				repository.SaveCollegeHockeyPlayerRecord(p, db)
				continue
			}
			// 1-100
			baseFloor := getTransferFloor(p.TransferLikeliness)
			// 10, 20, 40, 60, 70
			promiseModifier := getPromiseFloor(promise.PromiseWeight)
			difference := baseFloor - promiseModifier
			// In the future, add something like a bias modifier.
			// If the coach promises something
			// That does NOT match the bias of the player, it should not be as impactful.
			// However, this should be implemented after investigating how to make bias more impactful.

			diceRoll := util.GenerateIntFromRange(1, 100)

			// Lets say the difference is 40. 60-20.
			if diceRoll < difference {
				// If the dice roll is within the 40%. They leave.
				// Okay this makes sense.

				p.WillTransfer()

				// Create News Log
				message := "Breaking News! " + p.Team + " " + strconv.Itoa(int(p.Stars)) + " Star " + p.Position + " " + p.FirstName + " " + p.LastName + " has officially entered the transfer portal!"
				CreateNewsLog("CHL", message, "Transfer Portal", int(p.PreviousTeamID), ts, true)

				repository.SaveCollegeHockeyPlayerRecord(p, db)
				repository.DeleteCollegePromise(promise, db)
				continue
			}

			// Create News Log
			message := "Breaking News! " + p.Team + " " + strconv.Itoa(int(p.Stars)) + " Star " + p.Position + " " + p.FirstName + " " + p.LastName + " has withdrawn their name from the transfer portal!"
			CreateNewsLog("CHL", message, "Transfer Portal", int(p.PreviousTeamID), ts, true)

			promise.MakePromise()
			repository.SaveCollegePromiseRecord(promise, db)
			p.WillStay()
			repository.SaveCollegeHockeyPlayerRecord(p, db)
		}
	}

	ts.EnactPortalPhase()
	repository.SaveTimestamp(ts, db)
}

func AddTransferPlayerToBoard(transferPortalProfileDto structs.TransferPortalProfile) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	portalProfile := repository.FindTransferPortalProfileRecord(repository.TransferPortalQuery{
		CollegePlayerID: strconv.Itoa(int(transferPortalProfileDto.CollegePlayerID)),
		ProfileID:       strconv.Itoa(int(transferPortalProfileDto.ProfileID)),
	})

	// If Recruit Already Exists
	if portalProfile.CollegePlayerID != 0 && portalProfile.ProfileID != 0 {
		portalProfile.Reactivate()
		repository.SaveTransferPortalProfileRecord(portalProfile, db)
		return portalProfile
	}

	newProfileForRecruit := structs.TransferPortalProfile{
		SeasonID:           uint(transferPortalProfileDto.SeasonID),
		CollegePlayerID:    uint(transferPortalProfileDto.CollegePlayerID),
		ProfileID:          uint(transferPortalProfileDto.ProfileID),
		TeamAbbreviation:   transferPortalProfileDto.TeamAbbreviation,
		TotalPoints:        0,
		CurrentWeeksPoints: 0,
		SpendingCount:      0,
		RemovedFromBoard:   false,
	}

	repository.CreateTransferPortalProfileRecord(newProfileForRecruit, db)

	return newProfileForRecruit
}

func RemovePlayerFromTransferPortalBoard(dto structs.TransferPortalProfile) structs.TransferPortalProfile {
	db := dbprovider.GetInstance().GetDB()

	profile := repository.FindTransferPortalProfileRecord(repository.TransferPortalQuery{
		CollegePlayerID: strconv.Itoa(int(dto.CollegePlayerID)),
		ProfileID:       strconv.Itoa(int(dto.ProfileID)),
	})

	profile.Deactivate()
	pid := profile.PromiseID.Int64
	profile.RemovePromise()
	repository.SaveTransferPortalProfileRecord(profile, db)
	if pid > 0 && !profile.IsSigned {
		promiseID := strconv.Itoa(int(pid))
		promise := GetCollegePromiseByID(promiseID)
		promise.Deactivate()
		repository.DeleteCollegePromise(promise, db)
	}

	return profile
}

func AllocatePointsToTransferPlayer(updateTransferPortalBoardDto structs.UpdateTransferPortalBoard) {
	db := dbprovider.GetInstance().GetDB()

	var teamId = strconv.Itoa(updateTransferPortalBoardDto.TeamID)
	var profile = repository.FindTeamRecruitingProfile(teamId, false, false)
	var portalProfiles = repository.FindTransferPortalProfileRecords(repository.TransferPortalQuery{
		ProfileID: strconv.Itoa(updateTransferPortalBoardDto.TeamID),
		IsActive:  "Y"})
	var updatedPlayers = updateTransferPortalBoardDto.Players

	currentPoints := 0.0

	for i := 0; i < len(portalProfiles); i++ {
		updatedRecruit := GetPlayerFromTransferPortalList(int(portalProfiles[i].CollegePlayerID), updatedPlayers)

		if portalProfiles[i].CurrentWeeksPoints != updatedRecruit.CurrentWeeksPoints {

			// Allocate Points to Profile
			currentPoints += float64(updatedRecruit.CurrentWeeksPoints)
			profile.AllocateSpentPoints(float32(currentPoints))
			// If total not surpassed, allocate to the recruit and continue
			if profile.SpentPoints <= profile.WeeklyPoints {
				portalProfiles[i].AllocatePoints(updatedRecruit.CurrentWeeksPoints)
				fmt.Println("Saving recruit " + strconv.Itoa(int(portalProfiles[i].CollegePlayerID)))
			} else {
				panic("Error: Allocated more points for Profile " + strconv.Itoa(int(profile.TeamID)) + " than what is allowed.")
			}
			repository.SaveTransferPortalProfileRecord(portalProfiles[i], db)
		} else {
			currentPoints += float64(portalProfiles[i].CurrentWeeksPoints)
			profile.AllocateSpentPoints(float32(currentPoints))
		}
	}

	// Save profile
	repository.SaveTeamProfileRecord(db, profile)
}

func AICoachFillBoardsPhase() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	AITeams := repository.FindTeamRecruitingProfiles(true)
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})
	transferPortalPlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{
		TransferStatus: "2",
	})
	teamMap := GetCollegeTeamMap()
	standingsMap := GetCollegeStandingsMap(seasonID)
	profiles := []structs.TransferPortalProfile{}
	existingProfiles := repository.FindTransferPortalProfileRecords(repository.TransferPortalQuery{
		IsActive: "Y"})

	portalProfileMap := MakePortalProfileMapByTeamID(existingProfiles)

	collegeRosterMap := GetAllCollegePlayersMapByTeam()
	collegeTeamMap := GetCollegeTeamMap()

	for idx, teamProfile := range AITeams {
		if teamProfile.IsUserTeam {
			continue
		}
		fmt.Println("Iterating "+teamProfile.Team+" on IDX: ", idx)
		team := teamMap[teamProfile.ID]
		teamStandings := standingsMap[uint(teamProfile.TeamID)]
		teamID := strconv.Itoa(int(teamProfile.ID))
		portalProfileMap := portalProfileMap[teamProfile.ID]
		if len(portalProfileMap) > 100 {
			continue
		}
		roster := collegeRosterMap[teamProfile.ID]
		rosterSize := len(roster)
		teamCap := 34
		if rosterSize >= teamCap {
			continue
		}

		majorNeedsMap := getMajorNeedsMap()

		for _, r := range roster {
			if r.Overall > 18 && majorNeedsMap[r.Position] {
				majorNeedsMap[r.Position] = false
			}
		}
		profileCount := len(portalProfileMap)

		for _, tp := range transferPortalPlayers {
			if profileCount >= 100 {
				break
			}
			isBadFit := IsBadSchemeFit(teamProfile.OffensiveScheme, teamProfile.DefensiveScheme, tp.Archetype, tp.Position)
			if isBadFit || !majorNeedsMap[tp.Position] || tp.PreviousTeamID == uint8(team.ID) || portalProfileMap[tp.ID].CollegePlayerID == tp.ID || portalProfileMap[tp.ID].ID > 0 {
				continue
			}

			// Put together a player prestige rating to use as a qualifier on which teams will target specific players. Ideally more experienced coaches will be able to target higher rated players
			// playerPrestige := getPlayerPrestigeRating(tp.Stars, tp.Overall)
			// if coach.Prestige < playerPrestige {
			// 	continue
			// }
			biasMod := 0
			postSeasonStatus := teamStandings.PostSeasonStatus
			collegeTeam := collegeTeamMap[uint(teamProfile.TeamID)]
			programDiff := collegeTeam.ProgramPrestige - tp.ProgramPref
			if programDiff > 2 {
				// Get multiple season standings
				teamHistory := GetStandingsHistoryByTeamID(teamID)
				averageWins := getAverageWins(teamHistory)
				biasMod += averageWins
			} else if collegeTeam.SeasonMomentum > 7 && tp.SeasonMomentumPref > 5 {
				switch postSeasonStatus {
				case "Round of 16":
					biasMod += 10
				case "Frozen Four":
					biasMod += 15
				case "National Championship Participant":
					biasMod += 20
				case "National Champions":
					biasMod += 25
				}
			} else if collegeTeam.ConferencePrestige > 7 {
				if tp.ConferencePref > 5 {
					biasMod += 10
				} else if tp.ConferencePref > 7 {
					biasMod += 15
				} else if tp.ConferencePref > 9 {
					biasMod += 20
				} else if tp.ConferencePref < 4 {
					biasMod -= 15
				}
			} else if collegeTeam.Academics > 7 {
				if tp.AcademicsPref > 5 {
					biasMod += 10
				} else if tp.AcademicsPref > 7 {
					biasMod += 15
				} else if tp.AcademicsPref > 9 {
					biasMod += 20
				} else if tp.AcademicsPref < 4 {
					biasMod -= 15
				}
			} else if collegeTeam.Facilities > 7 {
				if tp.FacilitiesPref > 5 {
					biasMod += 10
				} else if tp.FacilitiesPref > 7 {
					biasMod += 15
				} else if tp.FacilitiesPref > 9 {
					biasMod += 20
				} else if tp.FacilitiesPref < 4 {
					biasMod -= 15
				}
			} else if collegeTeam.CoachRating > 7 {
				if tp.CoachPref > 5 {
					biasMod += 10
				} else if tp.CoachPref > 7 {
					biasMod += 15
				} else if tp.CoachPref > 9 {
					biasMod += 20
				} else if tp.CoachPref < 4 {
					biasMod -= 15
				}
			} else if collegeTeam.Traditions > 7 {
				if tp.TraditionsPref > 5 {
					biasMod += 10
				} else if tp.TraditionsPref > 7 {
					biasMod += 15
				} else if tp.TraditionsPref > 9 {
					biasMod += 20
				} else if tp.TraditionsPref < 4 {
					biasMod -= 15
				}
			} else if collegeTeam.ProgramPrestige > 7 {
				if tp.ProgramPref > 5 {
					biasMod += 10
				} else if tp.ProgramPref > 7 {
					biasMod += 15
				} else if tp.ProgramPref > 9 {
					biasMod += 20
				} else if tp.ProgramPref < 4 {
					biasMod -= 15
				}
			} else if collegeTeam.ProfessionalPrestige > 7 {
				if tp.ProfDevPref > 5 {
					biasMod += 10
				} else if tp.ProfDevPref > 7 {
					biasMod += 15
				} else if tp.ProfDevPref > 9 {
					biasMod += 20
				} else if tp.ProfDevPref < 4 {
					biasMod -= 15
				}
			} else {
				biasMod = 5
			}

			diceRoll := util.GenerateIntFromRange(1, 50)
			if diceRoll < biasMod {
				profileCount++
				portalProfile := structs.TransferPortalProfile{
					ProfileID:        teamProfile.ID,
					CollegePlayerID:  tp.ID,
					SeasonID:         uint(ts.SeasonID),
					TeamAbbreviation: teamProfile.Team,
				}
				profiles = append(profiles, portalProfile)
			}
		}

	}

	repository.CreateTransferPortalProfileRecordsBatch(db, profiles, 500)
}

func AICoachAllocateAndPromisePhase() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	AITeams := repository.FindTeamRecruitingProfiles(true)
	portalProfiles := repository.FindTransferPortalProfileRecords(repository.TransferPortalQuery{
		IsActive: "Y"})
	// Get Maps for easy access
	profileMapByTeamID := MakePortalProfileMapByTeamID(portalProfiles)
	profileMapByPlayerID := MakePortalProfileMapByPlayerID(portalProfiles)
	collegePlayers := GetAllCollegePlayers()
	transferPortalPlayerMap := MakeCollegePlayerMap(collegePlayers)
	collegePlayerMap := MakeCollegePlayerMapByTeamID(collegePlayers)
	collegeTeamMap := GetCollegeTeamMap()
	// regionMap := util.GetRegionMap()
	// Shuffles the list of AI teams so that it's not always iterating from A-Z. Gives the teams at the lower end of the list a chance to recruit other croots
	rand.Shuffle(len(AITeams), func(i, j int) {
		AITeams[i], AITeams[j] = AITeams[j], AITeams[i]
	})

	for _, teamProfile := range AITeams {
		if teamProfile.IsUserTeam {
			continue
		}

		roster := collegePlayerMap[teamProfile.ID]
		teamCap := 34
		if len(roster) >= teamCap {
			continue
		}

		majorNeedsMap := getMajorNeedsMap()

		for _, r := range roster {
			if r.Overall > 18 && majorNeedsMap[r.Position] {
				majorNeedsMap[r.Position] = false
			}
		}

		teamProfile.ResetSpentPoints()
		points := 0.0

		portalProfiles := profileMapByTeamID[teamProfile.ID]
		for _, profile := range portalProfiles {
			if points >= float64(teamProfile.WeeklyPoints) {
				break
			}
			if profile.CurrentWeeksPoints > 0 || profile.RemovedFromBoard {
				continue
			}
			tp := transferPortalPlayerMap[profile.CollegePlayerID]
			// If player has already signed or if the position has been fulfilled
			isBadFit := IsBadSchemeFit(teamProfile.OffensiveScheme, teamProfile.DefensiveScheme, tp.Archetype, tp.Position)
			if isBadFit || tp.TeamID > 0 || tp.TransferStatus == 0 || tp.ID == 0 || !majorNeedsMap[tp.Position] {
				profile.Deactivate()
				repository.SaveTransferPortalProfileRecord(profile, db)
				continue
			}
			pointsRemaining := teamProfile.WeeklyPoints - teamProfile.SpentPoints
			if teamProfile.SpentPoints >= teamProfile.WeeklyPoints || pointsRemaining <= 0 || (pointsRemaining < 1 && pointsRemaining > 0) {
				break
			}

			removePlayerFromBoard := false
			num := 0.0

			profiles := profileMapByPlayerID[profile.CollegePlayerID]
			leadingTeamVal := IsAITeamContendingForPortalPlayer(profiles)
			if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) >= float64(leadingTeamVal)*0.66 {
				// continue, leave everything alone
				points += float64(profile.CurrentWeeksPoints)
				continue
			} else if profile.CurrentWeeksPoints > 0 && profile.TotalPoints+float64(profile.CurrentWeeksPoints) < float64(leadingTeamVal)*0.66 {
				profile.Deactivate()
				repository.SaveTransferPortalProfileRecord(profile, db)
				continue
			}

			maxChance := 2
			if ts.TransferPortalRound > 3 {
				maxChance = 4
			}
			chance := util.GenerateIntFromRange(1, maxChance)
			if (chance < 2 && ts.TransferPortalPhase <= 3) || (chance < 4 && ts.TransferPortalPhase > 3) {
				continue
			}

			min := teamProfile.AIMaxThreshold
			max := teamProfile.AIMinThreshold
			if max > 10 {
				max = 10
			}
			if !teamProfile.IsUserTeam && max > 5 {
				max = 5
			}
			if min >= max {
				min = 4
			}
			num = util.GenerateFloatFromRange(float64(min), float64(max))
			if num > float64(pointsRemaining) {
				num = float64(pointsRemaining)
			}

			if num+profile.TotalPoints < float64(leadingTeamVal)*0.66 {
				removePlayerFromBoard = true
			}
			if leadingTeamVal < 8 {
				removePlayerFromBoard = false
			}

			if removePlayerFromBoard {
				profile.Deactivate()
				repository.SaveTransferPortalProfileRecord(profile, db)
				continue
			}
			profile.AllocatePoints(int(num))
			points += num

			team := collegeTeamMap[uint(teamProfile.TeamID)]
			p := transferPortalPlayerMap[profile.CollegePlayerID]

			// Generate Promise based on coach bias
			if profile.PromiseID.Int64 == 0 && !profile.RolledOnPromise {
				promiseOdds := getBasePromiseOdds(team, p.PlayerPreferences)
				diceRoll := util.GenerateIntFromRange(1, 100)

				if diceRoll < promiseOdds {
					// Commit Promise
					promiseWeight := "Medium"
					promiseType := ""
					benchmarkStr := ""
					promiseBenchmark := 0

					promiseType = "Minutes"
					promiseBenchmark = 10

					if p.Overall > 20 {
						promiseBenchmark += 6
					} else if p.Overall > 16 {
						promiseBenchmark += 3
					} else if p.Overall < 10 {
						promiseBenchmark -= 3
					} else if p.Overall < 7 {
						promiseBenchmark -= 6
					}
					promiseWeight = getPromiseWeightByTimeOrWins("Time on Ice", promiseBenchmark)
					if p.SeasonMomentumPref > 6 || p.ProgramPref > 6 {
						// Promise based on wins
						promiseBenchmark = 15
						promiseType = "Wins"
						if team.SeasonMomentum > 8 {
							promiseBenchmark += 10
						} else if team.SeasonMomentum > 5 {
							promiseBenchmark += 6
						} else if team.SeasonMomentum < 5 {
							promiseBenchmark -= 6
						} else if team.SeasonMomentum < 3 {
							promiseBenchmark -= 10
						}
						promiseWeight = getPromiseWeightByTimeOrWins("Wins", promiseBenchmark)
					}

					if promiseType != "" {
						collegePromise := structs.CollegePromise{
							TeamID:          uint(teamProfile.TeamID),
							CollegePlayerID: tp.ID,
							PromiseType:     promiseType,
							PromiseWeight:   promiseWeight,
							Benchmark:       promiseBenchmark,
							BenchmarkStr:    benchmarkStr,
							IsActive:        true,
						}
						repository.CreateCollegePromiseRecord(collegePromise, db)
					}
				}

				profile.ToggleRolledOnPromise()
			}
			// Save Profile
			if profile.CurrentWeeksPoints > 0 {
				repository.SaveTransferPortalProfileRecord(profile, db)
			}
		}
		teamProfile.AIAllocateSpentPoints(float32(points))
		repository.SaveTeamProfileRecord(db, teamProfile)
	}
}

func SyncTransferPortal() {
	db := dbprovider.GetInstance().GetDB()
	//GetCurrentWeek
	ts := GetTimestamp()
	// Use IsRecruitingLocked to lock the TP when not in use
	teamProfiles := repository.FindTeamRecruitingProfiles(false)
	teamProfileMap := MakeTeamProfileMap(teamProfiles)
	transferPortalPlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{
		TransferStatus: "2",
	})
	// Only get players that are actively being recruited or have not been signed yet.
	//
	transferPortalProfiles := repository.FindTransferPortalProfileRecords(repository.TransferPortalQuery{
		IsActive: "Y"})
	// Get Maps for easy access
	transferPortalProfileMap := MakePortalProfileMapByPlayerID(transferPortalProfiles)
	rosterMap := GetAllCollegePlayersMapByTeam()
	collegePromises := GetAllCollegePromises()
	collegePromiseMap := MakeCollegePromiseMap(collegePromises)

	if !ts.IsRecruitingLocked {
		ts.ToggleLockRecruiting()
		repository.SaveTimestamp(ts, db)
	}

	for _, portalPlayer := range transferPortalPlayers {
		// Skip over players that have already transferred
		if portalPlayer.TransferStatus != 2 || portalPlayer.TeamID > 0 {
			continue
		}

		portalProfiles := transferPortalProfileMap[portalPlayer.ID]
		if len(portalProfiles) == 0 && ts.TransferPortalRound < uint(util.FinalPortalRound) {
			continue
		}

		// If no one has a profile on them during round 10
		if len(portalProfiles) == 0 && ts.TransferPortalRound == 10 && len(portalPlayer.TransferLikeliness) > 0 {
			roster := rosterMap[uint(portalPlayer.PreviousTeamID)]
			if len(roster) > util.MaxCollegeRosterSize {
				continue
			}
			rosterMap[uint(portalPlayer.PreviousTeamID)] = append(rosterMap[uint(portalPlayer.PreviousTeamID)], portalPlayer)
			portalPlayer.WillReturn()
			repository.SaveCollegeHockeyPlayerRecord(portalPlayer, db)
			continue
		}

		totalPointsOnPlayer := 0.0
		eligiblePointThreshold := 0.0
		readyToSign := false
		minSpendingCount := 100
		maxSpendingCount := 0
		signingMinimum := util.PortalSigningMinimum
		teamCount := 0
		eligibleTeams := []structs.TransferPortalProfile{}

		for i := range portalProfiles {
			if portalProfiles[i].RemovedFromBoard {
				continue
			}
			promise := collegePromiseMap[uint(portalProfiles[i].PromiseID.Int64)]

			multiplier := getMultiplier(promise)
			portalProfiles[i].AddPointsToTotal(multiplier)
		}

		sort.Slice(portalProfiles, func(i, j int) bool {
			return portalProfiles[i].TotalPoints > portalProfiles[j].TotalPoints
		})

		for i := range portalProfiles {
			// roster := rosterMap[portalProfiles[i].ProfileID]
			// tp := teamProfileMap[strconv.Itoa(int(portalProfiles[i].ProfileID))]
			// // if (len(roster) > 105 && tp.IsFBS) || (len(roster) > 80 && !tp.IsFBS) {
			// // 	continue
			// // }
			if eligiblePointThreshold == 0.0 {
				eligiblePointThreshold = portalProfiles[i].TotalPoints * signingMinimum
			}
			if portalProfiles[i].TotalPoints >= eligiblePointThreshold {
				if portalProfiles[i].SpendingCount < minSpendingCount {
					minSpendingCount = portalProfiles[i].SpendingCount
				}
				if portalProfiles[i].SpendingCount > maxSpendingCount {
					maxSpendingCount = portalProfiles[i].SpendingCount
				}
				eligibleTeams = append(eligibleTeams, portalProfiles[i])
				totalPointsOnPlayer += portalProfiles[i].TotalPoints
				teamCount += 1
			}

		}

		if (teamCount >= 1 && minSpendingCount >= 2) || (teamCount > 1 && minSpendingCount > 3) || (ts.TransferPortalRound >= uint(util.FinalPortalRound)) && totalPointsOnPlayer > 0 {
			// threshold met
			readyToSign = true
		}
		var winningTeamID uint = 0
		var odds float64 = 0
		if readyToSign {
			for winningTeamID == 0 && len(eligibleTeams) > 0 {
				percentageOdds := rand.Float64() * (totalPointsOnPlayer)
				currentProbability := 0.0
				for _, profile := range eligibleTeams {
					currentProbability += profile.TotalPoints
					if percentageOdds <= currentProbability {
						// WINNING TEAM
						winningTeamID = profile.ProfileID
						odds = profile.TotalPoints / totalPointsOnPlayer * 100
						break
					}
				}

				if winningTeamID > 0 {
					teamProfile := teamProfileMap[winningTeamID]
					currentRoster := rosterMap[teamProfile.ID]
					teamCap := util.MaxCollegeRosterSize
					if len(currentRoster) < teamCap {
						promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), strconv.Itoa(int(winningTeamID)))
						if promise.ID > 0 {
							promise.MakePromise()
							repository.SaveCollegePromiseRecord(promise, db)
						}
						portalPlayer.SignWithNewTeam(int(teamProfile.TeamID), teamProfile.Team, 1)
						message := portalPlayer.FirstName + " " + portalPlayer.LastName + ", " + strconv.Itoa(int(portalPlayer.Stars)) + " star " + portalPlayer.Position + " from " + portalPlayer.PreviousTeam + " has signed with " + portalPlayer.Team + " with " + strconv.Itoa(int(odds)) + " percent odds."
						CreateNewsLog("CHL", message, "Transfer Portal", int(winningTeamID), ts, true)
						fmt.Println("Created new log!")
						// Add player to existing roster map
						rosterMap[teamProfile.ID] = append(rosterMap[teamProfile.ID], portalPlayer)
						for i := range portalProfiles {
							if portalProfiles[i].ID == winningTeamID {
								portalProfiles[i].SignPlayer()
								break
							}
						}

					} else {
						// Filter out profile
						eligibleTeams = FilterOutPortalProfile(eligibleTeams, winningTeamID)
						winningTeamID = 0
						if len(eligibleTeams) == 0 {
							break
						}

						totalPointsOnPlayer = 0
						for _, p := range eligibleTeams {
							totalPointsOnPlayer += p.TotalPoints
						}
					}

				}
			}

		}
		for _, p := range portalProfiles {
			if winningTeamID > 0 && p.ID != winningTeamID {
				p.RemovePromise()
				p.Lock()
			}
			if winningTeamID > 0 || p.SpendingCount > 0 {
				repository.SaveTransferPortalProfileRecord(p, db)
			}
			fmt.Println("Save transfer portal profile from " + portalPlayer.Team + " towards " + portalPlayer.FirstName + " " + portalPlayer.LastName)
			if winningTeamID > 0 && p.ProfileID != winningTeamID {
				promise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(portalPlayer.ID)), strconv.Itoa(int(p.ProfileID)))
				if promise.ID > 0 {
					repository.DeleteCollegePromise(promise, db)
				}
			}
		}
		// Save Recruit
		if portalPlayer.TeamID > 0 {
			repository.SaveCollegeHockeyPlayerRecord(portalPlayer, db)
		}
	}

	ts.IncrementTransferPortalRound()
	repository.SaveTimestamp(ts, db)
}

// Portal Helper Functions
func getSchemeMod(tp *structs.RecruitingTeamProfile, p structs.CollegePlayer, drop, gain float64) float64 {
	schemeMod := 0.0
	if tp.OffensiveScheme == "" && tp.DefensiveScheme == "" {
		fmt.Println("PING!")
	}
	goodFit := IsGoodSchemeFit(tp.OffensiveScheme, tp.DefensiveScheme, p.Archetype, p.Position)
	badFit := IsBadSchemeFit(tp.OffensiveScheme, tp.DefensiveScheme, p.Archetype, p.Position)
	if goodFit {
		schemeMod += drop
	} else if badFit {
		schemeMod += gain
	}

	return schemeMod
}

func IsGoodSchemeFit(offensiveScheme, defensiveScheme, arch, position string) bool {
	archType := arch + " " + position
	offensiveSchemeList := GetFitsByScheme(offensiveScheme, false)
	defensiveSchemeList := GetFitsByScheme(defensiveScheme, false)
	totalFitList := append(offensiveSchemeList, defensiveSchemeList...)

	return CheckPlayerFits(archType, totalFitList)
}

func IsBadSchemeFit(offensiveScheme, defensiveScheme, arch, position string) bool {
	archType := arch + " " + position
	offensiveSchemeList := GetFitsByScheme(offensiveScheme, true)
	defensiveSchemeList := GetFitsByScheme(defensiveScheme, true)
	totalFitList := append(offensiveSchemeList, defensiveSchemeList...)

	return CheckPlayerFits(archType, totalFitList)
}

// Update this for hockey schemes later this offseason
func GetFitsByScheme(scheme string, isBadFit bool) []string {
	fullMap := map[string]structs.SchemeFits{
		"Power Run":                  {GoodFits: []string{"Power RB", "Blocking FB", "Blocking TE", "Red Zone Threat WR", "Run Blocking OG", "Run Blocking OT", "Run Blocking C"}, BadFits: []string{"Speed RB", "Receiving RB", "Receiving FB", "Receiving TE", "Vertical Threat TE", "Pass Blocking OG", "Pass Blocking OT", "Pass Blocking C"}},
		"Vertical":                   {GoodFits: []string{"Pocket QB", "Receiving RB", "Receiving TE", "Vertical Threat TE", "Route Runner WR", "Speed WR", "Pass Blocking OG", "Pass Blocking OT", "Pass Blocking C"}, BadFits: []string{"Balanced QB", "Scrambler QB", "Field General QB", "Power RB", "Blocking FB", "Rushing FB", "Blocking TE", "Red Zone Threat WR", "Run Blocking OG", "Run Blocking OT", "Run Blocking C"}},
		"West Coast":                 {GoodFits: []string{"Field General QB", "Balanced FB", "Receiving FB", "Receiving TE", "Route Runner WR", "Possession WR", "Line Captain C"}, BadFits: []string{"Blocking FB", "Red Zone Threat WR"}},
		"I-Option":                   {GoodFits: []string{"Scrambler QB", "Power RB", "Rushing FB", "Blocking TE", "Possession WR"}, BadFits: []string{"Pocket QB", "Speed RB", "Receiving RB", "Receiving FB", "Receiving TE", "Vertical Threat TE"}},
		"Run and Shoot":              {GoodFits: []string{"Field General QB", "Speed RB", "Receiving RB", "Speed WR", "Line Captain C"}, BadFits: []string{"Balanced RB", "Power RB", "Blocking FB", "Rushing FB", "Blocking TE", "Possession WR"}},
		"Air Raid":                   {GoodFits: []string{"Pocket QB", "Receiving RB", "Receiving FB", "Receiving TE", "Vertical Threat TE", "Speed WR", "Pass Blocking OG", "Pass Blocking OT", "Pass Blocking C"}, BadFits: []string{"Balanced QB", "Scrambler QB", "Field General QB", "Power RB", "Blocking FB", "Rushing FB", "Blocking TE", "Run Blocking OG", "Run Blocking OT", "Run Blocking C"}},
		"Pistol":                     {GoodFits: []string{"Balanced QB", "Pocket QB", "Balanced RB", "Rushing FB", "Vertical Threat TE", "Route Runner WR", "Possession WR"}, BadFits: []string{"Balanced FB", "Line Captain C"}},
		"Spread Option":              {GoodFits: []string{"Scrambler QB", "Speed RB", "Receiving FB", "Route Runner WR", "Possession WR"}, BadFits: []string{"Balanced RB", "Balanced FB"}},
		"Wing-T":                     {GoodFits: []string{"Balanced QB", "Balanced RB", "Balanced FB", "Speed WR"}, BadFits: []string{}},
		"Double Wing":                {GoodFits: []string{"Power RB", "Blocking FB", "Rushing FB", "Blocking TE", "Red Zone Threat WR", "Run Blocking OG", "Run Blocking OT", "Run Blocking C"}, BadFits: []string{"Pocket QB", "Speed RB", "Receiving RB", "Receiving FB", "Receiving TE", "Vertical Threat TE", "Pass Blocking OG", "Pass Blocking OT", "Pass Blocking C"}},
		"Wishbone":                   {GoodFits: []string{"Balanced QB", "Field General QB", "Balanced RB", "Red Zone Threat WR"}, BadFits: []string{"Balanced FB", "Route Runner WR", "Line Captain C"}},
		"Flexbone":                   {GoodFits: []string{"Scrambler QB", "Speed RB", "Balanced FB", "Red Zone Threat WR"}, BadFits: []string{"Balanced RB", "Speed WR", "Possession WR"}},
		"Old School":                 {GoodFits: []string{"Run Stopper DE", "Run Stopper OLB", "Run Stopper ILB", "Field General ILB"}, BadFits: []string{"Nose Tackle DT", "Coverage OLB", "Coverage ILB"}},
		"2-Gap":                      {GoodFits: []string{"Run Stopper DE", "Nose Tackle DT", "Run Stopper OLB", "Pass Rush OLB", "Run Stopper ILB"}, BadFits: []string{"Speed Rusher DE", "Pass Rusher DT", "Speed OLB", "Speed ILB"}},
		"4-man Front Spread Stopper": {GoodFits: []string{"Speed Rusher DE", "Pass Rusher DT", "Coverage OLB", "Coverage ILB"}, BadFits: []string{"Run Stopper DE", "Nose Tackle DT", "Run Stoppper OLB", "Run Stopper ILB", "Run Stopper FS", "Run Stopper SS"}},
		"3-man Front Spread Stopper": {GoodFits: []string{"Nose Tackle DT", "Pash Rush OLB", "Coverage ILB"}, BadFits: []string{"Nose Tackle DT", "Run Stopper OLB", "Run Stopper ILB", "Run Stopper FS", "Run Stopper SS", "Speed OLB", "Speed ILB", "Field General ILB"}},
		"Speed":                      {GoodFits: []string{"Speed Rusher DE", "Pass Rusher DT", "Coverage OLB", "Speed OLB", "Speed ILB"}, BadFits: []string{"Run Stopper DE", "Nose Tackle DT", "Pass Rush OLB", "Field General ILB"}},
		"Multiple":                   {GoodFits: []string{"Run Stopper DE", "Speed OLB", "Speed ILB", "Field General ILB", "Run Stopper FS", "Run Stopper SS"}, BadFits: []string{"Speed Rusher DE", "Pass Rusher DT", "Coverage OLB", "Coverage ILB"}},
	}
	if schemeFits, ok := fullMap[scheme]; ok {
		if isBadFit {
			return schemeFits.BadFits
		}
		return schemeFits.GoodFits
	}
	return []string{}
}

func CheckPlayerFits(player string, fits []string) bool {
	for _, fit := range fits {
		if player == fit {
			return true
		}
	}
	return false
}

func filterRosterByPosition(roster []structs.CollegePlayer, pos string) []structs.CollegePlayer {
	estimatedSize := len(roster) / 5 // Adjust this based on your data
	filteredList := make([]structs.CollegePlayer, 0, estimatedSize)
	for _, p := range roster {
		if p.Position != pos || (p.Year == 5 || (p.Year == 4 && p.IsRedshirt)) {
			continue
		}
		filteredList = append(filteredList, p)
	}
	sort.Slice(filteredList, func(i, j int) bool {
		return filteredList[i].Overall > filteredList[j].Overall
	})

	return filteredList
}

// GetTransferFloor -- Get the Base Floor to determine if a player will transfer or not based on a promise
func getTransferFloor(likeliness string) int {
	min := 25
	max := 100
	switch likeliness {
	case "Low":
		max = 40
	case "Medium":
		min = 45
		max = 70
	default:
		min = 75
	}

	return util.GenerateIntFromRange(min, max)
}

// getPromiseFloor -- Get the modifier towards the floor value above
func getPromiseFloor(weight string) int {
	if weight == "Very Low" {
		return 10
	}
	if weight == "Low" {
		return 20
	}
	if weight == "Medium" {
		return 40
	}
	if weight == "High" {
		return 60
	}
	return util.GenerateIntFromRange(70, 80)
}

func getBasePromiseOdds(profile structs.CollegeTeam, pref structs.PlayerPreferences) int {
	promiseOdds := 50

	profPrestigeOdds := getPromiseLevelDifference(int(profile.ProfessionalPrestige), int(pref.ProfDevPref))
	programPrestige := getPromiseLevelDifference(int(profile.ProgramPrestige), int(pref.ProgramPref))
	academicsPrestige := getPromiseLevelDifference(int(profile.Academics), int(pref.AcademicsPref))
	traditionsPrestige := getPromiseLevelDifference(int(profile.Traditions), int(pref.TraditionsPref))
	atmospherePrestige := getPromiseLevelDifference(int(profile.Atmosphere), int(pref.AtmospherePref))
	conferencePrestige := getPromiseLevelDifference(int(profile.ConferencePrestige), int(pref.ConferencePref))
	coachRatingPrestige := getPromiseLevelDifference(int(profile.CoachRating), int(pref.CoachPref))
	seasonMomentum := getPromiseLevelDifference(int(profile.SeasonMomentum), int(pref.SeasonMomentumPref))

	promiseOdds += profPrestigeOdds + programPrestige + academicsPrestige + traditionsPrestige + atmospherePrestige + conferencePrestige + coachRatingPrestige + seasonMomentum

	return promiseOdds
}

func getPromiseLevelDifference(teamPrestige, playerPref int) int {
	diff := teamPrestige - playerPref
	if diff >= 2 {
		return 7
	}
	if diff <= -2 {
		return -7
	}
	return 0
}

func getPromiseWeightByTimeOrWins(category string, benchmark int) string {
	if benchmark == 0 {
		return "Very Low"
	}
	weight := "Medium"
	if category == "Wins" {
		if benchmark <= 8 {
			weight = "Very Low"
		}
		if benchmark <= 14 {
			weight = "Low"
		}
		if benchmark <= 21 {
			weight = "Medium"
		}
		if benchmark <= 27 {
			weight = "High"
		}
		if benchmark <= 31 {
			weight = "Very High"
		}

	}
	if category == "Time On Ice" {
		if benchmark <= 3 {
			weight = "Very Low"
		}
		if benchmark <= 7 {
			weight = "Low"
		}
		if benchmark <= 10 {
			weight = "Medium"
		}
		if benchmark <= 13 {
			weight = "High"
		}
		if benchmark <= 20 {
			weight = "Very High"
		}

	}
	return weight
}

func GetPlayerFromTransferPortalList(id int, profiles []structs.TransferPortalProfileResponse) structs.TransferPortalProfileResponse {
	var profile structs.TransferPortalProfileResponse

	for i := 0; i < len(profiles); i++ {
		if profiles[i].CollegePlayerID == uint(id) {
			profile = profiles[i]
			break
		}
	}

	return profile
}

func getMajorNeedsMap() map[string]bool {
	majorNeedsMap := make(map[string]bool)

	if _, ok := majorNeedsMap["C"]; !ok {
		majorNeedsMap["C"] = true
	}

	if _, ok := majorNeedsMap["F"]; !ok {
		majorNeedsMap["F"] = true
	}

	if _, ok := majorNeedsMap["D"]; !ok {
		majorNeedsMap["D"] = true
	}

	if _, ok := majorNeedsMap["G"]; !ok {
		majorNeedsMap["G"] = true
	}

	return majorNeedsMap
}

func getAverageWins(standings []structs.CollegeStandings) int {
	wins := 0
	for _, s := range standings {
		wins += int(s.TotalWins)
	}

	totalStandings := len(standings)
	if totalStandings > 0 {
		wins = wins / len(standings)
	}

	return wins
}

func IsAITeamContendingForPortalPlayer(profiles []structs.TransferPortalProfile) int {
	if len(profiles) == 0 {
		return 0
	}
	leadingVal := 0
	for _, profile := range profiles {
		if profile.TotalPoints != 0 && profile.TotalPoints > float64(leadingVal) {
			leadingVal = int(profile.TotalPoints)
		}
	}

	return leadingVal
}

func getMultiplier(pr structs.CollegePromise) float64 {
	if pr.ID == 0 || !pr.IsActive {
		return 1
	}
	weight := pr.PromiseWeight
	switch weight {
	case "Very Low":
		return 1.05
	case "Low":
		return 1.1
	case "Medium":
		return 1.3
	case "High":
		return 1.5
	}
	// Very High
	return 1.75
}

func FilterOutPortalProfile(profiles []structs.TransferPortalProfile, ID uint) []structs.TransferPortalProfile {
	var rp []structs.TransferPortalProfile

	for _, profile := range profiles {
		if profile.ProfileID != ID {
			rp = append(rp, profile)
		}
	}

	return rp
}
