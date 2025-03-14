package util

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetPositionList() []string {
	return []string{
		Center, Forward, Defender, Goalie,
	}
}

func GetStarRating(isCustom bool) int {
	roll := GenerateIntFromRange(1, 1000)
	if isCustom {
		roll -= 100
	}
	if roll < 0 {
		roll = 1
	}
	if roll < 2 && !isCustom {
		return 6
	}
	if roll < 42 {
		return 5
	}
	if roll < 122 {
		return 4
	}
	if roll < 352 {
		return 3
	}
	if roll < 652 {
		return 2
	}
	return 1
}

// Pick a US state or Canadian province for which the player is from
func PickState() string {
	states := []Locale{
		{Name: "MN", Weight: 50},
		{Name: "MI", Weight: 45},
		{Name: "MA", Weight: 40},
		{Name: "NY", Weight: 40},
		{Name: "IL", Weight: 30},
		{Name: "WI", Weight: 29},
		{Name: "PA", Weight: 29},
		{Name: "ND", Weight: 27},
		{Name: "CO", Weight: 26},
		{Name: "OH", Weight: 25},
		{Name: "CT", Weight: 19},
		{Name: "VT", Weight: 19},
		{Name: "AK", Weight: 19},
		{Name: "NH", Weight: 10}, // Collective weight for less prominent states
		{Name: "RI", Weight: 10}, // Collective weight for less prominent states
		{Name: "ME", Weight: 10}, // Collective weight for less prominent states
		{Name: "NJ", Weight: 10}, // Collective weight for less prominent states
		{Name: "IN", Weight: 4},  // Collective weight for less prominent states
		{Name: "DE", Weight: 3},  // Collective weight for less prominent states
		{Name: "NE", Weight: 2},  // Collective weight for less prominent states
		{Name: "MT", Weight: 3},  // Collective weight for less prominent states
		{Name: "MD", Weight: 1},  // Collective weight for less prominent states
		{Name: "VA", Weight: 1},  // Collective weight for less prominent states
		{Name: "MO", Weight: 3},  // Collective weight for less prominent states
		{Name: "CA", Weight: 3},  // Collective weight for less prominent states
		{Name: "OR", Weight: 1},  // Collective weight for less prominent states
		{Name: "WA", Weight: 3},  // Collective weight for less prominent states
		{Name: "IA", Weight: 1},  // Collective weight for less prominent states
		{Name: "NM", Weight: 1},  // Collective weight for less prominent states
		{Name: "SD", Weight: 3},  // Collective weight for less prominent states
		{Name: "AZ", Weight: 2},  // Collective weight for less prominent states
		{Name: "UT", Weight: 1},  // Collective weight for less prominent states
		{Name: "WY", Weight: 1},  // Collective weight for less prominent states
		{Name: "ID", Weight: 1},  // Collective weight for less prominent states
		{Name: "NV", Weight: 1},  // Collective weight for less prominent states
		{Name: "TN", Weight: 3},  // Collective weight for less prominent states
		{Name: "NC", Weight: 2},  // Collective weight for less prominent states
		{Name: "TX", Weight: 2},  // Collective weight for less prominent states
		{Name: "KY", Weight: 1},  // Collective weight for less prominent states
		{Name: "WV", Weight: 1},  // Collective weight for less prominent states
		{Name: "SC", Weight: 1},  // Collective weight for less prominent states
		{Name: "GA", Weight: 1},  // Collective weight for less prominent states
		{Name: "AL", Weight: 1},  // Collective weight for less prominent states
		{Name: "MS", Weight: 1},  // Collective weight for less prominent states
		{Name: "FL", Weight: 2},  // Collective weight for less prominent states
		{Name: "AR", Weight: 1},  // Collective weight for less prominent states
		{Name: "LA", Weight: 1},  // Collective weight for less prominent states
		{Name: "OK", Weight: 1},  // Collective weight for less prominent states
		{Name: "KS", Weight: 1},  // Collective weight for less prominent states
		{Name: "HI", Weight: 1},  // Collective weight for less prominent states
	}

	totalWeight := 0
	for _, state := range states {
		totalWeight += state.Weight
	}

	randomWeight := GenerateIntFromRange(0, totalWeight)
	for _, state := range states {
		if randomWeight < state.Weight {
			return state.Name
		}
		randomWeight -= state.Weight
	}
	return PickFromStringList([]string{"MN", "MI", "NY", "MA"})
}

