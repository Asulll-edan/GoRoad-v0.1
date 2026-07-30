package response

type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"has_more"`
	Total   *int64 `json:"total,omitempty"`
}

type PaginatedResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

type AuthResponse struct {
	User         UserResponse   `json:"user"`
	Tokens       TokenResponse  `json:"tokens"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Role        string `json:"role"`
	IsVerified  bool   `json:"is_verified"`
	CreatedAt   string `json:"created_at"`
}

type RoomResponse struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description,omitempty"`
	CoverURL        string        `json:"cover_url,omitempty"`
	Category        string        `json:"category,omitempty"`
	IsPrivate       bool          `json:"is_private"`
	MemberCount     int           `json:"member_count"`
	MaxMembers      int           `json:"max_members"`
	OriginCity      string        `json:"origin_city,omitempty"`
	DestinationCity string        `json:"destination_city,omitempty"`
	Status          string        `json:"status"`
	CreatedBy       string        `json:"created_by"`
	CreatedAt       string        `json:"created_at"`
	Members         []MemberResponse `json:"members,omitempty"`
}

type MemberResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at"`
}

type MotorResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model,omitempty"`
	Year         int     `json:"year"`
	PlateNumber  string  `json:"plate_number,omitempty"`
	Color        string  `json:"color,omitempty"`
	EngineCc     float64 `json:"engine_cc,omitempty"`
	FuelType     string  `json:"fuel_type,omitempty"`
	IsPrimary    bool    `json:"is_primary"`
	CreatedAt    string  `json:"created_at"`
}

type ChatMessageResponse struct {
	ID          string `json:"id"`
	RoomID      string `json:"room_id"`
	UserID      string `json:"user_id"`
	Username    string `json:"username,omitempty"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	CreatedAt   string `json:"created_at"`
}
