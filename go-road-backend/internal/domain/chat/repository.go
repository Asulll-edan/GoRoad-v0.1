package chat

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, msg *Message) error
	FindByID(ctx context.Context, id uuid.UUID) (*Message, error)
	FindByRoomID(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]Message, string, bool, error)
	Update(ctx context.Context, msg *Message) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Pin(ctx context.Context, id uuid.UUID, pinned bool) error
	MarkRead(ctx context.Context, msgID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, roomID, userID uuid.UUID) (int64, error)
}
