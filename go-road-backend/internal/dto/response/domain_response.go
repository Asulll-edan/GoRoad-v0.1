package response

type RouteResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	DistanceKm  float64 `json:"distance_km"`
	DurationH   float64 `json:"duration_hours,omitempty"`
	ElevationG  float64 `json:"elevation_gain,omitempty"`
	Polyline    string  `json:"polyline,omitempty"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	Waypoints   []WaypointResponse `json:"waypoints,omitempty"`
}

type WaypointResponse struct {
	ID        string  `json:"id"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Name      string  `json:"name,omitempty"`
	OrderIdx  int     `json:"order_index"`
}

type ConvoyResponse struct {
	ID            string `json:"id"`
	RoomID        string `json:"room_id"`
	FormationType string `json:"formation_type"`
	Status        string `json:"status"`
	LeadUserID    string `json:"lead_user_id"`
	SweepUserID   string `json:"sweep_user_id,omitempty"`
}

type EmergencyResponse struct {
	ID          string  `json:"id"`
	RoomID      string  `json:"room_id"`
	UserID      string  `json:"user_id"`
	Type        string  `json:"type"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Description string  `json:"description,omitempty"`
	Severity    string  `json:"severity"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type VotingResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	Options     []VotingOptionResponse `json:"options"`
	IsActive    bool              `json:"is_active"`
	IsAnonymous bool              `json:"is_anonymous"`
	CreatedAt   string            `json:"created_at"`
	RoomID      string            `json:"room_id"`
}

type VotingOptionResponse struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Votes  int    `json:"votes,omitempty"`
}

type FuelLogResponse struct {
	ID       string  `json:"id"`
	MotorID  string  `json:"motor_id"`
	Liters   float64 `json:"liters"`
	Price    float64 `json:"price"`
	Odometer float64 `json:"odometer,omitempty"`
	FuelType string  `json:"fuel_type"`
	Station  string  `json:"station,omitempty"`
	FilledAt string  `json:"filled_at"`
}

type ExpenseResponse struct {
	ID          string  `json:"id"`
	RoomID      string  `json:"room_id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type NotificationResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Type    string `json:"type"`
	IsRead  bool   `json:"is_read"`
	Data    string `json:"data,omitempty"`
	CreatedAt string `json:"created_at"`
}

type BadgeResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	AwardedAt   string `json:"awarded_at,omitempty"`
}

type SocialPostResponse struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	Username     string   `json:"username,omitempty"`
	Content      string   `json:"content"`
	Images       []string `json:"images,omitempty"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
	IsLiked      bool     `json:"is_liked"`
	CreatedAt    string   `json:"created_at"`
}

type PoiResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Address  string  `json:"address,omitempty"`
	Phone    string  `json:"phone,omitempty"`
	Distance float64 `json:"distance_km,omitempty"`
}

type ChecklistTemplateResponse struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Items    []ChecklistItemResponse `json:"items"`
}

type ChecklistItemResponse struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	IsDone   bool   `json:"is_done"`
}

type QRResponse struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Data     string `json:"data"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type ServiceReminderResponse struct {
	ID       string `json:"id"`
	MotorID  string `json:"motor_id"`
	Title    string `json:"title"`
	DueDate  string `json:"due_date"`
	Status   string `json:"status"`
	Notes    string `json:"notes,omitempty"`
}

type UploadResponse struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type AdminDashboardResponse struct {
	TotalUsers      int64                         `json:"total_users"`
	ActiveUsers24h  int64                         `json:"active_users_24h"`
	TotalRooms      int64                         `json:"total_rooms"`
	ActiveRooms     int64                         `json:"active_rooms"`
	TotalRoutes     int64                         `json:"total_routes"`
	TotalDistanceKm float64                       `json:"total_distance_km"`
	Emergency24h    int64                         `json:"emergency_events_24h"`
	PendingReports  int64                         `json:"reports_pending"`
	UserGrowth      []ChartDataPoint              `json:"user_growth"`
	RoomActivity    []ChartDataPoint              `json:"room_activity"`
	TopRooms        []TopRoomResponse             `json:"top_rooms"`
}

type ChartDataPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type TopRoomResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"member_count"`
}
