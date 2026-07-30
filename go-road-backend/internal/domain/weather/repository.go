package weather

import "context"

type Repository interface {
	// Weather is fetched from external API via Python service
	// Repository is a pass-through to cache
	GetCached(ctx context.Context, key string) (string, error)
	SetCached(ctx context.Context, key string, value string, ttl int) error
}
