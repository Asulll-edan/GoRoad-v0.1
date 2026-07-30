package room

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, room *TouringRoom) error
	FindByID(ctx context.Context, id uuid.UUID) (*TouringRoom, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]TouringRoom, error)
	Update(ctx context.Context, room *TouringRoom) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, cursor string, limit int, filters map[string]interface{}) ([]TouringRoom, string, bool, error)
	Discover(ctx context.Context, cursor string, limit int) ([]TouringRoom, string, bool, error)

	AddMember(ctx context.Context, member *RoomMember) error
	RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, roomID, userID uuid.UUID, role string) error
	ListMembers(ctx context.Context, roomID uuid.UUID, cursor string, limit int) ([]RoomMember, string, bool, error)
	GetMember(ctx context.Context, roomID, userID uuid.UUID) (*RoomMember, error)
	CountMembers(ctx context.Context, roomID uuid.UUID) (int64, error)

	AddRoleHistory(ctx context.Context, history *RoomRoleHistory) error
}
