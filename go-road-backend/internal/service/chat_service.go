package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/chat"
	"go-road-backend/internal/repository/redis"
)

type chatService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewChatService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &chatService{repo: repo, cache: cache, logger: logger}
}

func (s *chatService) SendMessage(ctx context.Context, req domain.SendMessageRequest, senderID uuid.UUID) (*domain.Message, error) {
	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		return nil, errors.New("invalid room_id")
	}

	msg := &domain.Message{
		RoomID:      roomID,
		SenderID:    senderID,
		MessageType: req.MessageType,
		Content:     req.Content,
		SentAt:      time.Now(),
	}

	if req.ReplyToID != "" {
		replyID, err := uuid.Parse(req.ReplyToID)
		if err != nil {
			return nil, errors.New("invalid reply_to_id")
		}
		msg.ReplyToID = &replyID
	}

	if msg.MessageType == "" {
		msg.MessageType = "text"
	}

	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	data, _ := json.Marshal(msg)
	pubKey := fmt.Sprintf("room:%s:message", roomID.String())
	if err := s.cache.Publish(ctx, pubKey, string(data)); err != nil {
		s.logger.Warn("failed to publish message", zap.String("room_id", roomID.String()), zap.Error(err))
	}

	cacheKey := fmt.Sprintf("cache:room:%s:messages", roomID.String())
	s.cache.Delete(ctx, cacheKey)

	return msg, nil
}

func (s *chatService) GetMessage(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	msg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("message not found")
	}
	return msg, nil
}

func (s *chatService) EditMessage(ctx context.Context, id uuid.UUID, content string, userID uuid.UUID) error {
	msg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("message not found")
	}
	if msg.SenderID != userID {
		return errors.New("unauthorized")
	}

	msg.Content = content
	now := time.Now()
	msg.EditedAt = &now

	if err := s.repo.Update(ctx, msg); err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}

func (s *chatService) DeleteMessage(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	msg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("message not found")
	}
	if msg.SenderID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	cacheKey := fmt.Sprintf("cache:room:%s:messages", msg.RoomID.String())
	s.cache.Delete(ctx, cacheKey)
	return nil
}

func (s *chatService) PinMessage(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	msg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("message not found")
	}

	if err := s.repo.Pin(ctx, id, !msg.IsPinned); err != nil {
		return fmt.Errorf("failed to toggle pin: %w", err)
	}

	cacheKey := fmt.Sprintf("cache:room:%s:messages", msg.RoomID.String())
	s.cache.Delete(ctx, cacheKey)
	return nil
}

func (s *chatService) ListMessages(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]domain.Message, string, bool, error) {
	cacheKey := fmt.Sprintf("cache:room:%s:messages", roomID.String())
	var cached struct {
		Messages []domain.Message `json:"messages"`
		Cursor   string           `json:"cursor"`
		HasMore  bool             `json:"has_more"`
	}
	if cursor == "" {
		if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
			return cached.Messages, cached.Cursor, cached.HasMore, nil
		}
	}

	messages, nextCursor, hasMore, err := s.repo.FindByRoomID(ctx, roomID, cursor, limit)
	if err != nil {
		return nil, "", false, err
	}

	if cursor == "" {
		s.cache.SetJSON(ctx, cacheKey, map[string]interface{}{
			"messages": messages,
			"cursor":   nextCursor,
			"has_more": hasMore,
		}, 2*time.Minute)
	}

	return messages, nextCursor, hasMore, nil
}

func (s *chatService) MarkRead(ctx context.Context, msgID, userID uuid.UUID) error {
	if err := s.repo.MarkRead(ctx, msgID, userID); err != nil {
		return fmt.Errorf("failed to mark read: %w", err)
	}

	msg, err := s.repo.FindByID(ctx, msgID)
	if err != nil {
		return nil
	}

	pubKey := fmt.Sprintf("room:%s:read", msg.RoomID.String())
	receipt := map[string]interface{}{
		"message_id": msgID.String(),
		"user_id":    userID.String(),
		"read_at":    time.Now(),
	}
	data, _ := json.Marshal(receipt)
	s.cache.Publish(ctx, pubKey, string(data))

	return nil
}
