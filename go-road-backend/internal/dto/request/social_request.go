package request

type CreatePostRequest struct {
	Content string   `json:"content" validate:"required,max=5000"`
	Images  []string `json:"images,omitempty"`
	RoomID  string   `json:"room_id,omitempty"`
	IsPublic bool    `json:"is_public,omitempty"`
}

type UpdatePostRequest struct {
	Content string `json:"content,omitempty" validate:"max=5000"`
}

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
}

type FollowRequest struct {
	FollowUserID string `json:"follow_user_id" validate:"required"`
}

type ReportRequest struct {
	ReportedID string `json:"reported_id" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=user post comment message"`
	Reason     string `json:"reason" validate:"required,max=500"`
}
