package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/checklist"
)

type checklistRepository struct {
	db *Database
}

func NewChecklistRepository(db *Database) domain.Repository {
	return &checklistRepository{db: db}
}

func (r *checklistRepository) CreateTemplate(ctx context.Context, t *domain.ChecklistTemplate) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *checklistRepository) FindTemplateByID(ctx context.Context, id uuid.UUID) (*domain.ChecklistTemplate, error) {
	var t domain.ChecklistTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error
	return &t, err
}

func (r *checklistRepository) ListTemplates(ctx context.Context, cursor string, limit int) ([]domain.ChecklistTemplate, string, bool, error) {
	var templates []domain.ChecklistTemplate
	db := r.db.WithContext(ctx)
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&templates).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(templates) > limit
	if hasMore {
		templates = templates[:limit]
	}
	var nextCursor string
	if len(templates) > 0 {
		last := templates[len(templates)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return templates, nextCursor, hasMore, nil
}

func (r *checklistRepository) AddItem(ctx context.Context, item *domain.ChecklistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *checklistRepository) ListItems(ctx context.Context, templateID uuid.UUID) ([]domain.ChecklistItem, error) {
	var items []domain.ChecklistItem
	err := r.db.WithContext(ctx).Where("template_id = ?", templateID).
		Order("order_index ASC").Find(&items).Error
	return items, err
}

func (r *checklistRepository) CreateTouringChecklist(ctx context.Context, tc *domain.TouringChecklist) error {
	return r.db.WithContext(ctx).Create(tc).Error
}

func (r *checklistRepository) GetTouringChecklist(ctx context.Context, roomID, userID uuid.UUID) ([]domain.TouringChecklist, error) {
	var items []domain.TouringChecklist
	err := r.db.WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).
		Find(&items).Error
	return items, err
}

func (r *checklistRepository) ToggleItem(ctx context.Context, id uuid.UUID, isChecked bool) error {
	updates := map[string]interface{}{"is_checked": isChecked}
	if isChecked {
		updates["checked_at"] = time.Now()
	} else {
		updates["checked_at"] = nil
	}
	return r.db.WithContext(ctx).Model(&domain.TouringChecklist{}).
		Where("id = ?", id).Updates(updates).Error
}
