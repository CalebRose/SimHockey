package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

func BootstrapHockeyData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collegeID := vars["collegeID"]
	proID := vars["proID"]
	data := managers.GetBootstrapData(collegeID, proID)
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Bootstrap Data: %+v", data)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
