package managers

import (
	"fmt"
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"golang.org/x/exp/rand"
)

func CollegeProgressionMain() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(int(ts.SeasonID))
	collegePlayerGameStats := repository.FindCollegePlayerGameStatsRecords(SeasonID, "", "", "")
	gameStatMap := MakeCollegePlayerGameStatsMap(collegePlayerGameStats)
	collegeTeams := GetAllCollegeTeams()

	canadianHockeyTeams := GetAllCanadianCHLTeams()

	graduatingPlayers := []structs.ProfessionalPlayer{}
	draftablePlayers := []structs.DraftablePlayer{}
	historicRecords := []structs.HistoricCollegePlayer{}
	collegePlayerIDs := []string{}

	for _, team := range collegeTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetCollegePlayersByTeamID(teamID)
		croots := repository.FindAllRecruits(false, true, true, false, false, teamID)

		if !team.PlayersProgressed {
			for _, player := range roster {
				if player.HasProgressed {
					continue
				}
				id := strconv.Itoa(int(player.ID))

				stats := gameStatMap[player.ID]
				player = ProgressCollegePlayer(player, SeasonID, stats, false)
				willDeclare := (player.Year > 4 && !player.IsRedshirt) || (player.Year > 5 && player.IsRedshirt)
				if willDeclare && player.DraftedTeamID > 0 {
					historicRecord := structs.HistoricCollegePlayer{CollegePlayer: player}
					historicRecords = append(historicRecords, historicRecord)
					// Graduating players with draft rights become pro players
					// Create a new professional player record from the college player data
					professionalPlayer := structs.ProfessionalPlayer{
						Model:          player.Model,
						DraftedTeamID:  uint8(player.DraftedTeamID),
						BasePlayer:     player.BasePlayer,
						BasePotentials: player.BasePotentials,
						BaseInjuryData: player.BaseInjuryData,
						Year:           0,
					}
					collegePlayerIDs = append(collegePlayerIDs, id)
					// Assign their drafted team ID if they have been drafted
					professionalPlayer.AssignTeam(player.DraftedTeamID, player.DraftedTeam, 1)
					graduatingPlayers = append(graduatingPlayers, professionalPlayer)
				} else if willDeclare && player.DraftedTeamID == 0 {
					collegePlayerIDs = append(collegePlayerIDs, id)
					historicRecord := structs.HistoricCollegePlayer{CollegePlayer: player}
					historicRecords = append(historicRecords, historicRecord)
					// Graduate players with no draft rights become draftee records before UDFAs
					draftee := structs.DraftablePlayer{
						Model:          player.Model,
						BasePlayer:     player.BasePlayer,
						BasePotentials: player.BasePotentials,
						BaseInjuryData: player.BaseInjuryData,
						CollegeID:      uint(player.TeamID),
					}
					draftablePlayers = append(draftablePlayers, draftee)
				} else {
					repository.SaveCollegeHockeyPlayerRecord(player, db)
				}
			}

			team.TogglePlayersProgressed()
		}

		// Add Recruits
		if !team.RecruitsAdded {
			playersToAdd := []structs.CollegePlayer{}
			for _, croot := range croots {
				cp := structs.CollegePlayer{
					Model:          croot.Model,
					BasePlayer:     croot.BasePlayer,
					BasePotentials: croot.BasePotentials,
					BaseInjuryData: croot.BaseInjuryData,
					Year:           1,
				}
				playersToAdd = append(playersToAdd, cp)
			}
			repository.CreateCollegeHockeyPlayerRecordsBatch(db, playersToAdd, 10)

			team.ToggleRecruitsAdded()
		}

		repository.SaveCollegeTeamRecord(db, team)
	}

	for _, team := range canadianHockeyTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetCollegePlayersByTeamID(teamID)

		if !team.PlayersProgressed {
			for _, player := range roster {
				if player.HasProgressed {
					continue
				}
				id := strconv.Itoa(int(player.ID))

				stats := gameStatMap[player.ID]
				player = ProgressCollegePlayer(player, SeasonID, stats, false)
				willGraduateFromTeam := (player.Age > 20)
				if willGraduateFromTeam && player.DraftedTeamID > 0 {
					historicRecord := structs.HistoricCollegePlayer{CollegePlayer: player}
					historicRecords = append(historicRecords, historicRecord)
					// Graduating players with draft rights become pro players
					// Create a new professional player record from the college player data
					professionalPlayer := structs.ProfessionalPlayer{
						Model:          player.Model,
						DraftedTeamID:  uint8(player.DraftedTeamID),
						BasePlayer:     player.BasePlayer,
						BasePotentials: player.BasePotentials,
						BaseInjuryData: player.BaseInjuryData,
						Year:           0,
					}
					collegePlayerIDs = append(collegePlayerIDs, id)
					// Assign their drafted team ID if they have been drafted
					professionalPlayer.AssignTeam(player.DraftedTeamID, player.DraftedTeam, 1)
					graduatingPlayers = append(graduatingPlayers, professionalPlayer)
				} else if willGraduateFromTeam && player.DraftedTeamID == 0 {
					// If player is over the age of 20 and players in the Canadian leagues, they must graduate into the portal no matter what
					player.WillTransfer()
					repository.SaveCollegeHockeyPlayerRecord(player, db)
				} else {
					repository.SaveCollegeHockeyPlayerRecord(player, db)
				}
			}

			team.TogglePlayersProgressed()
		}

		repository.SaveCollegeTeamRecord(db, team)
	}

	// Unsigned Players
	unsignedPlayers := repository.FindAllCollegePlayers(repository.PlayerQuery{TeamID: "0"})

	for _, player := range unsignedPlayers {
		if player.HasProgressed {
			continue
		}
		id := strconv.Itoa(int(player.ID))

		stats := gameStatMap[player.ID]
		player = ProgressCollegePlayer(player, SeasonID, stats, false)
		willDeclare := (player.Year > 4 && !player.IsRedshirt) || (player.Year > 5 && player.IsRedshirt)
		isCanadian := player.Country == util.Canada && (willDeclare && player.Age > 24) // Phase out really old Canadian players
		isOther := player.Country != util.Canada && willDeclare && player.Age > 21
		if (isCanadian || isOther) && player.DraftedTeamID > 0 {
			historicRecord := structs.HistoricCollegePlayer{CollegePlayer: player}
			historicRecords = append(historicRecords, historicRecord)
			// Graduating players with draft rights become pro players
			// Create a new professional player record from the college player data
			professionalPlayer := structs.ProfessionalPlayer{
				Model:          player.Model,
				DraftedTeamID:  uint8(player.DraftedTeamID),
				BasePlayer:     player.BasePlayer,
				BasePotentials: player.BasePotentials,
				BaseInjuryData: player.BaseInjuryData,
				Year:           0,
			}
			collegePlayerIDs = append(collegePlayerIDs, id)
			// Assign their drafted team ID if they have been drafted
			professionalPlayer.AssignTeam(player.DraftedTeamID, player.DraftedTeam, 1)
			graduatingPlayers = append(graduatingPlayers, professionalPlayer)
		} else if (isCanadian || isOther) && player.DraftedTeamID == 0 {
			collegePlayerIDs = append(collegePlayerIDs, id)
			historicRecord := structs.HistoricCollegePlayer{CollegePlayer: player}
			historicRecords = append(historicRecords, historicRecord)
			// Graduate players with no draft rights become draftee records before UDFAs
			draftee := structs.DraftablePlayer{
				Model:          player.Model,
				BasePlayer:     player.BasePlayer,
				BasePotentials: player.BasePotentials,
				BaseInjuryData: player.BaseInjuryData,
				CollegeID:      uint(player.TeamID),
			}
			draftablePlayers = append(draftablePlayers, draftee)
		} else {
			repository.SaveCollegeHockeyPlayerRecord(player, db)
		}
	}

	unsignedRecruits := repository.FindAllRecruits(false, false, false, false, false, "")
	collegePlayerBatch := []structs.CollegePlayer{}
	for _, croot := range unsignedRecruits {
		if croot.TeamID > 0 {
			continue
		}
		cp := structs.CollegePlayer{
			Model:          croot.Model,
			BasePlayer:     croot.BasePlayer,
			BasePotentials: croot.BasePotentials,
			BaseInjuryData: croot.BaseInjuryData,
			Year:           1,
		}
		collegePlayerBatch = append(collegePlayerBatch, cp)
	}

	repository.CreateCollegeHockeyPlayerRecordsBatch(db, collegePlayerBatch, 100)
	repository.CreateProHockeyPlayerRecordsBatch(db, graduatingPlayers, 200)
	repository.CreateDraftablePlayerRecordsBatch(db, draftablePlayers, 200)
	repository.CreateHistoricCollegePlayerRecordsBatch(db, historicRecords, 200)
	repository.MassDeleteCollegePlayerRecords(db, collegePlayerIDs)
}

