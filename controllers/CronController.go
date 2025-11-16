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
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver && !ts.IsPreseason {
		managers.FillAIRecruitingBoards()
	}

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.TransferPortalPhase == 3 {
		managers.AICoachFillBoardsPhase()
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver && !ts.IsPreseason {
		managers.AllocatePointsToAIBoards()
	}

	if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.TransferPortalPhase == 3 {
		// Portal Stuff
		managers.AICoachAllocateAndPromisePhase()
	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.CollegeSeasonOver && !ts.IsPreseason && !ts.IsOffSeason {
		// Sync Recruiting
		managers.SyncCollegeRecruiting()
	} else if ts.RunCron && ts.IsOffSeason && ts.CollegeSeasonOver && ts.Week > 24 && ts.TransferPortalPhase == 3 {
		// Sync Transfer Portal
		// managers.SyncTransferPortal()
	}
}

func SyncFreeAgencyViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		managers.SyncAIOffers()
		managers.SyncFreeAgencyOffers()
		managers.AllocateCapsheets()
		if ts.Week > 17 {
			managers.SyncExtensionOffers()
		}
	}
}

func SyncToNextWeekViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron {
		// Move up Week
		managers.MoveUpWeek()
	}
}

func RunAIGameplanViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && !ts.CollegeSeasonOver {
		managers.RunLineupsForAICollegeTeams()
	}
	if ts.RunCron && !ts.IsOffSeason && !ts.NHLSeasonOver {
		managers.RunLineupsForAIProTeams()
	}
}

func RunTheGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && !ts.IsOffSeason && ts.RunGames {
		// Run the Games!
		managers.RunGames()
	}
}

func ShowResultsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.RunGames && !ts.IsOffSeason {
		// Reveal Results+
		managers.ShowGames()
	}
}

func RunPostSeasonMigrationViaCron() {
	managers.HandlePostSeasonMigration()
}
