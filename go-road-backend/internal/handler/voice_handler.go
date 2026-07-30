package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"go-road-backend/internal/config"
)

func handleVoiceToken(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roomID := c.Query("room_id")

	var req struct {
		RoomID string `json:"room_id"`
	}
	if err := c.Bind().JSON(&req); err == nil && req.RoomID != "" {
		roomID = req.RoomID
	}

	if roomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "room_id is required",
		})
	}

	cfg := c.Locals("config").(*config.Config)
	now := time.Now()

	claims := jwt.MapClaims{
		"iss": cfg.LiveKitAPIKey,
		"sub": userID,
		"exp": now.Add(6 * time.Hour).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"video": map[string]interface{}{
			"room":         roomID,
			"roomJoin":     true,
			"canPublish":   true,
			"canSubscribe": true,
		},
		"metadata": map[string]string{
			"user_id": userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfg.LiveKitAPISecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token":    tokenStr,
		"host":     cfg.LiveKitHost,
		"room_id":  roomID,
		"identity": userID,
	})
}
