package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/chat"
)

type chatRepository struct {
	db *Database
}

func NewChatRepository(db *Database) domain.Repository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(ctx context.Context, msg *domain.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *chatRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	var msg domain.Message
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = false", id).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *chatRepository) FindByRoomID(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]domain.Message, string, bool, error) {
	var msgs []domain.Message
	db := r.db.WithContext(ctx).Where("room_id = ? AND is_deleted = false", roomID)

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				SentAt time.Time `json:"sent_at"`
				ID     string    `json:"id"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(sent_at, id) < (?, ?)", c.SentAt, c.ID)
			}
		}
	}

	db = db.Order("sent_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&msgs).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(msgs) > limit
	if hasMore {
		msgs = msgs[:limit]
	}

	var nextCursor string
	if len(msgs) > 0 {
		last := msgs[len(msgs)-1]
		cursorData, _ := json.Marshal(struct {
			SentAt time.Time `json:"sent_at"`
			ID     string    `json:"id"`
		}{SentAt: last.SentAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return msgs, nextCursor, hasMore, nil
}

func (r *chatRepository) Update(ctx context.Context, msg *domain.Message) error {
	return r.db.WithContext(ctx).Save(msg).Error
}

func (r *chatRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Message{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *chatRepository) Pin(ctx context.Context, id uuid.UUID, pinned bool) error {
	return r.db.WithContext(ctx).Model(&domain.Message{}).Where("id = ?", id).Update("is_pinned", pinned).Error
}

func (r *chatRepository) MarkRead(ctx context.Context, msgID, userID uuid.UUID) error {
	read := domain.MessageRead{
		MessageID: msgID,
		UserID:    userID,
		ReadAt:    time.Now(),
	}
	return r.db.WithContext(ctx).Create(&read).Error
}

func (r *chatRepository) GetUnreadCount(ctx context.Context, roomID, userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&domain.Message{}).
		Where("room_id = ? AND sender_id != ? AND is_deleted = false AND id NOT IN (SELECT message_id FROM message_reads WHERE user_id = ?)", roomID, userID, userID).
		Count(&total).Error
	return total, err
}
