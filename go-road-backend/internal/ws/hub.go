package ws

import (
	"sync"
)

type Client struct {
	UserID string
	RoomID string
	Send   chan []byte
	Hub    *Hub
}

type Hub struct {
	mu      sync.RWMutex
	rooms   map[string]map[*Client]bool
	userMap map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:   make(map[string]map[*Client]bool),
		userMap: make(map[string]*Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.userMap[client.UserID] = client
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client.RoomID != "" {
		if members, ok := h.rooms[client.RoomID]; ok {
			delete(members, client)
			if len(members) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}
	}
	if h.userMap[client.UserID] == client {
		delete(h.userMap, client.UserID)
	}
}

func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client.RoomID != "" {
		if members, ok := h.rooms[client.RoomID]; ok {
			delete(members, client)
		}
	}
	client.RoomID = roomID
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
}

func (h *Hub) LeaveRoom(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client.RoomID != "" {
		if members, ok := h.rooms[client.RoomID]; ok {
			delete(members, client)
			if len(members) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}
		client.RoomID = ""
	}
}

func (h *Hub) BroadcastToRoom(roomID string, message []byte, excludeUserID ...string) {
	h.mu.RLock()
	members, ok := h.rooms[roomID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	exclude := make(map[string]bool)
	for _, uid := range excludeUserID {
		exclude[uid] = true
	}
	for client := range members {
		if exclude[client.UserID] {
			continue
		}
		select {
		case client.Send <- message:
		default:
		}
	}
}

func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	client, ok := h.userMap[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	select {
	case client.Send <- message:
	default:
	}
}

func (h *Hub) GetRoomClients(roomID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	members, ok := h.rooms[roomID]
	if !ok {
		return nil
	}
	users := make([]string, 0, len(members))
	for client := range members {
		users = append(users, client.UserID)
	}
	return users
}

func (h *Hub) RoomCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}
