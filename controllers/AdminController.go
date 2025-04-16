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

func GenerateProTeams(w http.ResponseWriter, r *http.Request) {
	managers.ImportProTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateInitialRosters(w http.ResponseWriter, r *http.Request) {
	managers.GenerateInitialRosters()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateProTestData(w http.ResponseWriter, r *http.Request) {
	managers.GenerateInitialProPool()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func RunAICollegeLineups(w http.ResponseWriter, r *http.Request) {
	managers.RunLineupsForAICollegeTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func RunAIProLineups(w http.ResponseWriter, r *http.Request) {
	managers.RunLineupsForAIProTeams()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateCroots(w http.ResponseWriter, r *http.Request) {
	managers.GenerateCroots()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func ImportProRosters(w http.ResponseWriter, r *http.Request) {
	managers.ImportProRosters()
	managers.ImportStandingsForNewSeason()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func ImportTeamProfileRecords(w http.ResponseWriter, r *http.Request) {
	managers.ImportTeamRecruitingProfiles()
	json.NewEncoder(w).Encode("Data Generated ran!")
}

func GenerateCapsheets(w http.ResponseWriter, r *http.Request) {
	managers.AllocateCapsheets()

	json.NewEncoder(w).Encode("Data Generated ran!")
}

func AddFAPreferences(w http.ResponseWriter, r *http.Request) {
	managers.AddFAPreferences()

	json.NewEncoder(w).Encode("Data Generated ran!")
}
