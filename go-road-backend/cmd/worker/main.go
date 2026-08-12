package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go-road-backend/internal/config"
	"go-road-backend/internal/repository/redis"
	"go-road-backend/internal/worker"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			func(cfg *config.Config) (*nats.Conn, error) {
				return nats.Connect(cfg.NATSUrl())
			},
			func(cfg *config.Config) (*pgxpool.Pool, error) {
				return pgxpool.New(context.Background(), cfg.DBDSN())
			},
			func(cfg *config.Config) (redis.CacheRepository, error) {
				return redis.NewCacheRepository("redis://" + cfg.RedisURL + "/1")
			},
			func(pool *pgxpool.Pool, nc *nats.Conn) *worker.LocationAggregator {
				return worker.NewLocationAggregator(pool, nc)
			},
			func(nc *nats.Conn, cache redis.CacheRepository) *worker.SmartDetection {
				return worker.NewSmartDetection(nc, cache)
			},
			func(nc *nats.Conn) *worker.NotificationSender {
				return worker.NewNotificationSender(nc)
			},
			zap.NewProduction,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, aggregator *worker.LocationAggregator, detector *worker.SmartDetection, sender *worker.NotificationSender, logger *zap.Logger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go aggregator.Start(ctx)
						go detector.Start(ctx)
						go sender.Start(ctx)
						logger.Info("Worker started")
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Worker stopped")
						return nil
					},
				})
			},
		),
	).Run()
}
