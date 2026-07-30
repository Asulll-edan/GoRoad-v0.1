package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleListTemplates(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChecklistService(postgres.NewChecklistRepository(c.Locals("db").(*postgres.Database)), logger)

	templates, cursor, hasMore, err := svc.ListTemplates(c.Context(), params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: templates, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleCreateTemplate(c fiber.Ctx) error {
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Items       []string `json:"items"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChecklistService(postgres.NewChecklistRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	template, err := svc.CreateTemplate(c.Context(), req.Name, req.Description, req.Category, req.Items, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(template)
}

func handleGetTemplate(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChecklistService(postgres.NewChecklistRepository(c.Locals("db").(*postgres.Database)), logger)

	tid, _ := uuid.Parse(id)
	template, err := svc.GetTemplate(c.Context(), tid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "template not found"})
	}
	return c.JSON(template)
}

func handleCreateTouringChecklist(c fiber.Ctx) error {
	roomID := c.Params("roomId")
	var req struct{ TemplateID string `json:"template_id"` }
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChecklistService(postgres.NewChecklistRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(roomID)
	tid, _ := uuid.Parse(req.TemplateID)
	uid, _ := uuid.Parse(userID)
	if err := svc.CreateTouringChecklist(c.Context(), rid, tid, uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "checklist created"})
}

func handleGetTouringChecklist(c fiber.Ctx) error {
	roomID := c.Params("roomId")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChecklistService(postgres.NewChecklistRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(roomID)
	uid, _ := uuid.Parse(userID)
	items, err := svc.GetTouringChecklist(c.Context(), rid, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func handleToggleChecklistItem(c fiber.Ctx) error {
	itemID := c.Params("itemId")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChecklistService(postgres.NewChecklistRepository(c.Locals("db").(*postgres.Database)), logger)

	iid, _ := uuid.Parse(itemID)
	uid, _ := uuid.Parse(userID)
	if err := svc.ToggleItem(c.Context(), iid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "toggled"})
}
