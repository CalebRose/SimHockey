package controllers

import (
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
)

func AddCompensationPicksMidRound(w http.ResponseWriter, r *http.Request) {
	managers.AddCompensationPicksMidRound()
	w.WriteHeader(http.StatusOK)
}

func FixRecruitProfileModifiers(w http.ResponseWriter, r *http.Request) {
	// managers.AddPreferenceValuesToRecruits()
	managers.FixModifierValuesForRecruits()
	w.WriteHeader(http.StatusOK)
}
