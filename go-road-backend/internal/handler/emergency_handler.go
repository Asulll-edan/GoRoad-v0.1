package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/emergency"
	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleReportEmergency(c fiber.Ctx) error {
	var req domain.CreateEmergencyRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewEmergencyService(postgres.NewEmergencyRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	event, err := svc.ReportEmergency(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(event)
}

func handleListEmergencies(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	status := c.Query("status")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewEmergencyService(postgres.NewEmergencyRepository(c.Locals("db").(*postgres.Database)), logger)

	events, cursor, hasMore, err := svc.ListEmergencies(c.Context(), params.Cursor, params.Limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: events, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleGetEmergency(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewEmergencyService(postgres.NewEmergencyRepository(c.Locals("db").(*postgres.Database)), logger)

	eid, _ := uuid.Parse(id)
	event, err := svc.GetEmergency(c.Context(), eid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "event not found"})
	}
	return c.JSON(event)
}

func handleAcknowledgeEmergency(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewEmergencyService(postgres.NewEmergencyRepository(c.Locals("db").(*postgres.Database)), logger)

	eid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.AcknowledgeEmergency(c.Context(), eid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "acknowledged"})
}

func handleResolveEmergency(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewEmergencyService(postgres.NewEmergencyRepository(c.Locals("db").(*postgres.Database)), logger)

	eid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.ResolveEmergency(c.Context(), eid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "resolved"})
}

func handleTriggerSOS(c fiber.Ctx) error {
	var req struct {
		Lat    float64 `json:"lat"`
		Lng    float64 `json:"lng"`
		RoomID string  `json:"room_id,omitempty"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewEmergencyService(postgres.NewEmergencyRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	var rid *uuid.UUID
	if req.RoomID != "" {
		if parsed, err := uuid.Parse(req.RoomID); err == nil {
			rid = &parsed
		}
	}
	sos, err := svc.TriggerSOS(c.Context(), uid, req.Lat, req.Lng, rid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(sos)
}

func handleDismissSOS(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewEmergencyService(postgres.NewEmergencyRepository(c.Locals("db").(*postgres.Database)), logger)

	sid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.DismissSOS(c.Context(), sid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "dismissed"})
}
