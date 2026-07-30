package postgres

import (
	"context"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/qr"
)

type qrRepository struct {
	db *Database
}

func NewQRRepository(db *Database) domain.Repository {
	return &qrRepository{db: db}
}

func (r *qrRepository) Create(ctx context.Context, card *domain.QRCard) error {
	return r.db.WithContext(ctx).Create(card).Error
}

func (r *qrRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.QRCard, error) {
	var card domain.QRCard
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&card).Error
	return &card, err
}

func (r *qrRepository) FindByCode(ctx context.Context, code string) (*domain.QRCard, error) {
	var card domain.QRCard
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&card).Error
	return &card, err
}

func (r *qrRepository) IncrementScanCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.QRCard{}).
		Where("id = ?", id).
		Update("scan_count", gorm.Expr("scan_count + 1")).Error
}
