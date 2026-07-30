package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/notification"
)

type notificationRepository struct {
	db *Database
}

func NewNotificationRepository(db *Database) domain.Repository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notif *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Notification, string, bool, error) {
	var notifs []domain.Notification
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&notifs).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(notifs) > limit
	if hasMore {
		notifs = notifs[:limit]
	}
	var nextCursor string
	if len(notifs) > 0 {
		last := notifs[len(notifs)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return notifs, nextCursor, hasMore, nil
}

func (r *notificationRepository) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": time.Now()}).Error
}

func (r *notificationRepository) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": time.Now()}).Error
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}

func (r *notificationRepository) UpsertPreferences(ctx context.Context, prefs *domain.NotificationPreference) error {
	return r.db.WithContext(ctx).Save(prefs).Error
}

func (r *notificationRepository) GetPreferences(ctx context.Context, userID uuid.UUID) (*domain.NotificationPreference, error) {
	var p domain.NotificationPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error
	return &p, err
}
