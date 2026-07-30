package expense

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, expense *Expense) error
	FindByID(ctx context.Context, id uuid.UUID) (*Expense, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Expense, string, bool, error)
	Update(ctx context.Context, expense *Expense) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
