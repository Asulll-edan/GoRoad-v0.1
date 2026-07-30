package convoy

import (
	"time"
	"github.com/google/uuid"
)

type Formation struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID          uuid.UUID `json:"room_id" gorm:"not null;unique"`
	Name            string    `json:"name" gorm:"not null"`
	FormationType   string    `json:"formation_type" gorm:"default:single_line"`
	MemberOrder     []string  `json:"member_order" gorm:"type:uuid[]"`
	SpeedLimitKmh   float64   `json:"speed_limit_kmh,omitempty"`
	SafeDistanceM   float64   `json:"safe_distance_m,omitempty"`
	IsActive        bool      `json:"is_active" gorm:"default:false"`
	CreatedBy       uuid.UUID `json:"created_by" gorm:"not null"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RiderPosition struct {
	RoomID      uuid.UUID `json:"room_id" gorm:"not null"`
	UserID      uuid.UUID `json:"user_id" gorm:"not null"`
	Lat         float64   `json:"lat" gorm:"not null"`
	Lng         float64   `json:"lng" gorm:"not null"`
	SpeedKmh    float64   `json:"speed_kmh,omitempty"`
	Heading     float64   `json:"heading,omitempty"`
	Altitude    float64   `json:"altitude,omitempty"`
	Battery     float64   `json:"battery,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}
