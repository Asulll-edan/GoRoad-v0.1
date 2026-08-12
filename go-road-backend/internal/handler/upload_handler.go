package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleUploadFile(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	header, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file not provided"})
	}
	category := c.FormValue("category", "general")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewUploadService(postgres.NewUploadRepository(c.Locals("db").(*postgres.Database)), logger)

	f, _ := header.Open()
	defer f.Close()
	uid, _ := uuid.Parse(userID)
	upload, err := svc.UploadFile(c.Context(), uid, f, header, category)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(upload)
}

func handleUploadAvatar(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	header, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "avatar not provided"})
	}
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewUploadService(postgres.NewUploadRepository(c.Locals("db").(*postgres.Database)), logger)

	f, _ := header.Open()
	defer f.Close()
	uid, _ := uuid.Parse(userID)
	upload, err := svc.UploadAvatar(c.Context(), uid, f, header)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(upload)
}

func handleUploadPhoto(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	header, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "photo not provided"})
	}
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewUploadService(postgres.NewUploadRepository(c.Locals("db").(*postgres.Database)), logger)

	f, _ := header.Open()
	defer f.Close()
	uid, _ := uuid.Parse(userID)
	upload, err := svc.UploadPhoto(c.Context(), uid, f, header)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(upload)
}

func handleDeleteFile(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewUploadService(postgres.NewUploadRepository(c.Locals("db").(*postgres.Database)), logger)

	fid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.DeleteFile(c.Context(), fid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}
