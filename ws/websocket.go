package ws

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/CalebRose/SimHockey/structs"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	sendBuffer = 4
)

// client owns a single connection. Only writePump writes to conn,
// which satisfies gorilla's one-concurrent-writer rule.
type client struct {
	conn *websocket.Conn
	send chan structs.Timestamp
}

// Upgrader configures the WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

var (
	clients = make(map[*client]struct{})
	mu      sync.Mutex
)

// checkOrigin allows requests with no Origin header (non-browser clients)
// and origins listed in ORIGIN_ALLOWED (comma-separated, or "*").
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	for _, allowed := range strings.Split(os.Getenv("ORIGIN_ALLOWED"), ",") {
		allowed = strings.TrimSpace(allowed)
		if allowed == "*" || strings.EqualFold(allowed, origin) {
			return true
		}
	}
	return false
}

func removeClient(c *client) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := clients[c]; ok {
		delete(clients, c)
		close(c.send)
	}
}

// WebSocketHandler handles WebSocket connection requests
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	c := &client{conn: conn, send: make(chan structs.Timestamp, sendBuffer)}

	// Queue the latest timestamp before registering so it is the first message sent.
	c.send <- currentTimestamp()

	mu.Lock()
	clients[c] = struct{}{}
	mu.Unlock()
	log.Println("New WebSocket client connected")

	go c.writePump()
	c.readPump()
}

// readPump discards incoming messages; it exists to process control frames
// (pong/close) and detect dead connections via the read deadline.
func (c *client) readPump() {
	defer func() {
		removeClient(c)
		c.conn.Close()
		log.Println("WebSocket client disconnected")
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("WebSocket client error:", err)
			}
			return
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case ts, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(ts); err != nil {
				log.Println("Error writing WebSocket message:", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
