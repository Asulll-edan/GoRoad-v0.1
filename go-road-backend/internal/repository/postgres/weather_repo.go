package postgres

import (
	"context"

	domain "go-road-backend/internal/domain/weather"
	"go-road-backend/internal/repository/redis"
)

type weatherRepository struct {
	cache redis.CacheRepository
}

func NewWeatherRepository(cache redis.CacheRepository) domain.Repository {
	return &weatherRepository{cache: cache}
}

func (r *weatherRepository) GetCached(ctx context.Context, key string) (string, error) {
	return r.cache.Get(ctx, key)
}

func (r *weatherRepository) SetCached(ctx context.Context, key string, value string, ttl int) error {
	return r.cache.Set(ctx, key, value, 0)
}
