package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/chat"
	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleListMessages(c fiber.Ctx) error {
	roomID := c.Params("id")
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChatService(postgres.NewChatRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(roomID)
	msgs, cursor, hasMore, err := svc.ListMessages(c.Context(), rid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: msgs, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleSendMessage(c fiber.Ctx) error {
	var req domain.SendMessageRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChatService(postgres.NewChatRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	msg, err := svc.SendMessage(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(msg)
}

func handleGetMessage(c fiber.Ctx) error {
	msgID := c.Params("messageId")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewChatService(postgres.NewChatRepository(c.Locals("db").(*postgres.Database)), logger)

	mid, _ := uuid.Parse(msgID)
	msg, err := svc.GetMessage(c.Context(), mid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
	}
	return c.JSON(msg)
}

func handleEditMessage(c fiber.Ctx) error {
	msgID := c.Params("messageId")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	var req struct{ Content string `json:"content"` }
	c.Bind().JSON(&req)

	svc := service.NewChatService(postgres.NewChatRepository(c.Locals("db").(*postgres.Database)), logger)
	mid, _ := uuid.Parse(msgID)
	uid, _ := uuid.Parse(userID)
	if err := svc.EditMessage(c.Context(), mid, req.Content, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "edited"})
}

func handleDeleteMessage(c fiber.Ctx) error {
	msgID := c.Params("messageId")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	svc := service.NewChatService(postgres.NewChatRepository(c.Locals("db").(*postgres.Database)), logger)
	mid, _ := uuid.Parse(msgID)
	uid, _ := uuid.Parse(userID)
	if err := svc.DeleteMessage(c.Context(), mid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

func handlePinMessage(c fiber.Ctx) error {
	msgID := c.Params("messageId")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	svc := service.NewChatService(postgres.NewChatRepository(c.Locals("db").(*postgres.Database)), logger)
	mid, _ := uuid.Parse(msgID)
	uid, _ := uuid.Parse(userID)
	if err := svc.PinMessage(c.Context(), mid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "pinned"})
}

func handleMarkRead(c fiber.Ctx) error {
	msgID := c.Params("messageId")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	svc := service.NewChatService(postgres.NewChatRepository(c.Locals("db").(*postgres.Database)), logger)
	mid, _ := uuid.Parse(msgID)
	uid, _ := uuid.Parse(userID)
	svc.MarkRead(c.Context(), mid, uid)
	return c.JSON(fiber.Map{"message": "read"})
}
