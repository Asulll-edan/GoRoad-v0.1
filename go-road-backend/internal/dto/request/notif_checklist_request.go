package request

type NotificationPreferenceRequest struct {
	PushEnabled  *bool `json:"push_enabled,omitempty"`
	ChatEnabled  *bool `json:"chat_enabled,omitempty"`
	EmailEnabled *bool `json:"email_enabled,omitempty"`
}

type CreateChecklistTemplateRequest struct {
	Title       string   `json:"title" validate:"required,max=200"`
	Category    string   `json:"category" validate:"oneof=gear safety document vehicle"`
	Items       []string `json:"items" validate:"required,min=1"`
}

type CreateTouringChecklistRequest struct {
	TemplateID string `json:"template_id" validate:"required"`
	RoomID     string `json:"room_id" validate:"required"`
}
