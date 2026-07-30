package main

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go-road-backend/internal/config"
	"go-road-backend/internal/middleware"
	"go-road-backend/internal/handler"
	"go-road-backend/internal/repository/redis"
	"go-road-backend/internal/repository/postgres"
	authRepo "go-road-backend/internal/domain/auth"
	authSvc "go-road-backend/internal/service"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			func(cfg *config.Config) (redis.CacheRepository, error) {
				return redis.NewCacheRepository("redis://" + cfg.RedisURL + "/0")
			},
			func(cfg *config.Config) (*postgres.Database, error) {
				return postgres.NewDatabase(cfg.DBDSN())
			},
			func(db *postgres.Database) authRepo.Repository {
				return postgres.NewAuthRepository(db)
			},
			func(repo authRepo.Repository, cache redis.CacheRepository, cfg *config.Config, logger *zap.Logger) authRepo.Service {
				return authSvc.NewAuthService(repo, cache, cfg, logger)
			},
			zap.NewProduction,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, cfg *config.Config, cache redis.CacheRepository, db *postgres.Database, logger *zap.Logger, authService authRepo.Service) {
				app := fiber.New(fiber.Config{
					AppName:       "Go Road API v3",
					CaseSensitive: true,
				})

				middleware.SetupMiddleware(app, cfg.IsProduction(), logger)
				handler.SetupRoutes(app, cfg, cache, db, logger, authService)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Info("Starting Go Road API server", zap.String("port", cfg.APIPort))
						go app.Listen(":" + cfg.APIPort)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Shutting down Go Road API server")
						return app.Shutdown()
					},
				})
			},
		),
	).Run()
}
