package managers

import (
	"strings"
	"testing"

	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
)

func TestIncompleteRostersCreateReports(t *testing.T) {
	var collegeRoster []structs.CollegePlayer
	var proRoster []structs.ProfessionalPlayer
	for position, count := range map[string]int{"C": 4, "F": 8, "D": 6, "G": 2} {
		for i := 0; i < count; i++ {
			player := structs.BasePlayer{Position: position}
			collegeRoster = append(collegeRoster, structs.CollegePlayer{BasePlayer: player})
			proRoster = append(proRoster, structs.ProfessionalPlayer{BasePlayer: player})
		}
	}

	collegeTeam := structs.CollegeTeam{Model: gorm.Model{ID: 7}, BaseTeam: structs.BaseTeam{Abbreviation: "COL"}}
	proTeam := structs.ProfessionalTeam{Model: gorm.Model{ID: 8}, BaseTeam: structs.BaseTeam{Abbreviation: "PRO"}}
	var league, message string
	var teamID uint
	calls := 0
	report := func(l, m string, id uint) {
		league, message, teamID = l, m, id
		calls++
	}

	if !validateCollegeRoster(collegeTeam, collegeRoster, false, false, report) ||
		!validateProRoster(proTeam, proRoster, false, false, report) || calls != 0 {
		t.Fatal("complete rosters must not create reports")
	}

	for i := range collegeRoster {
		if collegeRoster[i].Position == "C" {
			collegeRoster[i].IsRedshirt = true
			break
		}
	}
	if validateCollegeRoster(collegeTeam, collegeRoster, false, false, report) ||
		calls != 1 || league != "CHL" || teamID != 7 || !strings.Contains(message, "1 centers") {
		t.Errorf("unexpected college report: calls=%d league=%s team=%d message=%q", calls, league, teamID, message)
	}

	for i := range proRoster {
		if proRoster[i].Position == "G" {
			proRoster[i].IsInjured = true
			break
		}
	}
	if validateProRoster(proTeam, proRoster, false, false, report) ||
		calls != 2 || league != "PHL" || teamID != 8 || !strings.Contains(message, "1 goalies") {
		t.Errorf("unexpected pro report: calls=%d league=%s team=%d message=%q", calls, league, teamID, message)
	}
}
