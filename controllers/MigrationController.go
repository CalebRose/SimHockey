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
	json.NewEncoder(w).Encode("Recruiting Sync Complete")
}

func ImportPHLDraftOrder(w http.ResponseWriter, r *http.Request) {
	managers.ImportPhlDraftOrder()
	json.NewEncoder(w).Encode("Recruiting Sync Complete")
}
