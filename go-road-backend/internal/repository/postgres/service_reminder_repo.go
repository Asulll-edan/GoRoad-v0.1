package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/service_reminder"
)

type serviceReminderRepository struct {
	db *Database
}

func NewServiceReminderRepository(db *Database) domain.Repository {
	return &serviceReminderRepository{db: db}
}

func (r *serviceReminderRepository) Create(ctx context.Context, reminder *domain.ServiceReminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

func (r *serviceReminderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ServiceReminder, error) {
	var rem domain.ServiceReminder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rem).Error
	return &rem, err
}

func (r *serviceReminderRepository) FindByUserID(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.ServiceReminder, string, bool, error) {
	var reminders []domain.ServiceReminder
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)
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
	if err := db.Find(&reminders).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(reminders) > limit
	if hasMore {
		reminders = reminders[:limit]
	}
	var nextCursor string
	if len(reminders) > 0 {
		last := reminders[len(reminders)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return reminders, nextCursor, hasMore, nil
}

func (r *serviceReminderRepository) Update(ctx context.Context, reminder *domain.ServiceReminder) error {
	return r.db.WithContext(ctx).Save(reminder).Error
}

func (r *serviceReminderRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.ServiceReminder{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *serviceReminderRepository) Complete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.ServiceReminder{}).Where("id = ?", id).
		Updates(map[string]interface{}{"completed_at": time.Now()}).Error
}

func (r *serviceReminderRepository) FindDueReminders(ctx context.Context, dueBefore time.Time) ([]domain.ServiceReminder, error) {
	var reminders []domain.ServiceReminder
	err := r.db.WithContext(ctx).
		Where("due_date <= ? AND completed_at IS NULL AND deleted_at IS NULL", dueBefore).
		Find(&reminders).Error
	return reminders, err
}
