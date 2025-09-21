package managers

import (
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

	graduatingPlayers := []structs.ProfessionalPlayer{}
	draftablePlayers := []structs.DraftablePlayer{}
	historicRecords := []structs.HistoricCollegePlayer{}

	for _, team := range collegeTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetCollegePlayersByTeamID(teamID)
		croots := repository.FindAllRecruits(false, true, true, false, false, teamID)

		if !team.PlayersProgressed {
			for _, player := range roster {
				if player.HasProgressed {
					continue
				}
				stats := gameStatMap[player.ID]
				player = ProgressCollegePlayer(player, SeasonID, stats)
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
					// Assign their drafted team ID if they have been drafted
					professionalPlayer.AssignTeam(player.DraftedTeamID, player.DraftedTeam)
					graduatingPlayers = append(graduatingPlayers, professionalPlayer)
				} else if willDeclare && player.DraftedTeamID == 0 {
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

	repository.CreateProHockeyPlayerRecordsBatch(db, graduatingPlayers, 200)
	repository.CreateDraftablePlayerRecordsBatch(db, draftablePlayers, 200)
	repository.CreateHistoricCollegePlayerRecordsBatch(db, historicRecords, 200)
}

func ProfessionalProgressionMain() {

}

func ProgressCollegePlayer(player structs.CollegePlayer, SeasonID string, stats []structs.CollegePlayerGameStats) structs.CollegePlayer {
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

	// Attributes
	agility := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Agility", growth, metMinutes, player.IsRedshirt)
	faceoffs := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Faceoffs", growth, metMinutes, player.IsRedshirt)
	CloseShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotAccuracy", growth, metMinutes, player.IsRedshirt)
	CloseShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotPower", growth, metMinutes, player.IsRedshirt)
	LongShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotAccuracy", growth, metMinutes, player.IsRedshirt)
	LongShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotPower", growth, metMinutes, player.IsRedshirt)
	passing := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Passing", growth, metMinutes, player.IsRedshirt)
	puckHandling := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "PuckHandling", growth, metMinutes, player.IsRedshirt)
	strength := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Strength", growth, metMinutes, player.IsRedshirt)
	bodyChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "BodyChecking", growth, metMinutes, player.IsRedshirt)
	stickChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "StickChecking", growth, metMinutes, player.IsRedshirt)
	shotBlocking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "ShotBlocking", growth, metMinutes, player.IsRedshirt)
	goalkeeping := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Goalkeeping", growth, metMinutes, player.IsRedshirt)
	goalieVision := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieVision", growth, metMinutes, player.IsRedshirt)
	goalieRebound := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieRebound", growth, metMinutes, player.IsRedshirt)

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
