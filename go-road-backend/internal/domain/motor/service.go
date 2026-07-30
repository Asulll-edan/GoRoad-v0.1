package motor

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateMotorRequest, userID uuid.UUID) (*Motor, error)
	Get(ctx context.Context, id uuid.UUID) (*Motor, error)
	List(ctx context.Context, userID uuid.UUID) ([]Motor, error)
	Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) (*Motor, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	SetPrimary(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
