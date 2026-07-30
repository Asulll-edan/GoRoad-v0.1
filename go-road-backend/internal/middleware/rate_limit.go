package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"go-road-backend/internal/repository/redis"
)

type RateLimitConfig struct {
	MaxRequests int
	WindowTTL   time.Duration
	KeyPrefix   string
}

var defaultRateLimits = map[string]RateLimitConfig{
	"/api/v1/auth/":    {MaxRequests: 10, WindowTTL: time.Minute, KeyPrefix: "rate:ip"},
	"/api/v1/ai/":      {MaxRequests: 20, WindowTTL: time.Hour, KeyPrefix: "rate:user"},
	"/api/v1/rooms/":   {MaxRequests: 60, WindowTTL: time.Minute, KeyPrefix: "rate:user"},
	"/api/v1/chat/":    {MaxRequests: 120, WindowTTL: time.Minute, KeyPrefix: "rate:user"},
	"/api/v1/location": {MaxRequests: 300, WindowTTL: time.Minute, KeyPrefix: "rate:user"},
	"default":          {MaxRequests: 100, WindowTTL: time.Minute, KeyPrefix: "rate:ip"},
}

func RateLimitMiddleware(cache redis.CacheRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()
		ip := c.IP()

		userID, _ := c.Locals("user_id").(string)

		config := getRateLimitConfig(path)

		var key string
		if config.KeyPrefix == "rate:user" && userID != "" {
			key = "rate:user:" + userID + ":" + path
		} else {
			key = "rate:ip:" + ip + ":" + path
		}

		count, err := cache.IncrWithExpiry(c.Context(), key, config.WindowTTL)
		if err != nil {
			return c.Next()
		}

		c.Set("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(max(0, config.MaxRequests-int(count))))
		c.Set("X-RateLimit-Reset", strconv.Itoa(int(config.WindowTTL.Seconds())))

		if int(count) > config.MaxRequests {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		}

		return c.Next()
	}
}

func getRateLimitConfig(path string) RateLimitConfig {
	for prefix, config := range defaultRateLimits {
		if strings.HasPrefix(path, prefix) {
			return config
		}
	}
	return defaultRateLimits["default"]
}
