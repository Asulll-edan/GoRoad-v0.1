package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/admin"
	authDomain "go-road-backend/internal/domain/auth"
	socialDomain "go-road-backend/internal/domain/social"
)

type adminRepository struct {
	db *Database
}

func NewAdminRepository(db *Database) domain.Repository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	r.db.WithContext(ctx).Model(&authDomain.User{}).Where("deleted_at IS NULL").Count(&stats.TotalUsers)
	r.db.WithContext(ctx).Model(&authDomain.User{}).Where("is_active = true AND deleted_at IS NULL").Count(&stats.ActiveUsers)

	var rooms [5]struct{ Count int64 }
	r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM touring_rooms WHERE deleted_at IS NULL`).Scan(&rooms[0])
	stats.TotalRooms = rooms[0].Count

	r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM touring_rooms WHERE status = 'touring' AND deleted_at IS NULL`).Scan(&rooms[1])
	stats.ActiveTourings = rooms[1].Count

	r.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(distance_km), 0) FROM touring_rooms WHERE deleted_at IS NULL`).Scan(&rooms[2])
	stats.TotalDistance = float64(rooms[2].Count)

	var emergencyCount int64
	r.db.WithContext(ctx).Model(&struct{}{}).Table("emergency_events").Where("status = 'active'").Count(&emergencyCount)
	stats.EmergencyEvents = emergencyCount

	var reportsPending int64
	r.db.WithContext(ctx).Model(&socialDomain.Report{}).Where("status = 'pending'").Count(&reportsPending)
	stats.ReportsPending = reportsPending

	today := time.Now().Truncate(24 * time.Hour)
	r.db.WithContext(ctx).Model(&authDomain.User{}).Where("created_at >= ?", today).Count(&stats.NewUsersToday)

	var newRooms int64
	r.db.WithContext(ctx).Table("touring_rooms").Where("created_at >= ?", today).Count(&newRooms)
	stats.NewRoomsToday = newRooms

	stats.GeneratedAt = time.Now()
	return stats, nil
}

func (r *adminRepository) ListUsers(ctx context.Context, cursor string, limit int, search string) ([]domain.UserManagementRow, string, bool, error) {
	var rows []domain.UserManagementRow
	db := r.db.WithContext(ctx).Table("users u").
		Select("u.id, u.username, u.full_name, u.email, u.is_verified, u.is_banned, u.total_points, u.created_at").
		Where("u.deleted_at IS NULL")

	if search != "" {
		db = db.Where("u.username ILIKE ? OR u.full_name ILIKE ? OR u.email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(u.created_at, u.id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("u.created_at DESC, u.id DESC").Limit(limit + 1)
	if err := db.Find(&rows).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	var nextCursor string
	if len(rows) > 0 {
		last := rows[len(rows)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return rows, nextCursor, hasMore, nil
}

func (r *adminRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.UserManagementRow, error) {
	var row domain.UserManagementRow
	err := r.db.WithContext(ctx).Table("users").
		Select("id, username, full_name, email, is_verified, is_banned, total_points, created_at").
		Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	return &row, err
}

func (r *adminRepository) BanUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&authDomain.User{}).Where("id = ?", id).
		Update("is_banned", true).Error
}

func (r *adminRepository) UnbanUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&authDomain.User{}).Where("id = ?", id).
		Update("is_banned", false).Error
}

func (r *adminRepository) ListRooms(ctx context.Context, cursor string, limit int, status string) ([]domain.RoomManagementRow, string, bool, error) {
	var rows []domain.RoomManagementRow
	db := r.db.WithContext(ctx).Table("touring_rooms tr").
		Select("tr.id, tr.name, tr.status, tr.start_date::text as start_date, tr.distance_km, tr.created_by, tr.created_at").
		Where("tr.deleted_at IS NULL")
	if status != "" {
		db = db.Where("tr.status = ?", status)
	}
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(tr.created_at, tr.id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("tr.created_at DESC, tr.id DESC").Limit(limit + 1)
	if err := db.Find(&rows).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	var nextCursor string
	if len(rows) > 0 {
		last := rows[len(rows)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return rows, nextCursor, hasMore, nil
}

func (r *adminRepository) FindRoomByID(ctx context.Context, id uuid.UUID) (*domain.RoomManagementRow, error) {
	var row domain.RoomManagementRow
	var startDate *string
	err := r.db.WithContext(ctx).Table("touring_rooms").
		Select("id, name, status, start_date::text as start_date, distance_km, created_by, created_at").
		Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if startDate != nil {
		row.StartDate = startDate
	}
	return &row, err
}

func (r *adminRepository) ListReports(ctx context.Context, cursor string, limit int, status string) ([]domain.ReportRow, string, bool, error) {
	var rows []domain.ReportRow
	db := r.db.WithContext(ctx).Table("reports").
		Select("id, reporter_id, reported_type, reason, status, created_at")
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&rows).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	var nextCursor string
	if len(rows) > 0 {
		last := rows[len(rows)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return rows, nextCursor, hasMore, nil
}

func (r *adminRepository) ReviewReport(ctx context.Context, id uuid.UUID, status string, reviewedBy uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&socialDomain.Report{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"reviewed_by":  reviewedBy,
			"reviewed_at":  time.Now(),
		}).Error
}

func (r *adminRepository) ListEmergencyEvents(ctx context.Context, cursor string, limit int) ([]domain.EmergencyEventRow, string, bool, error) {
	var rows []domain.EmergencyEventRow
	db := r.db.WithContext(ctx).Table("emergency_events").
		Select("id, event_type, severity, status, created_at")
	if cursor != "" {
		data, _ := base64.URLEncoding.DecodeString(cursor)
		var c struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}
		json.Unmarshal(data, &c)
		db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
	}
	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&rows).Error; err != nil {
		return nil, "", false, err
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	var nextCursor string
	if len(rows) > 0 {
		last := rows[len(rows)-1]
		cData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cData)
	}
	return rows, nextCursor, hasMore, nil
}

func (r *adminRepository) GetAnalytics(ctx context.Context, period string) (map[string]interface{}, error) {
	stats, err := r.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"stats":  stats,
		"period": period,
	}, nil
}

func (r *adminRepository) GetLogs(ctx context.Context, cursor string, limit int, level string) ([]domain.LogRow, string, bool, error) {
	return nil, "", false, nil
}
