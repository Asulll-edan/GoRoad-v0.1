package convoy

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateFormation(ctx context.Context, f *Formation) error
	FindActiveFormation(ctx context.Context, roomID uuid.UUID) (*Formation, error)
	UpdateFormation(ctx context.Context, f *Formation) error

	SavePosition(ctx context.Context, pos *RiderPosition) error
	GetLatestPosition(ctx context.Context, roomID, userID uuid.UUID) (*RiderPosition, error)
	GetPositionsInRoom(ctx context.Context, roomID uuid.UUID) ([]RiderPosition, error)
	GetPositionsSince(ctx context.Context, roomID uuid.UUID, since time.Time) ([]RiderPosition, error)
}
