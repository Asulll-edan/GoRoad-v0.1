package motor

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, motor *Motor) error
	FindByID(ctx context.Context, id uuid.UUID) (*Motor, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Motor, error)
	Update(ctx context.Context, motor *Motor) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	SetPrimary(ctx context.Context, userID, motorID uuid.UUID) error
}
