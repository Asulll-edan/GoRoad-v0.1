package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/checklist"
)

type checklistService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewChecklistService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &checklistService{repo: repo, logger: logger}
}

func (s *checklistService) ListTemplates(ctx context.Context, cursor string, limit int) ([]domain.ChecklistTemplate, string, bool, error) {
	return s.repo.ListTemplates(ctx, cursor, limit)
}

func (s *checklistService) CreateTemplate(ctx context.Context, name, description, category string, items []string, userID uuid.UUID) (*domain.ChecklistTemplate, error) {
	t := &domain.ChecklistTemplate{
		CreatedBy:   userID,
		Name:        name,
		Description: description,
		IsPublic:    true,
		Category:    category,
	}
	if err := s.repo.CreateTemplate(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}
	for i, label := range items {
		item := &domain.ChecklistItem{
			TemplateID: t.ID,
			Label:      label,
			OrderIndex: i + 1,
			IsRequired: i < 3,
		}
		if err := s.repo.AddItem(ctx, item); err != nil {
			return nil, fmt.Errorf("failed to add item: %w", err)
		}
	}
	return t, nil
}

func (s *checklistService) GetTemplate(ctx context.Context, id uuid.UUID) (*domain.ChecklistTemplate, error) {
	return s.repo.FindTemplateByID(ctx, id)
}

func (s *checklistService) CreateTouringChecklist(ctx context.Context, roomID, templateID, userID uuid.UUID) error {
	items, err := s.repo.ListItems(ctx, templateID)
	if err != nil {
		return err
	}
	for _, item := range items {
		tc := &domain.TouringChecklist{
			RoomID: roomID,
			UserID: userID,
			ItemID: item.ID,
		}
		if err := s.repo.CreateTouringChecklist(ctx, tc); err != nil {
			return fmt.Errorf("failed to create checklist item: %w", err)
		}
	}
	return nil
}

func (s *checklistService) GetTouringChecklist(ctx context.Context, roomID, userID uuid.UUID) ([]domain.TouringChecklist, error) {
	return s.repo.GetTouringChecklist(ctx, roomID, userID)
}

func (s *checklistService) ToggleItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	items, err := s.repo.GetTouringChecklist(ctx, id, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		newState := !item.IsChecked
		if err := s.repo.ToggleItem(ctx, item.ID, newState); err != nil {
			return err
		}
		break
	}
	return nil
}
