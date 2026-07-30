package handler

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"go-road-backend/internal/service"
)

func handleAIChat(c fiber.Ctx) error {
	var req struct {
		RoomID  string `json:"room_id"`
		Message string `json:"message"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAIService(c.Locals("cache").(redis.CacheRepository), logger)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	ch, err := svc.Chat(c.Context(), userID, req.RoomID, req.Message)
	if err != nil {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": err.Error()})
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for resp := range ch {
			jsonData, _ := json.Marshal(resp)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			w.Flush()
		}
	})
	return nil
}

func handleGenerateItinerary(c fiber.Ctx) error {
	var req struct {
		RouteID     string   `json:"route_id"`
		StartLoc    string   `json:"start_location"`
		EndLoc      string   `json:"end_location"`
		DurationDays int     `json:"duration_days"`
		RiderCount  int      `json:"rider_count"`
		MotorIDs    []string `json:"motor_ids"`
		Preferences []string `json:"preferences"`
	}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAIService(c.Locals("cache").(redis.CacheRepository), logger)

	result, err := svc.GenerateItinerary(c.Context(), aiDomain.ItineraryRequest{
		RouteID:      req.RouteID,
		StartLoc:     req.StartLoc,
		EndLoc:       req.EndLoc,
		DurationDays: req.DurationDays,
		RiderCount:   req.RiderCount,
		MotorIDs:     req.MotorIDs,
		Preferences:  req.Preferences,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"itinerary": result})
}

func handleEstimateCost(c fiber.Ctx) error {
	var req struct {
		RouteID     string `json:"route_id"`
		RiderCount  int    `json:"rider_count"`
		DurationDays int   `json:"duration_days"`
		FuelType    string `json:"fuel_type"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	_ = userID
	svc := service.NewAIService(c.Locals("cache").(redis.CacheRepository), logger)

	result, err := svc.EstimateCost(c.Context(), req.RouteID, nil, req.RiderCount, req.DurationDays, req.FuelType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"estimate": result})
}

func handleRouteAdvice(c fiber.Ctx) error {
	var req struct {
		Origin      string   `json:"origin"`
		Destination string   `json:"destination"`
		Waypoints   []string `json:"waypoints"`
		Preferences []string `json:"preferences"`
	}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAIService(c.Locals("cache").(redis.CacheRepository), logger)

	result, err := svc.AdviseRoute(c.Context(), req.Origin, req.Destination, req.Waypoints, req.Preferences)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"advice": result})
}

func handlePackingList(c fiber.Ctx) error {
	var req struct {
		DurationDays    int    `json:"duration_days"`
		WeatherCondition string `json:"weather_condition"`
		TouringType     string `json:"touring_type"`
	}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAIService(c.Locals("cache").(redis.CacheRepository), logger)

	result, err := svc.GeneratePackingList(c.Context(), req.DurationDays, req.WeatherCondition, req.TouringType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"packing_list": result})
}

func handleSafetyAdvice(c fiber.Ctx) error {
	var req struct {
		RouteID     string `json:"route_id"`
		Weather     string `json:"weather_condition"`
		RiderCount  int    `json:"rider_count"`
		SkillLevel  string `json:"skill_level"`
	}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewAIService(c.Locals("cache").(redis.CacheRepository), logger)

	result, err := svc.AdviseSafety(c.Context(), req.RouteID, req.Weather, req.SkillLevel, req.RiderCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"safety": result})
}
