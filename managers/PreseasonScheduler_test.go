package managers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/CalebRose/SimHockey/structs"
)

func preseasonTestTeams(count int) map[uint]structs.CollegeTeam {
	teams := make(map[uint]structs.CollegeTeam, count)
	for i := 1; i <= count; i++ {
		teams[uint(i)] = structs.CollegeTeam{BaseTeam: structs.BaseTeam{
			TeamName:     fmt.Sprintf("Team %d", i),
			ConferenceID: uint8((i-1)%4 + 1),
		}}
	}
	return teams
}

func preseasonTestGame(home, away uint, day string) structs.CollegeGame {
	return structs.CollegeGame{BaseGame: structs.BaseGame{
		HomeTeamID: home, AwayTeamID: away, GameDay: day,
	}}
}

func TestGenerateCollegePreseasonSchedule(t *testing.T) {
	ts := structs.Timestamp{SeasonID: 25, WeekID: 2500, Week: 0}
	teams := preseasonTestTeams(64)
	existing := []structs.CollegeGame{
		preseasonTestGame(1, 2, "A"),
		preseasonTestGame(3, 4, "B"),
		preseasonTestGame(1, 5, "C"), // Requested conference games must also be preserved.
		preseasonTestGame(6, 1000, "C"),
		preseasonTestGame(7, 1001, "C"),
	}
	games, err := generateCollegePreseasonSchedule(ts, teams, existing)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 92 {
		t.Fatalf("got %d generated games, want 92", len(games))
	}

	occupied := make(map[string]map[uint]bool)
	opponents := make(map[[2]uint]bool)
	for _, day := range []string{"A", "B", "C"} {
		occupied[day] = make(map[uint]bool)
	}
	for _, game := range append(existing, games...) {
		for _, id := range []uint{game.HomeTeamID, game.AwayTeamID} {
			if occupied[game.GameDay][id] {
				t.Fatalf("team %d scheduled twice in %s", id, game.GameDay)
			}
			occupied[game.GameDay][id] = true
		}
		key := [2]uint{game.HomeTeamID, game.AwayTeamID}
		if key[0] > key[1] {
			key[0], key[1] = key[1], key[0]
		}
		if opponents[key] {
			t.Fatalf("repeated pairing %v", key)
		}
		opponents[key] = true
	}
	for _, day := range []string{"A", "B", "C"} {
		for id := range teams {
			if !occupied[day][id] {
				t.Errorf("team %d missing game in %s", id, day)
			}
		}
	}
	for _, game := range games {
		if !game.IsPreseason || game.SeasonID != ts.SeasonID || game.WeekID != ts.WeekID || game.Week != int(ts.Week) {
			t.Errorf("incorrect game metadata: %+v", game.BaseGame)
		}
	}
}

func TestGenerateCollegePreseasonScheduleExistingGames(t *testing.T) {
	ts := structs.Timestamp{}
	teams := preseasonTestTeams(4)
	tests := []struct {
		name     string
		existing []structs.CollegeGame
		wantErr  string
	}{
		{
			name: "conflicting requests",
			existing: []structs.CollegeGame{
				preseasonTestGame(1, 2, "A"),
				preseasonTestGame(3, 1, "A"),
			},
			wantErr: "already scheduled more than once",
		},
		{
			name:     "unmatchable slot",
			existing: []structs.CollegeGame{preseasonTestGame(1, 100, "A")},
			wantErr:  "unable to complete",
		},
		{
			name: "fully scheduled",
			existing: []structs.CollegeGame{
				preseasonTestGame(1, 2, "A"), preseasonTestGame(3, 4, "A"),
				preseasonTestGame(1, 3, "B"), preseasonTestGame(2, 4, "B"),
				preseasonTestGame(1, 4, "C"), preseasonTestGame(2, 3, "C"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games, err := generateCollegePreseasonSchedule(ts, teams, tt.existing)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got error %v, want %q", err, tt.wantErr)
				}
				if games != nil {
					t.Fatal("returned partial schedule on error")
				}
				return
			}
			if err != nil || games == nil || len(games) != 0 {
				t.Fatalf("got games %v and error %v; want an empty successful schedule", games, err)
			}
		})
	}
}

func TestGenerateCollegePreseasonScheduleSameConference(t *testing.T) {
	for _, conference := range []uint8{1, 7} {
		t.Run(fmt.Sprintf("conference-%d", conference), func(t *testing.T) {
			teams := preseasonTestTeams(4)
			for id, team := range teams {
				team.ConferenceID = conference
				teams[id] = team
			}
			games, err := generateCollegePreseasonSchedule(structs.Timestamp{}, teams, nil)
			if err != nil || len(games) != 6 {
				t.Fatalf("got %d games and error %v, want 6 games", len(games), err)
			}
		})
	}
}

func TestGenerateCollegePreseasonScheduleReservesFutureOpponents(t *testing.T) {
	teams := preseasonTestTeams(4)
	existing := []structs.CollegeGame{
		preseasonTestGame(1, 2, "A"),
		preseasonTestGame(1, 3, "C"),
	}
	games, err := generateCollegePreseasonSchedule(structs.Timestamp{}, teams, existing)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 4 {
		t.Fatalf("got %d games, want 4", len(games))
	}
	for _, game := range games {
		if game.GameDay == "B" && (game.HomeTeamID == 1 || game.AwayTeamID == 1) {
			if game.HomeTeamID != 4 && game.AwayTeamID != 4 {
				t.Fatalf("team 1 must play team 4 in B: %+v", game.BaseGame)
			}
			return
		}
	}
	t.Fatal("team 1 was not scheduled in B")
}
