package postgres

import (
	"context"

	domain "go-road-backend/internal/domain/ai"
	"go-road-backend/internal/repository/redis"
)

type aiRepository struct {
	cache redis.CacheRepository
}

func NewAIRepository(cache redis.CacheRepository) domain.Repository {
	return &aiRepository{cache: cache}
}

func (r *aiRepository) GetCached(ctx context.Context, key string) (string, error) {
	return r.cache.Get(ctx, key)
}

func (r *aiRepository) SetCached(ctx context.Context, key string, value string, ttl int) error {
	return r.cache.Set(ctx, key, value, 0)
}
