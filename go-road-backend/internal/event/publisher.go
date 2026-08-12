package event

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewPublisher(conn *nats.Conn) (*Publisher, error) {
	if conn == nil {
		return nil, errors.New("nil nats connection")
	}
	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}
	return &Publisher{conn: conn, js: js}, nil
}

func (p *Publisher) Publish(subject string, data interface{}) error {
	if p == nil || p.conn == nil {
		return errors.New("nil publisher")
	}
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.conn.Publish(subject, msg)
}

func (p *Publisher) PublishJS(subject string, data interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = p.js.Publish(subject, msg)
	return err
}

func (p *Publisher) PublishLocation(roomID, userID string, lat, lon, speed, heading float64) {
	data := map[string]interface{}{
		"user_id": userID,
		"room_id": roomID,
		"lat":     lat,
		"lon":     lon,
		"speed":   speed,
		"heading": heading,
	}
	subject := "room." + roomID + ".location"
	if err := p.Publish(subject, data); err != nil {
		log.Printf("publish location error: %v", err)
	}
}

func (p *Publisher) PublishChat(roomID, userID, message string) {
	data := map[string]interface{}{
		"user_id": userID,
		"room_id": roomID,
		"message": message,
	}
	subject := "room." + roomID + ".chat"
	if err := p.Publish(subject, data); err != nil {
		log.Printf("publish chat error: %v", err)
	}
}

func (p *Publisher) PublishEmergency(roomID, userID string, details interface{}) {
	subject := "emergency.alert"
	data := map[string]interface{}{
		"room_id": roomID,
		"user_id": userID,
		"details": details,
	}
	if err := p.PublishJS(subject, data); err != nil {
		log.Printf("publish emergency error: %v", err)
	}
}

func (p *Publisher) PublishNotification(userID string, notification interface{}) {
	subject := "notification.send"
	data := map[string]interface{}{
		"user_id":      userID,
		"notification": notification,
	}
	if err := p.PublishJS(subject, data); err != nil {
		log.Printf("publish notification error: %v", err)
	}
}

func (p *Publisher) PublishLeaderboardUpdate() {
	if err := p.Publish("leaderboard.updated", map[string]string{"status": "recompute"}); err != nil {
		log.Printf("publish leaderboard update error: %v", err)
	}
}

func (p *Publisher) JoinStream(name string) error {
	_, err := p.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: []string{name + ".>"},
	})
	return err
}
