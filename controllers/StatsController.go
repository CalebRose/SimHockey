package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

func GetCHLStatsPageContentForSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	response := managers.SearchCollegeStats(seasonID, weekID, viewType, gameType)
	json.NewEncoder(w).Encode(response)
}

func ExportCHLStatsPageContentForSeason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	managers.ExportCollegeStats(seasonID, weekID, viewType, gameType, w)
}

func GetProStatsPageContent(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	response := managers.SearchProStats(seasonID, weekID, viewType, gameType)

	json.NewEncoder(w).Encode(response)
}

func ExportProStatsPageContent(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	viewType := vars["viewType"]
	weekID := vars["weekID"]
	gameType := vars["gameType"]

	if len(viewType) == 0 {
		panic("User did not provide view type")
	}

	if len(seasonID) == 0 {
		panic("User did not provide TeamID")
	}

	managers.ExportProStats(seasonID, weekID, viewType, gameType, w)
}

func GetCollegeGameResultsByGameID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetCHLGameResultsByGameID(gameID)

	json.NewEncoder(w).Encode(player)
}

func GetProGameResultsByGameID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		panic("User did not provide a first name")
	}

	player := managers.GetPHLGameResultsByGameID(gameID)

	json.NewEncoder(w).Encode(player)
}
