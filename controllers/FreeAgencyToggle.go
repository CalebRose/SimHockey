package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/CalebRose/SimHockey/managers"
)

func ToggleFreeAgency(w http.ResponseWriter, r *http.Request) {
    managers.ToggleFreeAgency()
	json.NewEncoder(w).Encode("Free Agency Unlocked!")
}
