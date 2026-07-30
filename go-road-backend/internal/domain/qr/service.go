package qr

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	GetMyQRCard(ctx context.Context, userID uuid.UUID) (*QRCard, error)
	RegenerateQR(ctx context.Context, userID uuid.UUID) (*QRCard, error)
	ScanQR(ctx context.Context, code string) (*QRCard, error)
}
