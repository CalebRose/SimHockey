package util

import (
	"fmt"
	"strings"

	"github.com/CalebRose/SimHockey/structs"
)

func GetPositionList() []string {
	return []string{
		Center, Forward, Defender, Goalie,
	}
}

func GetStarRating() int {
	roll := GenerateIntFromRange(1, 1000)
	if roll < 40 {
		return 5
	}
	if roll < 120 {
		return 4
	}
	if roll < 350 {
		return 3
	}
	if roll < 650 {
		return 2
	}
	return 1
}

// Pick a US state or Canadian province for which the player is from
func PickState() string {
	states := []structs.Locale{
		{Name: "MN", Weight: 50},
		{Name: "MI", Weight: 45},
		{Name: "MA", Weight: 40},
		{Name: "NY", Weight: 40},
		{Name: "IL", Weight: 38},
		{Name: "WI", Weight: 37},
		{Name: "PA", Weight: 37},
		{Name: "ND", Weight: 35},
		{Name: "CO", Weight: 34},
		{Name: "OH", Weight: 33},
		{Name: "CT", Weight: 27},
		{Name: "VT", Weight: 27},
		{Name: "AK", Weight: 27},
		{Name: "NH", Weight: 18}, // Collective weight for less prominent states
		{Name: "RI", Weight: 18}, // Collective weight for less prominent states
		{Name: "ME", Weight: 18}, // Collective weight for less prominent states
		{Name: "NJ", Weight: 16}, // Collective weight for less prominent states
		{Name: "IN", Weight: 12}, // Collective weight for less prominent states
		{Name: "DE", Weight: 11}, // Collective weight for less prominent states
		{Name: "NE", Weight: 10}, // Collective weight for less prominent states
		{Name: "MT", Weight: 10}, // Collective weight for less prominent states
		{Name: "MD", Weight: 5},  // Collective weight for less prominent states
		{Name: "VA", Weight: 5},  // Collective weight for less prominent states
		{Name: "MO", Weight: 5},  // Collective weight for less prominent states
		{Name: "CA", Weight: 5},  // Collective weight for less prominent states
		{Name: "OR", Weight: 5},  // Collective weight for less prominent states
		{Name: "WA", Weight: 5},  // Collective weight for less prominent states
		{Name: "IA", Weight: 3},  // Collective weight for less prominent states
		{Name: "NM", Weight: 3},  // Collective weight for less prominent states
		{Name: "SD", Weight: 3},  // Collective weight for less prominent states
		{Name: "AZ", Weight: 3},  // Collective weight for less prominent states
		{Name: "UT", Weight: 3},  // Collective weight for less prominent states
		{Name: "WY", Weight: 3},  // Collective weight for less prominent states
		{Name: "ID", Weight: 3},  // Collective weight for less prominent states
		{Name: "NV", Weight: 3},  // Collective weight for less prominent states
		{Name: "TN", Weight: 3},  // Collective weight for less prominent states
		{Name: "NC", Weight: 3},  // Collective weight for less prominent states
		{Name: "TX", Weight: 3},  // Collective weight for less prominent states
		{Name: "KY", Weight: 1},  // Collective weight for less prominent states
		{Name: "WV", Weight: 1},  // Collective weight for less prominent states
		{Name: "SC", Weight: 1},  // Collective weight for less prominent states
		{Name: "GA", Weight: 1},  // Collective weight for less prominent states
		{Name: "AL", Weight: 1},  // Collective weight for less prominent states
		{Name: "MS", Weight: 1},  // Collective weight for less prominent states
		{Name: "FL", Weight: 1},  // Collective weight for less prominent states
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
	provinces := []structs.Locale{
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
	normalizedState := strings.Title(strings.ToLower(strings.TrimSpace(state)))

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
		return PickFromStringList([]string{Enforcer, Grinder, Playmaker, Power, Sniper, "Two-Way"})
	} else if pos == Defender {
		return PickFromStringList([]string{"Defensive", Enforcer, "Offensive", "Two-Way"})
	} else if pos == Goalie {
		return PickFromStringList([]string{"Stand-Up", "Hybrid", "Butterfly"})
	}
	return ""
}

func GeneratePotential(pos, arch, attr string) uint8 {
	mean := 50
	if pos == Center || pos == Forward {
		if arch == Enforcer && (attr == Agility || attr == Strength || attr == PuckHandling) {
			mean += 10
		}
		if arch == Enforcer && (attr == SlapshotPower || attr == WristShotPower) {
			mean -= 10
		}
		if arch == Grinder && (attr == BodyChecking || attr == Strength || attr == StickChecking || attr == Passing) {
			mean += 10
		}
		if arch == Grinder && (attr == SlapshotPower || attr == SlapshotAccuracy || attr == WristShotPower || attr == WristShotAccuracy || attr == PuckHandling) {
			mean -= 10
		}
		if arch == Playmaker && (attr == Passing || attr == PuckHandling) {
			mean += 10
		}
		if arch == Playmaker && (attr == Strength) {
			mean -= 10
		}
		if arch == Power && (attr == SlapshotPower || attr == Strength) {
			mean += 10
		}
		if arch == Power && (attr == BodyChecking || attr == StickChecking || attr == WristShotPower) {
			mean -= 10
		}
		if arch == Sniper && (attr == WristShotPower || attr == WristShotAccuracy || attr == Passing) {
			mean += 10
		}
		if arch == Sniper && (attr == BodyChecking || attr == StickChecking || attr == WristShotPower) {
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
		if arch == Defensive && (attr == WristShotAccuracy || attr == WristShotPower || attr == SlapshotAccuracy || attr == SlapshotPower || attr == PuckHandling) {
			mean -= 10
		}
		if arch == Offensive && (attr == WristShotAccuracy || attr == WristShotPower || attr == Passing || attr == StickChecking || attr == PuckHandling) {
			mean += 10
		}
		if arch == Offensive && (attr == BodyChecking || attr == Strength || attr == SlapshotAccuracy || attr == SlapshotPower || attr == ShotBlocking) {
			mean -= 10
		}
		if arch == TwoWay && (attr == BodyChecking || attr == Passing || attr == StickChecking) {
			mean += 10
		}
		if arch == TwoWay && (attr == PuckHandling || attr == Agility || attr == SlapshotAccuracy || attr == SlapshotPower) {
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
	vDiceRoll := GenerateIntFromRange(1, 10000)
	chance := getVenerableChance(pos)
	if vDiceRoll < chance {
		venerable = true
	}

	mean, stddev := getPositionMean(pos, venerable)

	age := GenerateNormalizedIntFromMeanStdev(mean, stddev)
	return int(age)
}

func getPositionMean(pos string, venerable bool) (float64, float64) {
	meanMap := getPositionMeanMap()
	return meanMap[venerable][pos][0], meanMap[venerable][pos][1]
}

func getPositionMeanMap() map[bool]map[string][]float64 {
	return map[bool]map[string][]float64{
		true: {
			Center:   []float64{39, 2},
			Forward:  []float64{32, 1},
			Defender: []float64{35, 1},
			Goalie:   []float64{35, 1},
		},
		false: {
			Center:   []float64{32, 2},
			Forward:  []float64{26, 0.67},
			Defender: []float64{26, 0.67},
			Goalie:   []float64{29, 1.33},
		},
	}
}

func getVenerableChance(pos string) int {
	return 10
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
