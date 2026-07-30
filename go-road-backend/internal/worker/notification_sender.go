package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type NotificationSender struct {
	nc *nats.Conn
}

type Notification struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

func NewNotificationSender(nc *nats.Conn) *NotificationSender {
	return &NotificationSender{nc: nc}
}

func (s *NotificationSender) Start(ctx context.Context) {
	_, err := s.nc.Subscribe("notification.send", func(msg *nats.Msg) {
		var notif Notification
		if err := json.Unmarshal(msg.Data, &notif); err != nil {
			log.Printf("notification parse error: %v", err)
			return
		}
		s.sendPush(notif)
	})
	if err != nil {
		log.Printf("notification subscribe error: %v", err)
	}
	<-ctx.Done()
}

func (s *NotificationSender) sendPush(notif Notification) {
	log.Printf("sending push notification to user %s: %s", notif.UserID, notif.Title)
}
