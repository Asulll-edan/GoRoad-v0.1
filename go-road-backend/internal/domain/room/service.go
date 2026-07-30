package room

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateRoomRequest, userID uuid.UUID) (*RoomResponse, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID, includes []string) (*TouringRoom, error)
	Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	List(ctx context.Context, cursor string, limit int, userID uuid.UUID) ([]RoomResponse, string, bool, error)
	Discover(ctx context.Context, cursor string, limit int) ([]RoomResponse, string, bool, error)

	Join(ctx context.Context, roomID, userID uuid.UUID) error
	Leave(ctx context.Context, roomID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, roomID, userID, targetUserID uuid.UUID, role string) error
	ListMembers(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]RoomMember, string, bool, error)

	StartTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	PauseTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	CompleteTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
	CancelTouring(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error
}
