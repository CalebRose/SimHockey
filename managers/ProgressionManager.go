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

	growth := GetGrowth(int(player.Age), int(player.PrimeAge), int(player.Regression), float64(player.DecayRate), true)

	metMinutesOrRedshirt := true

	// Attributes
	agility := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Agility", growth, metMinutesOrRedshirt)
	faceoffs := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Faceoffs", growth, metMinutesOrRedshirt)
	CloseShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotAccuracy", growth, metMinutesOrRedshirt)
	CloseShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotPower", growth, metMinutesOrRedshirt)
	LongShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotAccuracy", growth, metMinutesOrRedshirt)
	LongShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotPower", growth, metMinutesOrRedshirt)
	passing := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Passing", growth, metMinutesOrRedshirt)
	puckHandling := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "PuckHandling", growth, metMinutesOrRedshirt)
	strength := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Strength", growth, metMinutesOrRedshirt)
	bodyChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "BodyChecking", growth, metMinutesOrRedshirt)
	stickChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "StickChecking", growth, metMinutesOrRedshirt)
	shotBlocking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "ShotBlocking", growth, metMinutesOrRedshirt)
	goalkeeping := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Goalkeeping", growth, metMinutesOrRedshirt)
	goalieVision := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieVision", growth, metMinutesOrRedshirt)
	goalieRebound := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieRebound", growth, metMinutesOrRedshirt)

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
	// minutes := 0

	// for _, stat := range stats {
	// 	minutes += int(stat.TimeOnIce)
	// }

	// averageTimeOnIce := 0
	// if len(stats) > 0 {
	// 	averageTimeOnIce = minutes / 34 // Total regular season games
	// }

	growth := GetGrowth(int(player.Age), int(player.PrimeAge), int(player.Regression), float64(player.DecayRate), false)

	metMinutesOrRedshirt := true

	// Attributes
	agility := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Agility", growth, metMinutesOrRedshirt)
	faceoffs := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Faceoffs", growth, metMinutesOrRedshirt)
	CloseShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotAccuracy", growth, metMinutesOrRedshirt)
	CloseShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "CloseShotPower", growth, metMinutesOrRedshirt)
	LongShotAccuracy := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotAccuracy", growth, metMinutesOrRedshirt)
	LongShotPower := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "LongShotPower", growth, metMinutesOrRedshirt)
	passing := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Passing", growth, metMinutesOrRedshirt)
	puckHandling := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "PuckHandling", growth, metMinutesOrRedshirt)
	strength := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Strength", growth, metMinutesOrRedshirt)
	bodyChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "BodyChecking", growth, metMinutesOrRedshirt)
	stickChecking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "StickChecking", growth, metMinutesOrRedshirt)
	shotBlocking := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "ShotBlocking", growth, metMinutesOrRedshirt)
	goalkeeping := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "Goalkeeping", growth, metMinutesOrRedshirt)
	goalieVision := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieVision", growth, metMinutesOrRedshirt)
	goalieRebound := calculateAttributeGrowth(&player.BasePlayer, &player.BasePotentials, "GoalieRebound", growth, metMinutesOrRedshirt)

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

func calculateAttributeGrowth(player *structs.BasePlayer, pots *structs.BasePotentials, attribute string, baseGrowth int, metMinutesOrRedshirt bool) int {
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
	if metMinutesOrRedshirt && player.Age < player.PrimeAge {
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
