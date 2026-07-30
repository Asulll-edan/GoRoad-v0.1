package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleGetNearbyPOI(c fiber.Ctx) error {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	radius, _ := strconv.ParseFloat(c.Query("radius_km"), 64)
	typesStr := c.Query("types")

	if radius <= 0 {
		radius = 10
	}

	var types []string
	if typesStr != "" {
		types = strings.Split(typesStr, ",")
	}

	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewPOIService(postgres.NewPOIRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pois, err := svc.GetNearbyPOI(c.Context(), lat, lng, radius, types)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pois)
}

func handleGetPOIDetail(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewPOIService(postgres.NewPOIRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(id)
	poi, err := svc.GetPOIDetail(c.Context(), pid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "POI not found"})
	}
	return c.JSON(poi)
}

func handleReportPOI(c fiber.Ctx) error {
	var req struct {
		POIID       string `json:"poi_id"`
		ReportType  string `json:"report_type"`
		Description string `json:"description"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewPOIService(postgres.NewPOIRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(req.POIID)
	uid, _ := uuid.Parse(userID)
	if err := svc.ReportPOI(c.Context(), uid, pid, req.ReportType, req.Description); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "reported"})
}

func handleGetPOICategories(c fiber.Ctx) error {
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewPOIService(postgres.NewPOIRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	categories, err := svc.GetCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(categories)
}
