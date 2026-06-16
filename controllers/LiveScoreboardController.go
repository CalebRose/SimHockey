package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/mux"
)

// StreamCHLLiveGames handles the SSE connection for the CHL frontend
func StreamCHLLiveGames(w http.ResponseWriter, r *http.Request) {
	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Type assert the ResponseWriter to an http.Flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// The channel that our manager will push play-by-play events into
	playChannel := make(chan string)

	// Fire the manager, passing the request context so it knows when the user disconnects
	go managers.StartLiveScoreboardSession(r.Context(), "CHL", 8, playChannel)

	for {
		select {
		case playJSON := <-playChannel:
			// Send the JSON payload to the browser
			fmt.Fprintf(w, "data: %s\n\n", playJSON)
			// FIXED: Use the locally scoped flusher variable
			flusher.Flush()
		case <-r.Context().Done():
			// Browser closed the connection
			fmt.Println("Client disconnected from CHL Live Scoreboard")
			return
		}
	}
}

// StreamPHLLiveGames handles the SSE connection for the PHL frontend
func StreamPHLLiveGames(w http.ResponseWriter, r *http.Request) {
	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Type assert the ResponseWriter to an http.Flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	playChannel := make(chan string)

	go managers.StartLiveScoreboardSession(r.Context(), "PHL", 4, playChannel)

	for {
		select {
		case playJSON := <-playChannel:
			fmt.Fprintf(w, "data: %s\n\n", playJSON)
			// FIXED: Use the locally scoped flusher variable
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client disconnected from PHL Live Scoreboard")
			return
		}
	}
}

// GetLiveGamesHub returns the current state of games for the live rink hub
func GetLiveGamesHub(w http.ResponseWriter, r *http.Request) {
	isCollege := r.URL.Query().Get("isCollege") == "true"
	season := r.URL.Query().Get("season")
	week := r.URL.Query().Get("week")
	timeslot := r.URL.Query().Get("timeslot")

	response := managers.GetLiveGamesHubData(isCollege, season, week, timeslot)
	json.NewEncoder(w).Encode(response)
}

// GetBulkPlayByPlay returns the massive array of plays to feed the frontend spoofing loop
func GetBulkPlayByPlay(w http.ResponseWriter, r *http.Request) {
	isCollege := r.URL.Query().Get("isCollege") == "true"
	season := r.URL.Query().Get("season")
	week := r.URL.Query().Get("week")
	timeslot := r.URL.Query().Get("timeslot")

	response := managers.GetBulkPlayByPlayData(isCollege, season, week, timeslot)
	json.NewEncoder(w).Encode(response)
}

// GetLivePlays returns the ordered play-by-play slice for a single game from the
// database.  Used by the frontend streaming client to fetch plays once per game.
// No Firebase reads occur in this handler.
func GetLivePlays(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	gameID := vars["gameID"]

	var response interface{}
	if league == "chl" {
		response = managers.GetCHLLivePlays(gameID)
	} else {
		response = managers.GetPHLLivePlays(gameID)
	}
	json.NewEncoder(w).Encode(response)
}

// RunAdminGames manually triggers the game engine via POST from the Control Room
func RunAdminGames(w http.ResponseWriter, r *http.Request) {
	managers.RunGames()
	json.NewEncoder(w).Encode("Live Broadcast Engine Started!")
}

func TestCHLCronJob(w http.ResponseWriter, r *http.Request) {
	managers.StartCHLLiveStreamingCron()
}

func TestPHLCronJob(w http.ResponseWriter, r *http.Request) {
	managers.StartPHLLiveStreamingCron()
}
