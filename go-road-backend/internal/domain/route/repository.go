package route

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, route *Route) error
	FindByID(ctx context.Context, id uuid.UUID) (*Route, error)
	FindByRoomID(ctx context.Context, roomID uuid.UUID) ([]Route, error)
	Update(ctx context.Context, route *Route) error
	SoftDelete(ctx context.Context, id uuid.UUID) error

	AddWaypoint(ctx context.Context, wp *Waypoint) error
	UpdateWaypoint(ctx context.Context, wp *Waypoint) error
	ListWaypoints(ctx context.Context, routeID uuid.UUID) ([]Waypoint, error)
	DeleteWaypoint(ctx context.Context, id uuid.UUID) error

	ActivateRoute(ctx context.Context, routeID, roomID uuid.UUID) error
}
