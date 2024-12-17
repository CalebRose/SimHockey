package ws

import (
	"log"

	"github.com/CalebRose/SimHockey/structs"
)

func BroadcastTSUpdate(ts structs.Timestamp) error {
	mu.Lock()
	defer mu.Unlock()
	for conn := range clients {
		err := conn.WriteJSON(ts)
		if err != nil {
			log.Println("Error broadcasting to WebSocket client:", err)
			conn.Close()
			delete(clients, conn) // Remove client on error
		}
	}
	return nil
}
