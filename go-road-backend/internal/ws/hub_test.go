package ws

import (
	"testing"
)

func TestHubNew(t *testing.T) {
	h := NewHub()
	if h == nil {
		t.Fatal("expected hub to be non-nil")
	}
}

func TestHubRegisterUnregister(t *testing.T) {
	h := NewHub()
	c := &Client{UserID: "user1", Send: make(chan []byte, 1), Hub: h}

	h.Register(c)
	if _, ok := h.userMap["user1"]; !ok {
		t.Error("expected user1 to be registered")
	}

	h.Unregister(c)
	if _, ok := h.userMap["user1"]; ok {
		t.Error("expected user1 to be unregistered")
	}
}

func TestHubJoinRoom(t *testing.T) {
	h := NewHub()
	c := &Client{UserID: "user1", Send: make(chan []byte, 1), Hub: h}

	h.Register(c)
	h.JoinRoom(c, "room1")

	if c.RoomID != "room1" {
		t.Errorf("expected room1, got %s", c.RoomID)
	}
	if h.RoomCount("room1") != 1 {
		t.Errorf("expected 1 member, got %d", h.RoomCount("room1"))
	}
}

func TestHubLeaveRoom(t *testing.T) {
	h := NewHub()
	c := &Client{UserID: "user1", Send: make(chan []byte, 1), Hub: h}

	h.Register(c)
	h.JoinRoom(c, "room1")
	h.LeaveRoom(c)

	if c.RoomID != "" {
		t.Error("expected empty room after leave")
	}
	if h.RoomCount("room1") != 0 {
		t.Error("expected 0 members after leave")
	}
}

func TestHubBroadcastToRoom(t *testing.T) {
	h := NewHub()
	c1 := &Client{UserID: "user1", Send: make(chan []byte, 10), Hub: h}
	c2 := &Client{UserID: "user2", Send: make(chan []byte, 10), Hub: h}

	h.Register(c1)
	h.Register(c2)
	h.JoinRoom(c1, "room1")
	h.JoinRoom(c2, "room1")

	h.BroadcastToRoom("room1", []byte("hello"))

	select {
	case msg := <-c1.Send:
		if string(msg) != "hello" {
			t.Errorf("expected hello, got %s", msg)
		}
	default:
		t.Error("expected c1 to receive message")
	}

	select {
	case msg := <-c2.Send:
		if string(msg) != "hello" {
			t.Errorf("expected hello, got %s", msg)
		}
	default:
		t.Error("expected c2 to receive message")
	}
}

func TestHubGetRoomClients(t *testing.T) {
	h := NewHub()
	c := &Client{UserID: "user1", Send: make(chan []byte, 1), Hub: h}

	h.Register(c)
	h.JoinRoom(c, "room1")

	users := h.GetRoomClients("room1")
	if len(users) != 1 || users[0] != "user1" {
		t.Errorf("expected [user1], got %v", users)
	}
}

func TestHubSendToUser(t *testing.T) {
	h := NewHub()
	c := &Client{UserID: "user1", Send: make(chan []byte, 10), Hub: h}
	h.Register(c)

	h.SendToUser("user1", []byte("ping"))

	select {
	case msg := <-c.Send:
		if string(msg) != "ping" {
			t.Errorf("expected ping, got %s", msg)
		}
	default:
		t.Error("expected user1 to receive message")
	}
}
