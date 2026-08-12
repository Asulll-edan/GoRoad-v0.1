package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/weather"
	"go-road-backend/internal/repository/redis"
)

type weatherService struct {
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewWeatherService(cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &weatherService{cache: cache, logger: logger}
}

func (s *weatherService) GetCurrentWeather(ctx context.Context, lat, lng float64) (*domain.WeatherData, error) {
	cacheKey := fmt.Sprintf("cache:weather:current:%.4f:%.4f", lat, lng)
	var data domain.WeatherData
	if err := s.cache.GetJSON(ctx, cacheKey, &data); err == nil {
		return &data, nil
	}
	s.logger.Warn("weather data not in cache, delegating to Python", zap.Float64("lat", lat), zap.Float64("lng", lng))
	return nil, fmt.Errorf("weather data not available")
}

func (s *weatherService) GetForecast(ctx context.Context, lat, lng float64) (*domain.WeatherForecast, error) {
	cacheKey := fmt.Sprintf("cache:weather:forecast:%.4f:%.4f", lat, lng)
	var data domain.WeatherForecast
	if err := s.cache.GetJSON(ctx, cacheKey, &data); err == nil {
		return &data, nil
	}
	return nil, fmt.Errorf("forecast data not available")
}

func (s *weatherService) GetAlerts(ctx context.Context, lat, lng float64) ([]domain.WeatherAlert, error) {
	cacheKey := fmt.Sprintf("cache:weather:alerts:%.4f:%.4f", lat, lng)
	var alerts []domain.WeatherAlert
	if err := s.cache.GetJSON(ctx, cacheKey, &alerts); err == nil {
		return alerts, nil
	}
	return nil, nil
}

func (s *weatherService) GetRouteWeather(ctx context.Context, routeID string) ([]domain.WeatherData, error) {
	cacheKey := fmt.Sprintf("cache:weather:route:%s", routeID)
	var data []domain.WeatherData
	if err := s.cache.GetJSON(ctx, cacheKey, &data); err == nil {
		return data, nil
	}
	return nil, fmt.Errorf("route weather not available")
}

func UpdateWeatherCache(cache redis.CacheRepository, key string, data interface{}, ttl time.Duration) error {
	return cache.SetJSON(context.Background(), key, data, ttl)
}