func PickProvince() string {
	provinces := []Locale{
		{Name: "ON", Weight: 40},
		{Name: "QC", Weight: 20},
		{Name: "BC", Weight: 10},
		{Name: "AB", Weight: 10},
		{Name: "MB", Weight: 8},
		{Name: "SK", Weight: 5},
		{Name: "NS", Weight: 3},
		{Name: "NB", Weight: 2},
		{Name: "PE", Weight: 1},
		{Name: "NL", Weight: 1},
		{Name: "YT", Weight: 1},  // Yukon, NWT, Nunavut combined
		{Name: "NWT", Weight: 1}, // Yukon, NWT, Nunavut combined
		{Name: "NVT", Weight: 1}, // Yukon, NWT, Nunavut combined
	}

	totalWeight := 0
	for _, province := range provinces {
		totalWeight += province.Weight
	}

	randomWeight := GenerateIntFromRange(0, totalWeight)
	for _, province := range provinces {
		if randomWeight < province.Weight {
			return province.Name
		}
		randomWeight -= province.Weight
	}
	return "ON"
}

func PickSwedishRegion() string {
	provinces := []Locale{
		{Name: "Angermanland", Weight: 10},
		{Name: "Blekinge", Weight: 10},
		{Name: "Bohuslan", Weight: 10},
		{Name: "Dalarna", Weight: 10},
		{Name: "Dalsland", Weight: 5},
		{Name: "Gastrikland", Weight: 8},
		{Name: "Gotland", Weight: 5},
		{Name: "Halland", Weight: 8},
		{Name: "Halsingland", Weight: 6},
		{Name: "Harjedalen", Weight: 4},
		{Name: "Jamtland", Weight: 6},
		{Name: "Lappland", Weight: 5},
		{Name: "Medelpad", Weight: 8},
		{Name: "Narke", Weight: 7},
		{Name: "Norrbotten", Weight: 9},
		{Name: "Oland", Weight: 3},
		{Name: "Ostergotland", Weight: 10},
		{Name: "Skane", Weight: 10},
		{Name: "Smaland", Weight: 10},
		{Name: "Sodermanland", Weight: 10},
		{Name: "Uppland", Weight: 40},
		{Name: "Varmland", Weight: 8},
		{Name: "Vasterbotten", Weight: 9},
		{Name: "Vastergotland", Weight: 10},
		{Name: "Vastmanland", Weight: 7},
	}

	totalWeight := 0
	for _, province := range provinces {
		totalWeight += province.Weight
	}

	randomWeight := GenerateIntFromRange(0, totalWeight)
	for _, province := range provinces {
		if randomWeight < province.Weight {
			return province.Name
		}
		randomWeight -= province.Weight
	}
	return "Uppland"
}

