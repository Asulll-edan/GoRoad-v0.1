package notification

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	ListNotifications(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Notification, string, bool, error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
	UpdatePreferences(ctx context.Context, userID uuid.UUID, prefs map[string]interface{}) error
}
