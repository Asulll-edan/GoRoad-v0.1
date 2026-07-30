package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go-road-backend/internal/pkg/pagination"
	"go-road-backend/internal/repository/postgres"
	"go-road-backend/internal/service"
)

func handleCreateVoting(c fiber.Ctx) error {
	roomID := c.Params("id")
	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		VotingType  string   `json:"voting_type"`
		Answers     []string `json:"answers"`
	}
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewVotingService(postgres.NewVotingRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(roomID)
	uid, _ := uuid.Parse(userID)
	v, err := svc.CreateVoting(c.Context(), rid, uid, req.Title, req.Description, req.VotingType, req.Answers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(v)
}

func handleListVotings(c fiber.Ctx) error {
	roomID := c.Params("id")
	params := pagination.ParsePaginationParams(c)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewVotingService(postgres.NewVotingRepository(c.Locals("db").(*postgres.Database)), logger)

	rid, _ := uuid.Parse(roomID)
	votings, cursor, hasMore, err := svc.ListVotings(c.Context(), rid, params.Cursor, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pagination.PaginatedResponse{Data: votings, Meta: pagination.Meta{Cursor: cursor, HasMore: hasMore}})
}

func handleGetVoting(c fiber.Ctx) error {
	votingID := c.Params("votingId")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewVotingService(postgres.NewVotingRepository(c.Locals("db").(*postgres.Database)), logger)

	vid, _ := uuid.Parse(votingID)
	v, err := svc.GetVoting(c.Context(), vid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "voting not found"})
	}
	return c.JSON(v)
}

func handleSubmitVote(c fiber.Ctx) error {
	votingID := c.Params("votingId")
	var req struct{ AnswerID string `json:"answer_id"` }
	c.Bind().JSON(&req)
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewVotingService(postgres.NewVotingRepository(c.Locals("db").(*postgres.Database)), logger)

	vid, _ := uuid.Parse(votingID)
	aid, _ := uuid.Parse(req.AnswerID)
	uid, _ := uuid.Parse(userID)
	if err := svc.SubmitVote(c.Context(), vid, aid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "voted"})
}

func handleCloseVoting(c fiber.Ctx) error {
	votingID := c.Params("votingId")
	userID := c.Locals("user_id").(string)
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewVotingService(postgres.NewVotingRepository(c.Locals("db").(*postgres.Database)), logger)

	vid, _ := uuid.Parse(votingID)
	uid, _ := uuid.Parse(userID)
	if err := svc.CloseVoting(c.Context(), vid, uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "closed"})
}

func handleGetVotingResults(c fiber.Ctx) error {
	votingID := c.Params("votingId")
	logger := c.Locals("logger").(*zap.Logger)
	svc := service.NewVotingService(postgres.NewVotingRepository(c.Locals("db").(*postgres.Database)), logger)

	vid, _ := uuid.Parse(votingID)
	results, err := svc.GetResults(c.Context(), vid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}
