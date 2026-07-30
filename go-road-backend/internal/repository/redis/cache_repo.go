package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	GetJSON(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Hash operations
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// Set operations
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)

	// Sorted Set operations
	ZAdd(ctx context.Context, key string, members ...redis.Z) error
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZRevRank(ctx context.Context, key string, member string) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)

	// Rate limiting
	IncrWithExpiry(ctx context.Context, key string, ttl time.Duration) (int64, error)

	// Distributed lock
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error

	// Pub/Sub
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) *redis.PubSub

	// Close
	Close() error
}

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(redisURL string) (CacheRepository, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &cacheRepository{client: client}, nil
}

func (r *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *cacheRepository) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *cacheRepository) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *cacheRepository) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *cacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *cacheRepository) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *cacheRepository) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

func (r *cacheRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

func (r *cacheRepository) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

func (r *cacheRepository) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

func (r *cacheRepository) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

func (r *cacheRepository) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *cacheRepository) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.client.ZAdd(ctx, key, members...).Err()
}

func (r *cacheRepository) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func (r *cacheRepository) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return r.client.ZRevRank(ctx, key, member).Result()
}

func (r *cacheRepository) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}

func (r *cacheRepository) IncrWithExpiry(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		r.client.Expire(ctx, key, ttl)
	}
	return count, nil
}

func (r *cacheRepository) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, "lock:"+key, "1", ttl).Result()
}

func (r *cacheRepository) ReleaseLock(ctx context.Context, key string) error {
	return r.client.Del(ctx, "lock:"+key).Err()
}

func (r *cacheRepository) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.client.Publish(ctx, channel, message).Err()
}

func (r *cacheRepository) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}

func (r *cacheRepository) Close() error {
	return r.client.Close()
}
