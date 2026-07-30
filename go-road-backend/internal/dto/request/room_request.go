package request

type CreateRoomRequest struct {
	Name           string   `json:"name" validate:"required,min=3,max=100"`
	Description    string   `json:"description,omitempty" validate:"max=1000"`
	CoverURL       string   `json:"cover_url,omitempty"`
	Category       string   `json:"category,omitempty"`
	IsPrivate      bool     `json:"is_private"`
	MaxMembers     int      `json:"max_members" validate:"min=2,max=200"`
	OriginCity     string   `json:"origin_city,omitempty"`
	DestinationCity string  `json:"destination_city,omitempty"`
	Region         string   `json:"region,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	DepartureTime  string   `json:"departure_time,omitempty"`
}

type UpdateRoomRequest struct {
	Name           string   `json:"name,omitempty" validate:"min=3,max=100"`
	Description    string   `json:"description,omitempty" validate:"max=1000"`
	CoverURL       string   `json:"cover_url,omitempty"`
	Category       string   `json:"category,omitempty"`
	IsPrivate      *bool    `json:"is_private,omitempty"`
	MaxMembers     int      `json:"max_members,omitempty" validate:"min=2,max=200"`
	OriginCity     string   `json:"origin_city,omitempty"`
	DestinationCity string  `json:"destination_city,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

type UpdateMemberRoleRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"required,oneof=admin co_leader member"`
}

type UpdateRoomSettingsRequest struct {
	AllowGuest     *bool `json:"allow_guest,omitempty"`
	AutoApprove    *bool `json:"auto_approve,omitempty"`
	EnableVoting   *bool `json:"enable_voting,omitempty"`
	EnableChecklist *bool `json:"enable_checklist,omitempty"`
}

type RoomQueryParams struct {
	Cursor string `query:"cursor"`
	Limit  int    `query:"limit"`
	Status string `query:"status"`
	Region string `query:"region"`
	Search string `query:"q"`
}
