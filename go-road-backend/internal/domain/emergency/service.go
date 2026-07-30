package emergency

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	ReportEmergency(ctx context.Context, req CreateEmergencyRequest, userID uuid.UUID) (*EmergencyEvent, error)
	ListEmergencies(ctx context.Context, cursor string, limit int, status string) ([]EmergencyEvent, string, bool, error)
	GetEmergency(ctx context.Context, id uuid.UUID) (*EmergencyEvent, error)
	AcknowledgeEmergency(ctx context.Context, id, userID uuid.UUID) error
	ResolveEmergency(ctx context.Context, id, userID uuid.UUID) error

	TriggerSOS(ctx context.Context, userID uuid.UUID, lat, lng float64, roomID *uuid.UUID) (*SOSEvent, error)
	DismissSOS(ctx context.Context, id, userID uuid.UUID) error
}

type CreateEmergencyRequest struct {
	RoomID      string  `json:"room_id,omitempty"`
	EventType   string  `json:"event_type" validate:"required"`
	Severity    string  `json:"severity" validate:"required"`
	Lat         float64 `json:"lat,omitempty"`
	Lng         float64 `json:"lng,omitempty"`
	Description string  `json:"description,omitempty"`
}
