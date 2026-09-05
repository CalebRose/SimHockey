package managers

import (
	"encoding/csv"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

// =============================================================================
// Constants
// =============================================================================

const minutePositionFloor uint16 = 5

// starterSnapThreshold is the minimum snaps-per-season to qualify as a starter
// for playtime tag type determination.
const starterMinuteThreshold = 13

const goalieGamesPlayedFloor uint16 = 10
const goalieStarterGamesPlayedThreshold uint16 = 40

// =============================================================================
// Internal types
// =============================================================================

type playerGroupEntry struct {
	Player   structs.ProfessionalPlayer
	Contract structs.ProContract
	Stats    structs.ProfessionalPlayerSeasonStats
}

type rankedGroupEntry struct {
	Entry      playerGroupEntry
	AdjOverall int // overall with +4 adjustment for players <= age 25
}

// getTopTierCount returns the number of players in the "elite" tier for a group.
func getTopTierCount(group string) int {
	switch group {
	case "F":
		return 10
	case "D":
		return 8
	default:
		return 5
	}
}

// getAgeExclusionThreshold returns the minimum age (inclusive) at which a
// player is excluded from the custom mid-tier comparison set (step 5).
func getAgeExclusionThreshold(group string) int {
	switch group {
	case "G":
		return 30
	case "F", "D":
		return 30
	case "C":
		return 30
	}
	return 29
}

// getAgeAdjustmentFactor returns the decimal adjustment factor for a player's
// age and position group (e.g., +0.15 means +15%, -0.10 means −10%).
func getAgeAdjustmentFactor(age int) float64 {
	if age >= 34 {
		return -0.90
	}
	factors := map[int]float64{
		23: 0.15, 24: 0.10, 25: 0.05, 26: 0.00,
		27: -0.10, 28: -0.20, 29: -0.30, 30: -0.40,
		31: -0.50, 32: -0.60, 33: -0.75, 34: -0.90,
		35: -1.00, 36: -1.10, 37: -1.20, 38: -1.30, 39: -1.40,
		40: -1.50,
	}
	if f, ok := factors[age]; ok {
		return f
	}
	return 0.00
}

func avgFloats(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var total float64
	for _, v := range vals {
		total += v
	}
	return total / float64(len(vals))
}

// computeMidTierValue calculates a player's raw expected value using the
// mid-tier lookup: the highest signing value among comparably-or-lesser-
// rated peers who are below the group's age exclusion threshold. Falls back
// to the best comparable regardless of age if none qualify under the
// threshold, so groups whose veteran pool skews old (e.g. Centers) don't
// bottom out just because no contract passed the age filter.
func computeMidTierValue(actualOverall, ageThreshold int, entries []playerGroupEntry) float64 {
	var bestValue, bestValueAnyAge float64 = 0, 0

	for _, e := range entries {
		c := e.Contract
		if !c.IsActive || c.ContractType == "Rookie" || c.ContractType == "UDFA" {
			continue
		}
		if int(e.Player.Overall) > actualOverall {
			continue
		}
		if float64(c.SigningValue) > bestValueAnyAge {
			bestValueAnyAge = float64(c.SigningValue)
		}
		if int(e.Player.Age) >= ageThreshold {
			continue
		}
		if float64(c.SigningValue) > bestValue {
			bestValue = float64(c.SigningValue)
		}
	}

	if bestValue > 0 {
		return bestValue
	}
	return bestValueAnyAge // fallback: widen the search rather than return 0
}

// computeGroupExtensionValues calculates the expected minimum contract value
// and AAV for every player in a position group, following the steps in the
// TechDocs/extension_value_and_tag_calculations.md document.
func computeGroupExtensionValues(group string, entries []playerGroupEntry) []structs.ProfessionalPlayer {
	if len(entries) == 0 {
		return nil
	}

	topTierCount := getTopTierCount(group)
	ageThreshold := getAgeExclusionThreshold(group)

	// Step 2: compute adjusted overall (age <= 25 gets +4 bonus for ranking)
	ranked := make([]rankedGroupEntry, len(entries))
	for i, e := range entries {
		adj := int(e.Player.Overall)
		if int(e.Player.Age) <= 25 {
			adj += 4
		}
		ranked[i] = rankedGroupEntry{Entry: e, AdjOverall: adj}
	}

	// Step 3: sort by adjusted overall desc, age asc
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].AdjOverall != ranked[j].AdjOverall {
			return ranked[i].AdjOverall > ranked[j].AdjOverall
		}
		return ranked[i].Entry.Player.Age < ranked[j].Entry.Player.Age
	})

	// Determine the adjusted-overall floor for the top tier.
	// All players tied at the Nth place are included (per "tied at 5th" rule).
	topTierMinOverall := math.MinInt32 // default: all players are top tier
	if topTierCount < len(ranked) {
		topTierMinOverall = ranked[topTierCount-1].AdjOverall
	}

	// Step 4a: collect active contract signing values for top-tier reference
	var groupSigningValues []float64
	for _, e := range entries {
		if e.Contract.IsActive {
			groupSigningValues = append(groupSigningValues, float64(e.Contract.SigningValue))
		}
	}
	sort.Slice(groupSigningValues, func(i, j int) bool {
		return groupSigningValues[i] > groupSigningValues[j]
	})

	var maxContractValue float64
	if len(groupSigningValues) > 0 {
		maxContractValue = groupSigningValues[0]
		// Caruso-Bordewyk Rule: if the highest value is >= 150% of the second
		// highest, exclude it to avoid distorting the top-tier reference.
		if len(groupSigningValues) > 1 && maxContractValue >= groupSigningValues[1]*1.5 {
			maxContractValue = groupSigningValues[1]
		}
	}
	topTierExpectedValue := maxContractValue * 1.10

	// Steps 5–6: mid-tier reference value for every player (highest signing
	// value among comparably-or-lesser-rated peers). Computed for everyone,
	// including elite-tier members, since it doubles as the blend baseline
	// below (previously left at 0 for elite members, which silently broke
	// the weighted blend and let flat top-tier values leak through).
	midTierValues := make([]float64, len(ranked))
	for i, rp := range ranked {
		midTierValues[i] = computeMidTierValue(int(rp.Entry.Player.Overall), ageThreshold, entries)
	}

	// Find the real-overall range spanned by elite-tier members, so we can scale
	// each elite player's blend weight by how good they actually are, not just
	// whether their age-boosted rank cleared the cutoff.
	maxEliteOverall, minEliteOverall := 0, math.MaxInt32
	for _, rp := range ranked {
		if rp.AdjOverall >= topTierMinOverall {
			o := int(rp.Entry.Player.Overall)
			if o > maxEliteOverall {
				maxEliteOverall = o
			}
			if o < minEliteOverall {
				minEliteOverall = o
			}
		}
	}
	eliteOverallSpan := maxEliteOverall - minEliteOverall

	// Step 4: elite tier — blend max-contract-value+10% with the mid-tier
	// value, weighted by real overall position within the tier (floor of
	// 15% so tier membership is still worth something). Non-elite players
	// within a small band just below the tier's real-overall floor get a
	// tapering share of the same bonus, so values don't cliff hard at the
	// tier boundary (e.g. 5th vs 6th place separated by a single point).
	const boundaryBand = 2.0
	rawValues := make([]float64, len(ranked))
	for i, rp := range ranked {
		overall := int(rp.Entry.Player.Overall)
		weight := 0.0
		switch {
		case rp.AdjOverall >= topTierMinOverall:
			weight = 1.0
			if eliteOverallSpan > 0 {
				norm := float64(overall-minEliteOverall) / float64(eliteOverallSpan)
				weight = 0.15 + norm*0.85
			}
		case minEliteOverall != math.MaxInt32:
			dist := float64(minEliteOverall - overall)
			if dist > 0 && dist <= boundaryBand {
				weight = 0.15 * (1 - dist/boundaryBand)
			}
		}
		rawValues[i] = weight*topTierExpectedValue + (1-weight)*midTierValues[i]
	}

	// Step 7: smoothing — average values of players at (overall+1) and (overall-1)
	adjOverallValMap := make(map[int][]float64)
	for i, rp := range ranked {
		adjOverallValMap[rp.AdjOverall] = append(adjOverallValMap[rp.AdjOverall], rawValues[i])
	}

	result := make([]structs.ProfessionalPlayer, len(ranked))
	for i, rp := range ranked {
		rawVal := rawValues[i]

		// Build smoothing components from adjacent overall levels
		var smoothComponents []float64
		if higherVals, ok := adjOverallValMap[rp.AdjOverall+1]; ok {
			smoothComponents = append(smoothComponents, avgFloats(higherVals))
		}
		if lowerVals, ok := adjOverallValMap[rp.AdjOverall-1]; ok {
			smoothComponents = append(smoothComponents, avgFloats(lowerVals))
		}
		smoothedVal := avgFloats(smoothComponents)

		// Step 8: take the higher of raw and smoothed values
		bestVal := rawVal
		if smoothedVal > bestVal {
			bestVal = smoothedVal
		}

		// Step 9: apply age adjustment
		ageFactor := getAgeAdjustmentFactor(int(rp.Entry.Player.Age))
		adjustedVal := bestVal * (1.0 + ageFactor)
		if adjustedVal < 0.7 {
			adjustedVal = 0.7
		}
		// Ceiling value of 14 million
		if adjustedVal > 14 {
			adjustedVal = 14
		}

		p := rp.Entry.Player
		p.AssignCalculatedValues(adjustedVal)
		result[i] = p
	}

	return result
}

