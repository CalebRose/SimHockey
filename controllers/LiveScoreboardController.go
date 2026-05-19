package controllers

import (
	"fmt"
	"net/http"

	"github.com/CalebRose/SimHockey/managers"
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
