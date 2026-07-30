package upload

import (
	"time"
	"github.com/google/uuid"
)

type Upload struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID          uuid.UUID `json:"user_id" gorm:"not null;index"`
	FileName        string    `json:"file_name" gorm:"not null"`
	FileSize        int64     `json:"file_size" gorm:"not null"`
	MimeType        string    `json:"mime_type" gorm:"not null"`
	URL             string    `json:"url" gorm:"not null"`
	Bucket          string    `json:"bucket" gorm:"not null"`
	ObjectKey       string    `json:"object_key" gorm:"not null"`
	Category        string    `json:"category" gorm:"default:general"`
	IsPublic        bool      `json:"is_public" gorm:"default:false"`
	VirusScanStatus string   `json:"virus_scan_status" gorm:"default:pending"`
	CreatedAt       time.Time `json:"created_at"`
}
