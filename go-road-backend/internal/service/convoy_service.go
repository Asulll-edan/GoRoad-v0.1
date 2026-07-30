package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/convoy"
	"go-road-backend/internal/repository/redis"
)

type convoyService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewConvoyService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &convoyService{repo: repo, cache: cache, logger: logger}
}

func (s *convoyService) CreateFormation(ctx context.Context, roomID uuid.UUID, name, formationType string, userID uuid.UUID) (*domain.Formation, error) {
	f := &domain.Formation{
		RoomID:        roomID,
		Name:          name,
		FormationType: formationType,
		IsActive:      true,
		CreatedBy:     userID,
	}
	if err := s.repo.CreateFormation(ctx, f); err != nil {
		return nil, fmt.Errorf("failed to create formation: %w", err)
	}
	return f, nil
}

func (s *convoyService) UpdateFormation(ctx context.Context, formationID uuid.UUID, req map[string]interface{}, userID uuid.UUID) error {
	f, err := s.repo.FindActiveFormation(ctx, formationID)
	if err != nil {
		return fmt.Errorf("formation not found: %w", err)
	}
	for k, v := range req {
		switch k {
		case "name":
			f.Name = v.(string)
		case "formation_type":
			f.FormationType = v.(string)
		case "member_order":
			f.MemberOrder = v.([]string)
		case "speed_limit_kmh":
			f.SpeedLimitKmh = v.(float64)
		case "safe_distance_m":
			f.SafeDistanceM = v.(float64)
		}
	}
	return s.repo.UpdateFormation(ctx, f)
}

func (s *convoyService) GetActiveFormation(ctx context.Context, roomID uuid.UUID) (*domain.Formation, error) {
	return s.repo.FindActiveFormation(ctx, roomID)
}

func (s *convoyService) UpdateLocation(ctx context.Context, roomID, userID uuid.UUID, lat, lng, speed, heading float64) error {
	pos := &domain.RiderPosition{
		RoomID:    roomID,
		UserID:    userID,
		Lat:       lat,
		Lng:       lng,
		SpeedKmh:  speed,
		Heading:   heading,
		Timestamp: time.Now(),
	}

	cacheKey := fmt.Sprintf("pos:room:%s:rider:%s", roomID, userID)
	posJSON, _ := json.Marshal(pos)
	s.cache.Set(ctx, cacheKey, string(posJSON), 120*time.Second)

	s.cache.SAdd(ctx, fmt.Sprintf("pos:room:%s:all", roomID), userID.String())
	s.cache.Expire(ctx, fmt.Sprintf("pos:room:%s:all", roomID), 10*time.Second)

	channel := fmt.Sprintf("room:%s:location", roomID)
	s.cache.Publish(ctx, channel, string(posJSON))

	return s.repo.SavePosition(ctx, pos)
}

func (s *convoyService) GetLocations(ctx context.Context, roomID uuid.UUID) ([]domain.RiderPosition, error) {
	memberIDs, err := s.cache.SMembers(ctx, fmt.Sprintf("pos:room:%s:all", roomID))
	if err != nil || len(memberIDs) == 0 {
		return s.repo.GetPositionsInRoom(ctx, roomID)
	}

	positions := make([]domain.RiderPosition, 0, len(memberIDs))
	for _, mid := range memberIDs {
		uid, _ := uuid.Parse(mid)
		key := fmt.Sprintf("pos:room:%s:rider:%s", roomID, uid)
		data, err := s.cache.Get(ctx, key)
		if err != nil {
			continue
		}
		var pos domain.RiderPosition
		if json.Unmarshal([]byte(data), &pos) == nil {
			positions = append(positions, pos)
		}
	}
	return positions, nil
}

func (s *convoyService) GetTracking(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]domain.RiderPosition, string, bool, error) {
	return s.repo.GetPositionsSince(ctx, roomID, time.Now().Add(-5*time.Minute)), "", false, nil
}

func (s *convoyService) GetPositionsSince(ctx context.Context, roomID uuid.UUID, since time.Time) ([]domain.RiderPosition, error) {
	return s.repo.GetPositionsSince(ctx, roomID, since)
}
