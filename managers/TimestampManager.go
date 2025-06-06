package managers

import (
	"strconv"

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
	UpdateSeasonStats(ts, gameDay)
	ts.ToggleGames(gameDay)
	repository.SaveTimestamp(ts, db)
}

func MoveUpWeek() structs.Timestamp {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	if ts.Week < 21 || !ts.IsOffSeason {
		ResetCollegeStandingsRanks()
	}

	// Sync to Next Week
	RecoverPlayers()
	ts.SyncToNextWeek()

	if ts.Week < 21 && !ts.CollegeSeasonOver && !ts.IsOffSeason && !ts.IsPreseason {
		SyncCollegePollSubmissionForCurrentWeek(uint(ts.Week), uint(ts.WeekID), uint(ts.SeasonID))
		ts.TogglePollRan()
	}
	if ts.Week > 15 {
		// SyncExtensionOffers()
		AllocateCapsheets()
	}
	if ts.CollegeSeasonOver && ts.NHLSeasonOver {
		ts.MoveUpSeason()
		// Run Progressions
		if !ts.ProgressedCollegePlayers {

		}
		if !ts.ProgressedProfessionalPlayers {

		}
		//
	}
	repository.SaveTimestamp(ts, db)

	return ts
}

func ResetCollegeStandingsRanks() {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	seasonID := strconv.Itoa(int(ts.SeasonID))
	db.Model(&structs.CollegeStandings{}).Where("season_id = ?", seasonID).Updates(structs.CollegeStandings{Rank: 0})
}
