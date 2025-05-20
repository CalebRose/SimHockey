package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
)

func CreateFreeAgencyOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.FreeAgencyOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var offer = managers.CreateFAOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(offer)
}

func CancelFreeAgencyOffer(w http.ResponseWriter, r *http.Request) {
	var freeAgencyOfferDTO structs.FreeAgencyOfferDTO
	err := json.NewDecoder(r.Body).Decode(&freeAgencyOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CancelOffer(freeAgencyOfferDTO)

	json.NewEncoder(w).Encode(true)
}

func CreateWaiverWireOffer(w http.ResponseWriter, r *http.Request) {
	var waiverWireOfferDTO structs.WaiverOfferDTO
	err := json.NewDecoder(r.Body).Decode(&waiverWireOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var offer = managers.CreateWaiverOffer(waiverWireOfferDTO)

	json.NewEncoder(w).Encode(offer)
}

// CancelWaiverWireOffer
func CancelWaiverWireOffer(w http.ResponseWriter, r *http.Request) {
	var waiverWireOfferDTO structs.WaiverOfferDTO
	err := json.NewDecoder(r.Body).Decode(&waiverWireOfferDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CancelWaiverOffer(waiverWireOfferDTO)

	json.NewEncoder(w).Encode(true)
}

func TestFASync(w http.ResponseWriter, r *http.Request) {
	managers.SyncFreeAgencyOffers()

	json.NewEncoder(w).Encode(true)
}

func TestFAOffers(w http.ResponseWriter, r *http.Request) {
	managers.SyncAIOffers()

	json.NewEncoder(w).Encode(true)
}
