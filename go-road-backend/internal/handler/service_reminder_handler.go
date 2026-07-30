package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/service_reminder"
	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleCreateReminder(c fiber.Ctx) error {
	var req domain.CreateReminderRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewServiceReminderService(postgres.NewServiceReminderRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	reminder, err := svc.Create(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(reminder)
}

func handleListReminders(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewServiceReminderService(postgres.NewServiceReminderRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	reminders, cursor, hasMore, err := svc.List(c.Context(), uid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: reminders, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleUpdateReminder(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	var req map[string]interface{}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewServiceReminderService(postgres.NewServiceReminderRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.Update(c.Context(), rid, req, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

func handleDeleteReminder(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewServiceReminderService(postgres.NewServiceReminderRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.Delete(c.Context(), rid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

func handleCompleteReminder(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewServiceReminderService(postgres.NewServiceReminderRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.Complete(c.Context(), rid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "completed"})
}
