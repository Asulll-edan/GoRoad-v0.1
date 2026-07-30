package chat

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID      uuid.UUID              `json:"room_id" gorm:"not null;index"`
	SenderID    uuid.UUID              `json:"sender_id" gorm:"not null"`
	MessageType string                 `json:"message_type" gorm:"default:text"`
	Content     string                 `json:"content" gorm:"not null"`
	ReplyToID   *uuid.UUID             `json:"reply_to_id,omitempty"`
	IsPinned    bool                   `json:"is_pinned" gorm:"default:false"`
	IsDeleted   bool                   `json:"is_deleted" gorm:"default:false"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb;default:'{}'"`
	SentAt      time.Time              `json:"sent_at"`
	EditedAt    *time.Time             `json:"edited_at,omitempty"`
}

type MessageRead struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MessageID uuid.UUID `json:"message_id" gorm:"not null;uniqueIndex:idx_msg_user"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;uniqueIndex:idx_msg_user"`
	ReadAt    time.Time `json:"read_at"`
}

type SendMessageRequest struct {
	RoomID      string `json:"room_id" validate:"required"`
	MessageType string `json:"message_type,omitempty"`
	Content     string `json:"content" validate:"required"`
	ReplyToID   string `json:"reply_to_id,omitempty"`
}