func ProfessionalProgressionMain() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	SeasonID := strconv.Itoa(int(ts.SeasonID))
	proPlayerGameStats := repository.FindProPlayerGameStatsRecords(SeasonID, "", "", "")
	proPlayerSeasonStats := repository.FindProPlayerSeasonStatsRecords("", "2")
	seasonStatMap := MakeHistoricProPlayerSeasonStatMap(proPlayerSeasonStats)
	gameStatMap := MakeProPlayerGameStatsMap(proPlayerGameStats)
	proTeams := GetAllProfessionalTeams()
	freeAgents := GetAllFreeAgents()
	proContracts := repository.FindAllProContracts(true)
	proContractMap := MakeContractMap(proContracts)
	extensions := repository.FindAllProExtensions(true)
	extensionMap := MakeExtensionMap(extensions)

	for _, team := range proTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetProPlayersByTeamID(teamID)

		if !team.PlayersProgressed {
			for _, player := range roster {
				if player.ID == 18 {
					fmt.Println(("Stop here"))
				}
				if player.HasProgressed {
					continue
				}
				stats := gameStatMap[player.ID]

				if player.PrimeAge == 0 {
					pa := util.GetPrimeAge(player.Position, player.Archetype)
					player.PrimeAge = uint8(pa)
				}
				player = ProgressProPlayer(player, SeasonID, stats)
				willRetire := DetermineIfRetiring(player, seasonStatMap)
				// Get Contraact Info
				contract := proContractMap[player.ID]
				if contract.ID == 0 && !player.IsFreeAgent {
					player.ToggleIsFreeAgent()
				} else {
					contract.ProgressContract()
					if contract.ContractLength == 0 || contract.IsComplete {
						// Check for existing extension
						extension := extensionMap[player.ID]
						if extension.ID > 0 && extension.IsAccepted && extension.IsActive {
							// Apply Extension
							contract.MapExtension(extension)
							message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has official signed his extended offer with " + player.Team + " for $" + strconv.Itoa(int(contract.ContractValue)) + " Million Dollars!"
							CreateNewsLog("PHL", message, "Free Agency", int(player.TeamID), ts, true)
							repository.DeleteExtensionRecord(extension, db)
						} else {
							// Player becomes a free agent
							player.ToggleIsFreeAgent()
						}
					}
					if willRetire {
						contract.ToggleRetirement()
					}
					repository.SaveProContractRecord(contract, db)
				}

				// If Player Retires
				if !willRetire {
					player.ToggleHasProgressed()
					repository.SaveProPlayerRecord(player, db)
					continue
				}
				historicRecord := structs.RetiredPlayer{ProfessionalPlayer: player}
				repository.CreateRetiredPlayer(historicRecord, db)
				repository.DeleteProPlayerRecord(player, db)
				message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has decided to retire from SimPHL. He was drafted by " + player.DraftedTeam + " and last played with " + player.Team + " and " + player.PreviousTeam + ". We thank him for his wondrous, extensive career and hope he enjoys his retirement!"
				CreateNewsLog("PHL", message, "Retirement", int(player.TeamID), ts, true)
			}
			team.TogglePlayersProgressed()
		} else {
			continue
		}

		repository.SaveProTeamRecord(db, team)
	}

	for _, player := range freeAgents {
		if player.HasProgressed {
			continue
		}

		stats := gameStatMap[player.ID]

		if player.PrimeAge == 0 {
			pa := util.GetPrimeAge(player.Position, player.Archetype)
			player.PrimeAge = uint8(pa)
		}
		player = ProgressProPlayer(player, SeasonID, stats)
		// willRetire := DetermineIfRetiring(player, seasonStatMap)
		willRetire := false // Free agents do not retire automatically for now
		// If Player Retires
		if !willRetire {
			player.ToggleHasProgressed()
			repository.SaveProPlayerRecord(player, db)
			continue
		}
		historicRecord := structs.RetiredPlayer{ProfessionalPlayer: player}
		repository.CreateRetiredPlayer(historicRecord, db)
		repository.DeleteProPlayerRecord(player, db)
		message := "Breaking News: " + player.Position + " " + player.FirstName + " " + player.LastName + " has decided to retire from SimPHL. He was drafted by " + player.DraftedTeam + " and last played with " + player.Team + " and " + player.PreviousTeam + ". We thank him for his wondrous, extensive career and hope he enjoys his retirement!"
		CreateNewsLog("PHL", message, "Retirement", int(player.TeamID), ts, true)
	}

	ts.ToggleProfessionalProgression()
	repository.SaveTimestamp(ts, db)
}

