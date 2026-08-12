package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/notification"
	"go-road-backend/internal/repository/redis"
)

type notificationService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewNotificationService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &notificationService{repo: repo, cache: cache, logger: logger}
}

func (s *notificationService) ListNotifications(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Notification, string, bool, error) {
	return s.repo.FindByUserID(ctx, userID, cursor, limit)
}

func (s *notificationService) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	if err := s.repo.MarkRead(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	cacheKey := fmt.Sprintf("cache:notif:unread:%s", userID.String())
	s.cache.Delete(ctx, cacheKey)

	return nil
}

func (s *notificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.MarkAllRead(ctx, userID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	cacheKey := fmt.Sprintf("cache:notif:unread:%s", userID.String())
	s.cache.Delete(ctx, cacheKey)

	return nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	cacheKey := fmt.Sprintf("cache:notif:unread:%s", userID.String())

	var cached int64
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	count, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	if err := s.cache.SetJSON(ctx, cacheKey, count, 5*time.Minute); err != nil {
		s.logger.Warn("failed to cache unread count", zap.Error(err))
	}

	return count, nil
}

func (s *notificationService) UpdatePreferences(ctx context.Context, userID uuid.UUID, prefs map[string]interface{}) error {
	existing, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		existing = &domain.NotificationPreference{
			UserID: userID,
		}
	}

	if v, ok := prefs["push_enabled"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			parsed := make(map[string]bool)
			for k, val := range m {
				parsed[k] = val.(bool)
			}
			existing.PushEnabled = parsed
		}
	}
	if v, ok := prefs["email_enabled"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			parsed := make(map[string]bool)
			for k, val := range m {
				parsed[k] = val.(bool)
			}
			existing.EmailEnabled = parsed
		}
	}
	if v, ok := prefs["quiet_hours_from"]; ok {
		existing.QuietHoursFrom = v.(string)
	}
	if v, ok := prefs["quiet_hours_to"]; ok {
		existing.QuietHoursTo = v.(string)
	}

	return s.repo.UpsertPreferences(ctx, existing)
}
