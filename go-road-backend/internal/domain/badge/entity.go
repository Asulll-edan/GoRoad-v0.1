package badge

import (
	"time"
	"github.com/google/uuid"
)

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
