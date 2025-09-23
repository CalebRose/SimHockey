package managers

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	util "github.com/CalebRose/SimHockey/_util"
	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllFreeAgencyOffers() []structs.FreeAgencyOffer {
	return repository.FindAllFreeAgentOffers("", "", "", true)
}

func GetAllFreeAgents() []structs.ProfessionalPlayer {
	return repository.FindAllFreeAgents(true, false, false, false)
}

func GetAllWaiverWirePlayers() []structs.ProfessionalPlayer {
	return repository.FindAllFreeAgents(false, true, false, false)
}

func GetContractMap() map[uint]structs.ProContract {
	contracts := repository.FindAllProContracts(true)
	return MakeContractMap(contracts)
}

func GetExtensionMap() map[uint]structs.ExtensionOffer {
	extensions := repository.FindAllProExtensions(true)
	return MakeExtensionMap(extensions)
}

func GetAllAvailableProPlayers(TeamID string, ch chan<- structs.FreeAgencyResponse) {
	var wg sync.WaitGroup
	wg.Add(4)
	var (
		FAs           []structs.ProfessionalPlayer
		WaiverPlayers []structs.ProfessionalPlayer
		Offers        []structs.FreeAgencyOffer
		PracticeSquad []structs.ProfessionalPlayer
	)
	go func() {
		defer wg.Done()
		FAs = GetAllFreeAgentsWithOffers()
	}()
	go func() {
		defer wg.Done()
		WaiverPlayers = GetAllWaiverWirePlayersFAPage()
	}()
	go func() {
		defer wg.Done()
		Offers = GetFreeAgentOffersByTeamID(TeamID)
	}()
	go func() {
		defer wg.Done()
		PracticeSquad = GetAllAffiliatePlayers()

	}()
	wg.Wait()

	ch <- structs.FreeAgencyResponse{
		FreeAgents:    FAs,
		WaiverPlayers: WaiverPlayers,
		PracticeSquad: PracticeSquad,
		TeamOffers:    Offers,
	}
}

func GetAllFreeAgentsWithOffers() []structs.ProfessionalPlayer {
	freeAgents := repository.FindAllFreeAgents(true, false, true, true)

	sort.Slice(freeAgents[:], func(i, j int) bool {
		return freeAgents[i].Overall > freeAgents[j].Overall
	})

	return freeAgents
}

func GetAllWaiverWirePlayersFAPage() []structs.ProfessionalPlayer {
	waivedPlayers := repository.FindAllFreeAgents(false, true, false, true)
	sort.Slice(waivedPlayers[:], func(i, j int) bool {
		return waivedPlayers[i].Overall > waivedPlayers[j].Overall
	})

	return waivedPlayers
}

func GetFreeAgentOffersByTeamID(TeamID string) []structs.FreeAgencyOffer {
	offers := repository.FindAllFreeAgentOffers(TeamID, "", "", true)
	return offers
}

func GetAllAffiliatePlayers() []structs.ProfessionalPlayer {
	players := repository.FindAffiliatePlayers("", "", true, true)
	return players
}

func CreateFAOffer(offer structs.FreeAgencyOfferDTO) structs.FreeAgencyOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	freeAgentOffer := repository.FindFreeAgentOfferRecord("", "", strconv.Itoa(int(offer.ID)), true)
	players := repository.FindAllProPlayers(repository.PlayerQuery{PlayerIDs: []string{strconv.Itoa(int(offer.PlayerID))}})

	if len(players) == 0 {
		return structs.FreeAgencyOffer{}
	}
	player := players[0]
	if freeAgentOffer.ID == 0 {
		id := repository.FindLatestFreeAgentOfferID(db)
		freeAgentOffer.AssignID(id)
	}
	if ts.IsFreeAgencyLocked {
		return freeAgentOffer
	}

	freeAgentOffer.CalculateOffer(offer)

	// If the owning team is sending an offer to a player
	if player.IsAffiliatePlayer && int(player.TeamID) == int(offer.TeamID) {
		SignFreeAgent(freeAgentOffer, player, ts)
	} else {
		repository.SaveFreeAgentOfferRecord(freeAgentOffer, db)
		fmt.Println("Creating offer!")
	}

	if player.IsAffiliatePlayer && int(player.TeamID) != int(offer.TeamID) {
		// Notify team
		notificationMessage := offer.Team + " have placed an offer on " + player.Position + " " + player.FirstName + " " + player.LastName + " to pick up from the practice squad."
		CreateNotification("PHL", notificationMessage, "Affiliate Player Offer", uint(player.TeamID))
		message := offer.Team + " have placed an offer on " + player.Team + " " + player.Position + " " + player.FirstName + " " + player.LastName + " to pick up from the practice squad."
		CreateNewsLog("PHL", message, "Free Agency", int(player.TeamID), ts, true)
	}

	return freeAgentOffer
}

func CancelOffer(offer structs.FreeAgencyOfferDTO) {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	if ts.IsFreeAgencyLocked {
		return
	}

	OfferID := strconv.Itoa(int(offer.ID))

	freeAgentOffer := repository.FindFreeAgentOfferRecord("", "", OfferID, true)
	if freeAgentOffer.ID == 0 {
		return
	}
	freeAgentOffer.DeactivateOffer()

	db.Save(&freeAgentOffer)
}

