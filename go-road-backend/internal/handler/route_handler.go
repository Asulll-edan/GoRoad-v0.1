package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/route"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleCreateRoute(c fiber.Ctx) error {
	var req domain.CreateRouteRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewRouteService(postgres.NewRouteRepository(c.Locals("db").(*postgres.Database)), logger)

	uid, _ := uuid.Parse(userID)
	route, err := svc.Create(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(route)
}

func handleGetRoute(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewRouteService(postgres.NewRouteRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(id)
	route, err := svc.GetRoute(c.Context(), rid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "route not found"})
	}
	return c.JSON(route)
}

func handleUpdateRoute(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	var req map[string]interface{}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewRouteService(postgres.NewRouteRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.Update(c.Context(), rid, req, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

func handleDeleteRoute(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewRouteService(postgres.NewRouteRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.Delete(c.Context(), rid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

func handleAddWaypoint(c fiber.Ctx) error {
	routeID := c.Params("id")
	var input domain.WaypointInput
	c.Bind().JSON(&input)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewRouteService(postgres.NewRouteRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(routeID)
	wp, err := svc.AddWaypoint(c.Context(), rid, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(wp)
}

func handleListWaypoints(c fiber.Ctx) error {
	routeID := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewRouteService(postgres.NewRouteRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(routeID)
	wps, err := svc.ListWaypoints(c.Context(), rid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(wps)
}

func handleImportGPX(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleExportGPX(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleActivateRoute(c fiber.Ctx) error {
	routeID := c.Params("id")
	roomID := c.Query("room_id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewRouteService(postgres.NewRouteRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(routeID)
	rmid, _ := uuid.Parse(roomID)
	if err := svc.Activate(c.Context(), rid, rmid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "activated"})
}
