package managers

import (
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

// GetTimestamp -- Get the Timestamp
func GetTimestamp() structs.Timestamp {
	return repository.FindTimestamp()
}

func ShowGames() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	// UpdateStandings
	// Update Season Stats
	gameDay := ts.GetGameDay()
	UpdateStandings(ts, gameDay)
	// UpdateSeasonStats(ts, gameDay)
	ts.ToggleGames(gameDay)
	repository.SaveTimestamp(ts, db)
}
