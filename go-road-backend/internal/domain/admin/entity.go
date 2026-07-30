package admin

import "time"

type DashboardStats struct {
	TotalUsers      int64     `json:"total_users"`
	ActiveUsers     int64     `json:"active_users"`
	TotalRooms      int64     `json:"total_rooms"`
	ActiveTourings  int64     `json:"active_tourings"`
	TotalDistance   float64   `json:"total_distance"`
	EmergencyEvents int64     `json:"emergency_events"`
	ReportsPending  int64     `json:"reports_pending"`
	NewUsersToday   int64     `json:"new_users_today"`
	NewRoomsToday   int64     `json:"new_rooms_today"`
	Period          string    `json:"period"`
	GeneratedAt     time.Time `json:"generated_at"`
}

type UserManagementRow struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	IsVerified  bool      `json:"is_verified"`
	IsBanned    bool      `json:"is_banned"`
	TotalPoints int64     `json:"total_points"`
	RoomCount   int64     `json:"room_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type RoomManagementRow struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	MemberCount int       `json:"member_count"`
	CreatedBy   string    `json:"created_by"`
	StartDate   *string   `json:"start_date,omitempty"`
	DistanceKm  float64   `json:"distance_km,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
