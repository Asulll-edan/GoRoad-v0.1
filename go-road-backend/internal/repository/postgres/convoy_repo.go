package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/convoy"
)

type convoyRepository struct {
	db *Database
}

func NewConvoyRepository(db *Database) domain.Repository {
	return &convoyRepository{db: db}
}

func (r *convoyRepository) CreateFormation(ctx context.Context, f *domain.Formation) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *convoyRepository) FindActiveFormation(ctx context.Context, roomID uuid.UUID) (*domain.Formation, error) {
	var f domain.Formation
	err := r.db.WithContext(ctx).Where("room_id = ? AND is_active = true", roomID).First(&f).Error
	return &f, err
}

func (r *convoyRepository) UpdateFormation(ctx context.Context, f *domain.Formation) error {
	return r.db.WithContext(ctx).Save(f).Error
}

func (r *convoyRepository) SavePosition(ctx context.Context, pos *domain.RiderPosition) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO rider_locations (time, room_id, user_id, location, speed_kmh, heading, altitude)
		 VALUES (?, ?, ?, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?, ?, 0)`,
		pos.Timestamp, pos.RoomID, pos.UserID, pos.Lng, pos.Lat, pos.SpeedKmh, pos.Heading,
	).Error
}

func (r *convoyRepository) GetLatestPosition(ctx context.Context, roomID, userID uuid.UUID) (*domain.RiderPosition, error) {
	var pos domain.RiderPosition
	err := r.db.WithContext(ctx).Raw(`
		SELECT time as timestamp, room_id, user_id,
			ST_Y(location::geometry) as lat, ST_X(location::geometry) as lng,
			speed_kmh, heading
		FROM rider_locations
		WHERE room_id = ? AND user_id = ?
		ORDER BY time DESC LIMIT 1`, roomID, userID).Scan(&pos).Error
	return &pos, err
}

func (r *convoyRepository) GetPositionsInRoom(ctx context.Context, roomID uuid.UUID) ([]domain.RiderPosition, error) {
	var positions []domain.RiderPosition
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ON (user_id) time as timestamp, room_id, user_id,
			ST_Y(location::geometry) as lat, ST_X(location::geometry) as lng,
			speed_kmh, heading
		FROM rider_locations
		WHERE room_id = ?
		ORDER BY user_id, time DESC`, roomID).Scan(&positions).Error
	return positions, err
}

func (r *convoyRepository) GetPositionsSince(ctx context.Context, roomID uuid.UUID, since time.Time) ([]domain.RiderPosition, error) {
	var positions []domain.RiderPosition
	err := r.db.WithContext(ctx).Raw(`
		SELECT time as timestamp, room_id, user_id,
			ST_Y(location::geometry) as lat, ST_X(location::geometry) as lng,
			speed_kmh, heading
		FROM rider_locations
		WHERE room_id = ? AND time > ?
		ORDER BY time DESC`, roomID, since).Scan(&positions).Error
	return positions, err
}
