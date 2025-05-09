package managers

import (
	"strconv"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
	"gorm.io/gorm"
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

	proposal := repository.FindTradeProposalRecord(repository.TradeClauses{PreloadTradeOptions: true}, proposalID)

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

	repository.CreateNewsLog(newsLog, db)
	repository.SaveTradeProposalRecord(proposal, db)
}

func RejectTradeProposal(proposalID string) {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()

	proposal := repository.FindTradeProposalRecord(repository.TradeClauses{PreloadTradeOptions: false}, proposalID)

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

	proposal := repository.FindTradeProposalRecord(repository.TradeClauses{PreloadTradeOptions: true}, proposalID)
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
	proposals := repository.FindAllTradeProposalsRecords(repository.TradeClauses{PreloadTradeOptions: true})
	return MakeTradeProposalMap(proposals)
}

func SyncAcceptedTrade(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := repository.FindTradeProposalRecord(repository.TradeClauses{PreloadTradeOptions: true}, proposalID)
	SentOptions := proposal.TeamTradeOptions

	proTeamMap := GetProTeamMap()
	capsheetMap := GetProCapsheetMap()
	contractMap := GetContractMap()

	syncAcceptedOptions(db, SentOptions, proposal.TeamID, proposal.RecepientTeamID, proTeamMap, capsheetMap, contractMap)

	proposal.ToggleSyncStatus()

	repository.SaveTradeProposalRecord(proposal, db)
}

func syncAcceptedOptions(db *gorm.DB, options []structs.TradeOption, senderID uint, recepientID uint, proTeamMap map[uint]structs.ProfessionalTeam, capsheetMap map[uint]structs.ProCapsheet, contractMap map[uint]structs.ProContract) {
	sendingTeam := proTeamMap[senderID]
	receivingTeam := proTeamMap[recepientID]
	SendersCapsheet := capsheetMap[senderID]
	recepientCapsheet := capsheetMap[recepientID]
	salaryMinimum := 0.5
	for _, option := range options {
		// Contract
		percentage := option.SalaryPercentage
		if option.PlayerID > 0 {
			playerRecord := repository.FindProPlayer(strconv.Itoa(int(option.PlayerID)))
			contract := contractMap[playerRecord.ID]
			if playerRecord.TeamID == uint16(senderID) {
				sendersPercentage := percentage * 0.01
				receiversPercentage := (100 - percentage) * 0.01
				sendingTeamPay := float64(contract.Y1BaseSalary) * sendersPercentage
				receivingTeamPay := float64(contract.Y1BaseSalary) * receiversPercentage
				// If a team is eating the Y1 Salary for a player
				if sendersPercentage == 1 {
					sendingTeamPay -= salaryMinimum
					receivingTeamPay += salaryMinimum
				} else if contract.Y1BaseSalary == 0 {
					sendingTeamPay -= salaryMinimum
					receivingTeamPay += salaryMinimum
				}
				SendersCapsheet.SubtractFromCapsheetViaTrade(contract)
				SendersCapsheet.NegotiateSalaryDifference(contract.Y1BaseSalary, float32(sendingTeamPay))
				recepientCapsheet.AddContractViaTrade(contract, float32(receivingTeamPay))
				playerRecord.TradePlayer(recepientID, receivingTeam.Abbreviation)
				contract.TradePlayer(recepientID, receivingTeam.Abbreviation, float32(receiversPercentage))
			} else {
				receiversPercentage := percentage * 0.01
				sendersPercentage := (100 - percentage) * 0.01
				sendingTeamPay := float64(contract.Y1BaseSalary) * sendersPercentage
				receivingTeamPay := float64(contract.Y1BaseSalary) * receiversPercentage
				if sendersPercentage == 1 {
					receivingTeamPay -= salaryMinimum
					sendingTeamPay += salaryMinimum
				} else if contract.Y1BaseSalary == 0 {
					receivingTeamPay -= salaryMinimum
					sendingTeamPay += salaryMinimum
				}
				recepientCapsheet.SubtractFromCapsheetViaTrade(contract)
				recepientCapsheet.NegotiateSalaryDifference(contract.Y1BaseSalary, float32(receivingTeamPay))
				SendersCapsheet.AddContractViaTrade(contract, float32(sendingTeamPay))
				playerRecord.TradePlayer(senderID, sendingTeam.Abbreviation)
				contract.TradePlayer(senderID, sendingTeam.Abbreviation, float32(sendersPercentage))
			}

			repository.SaveProPlayerRecord(playerRecord, db)
			repository.SaveProContractRecord(contract, db)

		} else if option.DraftPickID > 0 {
			draftPick := repository.FindDraftPickRecord(strconv.Itoa(int(option.DraftPickID)))
			if draftPick.TeamID == senderID {
				draftPick.TradePick(recepientID, receivingTeam.Abbreviation)
			} else {
				draftPick.TradePick(senderID, sendingTeam.Abbreviation)
			}

			repository.SaveDraftPickRecord(draftPick, db)
		}

		db.Delete(&option)
	}
	repository.SaveProCapsheetRecord(SendersCapsheet, db)
	repository.SaveProCapsheetRecord(recepientCapsheet, db)
}

func VetoTrade(proposalID string) {
	db := dbprovider.GetInstance().GetDB()

	proposal := repository.FindTradeProposalRecord(repository.TradeClauses{PreloadTradeOptions: true}, proposalID)
	SentOptions := proposal.TeamTradeOptions

	deleteOptions(db, SentOptions)

	db.Delete(&proposal)
}

func deleteOptions(db *gorm.DB, options []structs.TradeOption) {
	// Delete Recepient Trade Options
	for _, option := range options {
		// Deletes the option
		db.Delete(&option)
	}
}

func RemoveRejectedTrades() {
	db := dbprovider.GetInstance().GetDB()

	rejectedProposals := repository.FindAllTradeProposalsRecords(repository.TradeClauses{IsRejected: true, PreloadTradeOptions: true})

	for _, proposal := range rejectedProposals {
		sentOptions := proposal.TeamTradeOptions
		deleteOptions(db, sentOptions)

		// Delete Proposal
		db.Delete(&proposal)
	}
}
