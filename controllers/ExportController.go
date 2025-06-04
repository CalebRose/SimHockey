package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

func ExportAllProPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	players := managers.GetAllProPlayers()
	managers.WriteProPlayersExport(w, players, "_phl_players.csv")

	json.NewEncoder(w).Encode("Players Exported")
}

func ExportAllCollegePlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	players := managers.GetAllCollegePlayers()
	managers.WriteCollegePlayersExport(w, players, "_chl_players.csv")

	json.NewEncoder(w).Encode("Players Exported")
}

func ExportCollegeRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	w.Header().Set("Content-Type", "text/csv")
	players := managers.GetCollegePlayersByTeamID(teamID)
	team := managers.GetCollegeTeamByTeamID(teamID)
	managers.WriteCollegePlayersExport(w, players, "_"+team.Abbreviation+"_roster.csv")

	json.NewEncoder(w).Encode("Players Exported")
}

func ExportProRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]

	w.Header().Set("Content-Type", "text/csv")
	players := managers.GetProPlayersByTeamID(teamID)
	team := managers.GetProTeamByTeamID(teamID)
	managers.WriteProPlayersExport(w, players, "_"+team.Abbreviation+"_roster.csv")

	json.NewEncoder(w).Encode("Players Exported")
}

func ExportCHLRecruits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	players := managers.GetAllRecruitRecords()
	managers.WriteCollegeRecruitsExport(w, players, "_toucans_secret_croot_list.csv")

	json.NewEncoder(w).Encode("Players Exported")
}

func ExportHCKGameResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	weekID := vars["weekID"]
	timeslot := vars["timeslot"]
	managers.ExportHCKGameResults(w, seasonID, weekID, timeslot)
}

func ExportCollegePlayByPlayResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		panic("User did not provide a first name")
	}

	managers.HandleCollegePlayByPlayExport(w, gameID)
}

func ExportProPlayByPlayResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if len(gameID) == 0 {
		panic("User did not provide a first name")
	}

	managers.HandleProPlayByPlayExport(w, gameID)
}