func CreateWaiverOffer(offer structs.WaiverOfferDTO) structs.WaiverOffer {
	db := dbprovider.GetInstance().GetDB()
	ts := GetTimestamp()
	waiverOffer := repository.FindWaiverWireOfferRecord("", "", strconv.Itoa(int(offer.ID)), true)

	if waiverOffer.ID == 0 {
		id := repository.FindLatestWaiverOfferID(db)
		waiverOffer.AssignID(id)
	}

	if ts.IsFreeAgencyLocked {
		return waiverOffer
	}

	waiverOffer.Map(offer)

	repository.SaveWaiverRecord(waiverOffer, db)

	fmt.Println("Creating offer!")

	return waiverOffer
}

func CancelWaiverOffer(offer structs.WaiverOfferDTO) {
	db := dbprovider.GetInstance().GetDB()

	OfferID := strconv.Itoa(int(offer.ID))
	waiverOffer := repository.FindWaiverWireOfferRecord("", "", OfferID, true)
	if waiverOffer.ID == 0 {
		return
	}

	repository.DeleteWaiverRecord(waiverOffer, db)
}

func SyncAIOffers() {
	db := dbprovider.GetInstance().GetDB()

	teams := repository.FindAllProTeams(repository.TeamClauses{})

	offers := GetAllFreeAgencyOffers()
	offerMapByTeamID := MakeFreeAgencyOfferMapByTeamID(offers)
	freeAgents := GetAllFreeAgents()
	players := repository.FindAllProPlayers(repository.PlayerQuery{})
	playerMap := MakeProfessionalPlayerMapByTeamID(players)

	for _, team := range teams {
		if len(team.Owner) > 0 && team.Owner != "AI" {
			continue
		}
		if team.Owner == "" && len(team.GM) > 0 && team.GM != "AI" {
			continue
		}

		offersByTeam := offerMapByTeamID[team.ID]
		if len(offersByTeam) > 7 {
			continue
		}
		freeAgentOfferMap := MakeFreeAgencyOfferMap(offersByTeam)
		roster := playerMap[team.ID]
		cCount := 0
		fCount := 0
		dCount := 0
		gCount := 0
		cBids := 0
		fBids := 0
		dBids := 0
		gBids := 0
		for _, p := range roster {
			switch p.Position {
			case Center:
				cCount++
			case Forward:
				fCount++
			case Defender:
				dCount++
			default:
				gCount++
			}
		}

		// Iterate through FA list to get bids
		for _, fa := range freeAgents {
			existingOffers := freeAgentOfferMap[fa.ID]
			if len(existingOffers) > 0 {
				switch fa.Position {
				case Center:
					cBids++
				case Forward:
					fBids++
				case Defender:
					dBids++
				default:
					gBids++
				}
			}
		}

		for _, fa := range freeAgents {
			existingOffers := freeAgentOfferMap[fa.ID]
			if len(existingOffers) > 0 {
				continue
			}
			if fa.Position == Center && (cCount > 4 || cBids > 1) {
				continue
			}
			if fa.Position == Forward && (fCount > 8 || fBids > 3) {
				continue
			}
			if fa.Position == Defender && (dCount > 6 || dBids > 2) {
				continue
			}
			if fa.Position == Goalie && (gCount > 2 || gBids > 0) {
				continue
			}
			coinFlip := util.GenerateIntFromRange(1, 2)
			if coinFlip == 2 {
				continue
			}

			// Okay, now we found an open player. Send a bid.
			basePay := float32(1.0)
			if fa.Age < 25 || fa.Overall < 19 {
				basePay = 0.7
			} else if fa.Overall > 24 {
				rangedPay := util.GenerateFloatFromRange(1, 3.5)
				basePay = RoundToFixedDecimalPlace(rangedPay, 2)
			}

			yearsOnContract := 2
			if fa.Overall > 24 {
				yearsOnContract = 3
			} else if fa.Overall < 19 {
				yearsOnContract = 1
			}
			y1 := basePay
			y2 := float32(0.0)
			y3 := float32(0.0)
			if yearsOnContract > 2 {
				y3 = basePay
			}
			if yearsOnContract > 1 {
				y2 = basePay
			}
			switch fa.Position {
			case Center:
				cBids++
			case Forward:
				fBids++
			case Defender:
				dBids++
			default:
				gBids++
			}
			offer := structs.FreeAgencyOffer{
				Y1BaseSalary:   y1,
				Y2BaseSalary:   y2,
				Y3BaseSalary:   y3,
				TotalSalary:    basePay * float32(yearsOnContract),
				ContractValue:  basePay,
				IsActive:       true,
				PlayerID:       fa.ID,
				TeamID:         team.ID,
				ContractLength: yearsOnContract,
			}

			repository.SaveFreeAgentOfferRecord(offer, db)
		}
	}
}

