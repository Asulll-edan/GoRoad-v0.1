package badge

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	ListBadges(ctx context.Context) ([]Badge, error)
	GetMyBadges(ctx context.Context, userID uuid.UUID) ([]Badge, error)
	GetBadgeProgress(ctx context.Context, userID uuid.UUID) ([]BadgeProgress, error)
	EvaluateBadges(ctx context.Context, userID, roomID uuid.UUID, touringData map[string]interface{}) ([]string, error)
}
