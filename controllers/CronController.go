package controllers

import (
	"fmt"

	"github.com/CalebRose/SimHockey/managers"
)

func CronTest() {
	fmt.Println("PING!")
}

func FillAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {

	}

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.TransferPortalPhase == 3 {

	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {

	}

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.TransferPortalPhase == 3 {

	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.CollegeSeasonOver && ts.Week > 0 && ts.Week < 18 {
		// Sync Recruiting
	} else if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.Week > 24 && ts.TransferPortalPhase == 3 {
		// Sync Transfer Portal
	}
}

func SyncFreeAgencyViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.SyncFreeAgencyOffers()
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		// Move up Week
	}
}

func RunAIGameplanViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {

	}
}

func RunTheGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && ts.RunGames {
		// Run the Games!

	}
}

func ShowResultsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.RunGames && !ts.IsOffSeason {
		// Reveal Results+
	}
}
