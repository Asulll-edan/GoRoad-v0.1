package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID string   `json:"user_id"`
	Role   string   `json:"role"`
	Skills []string `json:"skills"`
}

func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("user_role", claims["role"])

		return c.Next()
	}
}

func AdminMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals("user_role").(string)
		if !ok || role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "admin access required",
			})
		}
		return c.Next()
	}
}

func RoomRoleMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roomRole := c.Locals("room_role")
		if roomRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "room role not found",
			})
		}

		roleStr, ok := roomRole.(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "invalid room role",
			})
		}

		for _, allowed := range requiredRoles {
			if roleStr == allowed {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "insufficient room role",
		})
	}
}
