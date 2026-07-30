package upload

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, upload *Upload) error
	FindByID(ctx context.Context, id uuid.UUID) (*Upload, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Upload, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
