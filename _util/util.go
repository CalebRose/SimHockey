package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func GenerateIntFromRange(min int, max int) int {
	diff := max - min + 1
	if diff < 0 {
		diff = 1
	}
	return rand.Intn(diff) + min
}

func CoinFlip() int {
	return rand.Intn(2) + 1
}

func DiceRoll(mod, requirement float64) bool {
	dr := GenerateFloatFromRange(1, 20) + mod
	return dr >= requirement
}

func GenerateFloatFromRange(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GenerateNormalizedIntFromRange(min int, max int) int {
	mean := float64(min+max) / 2.0
	stdDev := float64(max-min) / 6.0 // This approximates the 3-sigma rule

	for {
		// Generate a number using normal distribution around the mean
		num := rand.NormFloat64()*stdDev + mean
		// Round to nearest integer and convert to int type
		intNum := int(num + 0.5) // Adding 0.5 before truncating simulates rounding
		// Check if the generated number is within bounds
		if intNum >= min && intNum <= max {
			return intNum
		}
	}
}

func GenerateNormalizedIntFromMeanStdev(mean, stdDev float64) float64 {
	num := rand.NormFloat64()*stdDev + mean
	// Round to nearest integer and convert to int type
	intNum := int(num + 0.5) // Adding 0.5 before truncating simulates rounding
	return float64(intNum)
}

func PickFromStringList(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[rand.Intn(len(list))]
}

func ConvertStringToInt(num string) int {
	val, err := strconv.Atoi(num)
	if err != nil {
		log.Fatalln("Could not convert string to int")
	}

	return val
}

func ConvertStringToFloat(num string) float64 {
	floatNum, error := strconv.ParseFloat(num, 64)
	if error != nil {
		fmt.Println("Could not convert string to float 64, resetting as 0.")
		return 0
	}
	return floatNum
}

// Reads specific CSV values as Boolean. If the value is "0" or "FALSE" or "False", it will be read as false. Anything else is considered True.
func ConvertStringToBool(str string) bool {
	if str == "NULL" || str == "0" || str == "FALSE" || str == "False" {
		return false
	}
	return true
}

func ConvertFloatToString(num float64) string {
	return fmt.Sprintf("%.3f", num)
}

func ReadJson(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	return content
}

func ReadCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to read input file "+path, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+path, err)
	}

	return rows
}

func GetLetterGrade(attr int, year int) string {
	if year < 3 {
		if attr > 18 {
			return "A"
		}
		if attr > 12 {
			return "B"
		}
		if attr > 7 {
			return "C"
		}
		if attr > 4 {
			return "D"
		}
		return "F"
	}
	if attr > 24 {
		return "A+"
	}
	if attr > 21 {
		return "A"
	}
	if attr > 18 {
		return "A-"
	}
	if attr > 16 {
		return "B+"
	}
	if attr > 14 {
		return "B"
	}
	if attr > 12 {
		return "B-"
	}
	if attr > 10 {
		return "C+"
	}
	if attr > 8 {
		return "C"
	}
	if attr > 6 {
		return "C-"
	}
	if attr > 4 {
		return "D"
	}
	return "F"
}

func GetPotentialGrade(potential int) string {
	if potential > 85 {
		return "A+"
	} else if potential > 75 {
		return "A"
	} else if potential > 70 {
		return "A-"
	} else if potential > 65 {
		return "B+"
	} else if potential > 60 {
		return "B"
	} else if potential > 55 {
		return "B-"
	} else if potential > 50 {
		return "C+"
	} else if potential > 40 {
		return "C"
	} else if potential > 35 {
		return "C-"
	} else if potential > 30 {
		return "D+"
	} else if potential > 25 {
		return "D"
	} else if potential > 20 {
		return "D-"
	} else {
		return "F"
	}
}

func GetCompetitivenessLabel(value int) string {
	switch value {
	case 1:
		return "Lazy"
	case 2:
		return "Slacker"
	case 3:
		return "Milquetoast"
	case 4:
		return "Needs Motivation"
	case 5:
		return "Average"
	case 6:
		return "Focused"
	case 7:
		return "Engrossed"
	case 8:
		return "Fierce"
	case 9:
		return "Cutthroat"
	default:
		return "Distracted"
	}
}

func GetTeamLoyaltyLabel(value int) string {
	switch value {
	case 1:
		return "Apathetic"
	case 2:
		return "Wavering"
	case 3:
		return "Fickle"
	case 4:
		return "Uninterested"
	case 5:
		return "Average"
	case 6:
		return "Dependable"
	case 7:
		return "Trusted"
	case 8:
		return "Devoted"
	case 9:
		return "Unwavering"
	default:
		return "Distracted"
	}
}

func GetPlaytimePreferenceLabel(value int) string {
	switch value {
	case 1:
		return "Complacent"
	case 2:
		return "Patient"
	case 3:
		return "Passive"
	case 4:
		return "Uninterested"
	case 5:
		return "Average"
	case 6:
		return "Avid"
	case 7:
		return "Driven"
	case 8:
		return "Ambitious"
	case 9:
		return "Zealous"
	default:
		return "Distracted"
	}
}

func GetFAMarketPrefLabel(value int) string {
	valueMap := map[int]string{
		1: "Average",
		2: "Close to Home",
		3: "Countrymen",
		4: "Large Market",
		5: "No Large Market",
		6: "Small Market",
		7: "No Small Market",
		8: "Loyal to Team",
		9: "Avoiding Prev. Team",
	}
	return valueMap[value]
}

func GetFACompetitivePrefLabel(value int) string {
	valueMap := map[int]string{
		1: "Average",
		2: "Seeking Mentorship",
		3: "Veteran Mentor",
		4: "First Line",
		5: "Second Line",
		6: "Competitive Team",
	}
	return valueMap[value]
}

func GetFAFinancialPrefLabel(value int) string {
	valueMap := map[int]string{
		1: "Average",
		2: "Short Contract",
		3: "Long Contract",
		4: "Large AAV",
	}
	return valueMap[value]
}

func GetRoundAbbreviation(str string) string {
	if str == "1" {
		return "1st Round"
	} else if str == "2" {
		return "2nd Round"
	} else if str == "3" {
		return "3rd Round"
	} else if str == "4" {
		return "4th Round"
	} else if str == "5" {
		return "5th Round"
	} else if str == "6" {
		return "6th Round"
	}
	return "7th Round"
}
