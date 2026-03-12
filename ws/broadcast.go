package ws

import (
	"github.com/CalebRose/SimHockey/structs"
)

func BroadcastTSUpdate(ts structs.Timestamp) error {
	mu.Lock()
	defer mu.Unlock()
	for _, send := range clients {
		select {
		case send <- ts:
		default:
			// Client send buffer is full; skip this update.
			// The ping/pong mechanism will detect and clean up dead connections.
		}
	}
	return nil
}