func ProgressCollegePlayer(player structs.CollegePlayer, SeasonID string, stats []structs.CollegePlayerGameStats, isInit bool) structs.CollegePlayer {
	minutes := 0

	for _, stat := range stats {
		// TimeOnIce is stored in seconds
		min := stat.TimeOnIce / 60
		minutes += int(min)
	}

	averageTimeOnIce := 0
	if len(stats) > 0 {
		averageTimeOnIce = minutes / 34 // Total regular season games
	}

	growth := GetGrowth(int(player.Age), int(player.PrimeAge), int(player.Regression), float64(player.DecayRate), true)
	metMinutes := averageTimeOnIce >= 12

	redshirtQualification := player.IsRedshirting || player.LeagueID > 1 || player.TeamID == 0 || isInit

	// Attributes
	agility := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Agility", growth, metMinutes, redshirtQualification)
	faceoffs := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Faceoffs", growth, metMinutes, redshirtQualification)
	CloseShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotAccuracy", growth, metMinutes, redshirtQualification)
	CloseShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotPower", growth, metMinutes, redshirtQualification)
	LongShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotAccuracy", growth, metMinutes, redshirtQualification)
	LongShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotPower", growth, metMinutes, redshirtQualification)
	passing := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Passing", growth, metMinutes, redshirtQualification)
	puckHandling := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "PuckHandling", growth, metMinutes, redshirtQualification)
	strength := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Strength", growth, metMinutes, redshirtQualification)
	bodyChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "BodyChecking", growth, metMinutes, redshirtQualification)
	stickChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "StickChecking", growth, metMinutes, redshirtQualification)
	shotBlocking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "ShotBlocking", growth, metMinutes, redshirtQualification)
	goalkeeping := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Goalkeeping", growth, metMinutes, redshirtQualification)
	goalieVision := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieVision", growth, metMinutes, redshirtQualification)
	goalieRebound := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieRebound", growth, metMinutes, redshirtQualification)

	progressions := structs.BasePlayerProgressions{
		Agility:              agility,
		Faceoffs:             faceoffs,
		CloseShotAccuracy:    CloseShotAccuracy,
		CloseShotPower:       CloseShotPower,
		LongShotAccuracy:     LongShotAccuracy,
		LongShotPower:        LongShotPower,
		Passing:              passing,
		PuckHandling:         puckHandling,
		Strength:             strength,
		BodyChecking:         bodyChecking,
		StickChecking:        stickChecking,
		ShotBlocking:         shotBlocking,
		Goalkeeping:          goalkeeping,
		GoalieVision:         goalieVision,
		GoalieReboundControl: goalieRebound,
	}

	player.ProgressPlayer(progressions)

	return player
}

