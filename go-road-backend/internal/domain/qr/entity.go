package qr

import (
	"time"
	"github.com/google/uuid"
)

type QRCard struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"uniqueIndex;not null"`
	Code      string    `json:"code" gorm:"uniqueIndex;not null"`
	Style     string    `json:"style" gorm:"default:default"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	ScanCount int       `json:"scan_count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
