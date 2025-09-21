package managers

import (
	"strconv"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
)

func CutCHLPlayer(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetCollegePlayerByID(playerId)
	player.WillTransfer()
	ts := GetTimestamp()
	if ts.IsOffSeason || ts.Week <= 1 || ts.Week >= 24 || ts.TransferPortalPhase == 3 {
		previousTeamID := strconv.Itoa(int(player.PreviousTeamID))
		deduction := 0
		promiseDeduction := 0
		if player.Stars > 2 {
			deduction = int(player.Stars) / 2
		}
		collegePromise := GetCollegePromiseByCollegePlayerID(strconv.Itoa(int(player.ID)), previousTeamID)
		if collegePromise.IsActive && collegePromise.PromiseMade {
			weight := collegePromise.PromiseWeight
			switch weight {
			case "Vew Low":
				promiseDeduction = 3
			case "Low":
				promiseDeduction = 8
			case "Medium":
				promiseDeduction = 13
			case "High":
				promiseDeduction = 23
			case "Very High":
				promiseDeduction = 28
			}
		}

		points := (-1 * deduction) - promiseDeduction
		teamProfile := repository.FindTeamRecruitingProfile(previousTeamID, false, false)
		teamProfile.IncrementClassSize()
		if player.Stars > 0 {
			teamProfile.AdjustPortalReputation(int8(points))
			repository.SaveTeamProfileRecord(db, teamProfile)
		}
	}
	repository.SaveCollegeHockeyPlayerRecord(player, db)
}

func RedshirtCHLPlayer(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := GetCollegePlayerByID(playerId)
	if player.IsRedshirt || player.IsRedshirting {
		return
	}
	player.RedshirtPlayer()
	repository.SaveCollegeHockeyPlayerRecord(player, db)
}

func CutProPlayer(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := repository.FindProPlayer(playerId)
	contract := repository.FindProContract(playerId)
	capsheetMap := GetProCapsheetMap()
	capsheet := capsheetMap[uint(player.TeamID)]
	ts := GetTimestamp()

	if player.Year < 4 && !ts.IsOffSeason && !player.IsAffiliatePlayer {
		player.WaivePlayer()
	} else {
		player.ToggleIsFreeAgent()
		contract.CutContract()
	}

	capsheet.CutPlayerFromCapsheet(contract)
	repository.SaveProContractRecord(contract, db)
	repository.SaveProPlayerRecord(player, db)
	repository.SaveProCapsheetRecord(capsheet, db)
}

func SendPHLPlayerToAffiliate(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := repository.FindProPlayer(playerId)
	player.ToggleAffiliation()
	repository.SaveProPlayerRecord(player, db)
}

func SendPHLPlayerToTradeBlock(playerId string) {
	db := dbprovider.GetInstance().GetDB()

	player := repository.FindProPlayer(playerId)
	player.ToggleTradeBlock()
	repository.SaveProPlayerRecord(player, db)
}
