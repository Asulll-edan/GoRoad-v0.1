package ws

import (
	"encoding/json"
	"log"
)

type MessageHandler struct {
	Hub          *Hub
	NatsPub      NATSPublisher
	RedisCache   RedisCache
	JWTSecret    string
}

type NATSPublisher interface {
	Publish(subject string, data interface{}) error
}

type RedisCache interface {
	HSet(ctx interface{}, key, field string, value interface{}) error
	HGetAll(ctx interface{}, key string) (map[string]string, error)
	HDel(ctx interface{}, key string, fields ...string) error
	Publish(ctx interface{}, channel string, message interface{}) error
	IncrWithExpiry(ctx interface{}, key string, ttl interface{}) (int64, error)
}

func NewMessageHandler(hub *Hub, natsPub NATSPublisher, redisCache RedisCache, jwtSecret string) *MessageHandler {
	return &MessageHandler{
		Hub:        hub,
		NatsPub:    natsPub,
		RedisCache: redisCache,
		JWTSecret:  jwtSecret,
	}
}

func (h *MessageHandler) HandleMessage(client *Client, msg Message) {
	switch msg.Type {
	case MsgTypeHeartbeat:
		h.handleHeartbeat(client, msg)
	case MsgTypeJoinRoom:
		h.handleJoinRoom(client, msg)
	case MsgTypeLeaveRoom:
		h.Hub.LeaveRoom(client)
	case MsgTypeLocation:
		h.handleLocation(client, msg)
	case MsgTypeChat:
		h.handleChat(client, msg)
	case MsgTypeChatTyping:
		h.handleTyping(client, msg)
	case MsgTypeVote:
		h.handleVote(client, msg)
	case MsgTypeConvoyUpdate:
		h.handleConvoyUpdate(client, msg)
	case MsgTypeEmergency:
		h.handleEmergency(client, msg)
	default:
		h.sendError(client, "unknown message type: "+msg.Type)
	}
}

func (h *MessageHandler) sendError(client *Client, errMsg string) {
	data, _ := json.Marshal(map[string]string{"type": MsgTypeError, "message": errMsg})
	select {
	case client.Send <- data:
	default:
	}
}

func (h *MessageHandler) sendJSON(client *Client, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	select {
	case client.Send <- data:
	default:
	}
}

func (h *MessageHandler) handleHeartbeat(client *Client, msg Message) {
	h.sendJSON(client, map[string]string{"type": MsgTypeHeartbeat, "status": "ok"})
}

func (h *MessageHandler) handleJoinRoom(client *Client, msg Message) {
	var payload struct {
		RoomID string `json:"room_id"`
	}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil || payload.RoomID == "" {
		h.sendError(client, "invalid room_id")
		return
	}
	h.Hub.JoinRoom(client, payload.RoomID)
	h.sendJSON(client, map[string]string{"type": MsgTypeJoinRoom, "room_id": payload.RoomID, "status": "joined"})
	h.Hub.BroadcastToRoom(payload.RoomID, []byte(
		`{"type":"presence","user_id":"`+client.UserID+`","status":"online"}`,
	), client.UserID)
}