func ProgressProPlayer(player structs.ProfessionalPlayer, SeasonID string, stats []structs.ProfessionalPlayerGameStats) structs.ProfessionalPlayer {
	minutes := 0

	for _, stat := range stats {
		// TimeOnIce is stored in seconds
		min := stat.TimeOnIce / 60
		minutes += int(min)
	}

	averageTimeOnIce := 0
	if len(stats) > 0 {
		averageTimeOnIce = minutes / 34 // Total regular season games
	}

	growth := GetGrowth(int(player.Age), int(player.PrimeAge), int(player.Regression), float64(player.DecayRate), false)

	metMinutes := averageTimeOnIce >= 12

	// Attributes
	agility := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Agility", growth, metMinutes, player.IsAffiliatePlayer)
	faceoffs := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Faceoffs", growth, metMinutes, player.IsAffiliatePlayer)
	CloseShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotAccuracy", growth, metMinutes, player.IsAffiliatePlayer)
	CloseShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotPower", growth, metMinutes, player.IsAffiliatePlayer)
	LongShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotAccuracy", growth, metMinutes, player.IsAffiliatePlayer)
	LongShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotPower", growth, metMinutes, player.IsAffiliatePlayer)
	passing := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Passing", growth, metMinutes, player.IsAffiliatePlayer)
	puckHandling := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "PuckHandling", growth, metMinutes, player.IsAffiliatePlayer)
	strength := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Strength", growth, metMinutes, player.IsAffiliatePlayer)
	bodyChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "BodyChecking", growth, metMinutes, player.IsAffiliatePlayer)
	stickChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "StickChecking", growth, metMinutes, player.IsAffiliatePlayer)
	shotBlocking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "ShotBlocking", growth, metMinutes, player.IsAffiliatePlayer)
	goalkeeping := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Goalkeeping", growth, metMinutes, player.IsAffiliatePlayer)
	goalieVision := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieVision", growth, metMinutes, player.IsAffiliatePlayer)
	goalieRebound := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieRebound", growth, metMinutes, player.IsAffiliatePlayer)

	progressions := structs.BasePlayerProgressions{
		Agility:              agility,
		Faceoffs:             faceoffs,
		CloseShotAccuracy:    CloseShotAccuracy,
		CloseShotPower:       CloseShotPower,
		LongShotAccuracy:     LongShotAccuracy,
		LongShotPower:        LongShotPower,
		Passing:              passing,
		PuckHandling:         puckHandling,
		Strength:             strength,
		BodyChecking:         bodyChecking,
		StickChecking:        stickChecking,
		ShotBlocking:         shotBlocking,
		Goalkeeping:          goalkeeping,
		GoalieVision:         goalieVision,
		GoalieReboundControl: goalieRebound,
	}

	player.ProgressPlayer(progressions)

	return player
}

