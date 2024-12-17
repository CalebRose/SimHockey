package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/websocket"
)

// Upgrader configures the WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should validate origins for security
		return true
	},
}

// Global map to keep track of connected WebSocket clients
var (
	clients = make(map[*websocket.Conn]bool)
	mu      sync.Mutex
)

// WebSocketHandler handles WebSocket connection requests
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Add new WebSocket connection to the clients map
	mu.Lock()
	clients[conn] = true
	mu.Unlock()
	log.Println("New WebSocket client connected")

	defer func() {
		conn.Close()
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		log.Println("WebSocket client disconnected")
	}()

	ts := managers.GetTimestamp()
	// Send the latest timestamp to the user immediately upon connection
	err = conn.WriteJSON(ts)
	if err != nil {
		log.Println("Error sending timestamp:", err)
		return
	}

	// Listen for messages (e.g., handle disconnects)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket client error:", err)
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Error writing WebSocket message:", err)
			break
		}
	}
}
