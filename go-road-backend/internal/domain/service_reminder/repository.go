package service_reminder

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, r *ServiceReminder) error
	FindByID(ctx context.Context, id uuid.UUID) (*ServiceReminder, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]ServiceReminder, string, bool, error)
	Update(ctx context.Context, r *ServiceReminder) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Complete(ctx context.Context, id uuid.UUID) error
	FindDueReminders(ctx context.Context, dueBefore time.Time) ([]ServiceReminder, error)
}