func calculateAttributeGrowth(player *structs.BasePlayer, pots *structs.BasePotentials, attribute string, baseGrowth int, metMinutes bool, redshirtedOrAffiliated bool) int {
	// Map attribute names to potentials
	potentialMap := map[string]uint8{
		"Agility":           pots.AgilityPotential,
		"Faceoffs":          pots.FaceoffsPotential,
		"CloseShotAccuracy": pots.CloseShotAccuracyPotential,
		"CloseShotPower":    pots.CloseShotPowerPotential,
		"LongShotAccuracy":  pots.LongShotAccuracyPotential,
		"LongShotPower":     pots.LongShotPowerPotential,
		"Passing":           pots.PassingPotential,
		"PuckHandling":      pots.PuckHandlingPotential,
		"Strength":          pots.StrengthPotential,
		"BodyChecking":      pots.BodyCheckingPotential,
		"StickChecking":     pots.StickCheckingPotential,
		"ShotBlocking":      pots.ShotBlockingPotential,
		"Goalkeeping":       pots.GoalkeepingPotential,
		"GoalieVision":      pots.GoalieVisionPotential,
		"GoalieRebound":     pots.GoalieReboundPotential,
	}

	// Attributes not attributing to a goalie should not be progressing
	if player.Position == util.Goalie && (attribute != util.Goalkeeping && attribute != util.GoalieVision && attribute != util.GoalieRebound && attribute != util.Passing && attribute != util.Strength && attribute != util.Agility) {
		return 0
	}

	// Fetch the potential for the specific attribute
	potential := potentialMap[attribute]

	g := baseGrowth
	if metMinutes && player.Age < player.PrimeAge {
		g += util.GenerateIntFromRange(1, 3)
	}
	if redshirtedOrAffiliated && player.Age < player.PrimeAge {
		g += util.GenerateIntFromRange(0, 2)
	}

	// Scale base growth by potential
	scaledGrowth := int(float64(g) * (float64(potential) / 100.0))

	return scaledGrowth
}

