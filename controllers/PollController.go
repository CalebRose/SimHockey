package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CalebRose/SimHockey/dbprovider"
	"github.com/CalebRose/SimHockey/managers"
	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/mux"
)

func CreatePollSubmission(w http.ResponseWriter, r *http.Request) {
	var dto structs.CollegePollSubmission
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate info from DTO
	if len(dto.Username) == 0 {
		log.Fatalln("ERROR: Cannot submit poll.")
	}

	poll := managers.CreatePoll(dto)
	json.NewEncoder(w).Encode(poll)
}

func GetPollSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	poll := managers.GetPollSubmissionByUsernameWeekAndSeason(username)

	res := structs.PollDataResponse{
		Poll: poll,
	}

	json.NewEncoder(w).Encode(res)
}

func SyncCollegePoll(w http.ResponseWriter, r *http.Request) {
	db := dbprovider.GetInstance().GetDB()
	ts := managers.GetTimestamp()
	managers.SyncCollegePollSubmissionForCurrentWeek(uint(ts.Week), uint(ts.WeekID), uint(ts.SeasonID))
	ts.TogglePollRan()
	db.Save(&ts)
}

func GetOfficialPollsBySeasonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seasonID := vars["seasonID"]
	if len(seasonID) == 0 {
		panic("User did not provide seasonID")
	}
	polls := managers.GetOfficialPollBySeasonID(seasonID)

	res := structs.PollDataResponse{
		OfficialPolls: polls,
	}

	json.NewEncoder(w).Encode(res)
}
