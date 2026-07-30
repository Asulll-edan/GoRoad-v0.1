package event

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

type EventHandler func(subject string, data []byte)

func NewSubscriber(conn *nats.Conn) (*Subscriber, error) {
	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}
	return &Subscriber{conn: conn, js: js}, nil
}

func (s *Subscriber) Subscribe(subject string, handler EventHandler) (*nats.Subscription, error) {
	return s.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
	})
}

func (s *Subscriber) SubscribeJS(subject string, handler EventHandler, durable string) (*nats.Subscription, error) {
	return s.js.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
		msg.Ack()
	}, nats.Durable(durable), nats.ManualAck())
}

func (s *Subscriber) SubscribeQueue(subject, queue string, handler EventHandler) (*nats.Subscription, error) {
	return s.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
	})
}

func HandleLocationUpdate(hubLocHandler interface{}) EventHandler {
	return func(subject string, data []byte) {
		log.Printf("location update: %s", subject)
	}
}

func HandleChatMessage(broadcastFn func(roomID string, msg []byte)) EventHandler {
	return func(subject string, data []byte) {
		var payload struct {
			RoomID  string `json:"room_id"`
			UserID  string `json:"user_id"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			log.Printf("chat message parse error: %v", err)
			return
		}
		if broadcastFn != nil {
			broadcastFn(payload.RoomID, data)
		}
	}
}
