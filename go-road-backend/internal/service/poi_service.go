package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/poi"
	"go-road-backend/internal/repository/redis"
)

type poiService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewPOIService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &poiService{repo: repo, cache: cache, logger: logger}
}

func (s *poiService) GetNearbyPOI(ctx context.Context, lat, lng, radiusKm float64, types []string) ([]domain.POI, error) {
	cacheKey := fmt.Sprintf("cache:poi:nearby:%.4f:%.4f:%v", lat, lng, types)
	var cached []domain.POI
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && len(cached) > 0 {
		return cached, nil
	}
	pois, err := s.repo.FindNearby(ctx, lat, lng, radiusKm, types, 50)
	if err != nil {
		return nil, err
	}
	s.cache.SetJSON(ctx, cacheKey, pois, 1*time.Hour)
	return pois, nil
}

func (s *poiService) GetPOIDetail(ctx context.Context, id uuid.UUID) (*domain.POI, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *poiService) ReportPOI(ctx context.Context, userID uuid.UUID, poiID uuid.UUID, reportType, description string) error {
	report := &domain.POIReport{
		UserID:      userID,
		POIID:       poiID,
		ReportType:  reportType,
		Description: description,
	}
	return s.repo.CreateReport(ctx, report)
}

func (s *poiService) GetCategories(ctx context.Context) ([]string, error) {
	return []string{
		"fuel", "food", "accommodation", "rest_area", "workshop",
		"hospital", "police", "atm", "mosque", "tourist_attraction",
		"scenic_view", "camping", "parking", "ferry",
	}, nil
}
