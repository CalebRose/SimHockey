package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

func ToggleNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notiID := vars["notiID"]
	managers.ToggleNotification(notiID)
	json.NewEncoder(w).Encode(true)
}

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notiID := vars["notiID"]
	managers.DeleteNotification(notiID)
	json.NewEncoder(w).Encode(true)
}
