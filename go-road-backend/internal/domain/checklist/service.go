package checklist

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	ListTemplates(ctx context.Context, cursor string, limit int) ([]ChecklistTemplate, string, bool, error)
	CreateTemplate(ctx context.Context, name, description, category string, items []string, userID uuid.UUID) (*ChecklistTemplate, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*ChecklistTemplate, error)
	CreateTouringChecklist(ctx context.Context, roomID, templateID, userID uuid.UUID) error
	GetTouringChecklist(ctx context.Context, roomID, userID uuid.UUID) ([]TouringChecklist, error)
	ToggleItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
