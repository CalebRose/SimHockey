package managers

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/CalebRose/SimHockey/engine"
	"github.com/CalebRose/SimHockey/structs"
)

func TestWriteBoxScoreIncludesArchetypes(t *testing.T) {
	player := func(position, name, archetype string) *engine.GamePlayer {
		return &engine.GamePlayer{
			BasePlayer: structs.BasePlayer{Position: position, FirstName: name, Archetype: archetype},
			Stats:      engine.PlayerStatsDTO{Goals: 1},
		}
	}
	game := engine.GameState{
		HomeStrategy: engine.GamePlaybook{
			Forwards:  []engine.LineStrategy{{Players: []*engine.GamePlayer{player("F", "HomeForward", "Sniper")}}},
			Defenders: []engine.LineStrategy{{Players: []*engine.GamePlayer{player("D", "HomeDefender", "Offensive")}}},
			Goalies:   []engine.LineStrategy{{Players: []*engine.GamePlayer{player("G", "HomeGoalie", "Hybrid")}}},
		},
		AwayStrategy: engine.GamePlaybook{
			Forwards:  []engine.LineStrategy{{Players: []*engine.GamePlayer{player("C", "AwayForward", "Playmaker")}}},
			Defenders: []engine.LineStrategy{{Players: []*engine.GamePlayer{player("D", "AwayDefender", "Defensive")}}},
			Goalies:   []engine.LineStrategy{{Players: []*engine.GamePlayer{player("G", "AwayGoalie", "Butterfly")}}},
		},
	}
	path := filepath.Join(t.TempDir(), "box.csv")
	if err := WriteBoxScoreFile(game, path); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"HomeForward": "Sniper", "HomeDefender": "Offensive", "HomeGoalie": "Hybrid",
		"AwayForward": "Playmaker", "AwayDefender": "Defensive", "AwayGoalie": "Butterfly",
	}
	var header []string
	for _, row := range rows {
		if len(row) == 1 && row[0] == "Injury Report" {
			break
		}
		if len(row) > 2 && row[0] == "Position" && row[1] == "Name" {
			header = row
			if header[len(header)-1] != "Archetype" {
				t.Fatalf("missing archetype header: %v", header)
			}
			continue
		}
		if len(row) < 2 {
			continue
		}
		name := row[1]
		if len(name) > 0 && name[len(name)-1] == ' ' {
			name = name[:len(name)-1]
		}
		if archetype, ok := want[name]; ok {
			if len(row) != len(header) || row[len(row)-1] != archetype {
				t.Errorf("incorrect archetype for %s: %v, header %v", name, row, header)
			}
			if row[0] != "G" && row[2] != "1" {
				t.Errorf("goal column shifted for %s: %v", name, row)
			}
			delete(want, name)
		}
	}
	if len(want) != 0 {
		t.Fatalf("players missing from export: %v", want)
	}
}
