package expense

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateExpenseRequest, userID uuid.UUID) (*Expense, error)
	Get(ctx context.Context, id uuid.UUID) (*Expense, error)
	List(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Expense, string, bool, error)
	Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) (*Expense, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
