package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/auth"
	"go-road-backend/internal/service"
)

func handleRegister(c fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService, ok := c.Locals("auth_service").(domain.Service)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "service not available",
		})
	}

	resp, err := authService.Register(c.Context(), req)
	if err != nil {
		code := fiber.StatusInternalServerError
		if err == service.ErrEmailAlreadyExists || err == service.ErrUsernameAlreadyExists {
			code = fiber.StatusConflict
		}

		logger, _ := c.Locals("logger").(*zap.Logger)
		if logger != nil {
			logger.Warn("register failed", zap.Error(err))
		}

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func handleLogin(c fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService, ok := c.Locals("auth_service").(domain.Service)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "service not available",
		})
	}

	resp, err := authService.Login(c.Context(), req)
	if err != nil {
		code := fiber.StatusUnauthorized
		if err == service.ErrInvalidCredentials {
			code = fiber.StatusUnauthorized
		}

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

func handleRefreshToken(c fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService := c.Locals("auth_service").(domain.Service)
	resp, err := authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

func handleForgotPassword(c fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService := c.Locals("auth_service").(domain.Service)
	authService.ForgotPassword(c.Context(), req.Email)

	return c.JSON(fiber.Map{
		"message": "if the email exists, a reset link has been sent",
	})
}

func handleResetPassword(c fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService := c.Locals("auth_service").(domain.Service)
	if err := authService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "password has been reset"})
}

func handleGetProfile(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)

	authService := c.Locals("auth_service").(domain.Service)
	uid, _ := uuid.Parse(userID)
	resp, err := authService.GetProfile(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

func handleUpdateProfile(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)

	var updates map[string]interface{}
	if err := c.Bind().JSON(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService := c.Locals("auth_service").(domain.Service)
	uid, _ := uuid.Parse(userID)
	resp, err := authService.UpdateProfile(c.Context(), uid, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

func handleChangePassword(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService := c.Locals("auth_service").(domain.Service)
	uid, _ := uuid.Parse(userID)
	if err := authService.ChangePassword(c.Context(), uid, req.OldPassword, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "password changed successfully"})
}

func handleRegisterDevice(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)

	var req struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	authService := c.Locals("auth_service").(domain.Service)
	uid, _ := uuid.Parse(userID)
	if err := authService.RegisterDevice(c.Context(), uid, req.Token, req.Platform); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "device registered"})
}

func handleUnregisterDevice(c fiber.Ctx) error {
	deviceID := c.Params("id")

	authService := c.Locals("auth_service").(domain.Service)
	did, _ := uuid.Parse(deviceID)
	if err := authService.UnregisterDevice(c.Context(), did); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "device unregistered"})
}
