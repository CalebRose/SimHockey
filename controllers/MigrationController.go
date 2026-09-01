package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
)

/*
	For Migrating data & Fixing data issues
*/

func FixSeasonStatTables(w http.ResponseWriter, r *http.Request) {
	managers.FixSeasonStatTables()
	json.NewEncoder(w).Encode("Recruiting Sync Complete")
}

func FixStandingsTables(w http.ResponseWriter, r *http.Request) {
	managers.FixStandingsTables()
	json.NewEncoder(w).Encode("Fix Sync Complete")
}

func ImportPHLDraftOrder(w http.ResponseWriter, r *http.Request) {
	managers.ImportPhlDraftOrder()
	json.NewEncoder(w).Encode("Fix Sync Complete")
}

func FixDraftablePlayersTable(w http.ResponseWriter, r *http.Request) {
	managers.FixDraftablePlayersTable()
	json.NewEncoder(w).Encode("Fix Sync Complete")
}

func FixAddingCroots(w http.ResponseWriter, r *http.Request) {
	managers.FixAddingRecruitsToCollege()
	json.NewEncoder(w).Encode("Fix Sync Complete")
}

func FixAddingRookieContracts(w http.ResponseWriter, r *http.Request) {
	managers.FixRookieContracts()
	json.NewEncoder(w).Encode("Fix Sync Complete")
}

func ExtendPHLPlayers(w http.ResponseWriter, r *http.Request) {
	managers.ExtendPHLPlayers()
	json.NewEncoder(w).Encode("Fix Sync Complete")
}

func FixSimPHLRetiredPlayers(w http.ResponseWriter, r *http.Request) {
	managers.FixSimPHLRetiredPlayers()
	json.NewEncoder(w).Encode("Fix Sync Complete")
}
