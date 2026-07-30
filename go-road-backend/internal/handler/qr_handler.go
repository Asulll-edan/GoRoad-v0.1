package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleGetMyQRCard(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewQRService(postgres.NewQRRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	card, err := svc.GetMyQRCard(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(card)
}

func handleRegenerateQR(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewQRService(postgres.NewQRRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	card, err := svc.RegenerateQR(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(card)
}

func handleScanQR(c fiber.Ctx) error {
	code := c.Params("code")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewQRService(postgres.NewQRRepository(c.Locals("db").(*postgres.Database)), logger)

	card, err := svc.ScanQR(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid QR code"})
	}
	return c.JSON(card)
}
