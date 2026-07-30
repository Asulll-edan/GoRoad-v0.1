package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/room"
	"go-road-backend/internal/repository/redis"
)

type roomService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewRoomService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &roomService{repo: repo, cache: cache, logger: logger}
}

func (s *roomService) Create(ctx context.Context, req domain.CreateRoomRequest, userID uuid.UUID) (*domain.RoomResponse, error) {
	room := &domain.TouringRoom{
		Name:          req.Name,
		Description:   req.Description,
		StartLocation: req.StartLocation,
		EndLocation:   req.EndLocation,
		MaxMembers:    req.MaxMembers,
		IsPublic:      req.IsPublic,
		TouringType:   req.TouringType,
		Difficulty:    req.Difficulty,
		CreatedBy:     userID,
		Status:        "planning",
	}

	if room.MaxMembers <= 0 {
		room.MaxMembers = 20
	}

	if err := s.repo.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	owner := &domain.RoomMember{
		RoomID: room.ID,
		UserID: userID,
		Role:   "owner",
	}
	if err := s.repo.AddMember(ctx, owner); err != nil {
		return nil, fmt.Errorf("failed to add owner: %w", err)
	}

	return s.toResponse(room), nil
}

func (s *roomService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID, includes []string) (*domain.TouringRoom, error) {
	cacheKey := fmt.Sprintf("cache:room:%s", id.String())
	var cached domain.TouringRoom
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("room not found")
	}

	_ = userID
	_ = includes

	s.cache.SetJSON(ctx, cacheKey, room, 5*time.Minute)
	return room, nil
}

func (s *roomService) Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) error {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("room not found")
	}

	member, _ := s.repo.GetMember(ctx, id, userID)
	if member == nil || (member.Role != "owner" && member.Role != "co_owner") {
		return errors.New("insufficient permissions")
	}

	for k, v := range req {
		switch k {
		case "name":
			room.Name = v.(string)
		case "description":
			room.Description = v.(string)
		case "status":
			room.Status = v.(string)
		case "is_public":
			room.IsPublic = v.(bool)
		}
	}

	if err := s.repo.Update(ctx, room); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	s.cache.Delete(ctx, fmt.Sprintf("cache:room:%s", id.String()))
	return nil
}

func (s *roomService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	member, _ := s.repo.GetMember(ctx, id, userID)
	if member == nil || member.Role != "owner" {
		return errors.New("only owner can delete room")
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	s.cache.Delete(ctx, fmt.Sprintf("cache:room:%s", id.String()))
	return nil
}

func (s *roomService) List(ctx context.Context, cursor string, limit int, userID uuid.UUID) ([]domain.RoomResponse, string, bool, error) {
	_ = userID
	rooms, nextCursor, hasMore, err := s.repo.List(ctx, cursor, limit, nil)
	if err != nil {
		return nil, "", false, err
	}

	responses := make([]domain.RoomResponse, len(rooms))
	for i, r := range rooms {
		responses[i] = *s.toResponse(&r)
	}

	return responses, nextCursor, hasMore, nil
}

func (s *roomService) Discover(ctx context.Context, cursor string, limit int) ([]domain.RoomResponse, string, bool, error) {
	cacheKey := fmt.Sprintf("cache:room:discover:%s:%d", cursor, limit)
	var cached struct {
		Rooms   []domain.RoomResponse `json:"rooms"`
		Cursor  string                `json:"cursor"`
		HasMore bool                  `json:"has_more"`
	}
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return cached.Rooms, cached.Cursor, cached.HasMore, nil
	}

	rooms, nextCursor, hasMore, err := s.repo.Discover(ctx, cursor, limit)
	if err != nil {
		return nil, "", false, err
	}

	responses := make([]domain.RoomResponse, len(rooms))
	for i, r := range rooms {
		responses[i] = *s.toResponse(&r)
	}

	s.cache.SetJSON(ctx, cacheKey, map[string]interface{}{
		"rooms":    responses,
		"cursor":   nextCursor,
		"has_more": hasMore,
	}, 5*time.Minute)

	return responses, nextCursor, hasMore, nil
}

