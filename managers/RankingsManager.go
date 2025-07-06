package managers

import (
	"math"
	"math/rand/v2"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func ESPNModifiers() map[string]map[string]float32 {
	return map[string]map[string]float32{
		"C": {
			"Height": 72.00,
			"Weight": 199.0,
		},
		"D": {
			"Height": 74.00,
			"Weight": 202.00,
		},
		"F": {
			"Height": 72.00,
			"Weight": 196,
		},
		"G": {
			"Height": 75.00,
			"Weight": 200,
		},
	}
}

func RivalsModifiers() []int {
	return []int{100, 83, 82, 81, 80,
		76, 75, 74, 73, 72,
		69, 68, 67, 66, 65,
		64, 63, 62, 61, 60,
		59, 58, 57, 56, 55,
		53, 53, 53, 53, 53,
		51, 51, 51, 51, 51,
		49, 49, 49, 49, 49,
		47, 47, 47, 47, 47,
		45, 45, 45, 45, 45,
		43, 43, 43, 43, 43,
		41, 41, 41, 41, 41,
		40, 40, 40, 40, 40,
		39, 39, 39, 39, 39,
		38, 38, 38, 38, 38,
		37, 37, 37, 37, 37,
		36, 36, 36, 36, 36,
		35, 35, 35, 35, 35,
		34, 34, 34, 34, 34,
		33, 33, 33, 33, 33,
		32, 32, 32, 32, 32,
		31, 31, 31, 31, 31,
		30, 30, 30, 30, 30,
		29, 29, 29, 29, 29,
		28, 28, 28, 28, 28,
		27, 27, 27, 27, 27,
		26, 26, 26, 26, 26,
		25, 25, 25, 25, 25,
		24, 24, 24, 24, 24,
		23, 23, 23, 23, 23,
		22, 22, 22, 22, 22,
		21, 21, 21, 21, 21,
		20, 20, 20, 20, 20,
		19, 19, 19, 19, 19,
		18, 18, 18, 18, 18,
		17, 17, 17, 17, 17,
		16, 16, 16, 16, 16,
		15, 15, 15, 15, 15,
		14, 14, 14, 14, 14,
		13, 13, 13, 13, 13,
		12, 12, 12, 12, 12,
		11, 11, 11, 11, 11,
		10, 10, 10, 10, 10,
		9, 9, 9, 9, 9,
		8, 8, 8, 8, 8,
		7, 7, 7, 7, 7,
		6, 6, 6, 6, 6,
		5, 5, 5, 5, 5,
		4, 4, 4, 4, 4,
		3, 3, 3, 3, 3,
	}
}

func AssignAllRecruitRanks() {
	db := dbprovider.GetInstance().GetDB()
	recruits := repository.FindAllRecruits(false, false, false, true, false, "")

	var rivalsMod float32 = 100.0

	for _, croot := range recruits {
		potentialGrade := GetAveragePotentialGrade(croot.BasePotentials)
		rank247 := Get247Ranking(potentialGrade, int(croot.Overall))

		espnRank := GetESPNRanking(croot, potentialGrade)

		var rivalsRank float32 = 0
		rivalsBonus := rivalsMod
		rivalsRank = GetRivalsRanking(int(croot.Stars), rivalsBonus)

		var r float32 = croot.TopRankModifier

		if croot.TopRankModifier == 0 || croot.TopRankModifier < 0.95 || croot.TopRankModifier > 1.05 {
			r = float32(0.95 + rand.Float64()*(1.05-0.95))
		}

		if croot.Stars == 0 {
			rank247 = 0.001
			espnRank = 0.001
			rivalsRank = 0.001
			r = 1
		}

		croot.AssignRankValues(rank247, espnRank, rivalsRank, r)

		recruitingModifier := getRecruitingModifier()

		croot.AssignRecruitingModifier(recruitingModifier)

		repository.SaveCollegeHockeyRecruitRecord(croot, db)
		if rivalsMod > 0.1 {
			rivalsMod -= 0.1
		}
	}
}

func GetAveragePotentialGrade(pots structs.BasePotentials) string {
	total := int(pots.AgilityPotential) + int(pots.BodyCheckingPotential) + int(pots.CloseShotAccuracyPotential) + int(pots.CloseShotPowerPotential) +
		int(pots.FaceoffsPotential) + int(pots.GoalieVisionPotential) + int(pots.GoalkeepingPotential) + int(pots.LongShotAccuracyPotential) +
		int(pots.LongShotPowerPotential) + int(pots.PassingPotential) + int(pots.PuckHandlingPotential) + int(pots.ShotBlockingPotential) + int(pots.StickCheckingPotential) +
		int(pots.StrengthPotential)

	avg := total / 14

	return util.GetPotentialGrade(avg)
}

func Get247Ranking(pg string, ovr int) float32 {
	mod := Get247PotentialModifier(pg)

	return float32(ovr) + (mod * 2)
}

func Get247PotentialModifier(pg string) float32 {
	if pg == "A+" {
		return 7.83
	} else if pg == "A" {
		return 7.06
	} else if pg == "A-" {
		return 6.77
	} else if pg == "B+" {
		return 6.33
	} else if pg == "B" {
		return 6.04
	} else if pg == "B-" {
		return 5.87
	} else if pg == "C+" {
		return 5.58
	} else if pg == "C" {
		return 5.43
	} else if pg == "C-" {
		return 5.31
	} else if pg == "D+" {
		return 5.03
	} else if pg == "D" {
		return 4.77
	} else if pg == "D-" {
		return 4.67
	}
	return 4.3
}

func GetESPNRanking(r structs.Recruit, pg string) float32 {
	// ESPN Ranking = Star Rank + Archetype Modifier + weight difference + height difference
	// + potential val, and then round.

	starRank := GetESPNStarRank(int(r.Stars))
	archMod := GetArchetypeModifier(r.Archetype)
	potentialMod := GetESPNPotentialModifier(pg)

	espnPositionMap := ESPNModifiers()
	heightMod := float32(r.Height) / espnPositionMap[r.Position]["Height"]
	weightMod := float32(r.Weight) / espnPositionMap[r.Position]["Weight"]
	espnRanking := math.Round(float64(starRank) + float64(archMod) + potentialMod + float64(heightMod) + float64(weightMod))

	return float32(espnRanking)
}

func GetESPNPotentialModifier(pg string) float64 {
	if pg == "A+" {
		return 1
	} else if pg == "A" {
		return 0.9
	} else if pg == "A-" {
		return 0.8
	} else if pg == "B+" {
		return 0.6
	} else if pg == "B" {
		return 0.4
	} else if pg == "B-" {
		return 0.2
	} else if pg == "C+" {
		return 0
	} else if pg == "C" {
		return -0.15
	} else if pg == "C-" {
		return -0.3
	} else if pg == "D+" {
		return -0.6
	} else if pg == "D" {
		return -0.75
	} else if pg == "D-" {
		return -0.9
	}
	return -1
}

func GetESPNStarRank(star int) int {
	if star == 5 {
		return 95
	} else if star == 4 {
		return 85
	} else if star == 3 {
		return 75
	} else if star == 2 {
		return 65
	}
	return 55
}

func GetArchetypeModifier(arch string) int {
	if arch == Power || arch == Sniper {
		return 1
	} else if arch == Enforcer || arch == Grinder || arch == TwoWay {
		return -1
	} else if arch == Playmaker || arch == Defensive || arch == Offensive {
		return 2
	}
	return 0
}
func GetRivalsRanking(stars int, bonus float32) float32 {
	return GetRivalsStarModifier(stars) + bonus
}
func GetRivalsStarModifier(stars int) float32 {
	if stars == 5 {
		return 6.1
	} else if stars == 4 {
		return RoundToFixedDecimalPlace(rand.Float64()*((6.0-5.8)+5.8), 1)
	} else if stars == 3 {
		return RoundToFixedDecimalPlace(rand.Float64()*((5.7-5.5)+5.5), 1)
	} else if stars == 2 {
		return RoundToFixedDecimalPlace(rand.Float64()*((5.4-5.2)+5.2), 1)
	} else {
		return 5
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundToFixedDecimalPlace(num float64, precision int) float32 {
	output := math.Pow(10, float64(precision))
	return float32(round(num*output)) / float32(output)
}

func getRecruitingModifier() float32 {
	diceRoll := util.GenerateFloatFromRange(1, 20)
	if diceRoll == 1 {
		return 0.02
	}
	num := util.GenerateIntFromRange(1, 100)
	mod := 0.0
	if num < 11 {
		mod = util.GenerateFloatFromRange(1.80, 2.00)
	} else if num < 31 {
		mod = util.GenerateFloatFromRange(1.50, 1.69)
	} else if num < 71 {
		mod = util.GenerateFloatFromRange(1.15, 1.49)
	} else if num < 91 {
		mod = util.GenerateFloatFromRange(0.90, 1.14)
	} else {
		mod = util.GenerateFloatFromRange(0.80, 0.89)
	}

	return float32(mod)
}
