package emergency

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *EmergencyEvent) error
	FindEventByID(ctx context.Context, id uuid.UUID) (*EmergencyEvent, error)
	ListEvents(ctx context.Context, cursor string, limit int, status string) ([]EmergencyEvent, string, bool, error)
	UpdateEvent(ctx context.Context, event *EmergencyEvent) error

	CreateSOS(ctx context.Context, sos *SOSEvent) error
	FindSOSByID(ctx context.Context, id uuid.UUID) (*SOSEvent, error)
	ListSOS(ctx context.Context, cursor string, limit int) ([]SOSEvent, string, bool, error)
	UpdateSOS(ctx context.Context, sos *SOSEvent) error
}
