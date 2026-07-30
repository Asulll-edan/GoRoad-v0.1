package ws

import (
	"encoding/json"
	"log"
	"time"
)

func (h *MessageHandler) handlePresence(client *Client) {
	presenceMsg, _ := json.Marshal(map[string]interface{}{
		"type":     MsgTypePresence,
		"user_id":  client.UserID,
		"status":   "online",
		"room_id":  client.RoomID,
		"timestamp": time.Now().Unix(),
	})
	h.Hub.BroadcastToRoom(client.RoomID, presenceMsg, client.UserID)
}

func StartPresenceChecker(hub *Hub, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		hub.mu.RLock()
		for _, client := range hub.userMap {
			if client.RoomID == "" {
				continue
			}
			presenceMsg, err := json.Marshal(map[string]interface{}{
				"type":     MsgTypePresence,
				"user_id":  client.UserID,
				"status":   "online",
				"room_id":  client.RoomID,
				"timestamp": time.Now().Unix(),
			})
			if err != nil {
				log.Printf("presence marshal error: %v", err)
				continue
			}
			hub.BroadcastToRoom(client.RoomID, presenceMsg, client.UserID)
		}
		hub.mu.RUnlock()
	}
}
