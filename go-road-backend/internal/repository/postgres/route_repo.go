package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/route"
)

type routeRepository struct {
	db *Database
}

func NewRouteRepository(db *Database) domain.Repository {
	return &routeRepository{db: db}
}

func (r *routeRepository) Create(ctx context.Context, route *domain.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *routeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Route, error) {
	var route domain.Route
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&route).Error
	return &route, err
}

func (r *routeRepository) FindByRoomID(ctx context.Context, roomID uuid.UUID) ([]domain.Route, error) {
	var routes []domain.Route
	err := r.db.WithContext(ctx).Where("room_id = ? AND deleted_at IS NULL", roomID).
		Order("created_at DESC").Find(&routes).Error
	return routes, err
}

func (r *routeRepository) Update(ctx context.Context, route *domain.Route) error {
	return r.db.WithContext(ctx).Save(route).Error
}

func (r *routeRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Route{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *routeRepository) AddWaypoint(ctx context.Context, wp *domain.Waypoint) error {
	return r.db.WithContext(ctx).Create(wp).Error
}

func (r *routeRepository) UpdateWaypoint(ctx context.Context, wp *domain.Waypoint) error {
	return r.db.WithContext(ctx).Save(wp).Error
}

func (r *routeRepository) ListWaypoints(ctx context.Context, routeID uuid.UUID) ([]domain.Waypoint, error) {
	var wps []domain.Waypoint
	err := r.db.WithContext(ctx).Where("route_id = ?", routeID).
		Order("order_index ASC").Find(&wps).Error
	return wps, err
}

func (r *routeRepository) DeleteWaypoint(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Waypoint{}, "id = ?", id).Error
}

func (r *routeRepository) ActivateRoute(ctx context.Context, routeID, roomID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	tx.Model(&domain.Route{}).Where("room_id = ?", roomID).Update("is_active", false)
	tx.Model(&domain.Route{}).Where("id = ?", routeID).Update("is_active", true)
	return tx.Commit().Error
}
