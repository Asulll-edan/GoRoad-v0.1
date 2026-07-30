package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleCreatePost(c fiber.Ctx) error {
	var req struct {
		Caption string   `json:"caption"`
		Photos  []string `json:"photos"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	userID, _ := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	post, err := svc.CreatePost(c.Context(), req.Caption, req.Photos, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(post)
}

func handleGetPost(c fiber.Ctx) error {
	id := c.Params("id")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(id)
	post, err := svc.GetPost(c.Context(), pid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
	}
	return c.JSON(post)
}

func handleDeletePost(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.DeletePost(c.Context(), pid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

func handleLikePost(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.LikePost(c.Context(), pid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "liked"})
}

func handleUnlikePost(c fiber.Ctx) error {
	id, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(id)
	uid, _ := uuid.Parse(userID)
	if err := svc.UnlikePost(c.Context(), pid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "unliked"})
}

func handleListComments(c fiber.Ctx) error {
	postID := c.Params("id")
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(postID)
	comments, cursor, hasMore, err := svc.ListComments(c.Context(), pid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: comments, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleCreateComment(c fiber.Ctx) error {
	postID := c.Params("id")
	userID := c.Locals("user_id").(string)
	var req struct {
		Content   string `json:"content"`
		ReplyToID string `json:"reply_to_id,omitempty"`
	}
	c.Bind().JSON(&req)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	pid, _ := uuid.Parse(postID)
	uid, _ := uuid.Parse(userID)
	var replyTo *uuid.UUID
	if req.ReplyToID != "" {
		if rid, err := uuid.Parse(req.ReplyToID); err == nil {
			replyTo = &rid
		}
	}
	comment, err := svc.CreateComment(c.Context(), pid, uid, req.Content, replyTo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(comment)
}

func handleDeleteComment(c fiber.Ctx) error {
	commentID, userID := c.Params("commentId"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	cid, _ := uuid.Parse(commentID)
	uid, _ := uuid.Parse(userID)
	if err := svc.DeleteComment(c.Context(), cid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

func handleReportContent(c fiber.Ctx) error {
	var req struct {
		ReportedType string `json:"reported_type"`
		ReportedID   string `json:"reported_id"`
		Reason       string `json:"reason"`
		Description  string `json:"description,omitempty"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	rid, _ := uuid.Parse(req.ReportedID)
	if err := svc.ReportContent(c.Context(), uid, req.ReportedType, rid, req.Reason, req.Description); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "reported"})
}

func handleGetFeed(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	posts, cursor, hasMore, err := svc.GetFeed(c.Context(), uid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: posts, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleExploreFeed(c fiber.Ctx) error {
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	posts, cursor, hasMore, err := svc.GetExploreFeed(c.Context(), params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: posts, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleGetFollowers(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	followers, cursor, hasMore, err := svc.GetFollowers(c.Context(), uid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: followers, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleGetFollowing(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	uid, _ := uuid.Parse(userID)
	following, cursor, hasMore, err := svc.GetFollowing(c.Context(), uid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: following, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleFollowUser(c fiber.Ctx) error {
	targetID, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	tuid, _ := uuid.Parse(targetID)
	uid, _ := uuid.Parse(userID)
	if err := svc.FollowUser(c.Context(), uid, tuid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "followed"})
}

func handleUnfollowUser(c fiber.Ctx) error {
	targetID, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	tuid, _ := uuid.Parse(targetID)
	uid, _ := uuid.Parse(userID)
	if err := svc.UnfollowUser(c.Context(), uid, tuid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "unfollowed"})
}

func handleBlockUser(c fiber.Ctx) error {
	targetID, userID := c.Params("id"), c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewSocialService(postgres.NewSocialRepository(c.Locals("db").(*postgres.Database)), c.Locals("cache").(redis.CacheRepository), logger)

	tuid, _ := uuid.Parse(targetID)
	uid, _ := uuid.Parse(userID)
	if err := svc.BlockUser(c.Context(), uid, tuid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "blocked"})
}

func handleUnblockUser(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "unblocked"})
}

func handleGetRiderProfile(c fiber.Ctx) error {
	return handleGetProfile(c)
}

func handleGetRiderBadges(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleGetRiderStats(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "not implemented"})
}

func handleGetRiderMotors(c fiber.Ctx) error {
	return handleListMotors(c)
}
