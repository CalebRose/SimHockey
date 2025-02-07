package managers

import (
	"sort"
	"sync"

	"github.com/CalebRose/SimHockey/repository"
	"github.com/CalebRose/SimHockey/structs"
)

func GetAllFreeAgents() []structs.ProfessionalPlayer {
	return repository.FindAllFreeAgents(true, false, false, false)
}

func GetAllWaiverWirePlayers() []structs.ProfessionalPlayer {
	return repository.FindAllFreeAgents(false, true, false, false)
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
