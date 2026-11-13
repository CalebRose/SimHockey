package managers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	fbsvc "github.com/CalebRose/SimHockey/firebase"
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

func SyncAcceptedTrade(proposalID string) error {
	db := dbprovider.GetInstance().GetDB()

	proposal := repository.FindTradeProposalRecord(repository.TradeClauses{PreloadTradeOptions: true}, proposalID)
	if proposal.ID == 0 {
		return fmt.Errorf("trade proposal %s was not found", proposalID)
	}
	SentOptions := proposal.TeamTradeOptions
	if proposal.IsSynced {
		return fmt.Errorf("trade proposal %s is already synced", proposalID)
	}
	if len(SentOptions) == 0 {
		return fmt.Errorf("trade proposal %s has no assets to sync", proposalID)
	}

	proTeamMap := GetProTeamMap()
	capsheetMap := GetProCapsheetMap()
	contractMap := GetContractMap()
	if err := validateTradeTeams(proposal, proTeamMap, capsheetMap); err != nil {
		return err
	}

	details, complete := syncAcceptedOptions(db, SentOptions, proposal.TeamID, proposal.RecepientTeamID, proTeamMap, capsheetMap, contractMap)

	proposal.ToggleSyncStatus()

	repository.SaveTradeProposalRecord(proposal, db)
	if !complete {
		return fmt.Errorf("trade proposal %s may be partially processed; check trade logs", proposalID)
	}
	reportProcessedTrade(proposal, proTeamMap, details, false)
	return nil
}

var ErrInvalidDraftTrade = errors.New("invalid draft trade")

func SyncAcceptedDraftTrade(proposal structs.TradeProposal) error {
	if proposal.TeamID == 0 || proposal.RecepientTeamID == 0 || proposal.TeamID == proposal.RecepientTeamID {
		return fmt.Errorf("%w: two distinct teams are required", ErrInvalidDraftTrade)
	}
	options, err := draftTradeOptions(proposal)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDraftTrade, err)
	}
	if len(options) == 0 {
		return fmt.Errorf("%w: at least one asset is required", ErrInvalidDraftTrade)
	}

	db := dbprovider.GetInstance().GetDB()
	proTeamMap := GetProTeamMap()
	capsheetMap := GetProCapsheetMap()
	contractMap := GetContractMap()
	if err := validateTradeTeams(proposal, proTeamMap, capsheetMap); err != nil {
		return err
	}

	details, complete := syncAcceptedOptions(db, options, proposal.TeamID, proposal.RecepientTeamID, proTeamMap, capsheetMap, contractMap)
	if !complete {
		return fmt.Errorf("draft trade may be partially processed; check trade logs")
	}
	reportProcessedTrade(proposal, proTeamMap, details, true)
	return nil
}

func validateTradeTeams(proposal structs.TradeProposal, teams map[uint]structs.ProfessionalTeam, capsheets map[uint]structs.ProCapsheet) error {
	if proposal.TeamID == 0 || proposal.RecepientTeamID == 0 || proposal.TeamID == proposal.RecepientTeamID ||
		teams[proposal.TeamID].ID == 0 || teams[proposal.RecepientTeamID].ID == 0 ||
		capsheets[proposal.TeamID].ID == 0 || capsheets[proposal.RecepientTeamID].ID == 0 {
		return fmt.Errorf("trade teams %d and %d must exist and have capsheets", proposal.TeamID, proposal.RecepientTeamID)
	}
	return nil
}

func draftTradeOptions(proposal structs.TradeProposal) ([]structs.TradeOption, error) {
	options := append([]structs.TradeOption(nil), proposal.TeamTradeOptions...)
	indices := make(map[string]int, len(options))
	combined := make([]structs.TradeOption, 0, len(options)+len(proposal.RecepientTeamTradeOptions))
	options = append(options, proposal.RecepientTeamTradeOptions...)
	for _, option := range options {
		if (option.PlayerID == 0) == (option.DraftPickID == 0) {
			return nil, fmt.Errorf("draft trade option must identify exactly one player or draft pick")
		}
		key := fmt.Sprintf("player:%d", option.PlayerID)
		if option.DraftPickID != 0 {
			key = fmt.Sprintf("pick:%d", option.DraftPickID)
		}
		if idx, exists := indices[key]; exists {
			previous := combined[idx]
			if previous.TeamID != option.TeamID || previous.SalaryPercentage != option.SalaryPercentage ||
				(previous.ID != 0 && option.ID != 0 && previous.ID != option.ID) {
				return nil, fmt.Errorf("conflicting draft trade options for %s", key)
			}
			if previous.ID == 0 && option.ID != 0 {
				combined[idx] = option
			}
			continue
		}
		indices[key] = len(combined)
		combined = append(combined, option)
	}
	return combined, nil
}

