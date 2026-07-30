package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func NewCORSConfig(isProduction bool) cors.Config {
	config := cors.Config{
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
			"X-Cursor",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	if isProduction {
		config.AllowOrigins = []string{
			"https://goroad.app",
			"https://admin.goroad.app",
		}
	} else {
		config.AllowOrigins = []string{
			"*",
		}
	}

	return config
}
