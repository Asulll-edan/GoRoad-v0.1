package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	RoomID  string          `json:"room_id,omitempty"`
	UserID  string          `json:"user_id,omitempty"`
}

type ClientHandler interface {
	HandleMessage(client *Client, msg Message)
}

type WSClient struct {
	*Client
	conn     *websocket.Conn
	handler  ClientHandler
	done     chan struct{}
}

func NewWSClient(hub *Hub, conn *websocket.Conn, userID string, handler ClientHandler) *WSClient {
	client := &Client{
		UserID: userID,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}
	return &WSClient{
		Client:  client,
		conn:    conn,
		handler: handler,
		done:    make(chan struct{}),
	}
}

func (c *WSClient) Start() {
	c.Hub.Register(c.Client)
	go c.writePump()
	go c.readPump()
}

func (c *WSClient) Stop() {
	close(c.done)
	c.conn.Close()
	c.Hub.Unregister(c.Client)
}

func (c *WSClient) readPump() {
	defer c.Stop()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("ws parse error: %v", err)
			continue
		}
		msg.UserID = c.UserID

		if c.handler != nil {
			c.handler.HandleMessage(c.Client, msg)
		}
	}
}

func (c *WSClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.done:
			return
		}
	}
}
