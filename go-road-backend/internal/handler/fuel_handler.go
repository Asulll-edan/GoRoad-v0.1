package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/fuel"
	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleCreateFuelLog(c fiber.Ctx) error {
	var req domain.CreateFuelLogRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewFuelService(postgres.NewFuelRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	log, err := svc.Create(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(log)
}

func handleListFuelLogs(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewFuelService(postgres.NewFuelRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	logs, cursor, hasMore, err := svc.List(c.Context(), uid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pagination.PaginatedResponse{
		Data: logs,
		Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore},
	})
}

func handleGetFuelLog(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewFuelService(postgres.NewFuelRepository(c.Locals("db").(*postgres.Database)), logger)

	lid, _ := uuid.Parse(id)
	log, err := svc.Get(c.Context(), lid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(log)
}

func handleUpdateFuelLog(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	var req map[string]interface{}
	c.Bind().JSON(&req)

	svc := service.NewFuelService(postgres.NewFuelRepository(c.Locals("db").(*postgres.Database)), logger)
	lid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	log, err := svc.Update(c.Context(), lid, req, uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(log)
}

func handleDeleteFuelLog(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	svc := service.NewFuelService(postgres.NewFuelRepository(c.Locals("db").(*postgres.Database)), logger)
	lid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := svc.Delete(c.Context(), lid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "fuel log deleted"})
}

func handleFuelAnalytics(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "delegated to Python analytics service"})
}
