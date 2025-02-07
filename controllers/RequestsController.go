package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

func ViewCHLTeamUponRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	team := managers.GetCHLTeamForAvailableTeamsPage(teamID)

	json.NewEncoder(w).Encode(team)
}

func ViewPHLTeamUponRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	if len(teamID) == 0 {
		panic("User did not provide TeamID")
	}

	team := managers.GetProTeamForAvailableTeamsPage(teamID)

	json.NewEncoder(w).Encode(team)
}
