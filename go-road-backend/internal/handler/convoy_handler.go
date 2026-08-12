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

func handleCreateFormation(c fiber.Ctx) error {
	roomID := c.Params("id")
	var req struct {
		Name          string `json:"name"`
		FormationType string `json:"formation_type"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewConvoyService(postgres.NewConvoyRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rid, _ := uuid.Parse(roomID)
	uid, _ := uuid.Parse(userID)
	f, err := svc.CreateFormation(c.Context(), rid, req.Name, req.FormationType, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(f)
}

func handleUpdateFormation(c fiber.Ctx) error {
	formationID := c.Params("formationId")
	var req map[string]interface{}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewConvoyService(postgres.NewConvoyRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	fid, _ := uuid.Parse(formationID)
	uid, _ := uuid.Parse(c.Locals("user_id").(string))
	if err := svc.UpdateFormation(c.Context(), fid, req, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

func handleGetActiveFormation(c fiber.Ctx) error {
	roomID := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewConvoyService(postgres.NewConvoyRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rid, _ := uuid.Parse(roomID)
	f, err := svc.GetActiveFormation(c.Context(), rid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no active formation"})
	}
	return c.JSON(f)
}

func handleUpdateLocation(c fiber.Ctx) error {
	roomID := c.Params("id")
	var req struct {
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
		Speed   float64 `json:"speed_kmh"`
		Heading float64 `json:"heading"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewConvoyService(postgres.NewConvoyRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rid, _ := uuid.Parse(roomID)
	uid, _ := uuid.Parse(userID)
	if err := svc.UpdateLocation(c.Context(), rid, uid, req.Lat, req.Lng, req.Speed, req.Heading); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "location updated"})
}

func handleGetLocations(c fiber.Ctx) error {
	roomID := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewConvoyService(postgres.NewConvoyRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rid, _ := uuid.Parse(roomID)
	positions, err := svc.GetLocations(c.Context(), rid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(positions)
}

func handleGetTracking(c fiber.Ctx) error {
	roomID := c.Params("id")
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewConvoyService(postgres.NewConvoyRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	rid, _ := uuid.Parse(roomID)
	positions, cursor, hasMore, err := svc.GetTracking(c.Context(), rid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: positions, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}
