package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
)

func ExportAllProPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	players := managers.GetAllProPlayers()
	managers.WriteProPlayersExport(w, players)

	json.NewEncoder(w).Encode("Players Exported")
}

func ExportAllCollegePlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	players := managers.GetAllCollegePlayers()
	managers.WriteCollegePlayersExport(w, players)

	json.NewEncoder(w).Encode("Players Exported")
}
