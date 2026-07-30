package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/emergency"
)

type emergencyRepository struct {
	db *Database
}

func NewEmergencyRepository(db *Database) domain.Repository {
	return &emergencyRepository{db: db}
}

func (r *emergencyRepository) CreateEvent(ctx context.Context, event *domain.EmergencyEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *emergencyRepository) FindEventByID(ctx context.Context, id uuid.UUID) (*domain.EmergencyEvent, error) {
	var e domain.EmergencyEvent
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&e).Error
	return &e, err
}

func (r *emergencyRepository) ListEvents(ctx context.Context, cursor string, limit int, status string) ([]domain.EmergencyEvent, string, bool, error) {
	var events []domain.EmergencyEvent
	db := r.db.WithContext(ctx)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&events).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(events) > limit
	if hasMore {
		events = events[:limit]
	}
	var nextCursor string
	if len(events) > 0 {
		last := events[len(events)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return events, nextCursor, hasMore, nil
}

func (r *emergencyRepository) UpdateEvent(ctx context.Context, event *domain.EmergencyEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *emergencyRepository) CreateSOS(ctx context.Context, sos *domain.SOSEvent) error {
	return r.db.WithContext(ctx).Create(sos).Error
}

func (r *emergencyRepository) FindSOSByID(ctx context.Context, id uuid.UUID) (*domain.SOSEvent, error) {
	var s domain.SOSEvent
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	return &s, err
}

func (r *emergencyRepository) ListSOS(ctx context.Context, cursor string, limit int) ([]domain.SOSEvent, string, bool, error) {
	var sos []domain.SOSEvent
	db := r.db.WithContext(ctx).Order("triggered_at DESC, id DESC").Limit(limit + 1)
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			TriggeredAt time.Time `json:"triggered_at"`
			ID          string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(triggered_at, id) < (?, ?)", c.TriggeredAt, c.ID)
	}
	if err := db.Find(&sos).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(sos) > limit
	if hasMore {
		sos = sos[:limit]
	}
	var nextCursor string
	if len(sos) > 0 {
		last := sos[len(sos)-1]
		cData, _ := json.Marshal(struct {
			TriggeredAt time.Time `json:"triggered_at"`
			ID          string    `json:"id"`
		}{TriggeredAt: last.TriggeredAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return sos, nextCursor, hasMore, nil
}

func (r *emergencyRepository) UpdateSOS(ctx context.Context, sos *domain.SOSEvent) error {
	return r.db.WithContext(ctx).Save(sos).Error
}
