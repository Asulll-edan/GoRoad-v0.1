package postgres

import (
	"context"

	"github.com/google/uuid"

	domain "go-road-backend/internal/domain/badge"
	badgeDomain "go-road-backend/internal/domain/badge"
)

type badgeRepository struct {
	db *Database
}

func NewBadgeRepository(db *Database) domain.Repository {
	return &badgeRepository{db: db}
}

func (r *badgeRepository) ListAllBadges(ctx context.Context) ([]badgeDomain.Badge, error) {
	var badges []badgeDomain.Badge
	err := r.db.WithContext(ctx).Order("tier ASC, code ASC").Find(&badges).Error
	return badges, err
}

func (r *badgeRepository) FindBadgeByID(ctx context.Context, id uuid.UUID) (*badgeDomain.Badge, error) {
	var b badgeDomain.Badge
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&b).Error
	return &b, err
}

func (r *badgeRepository) FindBadgeByCode(ctx context.Context, code string) (*badgeDomain.Badge, error) {
	var b badgeDomain.Badge
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&b).Error
	return &b, err
}

func (r *badgeRepository) AwardBadge(ctx context.Context, ub *badgeDomain.UserBadge) error {
	return r.db.WithContext(ctx).Create(ub).Error
}

func (r *badgeRepository) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]badgeDomain.UserBadge, error) {
	var badges []badgeDomain.UserBadge
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&badges).Error
	return badges, err
}

func (r *badgeRepository) HasBadge(ctx context.Context, userID uuid.UUID, badgeID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&badgeDomain.UserBadge{}).
		Where("user_id = ? AND badge_id = ?", userID, badgeID).Count(&count).Error
	return count > 0, err
}
