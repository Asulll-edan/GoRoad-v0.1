package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/badge"
	"go-road-backend/internal/repository/redis"
)

type badgeService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewBadgeService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &badgeService{repo: repo, cache: cache, logger: logger}
}

func (s *badgeService) ListBadges(ctx context.Context) ([]domain.Badge, error) {
	return s.repo.ListAllBadges(ctx)
}

func (s *badgeService) GetMyBadges(ctx context.Context, userID uuid.UUID) ([]domain.Badge, error) {
	cacheKey := fmt.Sprintf("cache:user:%s:badges", userID)
	var cached []domain.Badge
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}
	userBadges, err := s.repo.GetUserBadges(ctx, userID)
	if err != nil {
		return nil, err
	}
	badges := make([]domain.Badge, 0, len(userBadges))
	for _, ub := range userBadges {
		b, err := s.repo.FindBadgeByID(ctx, ub.BadgeID)
		if err == nil {
			badges = append(badges, *b)
		}
	}
	s.cache.SetJSON(ctx, cacheKey, badges, 1*time.Hour)
	return badges, nil
}

func (s *badgeService) GetBadgeProgress(ctx context.Context, userID uuid.UUID) ([]domain.BadgeProgress, error) {
	allBadges, err := s.repo.ListAllBadges(ctx)
	if err != nil {
		return nil, err
	}
	userBadges, _ := s.repo.GetUserBadges(ctx, userID)
	owned := make(map[uuid.UUID]bool)
	for _, ub := range userBadges {
		owned[ub.BadgeID] = true
	}
	progress := make([]domain.BadgeProgress, 0, len(allBadges))
	for _, b := range allBadges {
		p := domain.BadgeProgress{
			BadgeCode: b.Code,
			BadgeName: b.Name,
			Target:    1,
			IsAwarded: owned[b.ID],
		}
		if criteria, ok := b.Criteria["value"]; ok {
			switch v := criteria.(type) {
			case float64:
				p.Target = v
			}
		}
		progress = append(progress, p)
	}
	return progress, nil
}

func (s *badgeService) EvaluateBadges(ctx context.Context, userID, roomID uuid.UUID, touringData map[string]interface{}) ([]string, error) {
	s.logger.Info("delegating badge evaluation to Python analytics",
		zap.String("user_id", userID.String()),
		zap.String("room_id", roomID.String()))
	return nil, nil
}
