package poi

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	GetNearbyPOI(ctx context.Context, lat, lng, radiusKm float64, types []string) ([]POI, error)
	GetPOIDetail(ctx context.Context, id uuid.UUID) (*POI, error)
	ReportPOI(ctx context.Context, userID uuid.UUID, poiID uuid.UUID, reportType, description string) error
	GetCategories(ctx context.Context) ([]string, error)
}