func reportProcessedTrade(proposal structs.TradeProposal, teams map[uint]structs.ProfessionalTeam, details []string, draft bool) {
	ts := GetTimestamp()
	if ts.ID == 0 {
		log.Printf("TradeManager: cannot report trade %d: timestamp unavailable", proposal.ID)
		return
	}
	id, message, err := tradeReportIdentityAndMessage(proposal, teams, details, draft)
	if err != nil {
		log.Printf("TradeManager: cannot report trade %d: %v", proposal.ID, err)
		return
	}
	report, err := fbsvc.TradeTransactionReport(proposal.TeamID, message, ts)
	if err != nil {
		log.Printf("TradeManager: cannot report trade %d: %v", proposal.ID, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := fbsvc.CreateAdminReport(ctx, id, report); err != nil {
		log.Printf("TradeManager: failed to create report for trade %d: %v", proposal.ID, err)
	}
}

func tradeReportIdentityAndMessage(proposal structs.TradeProposal, teams map[uint]structs.ProfessionalTeam, details []string, draft bool) (string, string, error) {
	sender, senderOK := teams[proposal.TeamID]
	recipient, recipientOK := teams[proposal.RecepientTeamID]
	if !senderOK || !recipientOK || proposal.TeamID == 0 || proposal.RecepientTeamID == 0 || proposal.TeamID == proposal.RecepientTeamID {
		return "", "", fmt.Errorf("invalid trade teams %d and %d", proposal.TeamID, proposal.RecepientTeamID)
	}

	label := fmt.Sprintf("Trade proposal #%d", proposal.ID)
	id := fmt.Sprintf("trade-proposal-%d", proposal.ID)
	if draft {
		label = "Draft trade"
		if proposal.ID != 0 {
			label = fmt.Sprintf("Draft trade #%d", proposal.ID)
			id = fmt.Sprintf("draft-trade-%d", proposal.ID)
		} else {
			var randomID [16]byte
			if _, err := rand.Read(randomID[:]); err != nil {
				return "", "", fmt.Errorf("generate draft trade report ID: %w", err)
			}
			id = "draft-trade-" + hex.EncodeToString(randomID[:])
		}
	}
	message := fmt.Sprintf("%s processed between %s (team %d) and %s (team %d).\n%s",
		label, sender.Abbreviation, proposal.TeamID, recipient.Abbreviation, proposal.RecepientTeamID, strings.Join(details, "\n"))
	return id, message, nil
}

func syncAcceptedOptions(db *gorm.DB, options []structs.TradeOption, senderID uint, recepientID uint, proTeamMap map[uint]structs.ProfessionalTeam, capsheetMap map[uint]structs.ProCapsheet, contractMap map[uint]structs.ProContract) ([]string, bool) {
	sendingTeam := proTeamMap[senderID]
	receivingTeam := proTeamMap[recepientID]
	SendersCapsheet := capsheetMap[senderID]
	recepientCapsheet := capsheetMap[recepientID]
	salaryMinimum := 0.5
	details := make([]string, 0, len(options))
	complete := true
	for _, option := range options {
		if (option.PlayerID == 0) == (option.DraftPickID == 0) {
			log.Printf("TradeManager: trade option %d must identify exactly one player or draft pick", option.ID)
			complete = false
			continue
		}
		// Contract
		percentage := option.SalaryPercentage
		if option.PlayerID > 0 {
			playerRecord := repository.FindProPlayer(strconv.Itoa(int(option.PlayerID)))
			contract := contractMap[playerRecord.ID]
			if playerRecord.ID == 0 || contract.ID == 0 || (uint(playerRecord.TeamID) != senderID && uint(playerRecord.TeamID) != recepientID) {
				log.Printf("TradeManager: trade player %d has a missing player, contract, or valid team", option.PlayerID)
				complete = false
				continue
			}
			if playerRecord.TeamID == uint16(senderID) {
				details = append(details, fmt.Sprintf("%s %s (player %d): %s to %s (salary retention input %.2f%%)", playerRecord.FirstName, playerRecord.LastName, option.PlayerID, sendingTeam.Abbreviation, receivingTeam.Abbreviation, percentage))
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
				details = append(details, fmt.Sprintf("%s %s (player %d): %s to %s (salary retention input %.2f%%)", playerRecord.FirstName, playerRecord.LastName, option.PlayerID, receivingTeam.Abbreviation, sendingTeam.Abbreviation, percentage))
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

			if playerRecord.ID > 0 {
				repository.SaveProPlayerRecord(playerRecord, db)
			}
			if contract.ID > 0 {
				repository.SaveProContractRecord(contract, db)
			}

		} else if option.DraftPickID > 0 {
			draftPick := repository.FindDraftPickRecord(strconv.Itoa(int(option.DraftPickID)))
			if draftPick.ID == 0 || (draftPick.TeamID != senderID && draftPick.TeamID != recepientID) {
				log.Printf("TradeManager: trade draft pick %d is missing or has an invalid team", option.DraftPickID)
				complete = false
				continue
			}
			if draftPick.TeamID == senderID {
				details = append(details, fmt.Sprintf("%d round %d pick (pick %d): %s to %s", draftPick.Season, draftPick.DraftRound, option.DraftPickID, sendingTeam.Abbreviation, receivingTeam.Abbreviation))
				draftPick.TradePick(recepientID, receivingTeam.Abbreviation)
			} else {
				details = append(details, fmt.Sprintf("%d round %d pick (pick %d): %s to %s", draftPick.Season, draftPick.DraftRound, option.DraftPickID, receivingTeam.Abbreviation, sendingTeam.Abbreviation))
				draftPick.TradePick(senderID, sendingTeam.Abbreviation)
			}

			if draftPick.ID > 0 {
				repository.SaveDraftPickRecord(draftPick, db)
			}
		}

		if option.ID != 0 {
			if err := db.Delete(&option).Error; err != nil {
				log.Printf("TradeManager: failed to delete trade option %d: %v", option.ID, err)
				complete = false
			}
		}
	}
	repository.SaveProCapsheetRecord(SendersCapsheet, db)
	repository.SaveProCapsheetRecord(recepientCapsheet, db)
	return details, complete
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
