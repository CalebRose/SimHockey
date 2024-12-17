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
	// ts := GetTimestamp()
	// SeasonID := strconv.Itoa(int(ts.SeasonID))
	// statMap := nil
	collegeTeams := GetAllCollegeTeams()

	graduatingPlayers := []structs.ProfessionalPlayer{}

	for _, team := range collegeTeams {
		teamID := strconv.Itoa(int(team.ID))
		roster := GetCollegePlayersByTeamID(teamID)
		// croots := GetSignedRecruitsByTeamProfileID(teamID)

		if !team.PlayersProgressed {
			for _, player := range roster {
				if player.HasProgressed {
					continue
				}
				willDeclare := (player.Year > 4 && !player.IsRedshirt) || (player.Year > 5 && player.IsRedshirt)
				if willDeclare {
					professionalPlayer := structs.ProfessionalPlayer{
						BasePlayer:     player.BasePlayer,
						BasePotentials: player.BasePotentials,
						BaseInjuryData: player.BaseInjuryData,
						Year:           0,
					}
					graduatingPlayers = append(graduatingPlayers, professionalPlayer)
				} else {
					repository.SaveCollegeHockeyPlayerRecord(player, db)
				}
			}
		}
	}
}

func ProfessionalProgressionMain() {

}

func ProgressCollegePlayer(player structs.CollegePlayer, SeasonID string, stats []structs.CollegePlayerGameStats) structs.CollegePlayer {
	// minutes := 0

	// for _, stat := range stats {
	// 	minutes += int(stat.TimeOnIce)
	// }

	// averageTimeOnIce := 0
	// if len(stats) > 0 {
	// 	averageTimeOnIce = minutes / 34 // Total regular season games
	// }

	growth := GetGrowth(int(player.Age), int(player.PrimeAge), int(player.Regression), float64(player.DecayRate))

	metMinutesOrRedshirt := true

	// Attributes
	agility := calculateAttributeGrowth(&player, "Agility", growth, metMinutesOrRedshirt)
	faceoffs := calculateAttributeGrowth(&player, "Faceoffs", growth, metMinutesOrRedshirt)
	slapshotAccuracy := calculateAttributeGrowth(&player, "SlapshotAccuracy", growth, metMinutesOrRedshirt)
	slapshotPower := calculateAttributeGrowth(&player, "SlapshotPower", growth, metMinutesOrRedshirt)
	wristShotAccuracy := calculateAttributeGrowth(&player, "WristShotAccuracy", growth, metMinutesOrRedshirt)
	wristShotPower := calculateAttributeGrowth(&player, "WristShotPower", growth, metMinutesOrRedshirt)
	passing := calculateAttributeGrowth(&player, "Passing", growth, metMinutesOrRedshirt)
	puckHandling := calculateAttributeGrowth(&player, "PuckHandling", growth, metMinutesOrRedshirt)
	strength := calculateAttributeGrowth(&player, "Strength", growth, metMinutesOrRedshirt)
	bodyChecking := calculateAttributeGrowth(&player, "BodyChecking", growth, metMinutesOrRedshirt)
	stickChecking := calculateAttributeGrowth(&player, "StickChecking", growth, metMinutesOrRedshirt)
	shotBlocking := calculateAttributeGrowth(&player, "ShotBlocking", growth, metMinutesOrRedshirt)
	goalkeeping := calculateAttributeGrowth(&player, "Goalkeeping", growth, metMinutesOrRedshirt)
	goalieVision := calculateAttributeGrowth(&player, "GoalieVision", growth, metMinutesOrRedshirt)
	goalieRebound := calculateAttributeGrowth(&player, "GoalieRebound", growth, metMinutesOrRedshirt)

	progressions := structs.BasePlayerProgressions{
		Agility:              agility,
		Faceoffs:             faceoffs,
		SlapshotAccuracy:     slapshotAccuracy,
		SlapshotPower:        slapshotPower,
		WristShotAccuracy:    wristShotAccuracy,
		WristShotPower:       wristShotPower,
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

func calculateAttributeGrowth(player *structs.CollegePlayer, attribute string, baseGrowth int, metMinutesOrRedshirt bool) int {
	// Map attribute names to potentials
	potentialMap := map[string]uint8{
		"Agility":           player.AgilityPotential,
		"Faceoffs":          player.FaceoffsPotential,
		"SlapshotAccuracy":  player.SlapshotAccuracyPotential,
		"SlapshotPower":     player.SlapshotPowerPotential,
		"WristShotAccuracy": player.WristShotAccuracyPotential,
		"WristShotPower":    player.WristShotPowerPotential,
		"Passing":           player.PassingPotential,
		"PuckHandling":      player.PuckHandlingPotential,
		"Strength":          player.StrengthPotential,
		"BodyChecking":      player.BodyCheckingPotential,
		"StickChecking":     player.StickCheckingPotential,
		"ShotBlocking":      player.ShotBlockingPotential,
		"Goalkeeping":       player.GoalkeepingPotential,
		"GoalieVision":      player.GoalieVisionPotential,
		"GoalieRebound":     player.GoalieReboundPotential,
	}

	// Attributes not attributing to a goalie should not be progressing
	if player.Position == "G" && (attribute != "Goalkeeping" && attribute != "GoalieVision" && attribute != "GoalieRebound" && attribute != "Strength" && attribute != "Agility") {
		return 0
	}

	// Fetch the potential for the specific attribute
	potential := potentialMap[attribute]

	g := baseGrowth
	if metMinutesOrRedshirt {
		g += util.GenerateIntFromRange(0, 2)
	}

	// Scale base growth by potential
	scaledGrowth := int(float64(g) * (float64(potential) / 100.0))

	return scaledGrowth
}

func GetGrowth(age, primeage, regression int, decayRate float64) int {
	if age == primeage || age == primeage+1 {
		return rand.Intn(2)
	}
	if age < primeage {
		return GetPrePrimeGrowth(age, primeage)
	}
	return GetPostPrimeGrowth(age, primeage, regression, decayRate)
}

func GetPrePrimeGrowth(age, primeage int) int {
	baseGrowth := 7.0

	ageDifference := float64(age / primeage)
	ageMultiplier := 1 - ageDifference

	growth := baseGrowth * ageMultiplier
	return int(growth)
}

func GetPostPrimeGrowth(age, primeage, regression int, decayRate float64) int {
	yearsPastPrime := age - primeage
	postRegression := -1 * (regression + int(float64(yearsPastPrime)*decayRate))
	return postRegression
}
