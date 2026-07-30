package qr

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, card *QRCard) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*QRCard, error)
	FindByCode(ctx context.Context, code string) (*QRCard, error)
	IncrementScanCount(ctx context.Context, id uuid.UUID) error
}
