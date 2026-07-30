package route

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateRouteRequest, userID uuid.UUID) (*Route, error)
	GetRoute(ctx context.Context, id uuid.UUID) (*Route, error)
	Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	AddWaypoint(ctx context.Context, routeID uuid.UUID, wp WaypointInput) (*Waypoint, error)
	ListWaypoints(ctx context.Context, routeID uuid.UUID) ([]Waypoint, error)
	Activate(ctx context.Context, id, roomID uuid.UUID) error
}
