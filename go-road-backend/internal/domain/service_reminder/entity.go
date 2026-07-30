package service_reminder

import (
	"time"
	"github.com/google/uuid"
)

type ServiceReminder struct {
	ID                    uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID                uuid.UUID  `json:"user_id" gorm:"not null;index"`
	MotorID               uuid.UUID  `json:"motor_id" gorm:"not null"`
	ServiceType           string     `json:"service_type" gorm:"not null"`
	Title                 string     `json:"title" gorm:"not null"`
	Description           string     `json:"description,omitempty"`
	DueDate               time.Time  `json:"due_date" gorm:"not null;index"`
	DueOdometer           float64    `json:"due_odometer,omitempty"`
	CompletedAt           *time.Time `json:"completed_at,omitempty"`
	IsRecurring           bool       `json:"is_recurring" gorm:"default:false"`
	RecurringIntervalDays int        `json:"recurring_interval_days,omitempty"`
	RecurringIntervalKm   float64    `json:"recurring_interval_km,omitempty"`
	NotifiedH7            bool       `json:"notified_h7" gorm:"default:false"`
	NotifiedH1            bool       `json:"notified_h1" gorm:"default:false"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
