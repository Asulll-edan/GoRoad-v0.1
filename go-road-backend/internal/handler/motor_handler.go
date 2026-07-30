package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/motor"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleCreateMotor(c fiber.Ctx) error {
	var req domain.CreateMotorRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewMotorService(postgres.NewMotorRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	motor, err := svc.Create(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(motor)
}

func handleListMotors(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewMotorService(postgres.NewMotorRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	motors, err := svc.List(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(motors)
}

func handleGetMotor(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewMotorService(postgres.NewMotorRepository(c.Locals("db").(*postgres.Database)), logger)

	mid, _ := uuid.Parse(id)
	motor, err := svc.Get(c.Context(), mid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(motor)
}

func handleUpdateMotor(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	var req map[string]interface{}
	c.Bind().JSON(&req)

	svc := service.NewMotorService(postgres.NewMotorRepository(c.Locals("db").(*postgres.Database)), logger)
	mid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	motor, err := svc.Update(c.Context(), mid, req, uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(motor)
}

func handleDeleteMotor(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	svc := service.NewMotorService(postgres.NewMotorRepository(c.Locals("db").(*postgres.Database)), logger)
	mid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := svc.Delete(c.Context(), mid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "motor deleted"})
}

func handleSetPrimaryMotor(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)

	svc := service.NewMotorService(postgres.NewMotorRepository(c.Locals("db").(*postgres.Database)), logger)
	mid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := svc.SetPrimary(c.Context(), mid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "primary motor set"})
}
