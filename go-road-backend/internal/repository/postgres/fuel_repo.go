package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/fuel"
)

type fuelRepository struct {
	db *Database
}

func NewFuelRepository(db *Database) domain.Repository {
	return &fuelRepository{db: db}
}

func (r *fuelRepository) Create(ctx context.Context, log *domain.FuelLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *fuelRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.FuelLog, error) {
	var log domain.FuelLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *fuelRepository) FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.FuelLog, string, bool, error) {
	var logs []domain.FuelLog
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			LoggedAt time.Time `json:"logged_at"`
			ID       string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(logged_at, id) < (?, ?)", c.LoggedAt, c.ID)
	}

	db = db.Order("logged_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&logs).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit]
	}

	var nextCursor string
	if len(logs) > 0 {
		last := logs[len(logs)-1]
		cursorData, _ := json.Marshal(struct {
			LoggedAt time.Time `json:"logged_at"`
			ID       string    `json:"id"`
		}{LoggedAt: last.LoggedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return logs, nextCursor, hasMore, nil
}

func (r *fuelRepository) Update(ctx context.Context, log *domain.FuelLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *fuelRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.FuelLog{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
