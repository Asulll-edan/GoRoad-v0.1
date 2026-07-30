package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/emergency"
	"go-road-backend/internal/repository/redis"
)

type emergencyService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewEmergencyService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &emergencyService{repo: repo, cache: cache, logger: logger}
}

func (s *emergencyService) ReportEmergency(ctx context.Context, req domain.CreateEmergencyRequest, userID uuid.UUID) (*domain.EmergencyEvent, error) {
	event := &domain.EmergencyEvent{
		ReportedBy:  userID,
		EventType:   req.EventType,
		Severity:    req.Severity,
		Description: req.Description,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.RoomID != "" {
		roomID, err := uuid.Parse(req.RoomID)
		if err != nil {
			return nil, errors.New("invalid room_id")
		}
		event.RoomID = &roomID
	}

	if req.Lat != 0 && req.Lng != 0 {
		event.Location = fmt.Sprintf("POINT(%f %f)", req.Lng, req.Lat)
	}

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create emergency event: %w", err)
	}

	data, _ := json.Marshal(event)
	if err := s.cache.Publish(ctx, "emergency:new", string(data)); err != nil {
		s.logger.Warn("failed to publish emergency event", zap.Error(err))
	}

	s.cache.Delete(ctx, "cache:emergency:recent")

	return event, nil
}

func (s *emergencyService) ListEmergencies(ctx context.Context, cursor string, limit int, status string) ([]domain.EmergencyEvent, string, bool, error) {
	return s.repo.ListEvents(ctx, cursor, limit, status)
}

func (s *emergencyService) GetEmergency(ctx context.Context, id uuid.UUID) (*domain.EmergencyEvent, error) {
	event, err := s.repo.FindEventByID(ctx, id)
	if err != nil {
		return nil, errors.New("emergency event not found")
	}
	return event, nil
}

func (s *emergencyService) AcknowledgeEmergency(ctx context.Context, id, userID uuid.UUID) error {
	event, err := s.repo.FindEventByID(ctx, id)
	if err != nil {
		return errors.New("emergency event not found")
	}

	event.Status = "acknowledged"

	if err := s.repo.UpdateEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to acknowledge emergency: %w", err)
	}

	data, _ := json.Marshal(map[string]interface{}{
		"event_id": id.String(),
		"user_id":  userID.String(),
		"action":   "acknowledged",
	})
	s.cache.Publish(ctx, "emergency:update", string(data))

	return nil
}

func (s *emergencyService) ResolveEmergency(ctx context.Context, id, userID uuid.UUID) error {
	event, err := s.repo.FindEventByID(ctx, id)
	if err != nil {
		return errors.New("emergency event not found")
	}

	event.Status = "resolved"
	event.ResolvedBy = &userID
	now := time.Now()
	event.ResolvedAt = &now

	if err := s.repo.UpdateEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to resolve emergency: %w", err)
	}

	s.cache.Delete(ctx, "cache:emergency:recent")

	return nil
}

func (s *emergencyService) TriggerSOS(ctx context.Context, userID uuid.UUID, lat, lng float64, roomID *uuid.UUID) (*domain.SOSEvent, error) {
	sos := &domain.SOSEvent{
		UserID:      userID,
		RoomID:      roomID,
		Location:    fmt.Sprintf("POINT(%f %f)", lng, lat),
		TriggeredAt: time.Now(),
		Status:      "active",
	}

	if err := s.repo.CreateSOS(ctx, sos); err != nil {
		return nil, fmt.Errorf("failed to create SOS: %w", err)
	}

	data, _ := json.Marshal(sos)
	if err := s.cache.Publish(ctx, "emergency:sos", string(data)); err != nil {
		s.logger.Warn("failed to publish SOS", zap.Error(err))
	}

	return sos, nil
}

func (s *emergencyService) DismissSOS(ctx context.Context, id, userID uuid.UUID) error {
	sos, err := s.repo.FindSOSByID(ctx, id)
	if err != nil {
		return errors.New("SOS event not found")
	}
	if sos.UserID != userID {
		return errors.New("unauthorized")
	}

	sos.Status = "dismissed"
	return s.repo.UpdateSOS(ctx, sos)
}
