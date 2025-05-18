package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
)

func SaveCHLLineups(w http.ResponseWriter, r *http.Request) {
	req := structs.UpdateLineupsDTO{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := managers.SaveCHLLineup(req)

	json.NewEncoder(w).Encode(dto)
}

func SavePHLLineups(w http.ResponseWriter, r *http.Request) {
	req := structs.UpdateLineupsDTO{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := managers.SavePHLLineup(req)

	json.NewEncoder(w).Encode(dto)
}

func SaveCHLGameplan(w http.ResponseWriter, r *http.Request) {
	req := structs.CollegeGameplan{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gp := managers.SaveCollegeGameplanSettings(req)

	json.NewEncoder(w).Encode(gp)
}

func SavePHLGameplan(w http.ResponseWriter, r *http.Request) {
	req := structs.ProGameplan{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gp := managers.SaveProGameplanSettings(req)

	json.NewEncoder(w).Encode(gp)
}

func CreateGameplans(w http.ResponseWriter, r *http.Request) {
	managers.CreateGameplans()
	json.NewEncoder(w).Encode("Gameplans Created")
}
