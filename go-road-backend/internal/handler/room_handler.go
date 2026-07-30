package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/room"
	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/repository/redis"
	"go-road-backend/internal/service"
)

func handleListRooms(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	userID, _ := c.Locals("user_id").(string)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	uid, _ := uuid.Parse(userID)
	rooms, cursor, hasMore, err := roomService.List(c.Context(), params.Cursor, params.Limit, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pagination.PaginatedResponse{
		Data: rooms,
		Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore},
	})
}

func handleDiscoverRooms(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	rooms, cursor, hasMore, err := roomService.Discover(c.Context(), params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pagination.PaginatedResponse{
		Data: rooms,
		Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore},
	})
}

func handleCreateRoom(c fiber.Ctx) error {
	var req domain.CreateRoomRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, _ := c.Locals("user_id").(string)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	uid, _ := uuid.Parse(userID)
	resp, err := roomService.Create(c.Context(), req, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func handleGetRoom(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	uid, _ := uuid.Parse(userID)
	rid, _ := uuid.Parse(id)

	include := c.Query("include")
	var includes []string
	if include != "" {
		includes = []string{include}
	}

	room, err := roomService.Get(c.Context(), rid, uid, includes)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(room)
}

func handleUpdateRoom(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)

	var req map[string]interface{}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := roomService.Update(c.Context(), rid, req, uid); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "room updated"})
}

func handleDeleteRoom(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := roomService.Delete(c.Context(), rid, uid); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "room deleted"})
}

func handleJoinRoom(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := roomService.Join(c.Context(), rid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "joined room"})
}

func handleLeaveRoom(c fiber.Ctx) error {
	id := c.Params("id")
	userID, _ := c.Locals("user_id").(string)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	rid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)

	if err := roomService.Leave(c.Context(), rid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "left room"})
}

func handleListRoomMembers(c fiber.Ctx) error {
	id := c.Params("id")
	params := pagination.ParsePaginationParams(c)

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	rid, _ := uuid.Parse(id)
	members, cursor, hasMore, err := roomService.ListMembers(c.Context(), rid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(pagination.PaginatedResponse{
		Data: members,
		Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore},
	})
}

func handleUpdateMemberRole(c fiber.Ctx) error {
	roomID := c.Params("id")
	targetUserID := c.Params("userId")
	userID, _ := c.Locals("user_id").(string)

	var req struct {
		Role string `json:"role"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	db := c.Locals("db").(*postgres.Database)
	cache := c.Locals("cache").(redis.CacheRepository)
	logger := c.Locals("logger").(*zap.Logger)

	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, cache, logger)

	rid, _ := uuid.Parse(roomID)
	tuid, _ := uuid.Parse(targetUserID)
	uid, _ := uuid.Parse(userID)

	if err := roomService.UpdateMemberRole(c.Context(), rid, uid, tuid, req.Role); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "role updated"})
}

func handleGetRoomSettings(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleUpdateRoomSettings(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleStartTouring(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handlePauseTouring(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleCompleteTouring(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleCancelTouring(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}