func CalculatePlayerMinimumValues(w http.ResponseWriter) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	isTest := true
	csvRows := [][]string{}

	allProPlayers := repository.FindAllProPlayers(repository.PlayerQuery{})
	// Get Last Season Stats
	seasonStats := repository.FindProPlayerSeasonStatsRecords(strconv.Itoa(int(ts.SeasonID)-1), "2")
	seasonStatMap := MakeProPlayerSeasonStatMap(seasonStats)
	proContracts := repository.FindAllProContracts(true)
	proContractMap := MakeContractMap(proContracts)

	// Group all players by their position group
	groupEntries := make(map[string][]playerGroupEntry)
	for _, p := range allProPlayers {
		pos := p.Position
		groupEntries[pos] = append(groupEntries[pos], playerGroupEntry{
			Player:   p,
			Contract: proContractMap[p.ID],
			Stats:    seasonStatMap[p.ID],
		})
	}

	// Calculate extension values per group, then save
	saved := 0

	for group, entries := range groupEntries {
		updated := computeGroupExtensionValues(group, entries)
		for _, p := range updated {
			if !isTest {
				if err := db.Model(&p).Updates(map[string]interface{}{
					"minimum_value":          p.MinimumValue,
					"original_minimum_value": p.OriginalMinimumValue,
				}).Error; err != nil {
					log.Printf("ValuationManager: failed to save player %d (%s): %v", p.ID, group, err)
					continue
				}
			} else {
				originalContract := proContractMap[p.ID]
				row := []string{
					group,
					strconv.Itoa(int(p.ID)),
					p.FirstName,
					p.LastName,
					p.Position,
					p.Archetype,
					strconv.Itoa(int(p.Age)),
					strconv.Itoa(int(p.Overall)),
					strconv.FormatFloat(float64(p.OriginalMinimumValue), 'f', 2, 64),
					strconv.FormatFloat(float64(p.MinimumValue), 'f', 2, 64),
					strconv.FormatFloat(float64(originalContract.ContractValue), 'f', 2, 64),
				}
				csvRows = append(csvRows, row)
			}
			saved++
		}
	}

	if isTest {
		ExportMinimumValueAndTagUpdatesToCSV(w, csvRows)
	}
}

func ExportMinimumValueAndTagUpdatesToCSV(w http.ResponseWriter, csvRows [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	filename := "simhck_minimum_value_and_tag_updates_test.csv"
	w.Header().Set("Content-Disposition", "attachment;filename="+filename)

	writer := csv.NewWriter(w)
	defer writer.Flush()
	HeaderRow := []string{
		"Group", "ID", "First Name", "Last Name", "Position", "Archetype",
		"Age", "Overall", "Calculated Value", "Calculated AAV", "Original Contract Value",
	}
	err := writer.Write(HeaderRow)
	if err != nil {
		log.Fatal("Cannot write header row", err)
	}

	for _, row := range csvRows {
		if err := writer.Write(row); err != nil {
			log.Fatal("Cannot write row to CSV", err)
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			log.Fatal("Error while writing to file ::", err)
		}
	}
}