func PickRussianRegion() string {
	provinces := []Locale{
		{Name: "Adygea", Weight: 5},
		{Name: "Altai Krai", Weight: 5},
		{Name: "Altai Republic", Weight: 5},
		{Name: "Amur Oblast", Weight: 5},
		{Name: "Arkhangelsk Oblast", Weight: 5},
		{Name: "Astrakhan Oblast", Weight: 5},
		{Name: "Bashkortostan", Weight: 10},
		{Name: "Belgorod Oblast", Weight: 5},
		{Name: "Bryansk Oblast", Weight: 5},
		{Name: "Buryatia", Weight: 5},
		{Name: "Chechnya", Weight: 5},
		{Name: "Chelyabinsk Oblast", Weight: 10},
		{Name: "Chukotka Autonomous Okrug", Weight: 3},
		{Name: "Chuvashia", Weight: 5},
		{Name: "Dagestan", Weight: 5},
		{Name: "Ingushetia", Weight: 3},
		{Name: "Irkutsk Oblast", Weight: 5},
		{Name: "Ivanovo Oblast", Weight: 5},
		{Name: "Jewish Autonomous Oblast", Weight: 3},
		{Name: "Kabardino-Balkaria", Weight: 5},
		{Name: "Kaliningrad Oblast", Weight: 5},
		{Name: "Kalmykia", Weight: 3},
		{Name: "Kaluga Oblast", Weight: 5},
		{Name: "Kamchatka Krai", Weight: 3},
		{Name: "Karachay-Cherkessia", Weight: 5},
		{Name: "Karelia", Weight: 5},
		{Name: "Kemerovo Oblast", Weight: 8},
		{Name: "Khabarovsk Krai", Weight: 5},
		{Name: "Khakassia", Weight: 5},
		{Name: "Khanty-Mansi Autonomous Okrug", Weight: 8},
		{Name: "Kirov Oblast", Weight: 5},
		{Name: "Komi", Weight: 5},
		{Name: "Kostroma Oblast", Weight: 3},
		{Name: "Krasnodar Krai", Weight: 10},
		{Name: "Krasnoyarsk Krai", Weight: 13},
		{Name: "Kurgan Oblast", Weight: 5},
		{Name: "Kursk Oblast", Weight: 5},
		{Name: "Leningrad Oblast", Weight: 8},
		{Name: "Lipetsk Oblast", Weight: 5},
		{Name: "Magadan Oblast", Weight: 3},
		{Name: "Mari El", Weight: 3},
		{Name: "Mordovia", Weight: 5},
		{Name: "Moscow", Weight: 40},
		{Name: "Moscow Oblast", Weight: 30},
		{Name: "Murmansk Oblast", Weight: 5},
		{Name: "Nenets Autonomous Okrug", Weight: 2},
		{Name: "Nizhny Novgorod Oblast", Weight: 10},
		{Name: "North Ossetia-Alania", Weight: 3},
		{Name: "Novgorod Oblast", Weight: 5},
		{Name: "Novosibirsk Oblast", Weight: 13},
		{Name: "Omsk Oblast", Weight: 10},
		{Name: "Orenburg Oblast", Weight: 8},
		{Name: "Oryol Oblast", Weight: 5},
		{Name: "Penza Oblast", Weight: 5},
		{Name: "Perm Krai", Weight: 8},
		{Name: "Primorsky Krai", Weight: 8},
		{Name: "Pskov Oblast", Weight: 3},
		{Name: "Rostov Oblast", Weight: 8},
		{Name: "Ryazan Oblast", Weight: 5},
		{Name: "Sakha (Yakutia)", Weight: 3},
		{Name: "Sakhalin Oblast", Weight: 3},
		{Name: "Samara Oblast", Weight: 10},
		{Name: "Saratov Oblast", Weight: 5},
		{Name: "Saint Petersburg", Weight: 30},
		{Name: "Smolensk Oblast", Weight: 5},
		{Name: "Stavropol Krai", Weight: 5},
		{Name: "Sverdlovsk Oblast", Weight: 10},
		{Name: "Tambov Oblast", Weight: 5},
		{Name: "Tatarstan", Weight: 10},
		{Name: "Tomsk Oblast", Weight: 5},
		{Name: "Tuva", Weight: 3},
		{Name: "Tula Oblast", Weight: 5},
		{Name: "Tver Oblast", Weight: 5},
		{Name: "Tyumen Oblast", Weight: 8},
		{Name: "Udmurtia", Weight: 5},
		{Name: "Ulyanovsk Oblast", Weight: 5},
		{Name: "Vladimir Oblast", Weight: 5},
		{Name: "Volgograd Oblast", Weight: 5},
		{Name: "Vologda Oblast", Weight: 5},
		{Name: "Voronezh Oblast", Weight: 8},
		{Name: "Yamalo-Nenets Autonomous Okrug", Weight: 5},
		{Name: "Yaroslavl Oblast", Weight: 8},
		{Name: "Zabaykalsky Krai", Weight: 5},
	}

	totalWeight := 0
	for _, province := range provinces {
		totalWeight += province.Weight
	}

	randomWeight := GenerateIntFromRange(0, totalWeight)
	for _, province := range provinces {
		if randomWeight < province.Weight {
			return province.Name
		}
		randomWeight -= province.Weight
	}
	return "Moscow"
}