func (s *roomService) Join(ctx context.Context, roomID, userID uuid.UUID) error {
	member, _ := s.repo.GetMember(ctx, roomID, userID)
	if member != nil {
		return errors.New("already a member")
	}

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	count, _ := s.repo.CountMembers(ctx, roomID)
	if count >= int64(room.MaxMembers) {
		return errors.New("room is full")
	}

	newMember := &domain.RoomMember{
		RoomID: roomID,
		UserID: userID,
		Role:   "member",
	}
	return s.repo.AddMember(ctx, newMember)
}

func (s *roomService) Leave(ctx context.Context, roomID, userID uuid.UUID) error {
	member, _ := s.repo.GetMember(ctx, roomID, userID)
	if member == nil {
		return errors.New("not a member")
	}
	if member.Role == "owner" {
		return errors.New("owner cannot leave, transfer ownership first")
	}
	return s.repo.RemoveMember(ctx, roomID, userID)
}

func (s *roomService) UpdateMemberRole(ctx context.Context, roomID, userID, targetUserID uuid.UUID, role string) error {
	member, _ := s.repo.GetMember(ctx, roomID, userID)
	if member == nil || (member.Role != "owner" && member.Role != "co_owner") {
		return errors.New("insufficient permissions")
	}

	if err := s.repo.UpdateMemberRole(ctx, roomID, targetUserID, role); err != nil {
		return err
	}

	s.repo.AddRoleHistory(ctx, &domain.RoomRoleHistory{
		RoomID:    roomID,
		UserID:    targetUserID,
		NewRole:   role,
		ChangedBy: userID,
	})

	s.cache.Delete(ctx,
		fmt.Sprintf("cache:room:%s:members", roomID),
		fmt.Sprintf("cache:room:%s", roomID),
	)

	return nil
}

func (s *roomService) ListMembers(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]domain.RoomMember, string, bool, error) {
	cacheKey := fmt.Sprintf("cache:room:%s:members:%s:%d", roomID, cursor, limit)
	var cached struct {
		Members []domain.RoomMember `json:"members"`
		Cursor  string              `json:"cursor"`
		HasMore bool                `json:"has_more"`
	}
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return cached.Members, cached.Cursor, cached.HasMore, nil
	}

	members, nextCursor, hasMore, err := s.repo.ListMembers(ctx, roomID, cursor, limit)
	if err != nil {
		return nil, "", false, err
	}

	s.cache.SetJSON(ctx, cacheKey, map[string]interface{}{
		"members":  members,
		"cursor":   nextCursor,
		"has_more": hasMore,
	}, 2*time.Minute)

	return members, nextCursor, hasMore, nil
}

func (s *roomService) StartTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	member, _ := s.repo.GetMember(ctx, roomID, userID)
	if member == nil || (member.Role != "owner" && member.Role != "lead") {
		return errors.New("only lead can start touring")
	}

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}

	room.Status = "touring"
	return s.repo.Update(ctx, room)
}

func (s *roomService) PauseTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	_ = userID

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}

	room.Status = "paused"
	return s.repo.Update(ctx, room)
}

func (s *roomService) CompleteTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	member, _ := s.repo.GetMember(ctx, roomID, userID)
	if member == nil || (member.Role != "owner" && member.Role != "lead") {
		return errors.New("only lead can complete touring")
	}

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}

	room.Status = "completed"
	return s.repo.Update(ctx, room)
}

func (s *roomService) CancelTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	_ = userID

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}

	room.Status = "cancelled"
	return s.repo.Update(ctx, room)
}

func (s *roomService) toResponse(room *domain.TouringRoom) *domain.RoomResponse {
	return &domain.RoomResponse{
		ID:            room.ID,
		Name:          room.Name,
		Status:        room.Status,
		StartLocation: room.StartLocation,
		EndLocation:   room.EndLocation,
		DistanceKm:    room.DistanceKm,
		IsPublic:      room.IsPublic,
		TouringType:   room.TouringType,
		Difficulty:    room.Difficulty,
		CreatedBy:     room.CreatedBy,
		CreatedAt:     room.CreatedAt,
	}
}
