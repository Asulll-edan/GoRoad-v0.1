package ws

import (
	"encoding/json"
	"fmt"
	"log"
)

func (h *MessageHandler) handleVote(client *Client, msg Message) {
	var payload struct {
		VotingID string `json:"voting_id"`
		OptionID string `json:"option_id"`
	}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil || payload.VotingID == "" || payload.OptionID == "" {
		h.sendError(client, "invalid vote data")
		return
	}

	if h.NatsPub != nil {
		voteData, _ := json.Marshal(map[string]interface{}{
			"user_id":   client.UserID,
			"room_id":   client.RoomID,
			"voting_id": payload.VotingID,
			"option_id": payload.OptionID,
		})
		if err := h.NatsPub.Publish("vote.submitted", voteData); err != nil {
			log.Printf("nats publish vote error: %v", err)
		}
	}

	h.sendJSON(client, map[string]interface{}{
		"type":      MsgTypeVote,
		"voting_id": payload.VotingID,
		"status":    "recorded",
	})
}

func (h *MessageHandler) handleConvoyUpdate(client *Client, msg Message) {
	if h.NatsPub == nil {
		return
	}
	payload := map[string]interface{}{
		"user_id": client.UserID,
		"room_id": client.RoomID,
		"data":    msg.Payload,
	}
	data, _ := json.Marshal(payload)
	if err := h.NatsPub.Publish(fmt.Sprintf("convoy.%s.update", client.RoomID), data); err != nil {
		log.Printf("nats publish convoy error: %v", err)
	}
}

func (h *MessageHandler) handleEmergency(client *Client, msg Message) {
	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		h.sendError(client, "invalid emergency data")
		return
	}
	payload["user_id"] = client.UserID
	payload["room_id"] = client.RoomID

	if h.NatsPub != nil {
		data, _ := json.Marshal(payload)
		if err := h.NatsPub.Publish("emergency.alert", data); err != nil {
			log.Printf("nats publish emergency error: %v", err)
		}
	}

	emergencyMsg, _ := json.Marshal(map[string]interface{}{
		"type":    MsgTypeEmergency,
		"payload": payload,
	})
	h.Hub.BroadcastToRoom(client.RoomID, emergencyMsg)
}
