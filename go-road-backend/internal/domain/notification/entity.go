package notification

import (
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID              `json:"user_id" gorm:"not null;index"`
	Type      string                 `json:"type" gorm:"not null"`
	Title     string                 `json:"title" gorm:"not null"`
	Body      string                 `json:"body,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty" gorm:"type:jsonb;default:'{}'"`
	IsRead    bool                   `json:"is_read" gorm:"default:false;index"`
	ReadAt    *time.Time             `json:"read_at,omitempty"`
	CreatedAt time.Time              `json:"created_at" gorm:"index"`
}

type NotificationPreference struct {
	ID             uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID         uuid.UUID          `json:"user_id" gorm:"uniqueIndex;not null"`
	PushEnabled    map[string]bool    `json:"push_enabled" gorm:"type:jsonb;default:'{}'"`
	EmailEnabled   map[string]bool    `json:"email_enabled" gorm:"type:jsonb;default:'{}'"`
	QuietHoursFrom string             `json:"quiet_hours_from,omitempty"`
	QuietHoursTo   string             `json:"quiet_hours_to,omitempty"`
}
