package ws

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/CalebRose/SimHockey/managers"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// Upgrader configures the WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Global map to keep track of connected WebSocket clients
var (
	clients = make(map[*websocket.Conn]chan interface{})
	mu      sync.Mutex
)

// writePump serializes all writes for a single connection and sends periodic pings.
// All writes must go through this goroutine to satisfy gorilla/websocket's
// requirement of at most one concurrent writer per connection.
func writePump(conn *websocket.Conn, send <-chan interface{}) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("WebSocket ping error:", err)
				return
			}
		}
	}
}

// WebSocketHandler handles WebSocket connection requests
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// Send the initial timestamp before registering in the broadcast pool
	// to prevent a race with BroadcastTSUpdate writing to the same connection.
	ts := managers.GetTimestamp()
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteJSON(ts); err != nil {
		log.Println("Error sending initial timestamp:", err)
		conn.Close()
		return
	}

	send := make(chan interface{}, 256)

	mu.Lock()
	clients[conn] = send
	mu.Unlock()
	log.Println("New WebSocket client connected")

	defer func() {
		mu.Lock()
		if _, ok := clients[conn]; ok {
			delete(clients, conn)
			close(send)
		}
		mu.Unlock()
		conn.Close()
		log.Println("WebSocket client disconnected")
	}()

	go writePump(conn, send)

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
	}
}
