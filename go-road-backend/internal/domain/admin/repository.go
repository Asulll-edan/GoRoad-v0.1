package admin

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	ListUsers(ctx context.Context, cursor string, limit int, search string) ([]UserManagementRow, string, bool, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*UserManagementRow, error)
	BanUser(ctx context.Context, id uuid.UUID) error
	UnbanUser(ctx context.Context, id uuid.UUID) error
	ListRooms(ctx context.Context, cursor string, limit int, status string) ([]RoomManagementRow, string, bool, error)
	FindRoomByID(ctx context.Context, id uuid.UUID) (*RoomManagementRow, error)
	ListReports(ctx context.Context, cursor string, limit int, status string) ([]ReportRow, string, bool, error)
	ReviewReport(ctx context.Context, id uuid.UUID, status string, reviewedBy uuid.UUID) error
	ListEmergencyEvents(ctx context.Context, cursor string, limit int) ([]EmergencyEventRow, string, bool, error)
	GetAnalytics(ctx context.Context, period string) (map[string]interface{}, error)
	GetLogs(ctx context.Context, cursor string, limit int, level string) ([]LogRow, string, bool, error)
}

type ReportRow struct {
	ID           uuid.UUID `json:"id"`
	ReporterID   uuid.UUID `json:"reporter_id"`
	ReportedType string    `json:"reported_type"`
	Reason       string    `json:"reason"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type EmergencyEventRow struct {
	ID        uuid.UUID `json:"id"`
	EventType string    `json:"event_type"`
	Severity  string    `json:"severity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type LogRow struct {
	ID        uuid.UUID `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
