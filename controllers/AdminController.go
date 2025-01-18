package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
)

// GetTimeStamp
func GetCurrentTimestamp(w http.ResponseWriter, r *http.Request) {

	timestamp := managers.GetTimestamp()

	json.NewEncoder(w).Encode(timestamp)
}

func TestEngine(w http.ResponseWriter, r *http.Request) {
	managers.RunGames()

	json.NewEncoder(w).Encode("Game ran!")
}

func GenerateTestData(w http.ResponseWriter, r *http.Request) {
	managers.GenerateTestRosters()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateProTestData(w http.ResponseWriter, r *http.Request) {
	managers.GenerateTestProPool()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func RunAILineups(w http.ResponseWriter, r *http.Request) {
	managers.RunLineupsForAICollegeTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}
