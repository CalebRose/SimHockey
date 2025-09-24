package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/mux"
)

func CreatePromise(w http.ResponseWriter, r *http.Request) {

	var createPromiseDto structs.CollegePromise
	err := json.NewDecoder(r.Body).Decode(&createPromiseDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	promise := managers.CreatePromise(createPromiseDto)

	json.NewEncoder(w).Encode(promise)

	fmt.Fprintf(w, "New Promise Created")
}

func UpdatePromise(w http.ResponseWriter, r *http.Request) {

	var createPromiseDto structs.CollegePromise
	err := json.NewDecoder(r.Body).Decode(&createPromiseDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.UpdatePromise(createPromiseDto)

	fmt.Fprintf(w, "Promise Updated")
}

func CancelPromise(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	promiseID := vars["promiseID"]

	if len(promiseID) == 0 {
		panic("User did not provide Promise ID")
	}

	managers.CancelPromise(promiseID)

	fmt.Fprintf(w, "Promise Cancelled.")
}

func GetPromiseByPlayerID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["playerID"]
	teamID := vars["teamID"]
	if len(id) == 0 {
		panic("User did not provide proper IDs")
	}

	promise := managers.GetCollegePromiseByCollegePlayerID(id, teamID)

	encodedJson, err := json.Marshal(promise)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedJson)
}

// func GetTransferPortalPageData(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	teamID := vars["teamID"]
// 	if len(teamID) == 0 {
// 		panic("User did not provide proper IDs")
// 	}

// 	data := managers.GetTransferPortalData(teamID)

// 	encodedJson, err := json.Marshal(data)
// 	if err != nil {
// 		log.Printf("Error encoding JSON: %v", err)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(encodedJson)
// }

func AddTransferPlayerToBoard(w http.ResponseWriter, r *http.Request) {

	var transferPortalProfile structs.TransferPortalProfile
	err := json.NewDecoder(r.Body).Decode(&transferPortalProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := managers.AddTransferPlayerToBoard(transferPortalProfile)

	json.NewEncoder(w).Encode(profile)
}

func RemovePlayerFromTransferPortalBoard(w http.ResponseWriter, r *http.Request) {
	var transferPortalProfile structs.TransferPortalProfile
	err := json.NewDecoder(r.Body).Decode(&transferPortalProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := managers.RemovePlayerFromTransferPortalBoard(transferPortalProfile)
	json.NewEncoder(w).Encode(profile)
}

func SaveTransferBoard(w http.ResponseWriter, r *http.Request) {

	var transferPortalProfile structs.UpdateTransferPortalBoard
	err := json.NewDecoder(r.Body).Decode(&transferPortalProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.AllocatePointsToTransferPlayer(transferPortalProfile)
}

// func GetScoutingDataByTransfer(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	id := vars["id"]
// 	if len(id) == 0 {
// 		panic("User did not provide scout profile id")
// 	}

// 	data := managers.GetTransferScoutingDataByPlayerID(id)

// 	json.NewEncoder(w).Encode(data)
// }

func SyncTransferPortal(w http.ResponseWriter, r *http.Request) {
	managers.SyncTransferPortal()

	json.NewEncoder(w).Encode("AI Boards for Transfer Portal Complete.")
}

func FillUpTransferBoardsAI(w http.ResponseWriter, r *http.Request) {
	managers.AICoachFillBoardsPhase()

	json.NewEncoder(w).Encode("AI Boards for Transfer Portal Complete.")
}

func AllocateAndPromisePlayersAI(w http.ResponseWriter, r *http.Request) {

	managers.AICoachAllocateAndPromisePhase()

	json.NewEncoder(w).Encode("Allocated and promised.")
}
