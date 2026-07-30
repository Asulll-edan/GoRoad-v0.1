package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/voting"
)

type votingService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewVotingService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &votingService{repo: repo, logger: logger}
}

func (s *votingService) CreateVoting(ctx context.Context, roomID, userID uuid.UUID, title, description, votingType string, answerLabels []string) (*domain.Voting, error) {
	v := &domain.Voting{
		RoomID:      roomID,
		CreatedBy:   userID,
		Title:       title,
		Description: description,
		VotingType:  votingType,
		Status:      "active",
		StartsAt:    time.Now(),
	}
	if err := s.repo.CreateVoting(ctx, v); err != nil {
		return nil, fmt.Errorf("failed to create voting: %w", err)
	}
	for i, label := range answerLabels {
		a := &domain.VotingAnswer{
			VotingID:   v.ID,
			Label:      label,
			OrderIndex: i + 1,
		}
		if err := s.repo.AddAnswer(ctx, a); err != nil {
			return nil, fmt.Errorf("failed to add answer: %w", err)
		}
	}
	return v, nil
}

func (s *votingService) ListVotings(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]domain.Voting, string, bool, error) {
	return s.repo.ListVotingsByRoomID(ctx, roomID, cursor, limit)
}

func (s *votingService) GetVoting(ctx context.Context, id uuid.UUID) (*domain.Voting, error) {
	return s.repo.FindVotingByID(ctx, id)
}

func (s *votingService) SubmitVote(ctx context.Context, votingID, answerID, userID uuid.UUID) error {
	v, err := s.repo.FindVotingByID(ctx, votingID)
	if err != nil {
		return errors.New("voting not found")
	}
	if v.Status != "active" {
		return errors.New("voting is not active")
	}
	if v.VotingType == "single" {
		answers, _ := s.repo.ListAnswers(ctx, votingID)
		for _, a := range answers {
			existing, _ := s.repo.GetVote(ctx, votingID, a.ID, userID)
			if existing != nil {
				return errors.New("already voted")
			}
		}
	}
	vote := &domain.Vote{
		VotingID: votingID,
		AnswerID: answerID,
		UserID:   userID,
	}
	return s.repo.SubmitVote(ctx, vote)
}

func (s *votingService) CloseVoting(ctx context.Context, id, userID uuid.UUID) error {
	v, err := s.repo.FindVotingByID(ctx, id)
	if err != nil {
		return errors.New("voting not found")
	}
	if v.CreatedBy != userID {
		return errors.New("only creator can close voting")
	}
	return s.repo.CloseVoting(ctx, id)
}

func (s *votingService) GetResults(ctx context.Context, id uuid.UUID) ([]domain.VoteResult, error) {
	return s.repo.GetResults(ctx, id)
}
