package admin

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	GetDashboard(ctx context.Context) (*DashboardStats, error)
	ListUsers(ctx context.Context, cursor string, limit int, search string) ([]UserManagementRow, string, bool, error)
	GetUser(ctx context.Context, id uuid.UUID) (*UserManagementRow, error)
	BanUser(ctx context.Context, id uuid.UUID) error
	UnbanUser(ctx context.Context, id uuid.UUID) error
	ListRooms(ctx context.Context, cursor string, limit int, status string) ([]RoomManagementRow, string, bool, error)
	GetRoom(ctx context.Context, id uuid.UUID) (*RoomManagementRow, error)
	ListReports(ctx context.Context, cursor string, limit int, status string) ([]ReportRow, string, bool, error)
	ReviewReport(ctx context.Context, id uuid.UUID, status string, userID uuid.UUID) error
	ListEmergency(ctx context.Context, cursor string, limit int) ([]EmergencyEventRow, string, bool, error)
	GetAnalytics(ctx context.Context, period string) (map[string]interface{}, error)
	GetLogs(ctx context.Context, cursor string, limit int, level string) ([]LogRow, string, bool, error)
}
