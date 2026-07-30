package convoy

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	CreateFormation(ctx context.Context, roomID uuid.UUID, name string, formationType string, userID uuid.UUID) (*Formation, error)
	UpdateFormation(ctx context.Context, formationID uuid.UUID, req map[string]interface{}, userID uuid.UUID) error
	GetActiveFormation(ctx context.Context, roomID uuid.UUID) (*Formation, error)

	UpdateLocation(ctx context.Context, roomID, userID uuid.UUID, lat, lng, speed, heading float64) error
	GetLocations(ctx context.Context, roomID uuid.UUID) ([]RiderPosition, error)
	GetTracking(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]RiderPosition, string, bool, error)
}
