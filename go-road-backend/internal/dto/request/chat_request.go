package request

type SendMessageRequest struct {
	Message     string `json:"message" validate:"required,max=2000"`
	MessageType string `json:"message_type,omitempty" validate:"oneof=text image location voice"`
	ReplyToID   string `json:"reply_to_id,omitempty"`
}

type PinMessageRequest struct {
	MessageID string `json:"message_id" validate:"required"`
}

type EditMessageRequest struct {
	Message string `json:"message" validate:"required,max=2000"`
}
