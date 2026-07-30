package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (h *MessageHandler) handleChat(client *Client, msg Message) {
	var payload struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil || payload.Message == "" {
		h.sendError(client, "invalid chat message")
		return
	}

	chatMsg := map[string]interface{}{
		"type":       MsgTypeChat,
		"user_id":    client.UserID,
		"message":    payload.Message,
		"message_type": payload.Type,
		"room_id":    client.RoomID,
		"created_at": time.Now().Unix(),
	}

	if h.NatsPub != nil {
		data, _ := json.Marshal(chatMsg)
		if err := h.NatsPub.Publish(fmt.Sprintf("room.%s.chat", client.RoomID), data); err != nil {
			log.Printf("nats publish chat error: %v", err)
		}
	}

	broadcast, _ := json.Marshal(chatMsg)
	h.Hub.BroadcastToRoom(client.RoomID, broadcast)
}

func (h *MessageHandler) handleTyping(client *Client, msg Message) {
	var payload struct {
		IsTyping bool `json:"is_typing"`
	}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	typingMsg, _ := json.Marshal(map[string]interface{}{
		"type":      MsgTypeChatTyping,
		"user_id":   client.UserID,
		"is_typing": payload.IsTyping,
		"room_id":   client.RoomID,
	})
	h.Hub.BroadcastToRoom(client.RoomID, typingMsg, client.UserID)
}