// getStateAbbreviation returns the two-letter state abbreviation for a given state name.
func GetStateAbbreviation(state string) (string, error) {
	// Map of state names to their two-letter abbreviations
	stateAbbreviations := map[string]string{
		"Alabama":        "AL",
		"Alaska":         "AK",
		"Arizona":        "AZ",
		"Arkansas":       "AR",
		"California":     "CA",
		"Colorado":       "CO",
		"Connecticut":    "CT",
		"Delaware":       "DE",
		"Florida":        "FL",
		"Georgia":        "GA",
		"Hawai'i":        "HI",
		"Idaho":          "ID",
		"Illinois":       "IL",
		"Indiana":        "IN",
		"Iowa":           "IA",
		"Kansas":         "KS",
		"Kentucky":       "KY",
		"Louisiana":      "LA",
		"Maine":          "ME",
		"Maryland":       "MD",
		"Massachusetts":  "MA",
		"Michigan":       "MI",
		"Minnesota":      "MN",
		"Mississippi":    "MS",
		"Missouri":       "MO",
		"Montana":        "MT",
		"Nebraska":       "NE",
		"Nevada":         "NV",
		"New Hampshire":  "NH",
		"New Jersey":     "NJ",
		"New Mexico":     "NM",
		"New York":       "NY",
		"North Carolina": "NC",
		"North Dakota":   "ND",
		"Ohio":           "OH",
		"Oklahoma":       "OK",
		"Oregon":         "OR",
		"Pennsylvania":   "PA",
		"Rhode Island":   "RI",
		"South Carolina": "SC",
		"South Dakota":   "SD",
		"Tennessee":      "TN",
		"Texas":          "TX",
		"Utah":           "UT",
		"Vermont":        "VT",
		"Virginia":       "VA",
		"Washington":     "WA",
		"West Virginia":  "WV",
		"Wisconsin":      "WI",
		"Wyoming":        "WY",
	}

	// Normalize the input by trimming spaces and capitalizing the first letter of each word
	caser := cases.Title(language.English)
	normalizedState := caser.String(strings.ToLower(strings.TrimSpace(state)))

	// Check if the state exists in the map
	if abbreviation, ok := stateAbbreviations[normalizedState]; ok {
		return abbreviation, nil
	}

	return "", fmt.Errorf("state not found: %s", state)
}

func PickPosition() string {
	roll := GenerateIntFromRange(1, 100)
	if roll < 43 {
		return Forward
	}
	if roll < 70 {
		return Defender
	}
	if roll < 84 {
		return Center
	}
	return Goalie
}

func GetArchetype(pos string) string {
	if pos == Center || pos == Forward {
		return PickFromStringList([]string{Enforcer, Grinder, Playmaker, Power, Sniper, TwoWay})
	} else if pos == Defender {
		return PickFromStringList([]string{Defensive, Defensive, Enforcer, Offensive, TwoWay})
	} else if pos == Goalie {
		return PickFromStringList([]string{StandUp, Hybrid, Butterfly})
	}
	return ""
}

func GeneratePotential(pos, arch, attr string) uint8 {
	mean := 50
	if pos == Center || pos == Forward {
		if arch == Enforcer && (attr == Agility || attr == Strength || attr == PuckHandling) {
			mean += 10
		}
		if arch == Enforcer && (attr == CloseShotPower || attr == LongShotPower) {
			mean -= 10
		}
		if arch == Grinder && (attr == BodyChecking || attr == Strength || attr == StickChecking || attr == Passing) {
			mean += 10
		}
		if arch == Grinder && (attr == CloseShotPower || attr == CloseShotAccuracy || attr == LongShotPower || attr == LongShotAccuracy || attr == PuckHandling) {
			mean -= 10
		}
		if arch == Playmaker && (attr == Passing || attr == PuckHandling) {
			mean += 10
		}
		if arch == Playmaker && (attr == Strength) {
			mean -= 10
		}
		if arch == Power && (attr == CloseShotPower || attr == Strength) {
			mean += 10
		}
		if arch == Power && (attr == BodyChecking || attr == StickChecking || attr == LongShotPower) {
			mean -= 10
		}
		if arch == Sniper && (attr == LongShotPower || attr == LongShotAccuracy || attr == Passing) {
			mean += 10
		}
		if arch == Sniper && (attr == BodyChecking || attr == StickChecking || attr == LongShotPower) {
			mean -= 10
		}
		if arch == TwoWay && (attr == Passing || attr == BodyChecking || attr == StickChecking) {
			mean += 5
		}
	} else if pos == Defender {
		if arch == Enforcer && (attr == Agility || attr == Strength || attr == BodyChecking) {
			mean += 10
		}
		if arch == Enforcer && (attr == StickChecking || attr == PuckHandling) {
			mean -= 10
		}
		if arch == Defensive && (attr == BodyChecking || attr == Strength || attr == StickChecking || attr == ShotBlocking) {
			mean += 10
		}
		if arch == Defensive && (attr == LongShotAccuracy || attr == LongShotPower || attr == CloseShotAccuracy || attr == CloseShotPower || attr == PuckHandling) {
			mean -= 10
		}
		if arch == Offensive && (attr == LongShotAccuracy || attr == LongShotPower || attr == Passing || attr == StickChecking || attr == PuckHandling) {
			mean += 10
		}
		if arch == Offensive && (attr == BodyChecking || attr == Strength || attr == CloseShotAccuracy || attr == CloseShotPower || attr == ShotBlocking) {
			mean -= 10
		}
		if arch == TwoWay && (attr == BodyChecking || attr == Passing || attr == StickChecking) {
			mean += 10
		}
		if arch == TwoWay && (attr == PuckHandling || attr == Agility || attr == CloseShotAccuracy || attr == CloseShotPower) {
			mean -= 10
		}
	} else if pos == Goalie {
		if arch == StandUp && (attr == GoalieVision || attr == Strength) {
			mean += 10
		}
		if arch == StandUp && (attr == Goalkeeping || attr == Agility) {
			mean -= 10
		}
		if arch == Butterfly && (attr == Goalkeeping || attr == Agility) {
			mean += 10
		}
		if arch == Butterfly && (attr == GoalieVision || attr == Strength) {
			mean -= 10
		}
	}
	val := GenerateNormalizedIntFromMeanStdev(float64(mean), 15)
	if val > 100 {
		val = 100
	} else if val < 1 {
		val = 1
	}
	return uint8(val)
}

