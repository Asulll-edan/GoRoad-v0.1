package fuel

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateFuelLogRequest, userID uuid.UUID) (*FuelLog, error)
	Get(ctx context.Context, id uuid.UUID) (*FuelLog, error)
	List(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]FuelLog, string, bool, error)
	Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) (*FuelLog, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
