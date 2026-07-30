package request

type CreatePoiRequest struct {
	Name        string   `json:"name" validate:"required,max=200"`
	Category    string   `json:"category" validate:"oneof=gas_station restaurant mosque hotel workshop rest_area attraction"`
	Lat         float64  `json:"lat" validate:"required"`
	Lon         float64  `json:"lon" validate:"required"`
	Address     string   `json:"address,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	PhotoURLs   []string `json:"photo_urls,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type NearbyPoiQuery struct {
	Lat      float64  `query:"lat" validate:"required"`
	Lon      float64  `query:"lon" validate:"required"`
	RadiusKm float64  `query:"radius" validate:"min=0.1,max=100"`
	Category string   `query:"category"`
	Limit    int      `query:"limit"`
}

type CreateReminderRequest struct {
	MotorID  string `json:"motor_id" validate:"required"`
	Title    string `json:"title" validate:"required,max=200"`
	DueDate  string `json:"due_date" validate:"required"`
	IntervalKm int   `json:"interval_km,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

type AdminActionRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Reason string `json:"reason" validate:"required,max=500"`
	Action string `json:"action" validate:"required,oneof=ban unban warn"`
}

type ReviewReportRequest struct {
	Action string `json:"action" validate:"required,oneof=dismiss warn ban"`
	Notes  string `json:"notes,omitempty"`
}
