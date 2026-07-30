package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (h *MessageHandler) handleLocation(client *Client, msg Message) {
	var payload struct {
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		Speed    float64 `json:"speed"`
		Heading  float64 `json:"heading"`
		Altitude float64 `json:"altitude"`
		Accuracy float64 `json:"accuracy"`
		Battery  int     `json:"battery"`
	}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		h.sendError(client, "invalid location data")
		return
	}

	if h.NatsPub != nil {
		locData, _ := json.Marshal(map[string]interface{}{
			"user_id":  client.UserID,
			"room_id":  client.RoomID,
			"lat":      payload.Lat,
			"lon":      payload.Lon,
			"speed":    payload.Speed,
			"heading":  payload.Heading,
			"altitude": payload.Altitude,
			"accuracy": payload.Accuracy,
			"battery":  payload.Battery,
			"ts":       time.Now().Unix(),
		})
		if err := h.NatsPub.Publish(fmt.Sprintf("room.%s.location", client.RoomID), locData); err != nil {
			log.Printf("nats publish location error: %v", err)
		}
	}

	if h.RedisCache != nil {
		posKey := fmt.Sprintf("pos:room:%s:rider:%s", client.RoomID, client.UserID)
		posData, _ := json.Marshal(map[string]interface{}{
			"lat":     payload.Lat,
			"lon":     payload.Lon,
			"speed":   payload.Speed,
			"heading": payload.Heading,
			"ts":      time.Now().Unix(),
		})
		_ = h.RedisCache.HSet(nil, posKey, "latest", string(posData))
	}

	broadcastPayload, _ := json.Marshal(map[string]interface{}{
		"type":    MsgTypeLocation,
		"user_id": client.UserID,
		"lat":     payload.Lat,
		"lon":     payload.Lon,
		"speed":   payload.Speed,
		"heading": payload.Heading,
		"ts":      time.Now().Unix(),
	})
	h.Hub.BroadcastToRoom(client.RoomID, broadcastPayload, client.UserID)
}
