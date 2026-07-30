package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		correlationID := uuid.New().String()

		c.Set("X-Correlation-ID", correlationID)
		c.Locals("correlation_id", correlationID)

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()
		path := c.Path()
		method := c.Method()

		fields := []zap.Field{
			zap.String("correlation_id", correlationID),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		}

		if userID, ok := c.Locals("user_id").(string); ok {
			fields = append(fields, zap.String("user_id", userID))
		}

		if status >= 500 {
			logger.Error("server error", fields...)
		} else if status >= 400 {
			logger.Warn("client error", fields...)
		} else {
			logger.Info("request", fields...)
		}

		return err
	}
}

func SetupMiddleware(app *fiber.App, isProduction bool, logger *zap.Logger) {
	app.Use(LoggerMiddleware(logger))
	app.Use(cors.New(NewCORSConfig(isProduction)))
}
