package managers

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

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

	if len(players) == 0 || freeAgentOffer.ID == 0 {
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
		CreateNewsLog("PHL", message, "Free Agency", int(player.TeamID), ts)
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
	freeAgentOffer.CancelOffer()

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

func SignFreeAgent(offer structs.FreeAgencyOffer, FreeAgent structs.ProfessionalPlayer, ts structs.Timestamp) {
	db := dbprovider.GetInstance().GetDB()

	proTeam := repository.FindProTeamRecord(strconv.Itoa(int(offer.TeamID)))
	Contract := structs.ProContract{}
	messageStart := "FA "
	if !FreeAgent.IsAffiliatePlayer {
		Contract = structs.ProContract{
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
		Contract = repository.FindProContract(strconv.Itoa(int(FreeAgent.ID)))
		Contract.MapAffiliateOffer(offer)
		repository.SaveProContractRecord(Contract, db)
		messageStart = "PS "
	}
	FreeAgent.SignPlayer(uint(proTeam.ID), proTeam.Abbreviation, ts.Week > 18)
	repository.SaveProPlayerRecord(FreeAgent, db)

	// News Log
	message := messageStart + FreeAgent.Position + " " + FreeAgent.FirstName + " " + FreeAgent.LastName + " has signed with the " + proTeam.TeamName + " with a contract worth approximately $" + strconv.Itoa(int(Contract.ContractValue)) + " Million Dollars."
	CreateNewsLog("PHL", message, "Free Agency", int(offer.TeamID), ts)
}
