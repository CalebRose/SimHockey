package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/mux"
)

func GetAllHCKRequests(w http.ResponseWriter, r *http.Request) {
	requests := managers.GetAllHCKRequests()

	json.NewEncoder(w).Encode(requests)
}

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

func CreateCHLTeamRequest(w http.ResponseWriter, r *http.Request) {
	req := structs.CollegeTeamRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreateCHLTeamRequest(req)

	json.NewEncoder(w).Encode(req)
}

func CreatePHLTeamRequest(w http.ResponseWriter, r *http.Request) {
	req := structs.ProTeamRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.CreatePHLTeamRequest(req)

	json.NewEncoder(w).Encode(req)
}

func ApproveCHLTeamRequest(w http.ResponseWriter, r *http.Request) {
	req := structs.CollegeTeamRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.ApproveCHLTeamRequest(req)

	json.NewEncoder(w).Encode(req)
}

func ApprovePHLTeamRequest(w http.ResponseWriter, r *http.Request) {
	req := structs.ProTeamRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.ApprovePHLTeamRequest(req)

	json.NewEncoder(w).Encode(req)
}

func RejectCHLTeamRequest(w http.ResponseWriter, r *http.Request) {
	req := structs.CollegeTeamRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RejectCollegeTeamRequest(req)

	json.NewEncoder(w).Encode(req)
}

func RejectPHLTeamRequest(w http.ResponseWriter, r *http.Request) {
	req := structs.ProTeamRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	managers.RejectPHLTeamRequest(req)

	json.NewEncoder(w).Encode(req)
}
