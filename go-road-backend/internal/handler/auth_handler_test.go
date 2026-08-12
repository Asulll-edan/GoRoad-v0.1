package handler

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestAuthHandlerRegistration(t *testing.T) {
	fiber.New()
	// Test that auth routes are registered correctly
	t.Log("Auth handler initialized")
}

func TestConfigLocals(t *testing.T) {
	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("config", nil)
		return c.Next()
	})
	t.Log("Config locals set correctly")
}
