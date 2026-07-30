package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/expense"
)

type expenseRepository struct {
	db *Database
}

func NewExpenseRepository(db *Database) domain.Repository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *expenseRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error) {
	var expense domain.Expense
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&expense).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Expense, string, bool, error) {
	var expenses []domain.Expense
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			LoggedAt time.Time `json:"logged_at"`
			ID       string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(logged_at, id) < (?, ?)", c.LoggedAt, c.ID)
	}

	db = db.Order("logged_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&expenses).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(expenses) > limit
	if hasMore {
		expenses = expenses[:limit]
	}

	var nextCursor string
	if len(expenses) > 0 {
		last := expenses[len(expenses)-1]
		cursorData, _ := json.Marshal(struct {
			LoggedAt time.Time `json:"logged_at"`
			ID       string    `json:"id"`
		}{LoggedAt: last.LoggedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return expenses, nextCursor, hasMore, nil
}

func (r *expenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	return r.db.WithContext(ctx).Save(expense).Error
}

func (r *expenseRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Expense{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
