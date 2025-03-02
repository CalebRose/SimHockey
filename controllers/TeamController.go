package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/mux"
)

func RemoveUserFromCollegeTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamID"]
	managers.RemoveUserFromCollegeTeam(teamID)
	json.NewEncoder(w).Encode("Removed User")
}

func RemoveUserFromProTeam(w http.ResponseWriter, r *http.Request) {
	req := structs.ProTeamRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RemoveUserFromProTeam(req)

	json.NewEncoder(w).Encode(req)
}
