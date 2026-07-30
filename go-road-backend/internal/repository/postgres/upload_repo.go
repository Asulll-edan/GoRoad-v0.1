package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/upload"
)

type uploadRepository struct {
	db *Database
}

func NewUploadRepository(db *Database) domain.Repository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Create(ctx context.Context, upload *domain.Upload) error {
	return r.db.WithContext(ctx).Create(upload).Error
}

func (r *uploadRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Upload, error) {
	var u domain.Upload
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return &u, err
}

func (r *uploadRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Upload, error) {
	var uploads []domain.Upload
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Find(&uploads).Error
	return uploads, err
}

func (r *uploadRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Upload{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
