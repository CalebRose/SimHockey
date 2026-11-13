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
	if ts.RunCron && ts.CanRunAIFunctionsForCollegeRecruiting() {
		managers.FillAIRecruitingBoards()
	}

	if ts.RunCron && ts.CanRunAIFunctionsForTransferPortal() {
		managers.AICoachFillBoardsPhase()
	}
}

func SyncAIBoardsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CanRunAIFunctionsForCollegeRecruiting() {
		managers.AllocatePointsToAIBoards()
	}

	if ts.RunCron && ts.CanRunAIFunctionsForTransferPortal() {
		// Portal Stuff
		managers.AICoachAllocateAndPromisePhase()
	}
}

func SyncRecruitingViaCron() {
	ts := managers.GetTimestamp()
	if !ts.RunCron {
		return
	}
	if ts.CanRefillCHLRosters() {
		managers.RefillCHLRosters()
		managers.ResetScoutingProfilesForRecruiting()
	}
	if ts.CanRunRecruiting() {
		// Sync Recruiting
		managers.SyncCollegeRecruiting()
	} else if ts.CanSyncCollegePromises() {
		managers.SyncPromises()
		managers.ProcessTransferIntention()
	} else if ts.CanOpenTransferPortal() {
		managers.EnterTheTransferPortal()
	} else if ts.CanRunTransferPortal() {
		// Sync Transfer Portal
		managers.SyncTransferPortal()
		if ts.TransferPortalRound >= 10 {
			// Portal is closed, create a preseason ranking
			managers.CreatePreseasonRanking()
		}
	}
}

func SyncFreeAgencyViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.IsFreeAgencyOpen() {
		managers.SyncAIOffers()
		managers.SyncFreeAgencyOffers()
		managers.AllocateCapsheets()
		if ts.Week > 19 {
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
	// Preseason Prep: Select optimal systems and run lineups for AI college teams
	if ts.RunCron && ts.Phase == 10 {
		managers.SelectOffensiveAndDefensiveSystemsForAllTeams_Offseason()
		managers.RunLineupsForAICollegeTeams()
	}
	if ts.RunCron && ts.IsCollegeSeasonActive() {
		managers.RunLineupsForAICollegeTeams()
	}
	if ts.RunCron && ts.IsProfessionalSeasonActive() {
		managers.RunLineupsForAIProTeams()
	}
}

func RunTheGamesViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CanRunGames() {
		// Run the Games!
		managers.RunGames()
	}
}

func ShowResultsViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.CanRunGames() {
		// Reveal Results+
		managers.ShowGames()
	}
}

func RunPostSeasonMigrationViaCron() {
	ts := managers.GetTimestamp()
	if ts.RunCron && ts.IsBothSeasonsCompleted() {
		managers.HandlePostSeasonMigration()
	}
}

func StreamCHLGamesToInterfaceViaCron() {
	managers.StartCHLLiveStreamingCron()
}

func StreamPHLGamesToInterfaceViaCron() {
	managers.StartPHLLiveStreamingCron()
}
