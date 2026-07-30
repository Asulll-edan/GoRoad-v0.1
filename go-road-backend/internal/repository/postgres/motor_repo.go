package postgres

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	domain "go-road-backend/internal/domain/motor"
)

type motorRepository struct {
	db *Database
}

func NewMotorRepository(db *Database) domain.Repository {
	return &motorRepository{db: db}
}

func (r *motorRepository) Create(ctx context.Context, motor *domain.Motor) error {
	return r.db.WithContext(ctx).Create(motor).Error
}

func (r *motorRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Motor, error) {
	var motor domain.Motor
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&motor).Error
	if err != nil {
		return nil, err
	}
	return &motor, nil
}

func (r *motorRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Motor, error) {
	var motors []domain.Motor
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).Order("is_primary DESC, created_at DESC").Find(&motors).Error
	return motors, err
}

func (r *motorRepository) Update(ctx context.Context, motor *domain.Motor) error {
	return r.db.WithContext(ctx).Save(motor).Error
}

func (r *motorRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Motor{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *motorRepository) SetPrimary(ctx context.Context, userID, motorID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	tx.Model(&domain.Motor{}).Where("user_id = ?", userID).Update("is_primary", false)
	tx.Model(&domain.Motor{}).Where("id = ?", motorID).Update("is_primary", true)
	return tx.Commit().Error
}
