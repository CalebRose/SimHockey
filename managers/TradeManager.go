package managers

import (
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

/*
	Need to include proposals & preferences in bootstrap call
*/

func UpdateTradePreferences(pref structs.TradePreferencesDTO) {
	db := dbprovider.GetInstance().GetDB()

	preferences := repository.FindTradePreferencesByTeamID(strconv.Itoa(int(pref.TeamID)))

	preferences.UpdatePreferences(pref)

	repository.SaveTradePreferencesRecord(preferences, db)
}

func CreateTradeProposal(TradeProposal structs.TradeProposalDTO) {
	db := dbprovider.GetInstance().GetDB()
	latestID := repository.FindLatestProposalInDB(db)

	// Create Trade Proposal Object
	proposal := structs.TradeProposal{
		TeamID:          TradeProposal.TeamID,
		RecepientTeamID: TradeProposal.RecepientTeamID,
		IsTradeAccepted: false,
		IsTradeRejected: false,
	}
	proposal.AssignID(latestID)

	repository.CreateTradeProposalRecord(db, proposal)

	// Create Trade Options
	SentTradeOptions := TradeProposal.TeamTradeOptions
	ReceivedTradeOptions := TradeProposal.RecepientTeamTradeOptions

	optionsBatch := []structs.TradeOption{}

	for _, sentOption := range SentTradeOptions {
		tradeOption := structs.TradeOption{
			TradeProposalID:  latestID,
			TeamID:           TradeProposal.TeamID,
			PlayerID:         sentOption.PlayerID,
			DraftPickID:      sentOption.DraftPickID,
			SalaryPercentage: sentOption.SalaryPercentage,
			OptionType:       sentOption.OptionType,
		}
		optionsBatch = append(optionsBatch, tradeOption)
	}

	for _, recepientOption := range ReceivedTradeOptions {
		tradeOption := structs.TradeOption{
			TradeProposalID:  latestID,
			TeamID:           TradeProposal.RecepientTeamID,
			PlayerID:         recepientOption.PlayerID,
			DraftPickID:      recepientOption.DraftPickID,
			SalaryPercentage: recepientOption.SalaryPercentage,
			OptionType:       recepientOption.OptionType,
		}
		optionsBatch = append(optionsBatch, tradeOption)
	}

	repository.CreateTradeOptionRecordsBatch(db, optionsBatch, 20)
}

func AcceptTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	proTeamMap := GetProTeamMap()

	proposal := repository.FindTradeProposalRecord(proposalID)

	proposal.AcceptTrade()

	recepientTeam := proTeamMap[proposal.RecepientTeamID]
	sendingTeam := proTeamMap[proposal.TeamID]

	// Create News Log
	newsLogMessage := recepientTeam.Abbreviation + " has accepted a trade offer from " + sendingTeam.Abbreviation + " for trade the following players:\n\n"

	for _, options := range proposal.TeamTradeOptions {
		if options.TeamID == proposal.TeamID {
			if options.PlayerID > 0 {
				playerRecord := repository.FindProPlayer(strconv.Itoa(int(options.PlayerID)))
				ovrGrade := strconv.Itoa(int(playerRecord.Overall))
				ovr := playerRecord.Overall
				if playerRecord.Year > 1 {
					newsLogMessage += playerRecord.Position + " " + strconv.Itoa(int(ovr)) + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + recepientTeam.Abbreviation + "\n"
				} else {
					newsLogMessage += playerRecord.Position + " " + ovrGrade + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + recepientTeam.Abbreviation + "\n"
				}
			} else if options.DraftPickID > 0 {
				draftPick := repository.FindDraftPickRecord(strconv.Itoa(int(options.DraftPickID)))
				pickRound := strconv.Itoa(int(draftPick.DraftRound))
				roundAbbreviation := util.GetRoundAbbreviation(pickRound)
				season := strconv.Itoa(int(draftPick.Season))
				newsLogMessage += season + " " + roundAbbreviation + " pick to " + recepientTeam.Abbreviation + "\n"
			}
		} else {
			if options.PlayerID > 0 {
				playerRecord := repository.FindProPlayer(strconv.Itoa(int(options.PlayerID)))
				newsLogMessage += playerRecord.Position + " " + playerRecord.FirstName + " " + playerRecord.LastName + " to " + sendingTeam.Abbreviation + "\n"
			} else if options.DraftPickID > 0 {
				draftPick := repository.FindDraftPickRecord(strconv.Itoa(int(options.DraftPickID)))
				pickRound := strconv.Itoa(int(draftPick.DraftRound))
				roundAbbreviation := util.GetRoundAbbreviation(pickRound)
				season := strconv.Itoa(int(draftPick.Season))
				newsLogMessage += season + " " + roundAbbreviation + " pick to " + sendingTeam.Abbreviation + "\n"
			}
		}
	}
	newsLogMessage += "\n"

	newsLog := structs.NewsLog{
		TeamID:      0,
		WeekID:      ts.WeekID,
		Week:        ts.Week,
		SeasonID:    ts.SeasonID,
		League:      "",
		MessageType: "Trade",
		Message:     newsLogMessage,
	}

	db.Create(&newsLog)
	db.Save(&proposal)
}

func RejectTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	proposal := repository.FindTradeProposalRecord(proposalID)

	proTeamMap := GetProTeamMap()

	recepientTeam := proTeamMap[proposal.RecepientTeamID]
	sendingTeam := proTeamMap[proposal.TeamID]

	proposal.RejectTrade()
	newsLog := structs.NewsLog{
		TeamID:      0,
		WeekID:      ts.WeekID,
		Week:        ts.Week,
		SeasonID:    ts.SeasonID,
		League:      "",
		MessageType: "Trade",
		Message:     recepientTeam.Abbreviation + " has rejected a trade from " + sendingTeam.Abbreviation,
	}

	repository.CreateNewsLog(newsLog, db)
	repository.SaveTradeProposalRecord(proposal, db)
}

func CancelTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := repository.FindTradeProposalRecord(proposalID)
	options := proposal.TeamTradeOptions

	for _, option := range options {
		db.Delete(&option)
	}

	db.Delete(&proposal)
}

func GetTradePreferencesMap() map[uint]structs.TradePreferences {
	prefs := repository.FindAllTradePreferenceRecords()
	return MakeTradePreferencesMap(prefs)
}

func GetTradeProposalsMap() map[uint][]structs.TradeProposal {
	proposals := repository.FindAllTradeProposalsRecords()
	return MakeTradeProposalMap(proposals)
}
