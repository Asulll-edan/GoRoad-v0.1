package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/poi"
)

type poiRepository struct {
	db *Database
}

func NewPOIRepository(db *Database) domain.Repository {
	return &poiRepository{db: db}
}

func (r *poiRepository) Create(ctx context.Context, poi *domain.POI) error {
	return r.db.WithContext(ctx).Create(poi).Error
}

func (r *poiRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.POI, error) {
	var p domain.POI
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	return &p, err
}

func (r *poiRepository) FindNearby(ctx context.Context, lat, lng, radiusKm float64, types []string, limit int) ([]domain.POI, error) {
	var pois []domain.POI
	query := `SELECT * FROM pois 
		WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
		AND deleted_at IS NULL`
	args := []interface{}{lng, lat, radiusKm * 1000}

	if len(types) > 0 {
		query += ` AND category = ANY(?)`
		args = append(args, types)
	}
	query += ` ORDER BY ST_Distance(location, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) ASC LIMIT ?`
	args = append(args, lng, lat, limit)

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&pois).Error
	return pois, err
}

func (r *poiRepository) FindByCategory(ctx context.Context, category string, cursor string, limit int) ([]domain.POI, string, bool, error) {
	var pois []domain.POI
	db := r.db.WithContext(ctx).Where("category = ? AND deleted_at IS NULL", category)
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
	if err := db.Find(&pois).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(pois) > limit
	if hasMore {
		pois = pois[:limit]
	}
	var nextCursor string
	if len(pois) > 0 {
		last := pois[len(pois)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return pois, nextCursor, hasMore, nil
}

func (r *poiRepository) CreateReport(ctx context.Context, report *domain.POIReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *poiRepository) ListReports(ctx context.Context, cursor string, limit int) ([]domain.POIReport, string, bool, error) {
	var reports []domain.POIReport
	db := r.db.WithContext(ctx).Order("created_at DESC, id DESC").Limit(limit + 1)
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
	}
	if err := db.Find(&reports).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(reports) > limit
	if hasMore {
		reports = reports[:limit]
	}
	var nextCursor string
	if len(reports) > 0 {
		last := reports[len(reports)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return reports, nextCursor, hasMore, nil
}
