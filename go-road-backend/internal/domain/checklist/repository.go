package checklist

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTemplate(ctx context.Context, t *ChecklistTemplate) error
	FindTemplateByID(ctx context.Context, id uuid.UUID) (*ChecklistTemplate, error)
	ListTemplates(ctx context.Context, cursor string, limit int) ([]ChecklistTemplate, string, bool, error)
	AddItem(ctx context.Context, item *ChecklistItem) error
	ListItems(ctx context.Context, templateID uuid.UUID) ([]ChecklistItem, error)
	CreateTouringChecklist(ctx context.Context, tc *TouringChecklist) error
	GetTouringChecklist(ctx context.Context, roomID, userID uuid.UUID) ([]TouringChecklist, error)
	ToggleItem(ctx context.Context, id uuid.UUID, isChecked bool) error
}
