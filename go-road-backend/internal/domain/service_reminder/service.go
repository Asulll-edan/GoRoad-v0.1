package service_reminder

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateReminderRequest, userID uuid.UUID) (*ServiceReminder, error)
	List(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]ServiceReminder, string, bool, error)
	Get(ctx context.Context, id uuid.UUID) (*ServiceReminder, error)
	Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Complete(ctx context.Context, id, userID uuid.UUID) error
}

type CreateReminderRequest struct {
	MotorID           string `json:"motor_id" validate:"required"`
	ServiceType       string `json:"service_type" validate:"required"`
	Title             string `json:"title" validate:"required"`
	Description       string `json:"description,omitempty"`
	DueDate           string `json:"due_date" validate:"required"`
	DueOdometer       float64 `json:"due_odometer,omitempty"`
	IsRecurring       bool   `json:"is_recurring,omitempty"`
	RecurringIntervalDays int `json:"recurring_interval_days,omitempty"`
}
