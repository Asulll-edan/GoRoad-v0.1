package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/expense"
	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleCreateExpense(c fiber.Ctx) error {
	var req domain.CreateExpenseRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewExpenseService(postgres.NewExpenseRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	expense, err := svc.Create(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(expense)
}

func handleListExpenses(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewExpenseService(postgres.NewExpenseRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	expenses, cursor, hasMore, err := svc.List(c.Context(), uid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pagination.PaginatedResponse{
		Data: expenses,
		Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore},
	})
}

func handleGetExpense(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewExpenseService(postgres.NewExpenseRepository(c.Locals("db").(*postgres.Database)), logger)

	eid, _ := uuid.Parse(id)
	expense, err := svc.Get(c.Context(), eid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(expense)
}

func handleUpdateExpense(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	var req map[string]interface{}
	c.Bind().JSON(&req)

	svc := service.NewExpenseService(postgres.NewExpenseRepository(c.Locals("db").(*postgres.Database)), logger)
	eid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	expense, err := svc.Update(c.Context(), eid, req, uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(expense)
}

func handleDeleteExpense(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	svc := service.NewExpenseService(postgres.NewExpenseRepository(c.Locals("db").(*postgres.Database)), logger)
	eid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := svc.Delete(c.Context(), eid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "expense deleted"})
}