func SyncFreeAgencyOffers() {
	db := dbprovider.GetInstance().GetDB()

	ts := GetTimestamp()
	if !ts.IsFreeAgencyLocked {
		ts.ToggleFALock()
		repository.SaveTimestamp(ts, db)
	}

	freeAgents := GetAllFreeAgents()
	capsheetMap := GetProCapsheetMap()
	offers := GetAllFreeAgencyOffers()
	offerMap := MakeFreeAgencyOfferMap(offers)

	for _, FA := range freeAgents {
		if ts.IsOffSeason && !FA.IsAcceptingOffers {
			continue
		}
		offers := offerMap[FA.ID]

		if len(offers) == 0 {
			continue
		}

		maxDay := 1000
		for _, offer := range offers {
			if maxDay > int(offer.Syncs) {
				maxDay = int(offer.Syncs)
			}
		}
		if maxDay < 3 {
			for _, offer := range offers {
				offer.IncrementSyncs()
				repository.SaveFreeAgentOfferRecord(offer, db)
			}
		} else {
			// For Inaugural Season, don't worry about preference factors. These should be for next season.
			// Just sort by contract value & take the highest
			/*
				For next season
				Create new struct containing offer & them new modified AAV, based on the preferences & factors specified
				Run through all offers again to calculate the new AAV, and then resort the list
				Highest value should be the winning offer given the team isn't maxed out or anything
			*/
			sort.Slice(offers, func(i, j int) bool {
				return offers[i].ContractValue > offers[j].ContractValue
			})
			WinningOffer := structs.FreeAgencyOffer{}
			competingTeams := []structs.FreeAgencyOffer{}
			highestAAV := 0.0
			for _, offer := range offers {
				capsheet := capsheetMap[offer.TeamID]
				if capsheet.ID == 0 || !offer.IsActive {
					continue
				}
				if offer.ContractValue > float32(highestAAV) {
					highestAAV = float64(offer.ContractValue)
					competingTeams = []structs.FreeAgencyOffer{offer}
				} else if offer.ContractValue == float32(highestAAV) && highestAAV > 0 {
					competingTeams = append(competingTeams, offer)
				} else {
					break
				}
			}
			idx := 0
			// If there is more than one competing team
			if len(competingTeams) > 1 {
				idx = util.GenerateIntFromRange(0, len(competingTeams)-1)
			}
			WinningOffer = competingTeams[idx]
			// Cancel All Offers
			for _, offer := range offers {
				capsheet := capsheetMap[offer.TeamID]
				if capsheet.ID == 0 {
					continue
				}
				if offer.IsActive && offer.ID != WinningOffer.ID {
					offer.RejectOffer()
				} else if offer.IsActive && offer.ID == WinningOffer.ID {
					offer.DeactivateOffer()
				}

				repository.SaveFreeAgentOfferRecord(offer, db)
			}

			if WinningOffer.ID > 0 {
				SignFreeAgent(WinningOffer, FA, ts)
			} else if ts.IsOffSeason {
				FA.WaitUntilAfterDraft()
				repository.SaveProPlayerRecord(FA, db)
			}
		}
	}
	if ts.IsFreeAgencyLocked {
		ts.ToggleFALock()
		repository.SaveTimestamp(ts, db)
	}
}

func SignFreeAgent(offer structs.FreeAgencyOffer, FreeAgent structs.ProfessionalPlayer, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	proTeam := repository.FindProTeamRecord(strconv.Itoa(int(offer.TeamID)))
	messageStart := "FA "
	value := offer.ContractValue
	if !FreeAgent.IsAffiliatePlayer {
		Contract := structs.ProContract{
			PlayerID:       FreeAgent.ID,
			TeamID:         proTeam.ID,
			OriginalTeamID: proTeam.ID,
			ContractLength: offer.ContractLength,
			Y1BaseSalary:   offer.Y1BaseSalary,
			Y2BaseSalary:   offer.Y2BaseSalary,
			Y3BaseSalary:   offer.Y3BaseSalary,
			Y4BaseSalary:   offer.Y4BaseSalary,
			Y5BaseSalary:   offer.Y5BaseSalary,
			ContractValue:  offer.ContractValue,
			SigningValue:   offer.ContractValue,
			IsActive:       true,
			IsComplete:     false,
			IsExtended:     false,
		}
		repository.CreateProContractRecord(db, Contract)
	} else {
		Contract := repository.FindProContract(strconv.Itoa(int(FreeAgent.ID)))
		Contract.MapAffiliateOffer(offer)
		value = Contract.ContractValue
		repository.SaveProContractRecord(Contract, db)
		messageStart = "PS "
	}
	FreeAgent.SignPlayer(uint(proTeam.ID), proTeam.Abbreviation, ts.Week > 18, offer.ToAffiliate)
	repository.SaveProPlayerRecord(FreeAgent, db)

	// News Log
	message := messageStart + FreeAgent.Position + " " + FreeAgent.FirstName + " " + FreeAgent.LastName + " has signed with the " + proTeam.TeamName + " with a contract worth approximately $" + strconv.Itoa(int(value)) + " Million Dollars."
	CreateNewsLog("PHL", message, "Free Agency", int(offer.TeamID), ts, true)
}
