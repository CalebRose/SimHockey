package controllers

import (
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
)

func AddCompensationPicksMidRound(w http.ResponseWriter, r *http.Request) {
	managers.AddCompensationPicksMidRound()
	w.WriteHeader(http.StatusOK)
}