func GetGrowth(age, primeage, regression int, decayRate float64, isCollege bool) int {
	if age == primeage || age == primeage+1 {
		return rand.Intn(1)
	}
	if age < primeage {
		return GetPrePrimeGrowth(age, primeage, isCollege)
	}
	return GetPostPrimeGrowth(age, primeage, regression, decayRate)
}

func GetPrePrimeGrowth(age, primeage int, isCollege bool) int {
	baseGrowth := 4.5
	if isCollege {
		baseGrowth = 7.5
	}

	ageDifference := float64(age / primeage)
	ageMultiplier := 1 - ageDifference

	growth := baseGrowth * ageMultiplier
	return int(growth)
}

func GetPostPrimeGrowth(age, primeage, regression int, decayRate float64) int {
	yearsPastPrime := age - primeage
	postRegression := -1.15 * (float64(regression) + (float64(yearsPastPrime) * decayRate))
	return int(postRegression)
}

func DetermineIfRetiring(player structs.ProfessionalPlayer, statMap map[uint][]structs.ProfessionalPlayerSeasonStats) bool {
	if player.IsFreeAgent && player.Year > 1 {
		lastTwoSeasonStats := statMap[player.ID]
		totalMinutes := 0
		for _, stat := range lastTwoSeasonStats {
			totalMinutes += int(stat.TimeOnIce)
		}
		return totalMinutes == 0
	}

	if player.Age <= player.PrimeAge {
		return false
	}

	/*
		Thoughts - we could implement historic injuries into this somewhere, although we are impacting prime age upon injuries.
	*/
	benchmark := 0
	age := int(player.Age)
	primeAge := int(player.PrimeAge)
	retirementAge := primeAge + util.GenerateIntFromRange(3, 5)
	if age > retirementAge {
		benchmark += 50
	}
	if age > primeAge && player.Overall < 20 {
		benchmark += (15 * (age - primeAge))
	} else if age > primeAge && player.Overall < 30 {
		benchmark += (7 * (age - primeAge))
	} else if age > primeAge && player.Overall < 40 {
		benchmark += (4 * (age - primeAge))
	} else if age > primeAge && player.Overall < 50 {
		benchmark += (2 * (age - primeAge))
	}
	diceRoll := util.GenerateIntFromRange(1, 100)
	// If the roll is less than the benchmark, player will retire. Otherwise, they are staying.
	return diceRoll < benchmark
}
