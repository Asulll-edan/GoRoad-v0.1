package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleListBadges(c fiber.Ctx) error {
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewBadgeService(postgres.NewBadgeRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	badges, err := svc.ListBadges(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(badges)
}

func handleMyBadges(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewBadgeService(postgres.NewBadgeRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	badges, err := svc.GetMyBadges(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(badges)
}

func handleBadgeProgress(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewBadgeService(postgres.NewBadgeRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	progress, err := svc.GetBadgeProgress(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(progress)
}
