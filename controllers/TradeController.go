package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/mux"
)

func UpdateTradePreferences(w http.ResponseWriter, r *http.Request) {

	var tradePreferenceDTO structs.TradePreferencesDTO
	err := json.NewDecoder(r.Body).Decode(&tradePreferenceDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdateTradePreferences(tradePreferenceDTO)

	fmt.Fprintf(w, "Trade Preferences Updated")
}

// Create NFL Trade Proposal
func CreateTradeProposal(w http.ResponseWriter, r *http.Request) {

	var tradeProposalDTO structs.TradeProposalDTO
	err := json.NewDecoder(r.Body).Decode(&tradeProposalDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateTradeProposal(tradeProposalDTO)

	// recruitingProfile := managers.CreateRecruitingProfileForRecruit(tradeProposalDTO)
	fmt.Fprintf(w, "New Trade Proposal Created")
}

// Accept Trade Offer
func AcceptTradeOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.AcceptTradeProposal(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// Reject Trade Offer
func RejectTradeOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.RejectTradeProposal(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// Cancels Trade Offer
func CancelTradeOffer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.CancelTradeProposal(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// SyncAcceptedTrade -- Admin approve a trade
func SyncAcceptedTrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.SyncAcceptedTrade(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// SyncAcceptedTrade -- Admin approve a trade
func VetoAcceptedTrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	if len(proposalID) == 0 {
		panic("User did not provide a proposalID")
	}

	managers.VetoTrade(proposalID)

	json.NewEncoder(w).Encode("Proposal " + proposalID + " has been accepted.")
}

// CleanUpRejectedTrades -- Remove all rejected trades from the DB
func CleanUpRejectedTrades(w http.ResponseWriter, r *http.Request) {
	managers.RemoveRejectedTrades()

	json.NewEncoder(w).Encode("Removed all rejected trades from the interface.")
}
