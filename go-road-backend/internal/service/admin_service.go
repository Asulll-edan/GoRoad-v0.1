package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/admin"
	"go-road-backend/internal/repository/redis"
)

type adminService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewAdminService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &adminService{repo: repo, cache: cache, logger: logger}
}

func (s *adminService) GetDashboard(ctx context.Context) (*domain.DashboardStats, error) {
	var cached domain.DashboardStats
	if err := s.cache.GetJSON(ctx, "cache:admin:dashboard", &cached); err == nil {
		return &cached, nil
	}
	stats, err := s.repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.SetJSON(ctx, "cache:admin:dashboard", stats, 5*time.Minute)
	return stats, nil
}

func (s *adminService) ListUsers(ctx context.Context, cursor string, limit int, search string) ([]domain.UserManagementRow, string, bool, error) {
	return s.repo.ListUsers(ctx, cursor, limit, search)
}

func (s *adminService) GetUser(ctx context.Context, id uuid.UUID) (*domain.UserManagementRow, error) {
	return s.repo.FindUserByID(ctx, id)
}

func (s *adminService) BanUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.BanUser(ctx, id)
}

func (s *adminService) UnbanUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.UnbanUser(ctx, id)
}

func (s *adminService) ListRooms(ctx context.Context, cursor string, limit int, status string) ([]domain.RoomManagementRow, string, bool, error) {
	return s.repo.ListRooms(ctx, cursor, limit, status)
}

func (s *adminService) GetRoom(ctx context.Context, id uuid.UUID) (*domain.RoomManagementRow, error) {
	return s.repo.FindRoomByID(ctx, id)
}

func (s *adminService) ListReports(ctx context.Context, cursor string, limit int, status string) ([]domain.ReportRow, string, bool, error) {
	return s.repo.ListReports(ctx, cursor, limit, status)
}

func (s *adminService) ReviewReport(ctx context.Context, id uuid.UUID, status string, userID uuid.UUID) error {
	return s.repo.ReviewReport(ctx, id, status, userID)
}

func (s *adminService) ListEmergency(ctx context.Context, cursor string, limit int) ([]domain.EmergencyEventRow, string, bool, error) {
	return s.repo.ListEmergencyEvents(ctx, cursor, limit)
}

func (s *adminService) GetAnalytics(ctx context.Context, period string) (map[string]interface{}, error) {
	return s.repo.GetAnalytics(ctx, period)
}

func (s *adminService) GetLogs(ctx context.Context, cursor string, limit int, level string) ([]domain.LogRow, string, bool, error) {
	return s.repo.GetLogs(ctx, cursor, limit, level)
}
