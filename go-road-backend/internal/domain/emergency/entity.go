package emergency

import (
	"time"

	"github.com/google/uuid"
)

type EmergencyEvent struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID      *uuid.UUID `json:"room_id,omitempty"`
	ReportedBy  uuid.UUID  `json:"reported_by" gorm:"not null"`
	EventType   string     `json:"event_type" gorm:"not null"`
	Severity    string     `json:"severity" gorm:"not null"`
	Location    string     `json:"location,omitempty" gorm:"type:geography(point,4326)"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status" gorm:"default:active"`
	ResolvedBy  *uuid.UUID `json:"resolved_by,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type SOSEvent struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID          uuid.UUID  `json:"user_id" gorm:"not null"`
	RoomID          *uuid.UUID `json:"room_id,omitempty"`
	Location        string     `json:"location" gorm:"type:geography(point,4326);not null"`
	TriggeredAt     time.Time  `json:"triggered_at"`
	Status          string     `json:"status" gorm:"default:active"`
	AcknowledgedBy  *uuid.UUID `json:"acknowledged_by,omitempty"`
	AcknowledgedAt  *time.Time `json:"acknowledged_at,omitempty"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	Notes           string     `json:"notes,omitempty"`
}
