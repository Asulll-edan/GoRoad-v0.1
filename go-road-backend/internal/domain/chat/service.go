package chat

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	SendMessage(ctx context.Context, req SendMessageRequest, senderID uuid.UUID) (*Message, error)
	GetMessage(ctx context.Context, id uuid.UUID) (*Message, error)
	EditMessage(ctx context.Context, id uuid.UUID, content string, userID uuid.UUID) error
	DeleteMessage(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	PinMessage(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ListMessages(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]Message, string, bool, error)
	MarkRead(ctx context.Context, msgID, userID uuid.UUID) error
}
