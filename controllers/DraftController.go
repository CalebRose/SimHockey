package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/mux"
)

func GetDraftPageData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		http.Error(w, "User did not provide TeamID", http.StatusBadRequest)
		return
	}

	res := managers.GetDraftBootstrap()

	json.NewEncoder(w).Encode(res)
}

func AddPlayerToScoutBoard(w http.ResponseWriter, r *http.Request) {

	var scoutProfileDto structs.ScoutingProfileDTO
	err := json.NewDecoder(r.Body).Decode(&scoutProfileDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scoutingProfile := managers.CreateScoutingProfile(scoutProfileDto)

	json.NewEncoder(w).Encode(scoutingProfile)
}

func ExportDraftedPicks(w http.ResponseWriter, r *http.Request) {
	var draftPickDTO structs.ExportDraftPicksDTO
	err := json.NewDecoder(r.Body).Decode(&draftPickDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveComplete := managers.ExportDraftedPlayers(draftPickDTO.DraftPicks)

	json.NewEncoder(w).Encode(saveComplete)
}

func RevealScoutingAttribute(w http.ResponseWriter, r *http.Request) {
	var revealAttributeDTO structs.RevealAttributeDTO
	err := json.NewDecoder(r.Body).Decode(&revealAttributeDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saveComplete := managers.RevealScoutingAttribute(revealAttributeDTO)

	json.NewEncoder(w).Encode(saveComplete)
}

func RemovePlayerFromScoutBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		panic("User did not provide scout profile id")
	}

	managers.RemovePlayerFromScoutBoard(id)

	json.NewEncoder(w).Encode(true)
}

func GetScoutingDataByDraftee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		panic("User did not provide scout profile id")
	}

	data := managers.GetScoutingDataByPlayerID(id)

	json.NewEncoder(w).Encode(data)
}

func BringUpCollegePlayerToPros(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["draftPickID"]
	if len(id) == 0 {
		panic("User did not provide draftPick id")
	}

	data := managers.BringUpCollegePlayerToPros(id)

	json.NewEncoder(w).Encode(data)
}
