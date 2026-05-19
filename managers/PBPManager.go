package managers

import (
	"fmt"

	"github.com/CalebRose/SimHockey/structs"
)

// GeneratePlayByPlayText unifies PBP generation for both Pro and College.
// It accepts 'any' for playerMap to support both CollegePlayer and ProfessionalPlayer maps.
func GeneratePlayByPlayText(play structs.PbP, event, outcome string, playerMap any, team string) string {
	// This currently returns a basic string to ensure your project compiles.
	// You can expand this with your switch statements here once you are ready
	// to consolidate the logic from your old generateResultsString functions.
	return fmt.Sprintf("%s - %s", event, outcome)
}
