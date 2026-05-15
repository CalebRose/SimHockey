package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/mux"
)

// CreateCHLGameRequest accepts a CHLGameRequest body and persists it.
func CreateCHLGameRequest(w http.ResponseWriter, r *http.Request) {
	var request structs.CHLGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	managers.CreateCHLGameRequest(request)
	json.NewEncoder(w).Encode(true)
}

// AcceptCHLGameRequest marks the request as accepted.
func AcceptCHLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.AcceptCHLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// RejectCHLGameRequest deletes the request.
func RejectCHLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.RejectCHLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// ProcessCHLGameRequest converts an accepted request into a CollegeGame record.
func ProcessCHLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.ProcessCHLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}

// VetoCHLGameRequest deletes the request via admin veto.
func VetoCHLGameRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestID"]
	if len(requestID) == 0 {
		panic("User did not provide a requestID")
	}
	managers.VetoCHLGameRequest(requestID)
	json.NewEncoder(w).Encode(true)
}
