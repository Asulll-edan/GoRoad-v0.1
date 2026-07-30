package poi

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, poi *POI) error
	FindByID(ctx context.Context, id uuid.UUID) (*POI, error)
	FindNearby(ctx context.Context, lat, lng, radiusKm float64, types []string, limit int) ([]POI, error)
	FindByCategory(ctx context.Context, category string, cursor string, limit int) ([]POI, string, bool, error)
	CreateReport(ctx context.Context, report *POIReport) error
	ListReports(ctx context.Context, cursor string, limit int) ([]POIReport, string, bool, error)
}
