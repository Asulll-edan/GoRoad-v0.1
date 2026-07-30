package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"go-road-backend/internal/service"
)

func handleGetCurrentWeather(c fiber.Ctx) error {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewWeatherService(c.Locals("cache").(redis.CacheRepository), logger)

	data, err := svc.GetCurrentWeather(c.Context(), lat, lng)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "weather data not available"})
	}
	return c.JSON(data)
}

func handleGetWeatherForecast(c fiber.Ctx) error {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewWeatherService(c.Locals("cache").(redis.CacheRepository), logger)

	data, err := svc.GetForecast(c.Context(), lat, lng)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "forecast not available"})
	}
	return c.JSON(data)
}

func handleGetWeatherAlerts(c fiber.Ctx) error {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewWeatherService(c.Locals("cache").(redis.CacheRepository), logger)

	alerts, err := svc.GetAlerts(c.Context(), lat, lng)
	if err != nil {
		alerts = []interface{}{}
	}
	return c.JSON(alerts)
}

func handleGetRouteWeather(c fiber.Ctx) error {
	routeID := c.Query("route_id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewWeatherService(c.Locals("cache").(redis.CacheRepository), logger)

	data, err := svc.GetRouteWeather(c.Context(), routeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "route weather not available"})
	}
	return c.JSON(data)
}