func GenerateProfessionalPotential(pot int) int {
	floor := pot - 20
	ceil := pot + 20
	if floor < 0 {
		diff := 0 - floor
		floor = 0
		ceil += diff
	}
	if ceil > 100 {
		diff := ceil - 100
		ceil = 100
		floor += diff
	}
	return GenerateIntFromRange(floor, ceil)
}

func GetWeightedPotentialGrade(rating int) string {
	weightedRating := GenerateIntFromRange(rating-15, rating+15)
	if weightedRating > 100 {
		weightedRating = 99
	} else if weightedRating < 0 {
		weightedRating = 0
	}

	if weightedRating > 88 {
		return "A+"
	}
	if weightedRating > 80 {
		return "A"
	}
	if weightedRating > 74 {
		return "A-"
	}
	if weightedRating > 68 {
		return "B+"
	}
	if weightedRating > 62 {
		return "B"
	}
	if weightedRating > 56 {
		return "B-"
	}
	if weightedRating > 50 {
		return "C+"
	}
	if weightedRating > 44 {
		return "C"
	}
	if weightedRating > 38 {
		return "C-"
	}
	if weightedRating > 32 {
		return "D+"
	}
	if weightedRating > 26 {
		return "D"
	}
	if weightedRating > 20 {
		return "D-"
	}
	return "F"
}

func GetPrimeAge(pos, arch string) int {
	venerable := false
	vulnerable := false
	vDiceRoll := GenerateIntFromRange(1, 10000)
	chance := getVenerableChance(pos)
	cutShortChance := 9990
	if vDiceRoll < chance {
		venerable = true
	} else if vDiceRoll > cutShortChance {
		vulnerable = true
	}

	mean, stddev := getPositionMean(pos, venerable, vulnerable)

	age := GenerateNormalizedIntFromMeanStdev(mean, stddev)
	return int(age)
}

func getPositionMean(pos string, venerable, shortenedCareer bool) (float64, float64) {
	if shortenedCareer {
		shortenedMap := map[string][]float64{
			Center:   []float64{26, 2},
			Forward:  []float64{26, 2},
			Defender: []float64{26, 2},
			Goalie:   []float64{26, 2},
		}
		return shortenedMap[pos][0], shortenedMap[pos][1]
	}
	meanMap := getPositionMeanMap()
	return meanMap[venerable][pos][0], meanMap[venerable][pos][1]
}

func getPositionMeanMap() map[bool]map[string][]float64 {
	return map[bool]map[string][]float64{
		true: {
			Center:   []float64{35, 2},
			Forward:  []float64{33, 1},
			Defender: []float64{34, 1},
			Goalie:   []float64{34, 1},
		},
		false: {
			Center:   []float64{30, 2},
			Forward:  []float64{29, 0.67},
			Defender: []float64{29, 0.67},
			Goalie:   []float64{30, 1.33},
		},
	}
}

func getVenerableChance(pos string) int {
	return 20
}

func GetPersonality() string {
	chance := GenerateIntFromRange(1, 3)
	if chance < 3 {
		return "Average"
	}
	list := []string{"Reserved",
		"Eccentric",
		"Motivational",
		"Disloyal",
		"Cooperative",
		"Irrational",
		"Focused",
		"Book Worm",
		"Motivation",
		"Abrasive",
		"Absent Minded",
		"Uncooperative",
		"Introvert",
		"Disruptive",
		"Outgoing",
		"Tough",
		"Paranoid",
		"Stoic",
		"Dramatic",
		"Extroverted",
		"Selfish",
		"Impatient",
		"Reliable",
		"Frail",
		"Relaxed",
		"Average",
		"Flamboyant",
		"Perfectionist",
		"Popular",
		"Jokester",
		"Narcissist"}

	return PickFromStringList(list)
}
