package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/repository/redis"
	"go-road-backend/internal/service"
)

func handleListNotifications(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewNotificationService(postgres.NewNotificationRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	notifs, cursor, hasMore, err := svc.ListNotifications(c.Context(), uid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: notifs, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleMarkNotificationRead(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewNotificationService(postgres.NewNotificationRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	nid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.MarkRead(c.Context(), nid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "marked read"})
}

func handleMarkAllRead(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewNotificationService(postgres.NewNotificationRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	if err := svc.MarkAllRead(c.Context(), uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "all marked read"})
}

func handleUnreadCount(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewNotificationService(postgres.NewNotificationRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	count, err := svc.GetUnreadCount(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"unread_count": count})
}

func handleUpdateNotificationPreferences(c fiber.Ctx) error {
	var req map[string]interface{}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewNotificationService(postgres.NewNotificationRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	if err := svc.UpdatePreferences(c.Context(), uid, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "preferences updated"})
}
