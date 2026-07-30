package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/expense"
)

type expenseService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewExpenseService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &expenseService{repo: repo, logger: logger}
}

func (s *expenseService) Create(ctx context.Context, req domain.CreateExpenseRequest, userID uuid.UUID) (*domain.Expense, error) {
	expense := &domain.Expense{
		UserID:      userID,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		IsSplitBill: req.IsSplitBill,
		SplitWith:   req.SplitWith,
		LoggedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return expense, nil
}

func (s *expenseService) Get(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	expense, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("expense not found")
	}
	return expense, nil
}

func (s *expenseService) List(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Expense, string, bool, error) {
	return s.repo.FindByUserID(ctx, userID, cursor, limit)
}

func (s *expenseService) Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) (*domain.Expense, error) {
	expense, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("expense not found")
	}
	if expense.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	for k, v := range req {
		switch k {
		case "category":
			expense.Category = v.(string)
		case "amount":
			expense.Amount = v.(float64)
		case "description":
			expense.Description = v.(string)
		case "is_split_bill":
			expense.IsSplitBill = v.(bool)
		case "split_with":
			sw, ok := v.([]interface{})
			if ok {
				splitWith := make([]string, len(sw))
				for i, s := range sw {
					splitWith[i] = s.(string)
				}
				expense.SplitWith = splitWith
			}
		case "receipt_url":
			expense.ReceiptURL = v.(string)
		}
	}

	if err := s.repo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to update expense: %w", err)
	}

	return expense, nil
}

func (s *expenseService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	expense, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("expense not found")
	}
	if expense.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.SoftDelete(ctx, id)
}
