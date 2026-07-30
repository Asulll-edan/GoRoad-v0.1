package badge

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	ListAllBadges(ctx context.Context) ([]Badge, error)
	FindBadgeByID(ctx context.Context, id uuid.UUID) (*Badge, error)
	FindBadgeByCode(ctx context.Context, code string) (*Badge, error)
	AwardBadge(ctx context.Context, ub *UserBadge) error
	GetUserBadges(ctx context.Context, userID uuid.UUID) ([]UserBadge, error)
	HasBadge(ctx context.Context, userID uuid.UUID, badgeID uuid.UUID) (bool, error)
}
