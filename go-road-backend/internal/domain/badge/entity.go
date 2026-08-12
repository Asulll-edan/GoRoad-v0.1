package badge

import (
	"time"
	"github.com/google/uuid"
)

type Badge struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Code        string                 `json:"code" gorm:"uniqueIndex;not null"`
	Name        string                 `json:"name" gorm:"not null"`
	Description string                 `json:"description,omitempty"`
	IconURL     string                 `json:"icon_url,omitempty"`
	Category    string                 `json:"category" gorm:"not null"`
	Tier        string                 `json:"tier" gorm:"not null"`
	Criteria    map[string]interface{} `json:"criteria" gorm:"type:jsonb;not null"`
	IsHidden    bool                   `json:"is_hidden" gorm:"default:false"`
	CreatedAt   time.Time              `json:"created_at"`
}

type UserBadge struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"not null;uniqueIndex:idx_user_badge_composite"`
	BadgeID   uuid.UUID  `json:"badge_id" gorm:"not null;uniqueIndex:idx_user_badge_composite"`
	AwardedAt time.Time  `json:"awarded_at"`
	TouringID *uuid.UUID `json:"touring_id,omitempty"`
}

type BadgeProgress struct {
	UserID    uuid.UUID `json:"user_id"`
	BadgeCode string    `json:"badge_code"`
	BadgeName string    `json:"badge_name"`
	Progress  float64   `json:"progress"`
	Target    float64   `json:"target"`
	IsAwarded bool      `json:"is_awarded"`
}
