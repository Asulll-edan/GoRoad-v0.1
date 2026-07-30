package notification

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, notif *Notification) error
	FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Notification, string, bool, error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)

	UpsertPreferences(ctx context.Context, prefs *NotificationPreference) error
	GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreference, error)
}
