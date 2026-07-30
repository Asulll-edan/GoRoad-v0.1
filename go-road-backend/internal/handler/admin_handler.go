package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleAdminDashboard(c fiber.Ctx) error {
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	stats, err := svc.GetDashboard(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stats)
}

func handleAdminListUsers(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	search := c.Query("search")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	users, cursor, hasMore, err := svc.ListUsers(c.Context(), params.Cursor, params.Limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: users, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleAdminGetUser(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(id)
	user, err := svc.GetUser(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

func handleAdminBanUser(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(id)
	if err := svc.BanUser(c.Context(), uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "user banned"})
}

func handleAdminUnbanUser(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(id)
	if err := svc.UnbanUser(c.Context(), uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "user unbanned"})
}

func handleAdminListRooms(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	status := c.Query("status")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rooms, cursor, hasMore, err := svc.ListRooms(c.Context(), params.Cursor, params.Limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: rooms, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleAdminGetRoom(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rid, _ := uuid.Parse(id)
	room, err := svc.GetRoom(c.Context(), rid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "room not found"})
	}
	return c.JSON(room)
}

func handleAdminListReports(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	status := c.Query("status")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	reports, cursor, hasMore, err := svc.ListReports(c.Context(), params.Cursor, params.Limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: reports, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleAdminReviewReport(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct{ Status string `json:"status"` }
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.ReviewReport(c.Context(), rid, req.Status, uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "report reviewed"})
}

func handleAdminListEmergency(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	events, cursor, hasMore, err := svc.ListEmergency(c.Context(), params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: events, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleAdminAnalytics(c fiber.Ctx) error {
	period := c.Query("period", "today")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	analytics, err := svc.GetAnalytics(c.Context(), period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(analytics)
}

func handleAdminLogs(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	level := c.Query("level")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAdminService(postgres.NewAdminRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	logs, cursor, hasMore, err := svc.GetLogs(c.Context(), params.Cursor, params.Limit, level)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: logs, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}
