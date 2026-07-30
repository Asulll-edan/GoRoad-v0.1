package fuel

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, log *FuelLog) error
	FindByID(ctx context.Context, id uuid.UUID) (*FuelLog, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]FuelLog, string, bool, error)
	Update(ctx context.Context, log *FuelLog) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
