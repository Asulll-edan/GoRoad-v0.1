package room

import (
	"time"

	"github.com/google/uuid"
)

type TouringRoom struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string     `json:"name" gorm:"not null"`
	Description     string     `json:"description,omitempty"`
	CoverPhotoURL   string     `json:"cover_photo_url,omitempty"`
	Status          string     `json:"status" gorm:"default:planning"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	GatheringPoint  string     `json:"gathering_point,omitempty" gorm:"type:geography(point,4326)"`
	GatheringAddress string    `json:"gathering_address,omitempty"`
	GatheringTime   *time.Time `json:"gathering_time,omitempty"`
	StartLocation   string     `json:"start_location,omitempty"`
	EndLocation     string     `json:"end_location,omitempty"`
	MaxMembers      int        `json:"max_members" gorm:"default:20"`
	IsPublic        bool       `json:"is_public" gorm:"default:true"`
	RequiresApproval bool      `json:"requires_approval" gorm:"default:false"`
	TouringType     string     `json:"touring_type" gorm:"default:fun_tour"`
	Difficulty      string     `json:"difficulty" gorm:"default:easy"`
	DistanceKm      float64    `json:"distance_km,omitempty"`
	CreatedBy       uuid.UUID  `json:"created_by" gorm:"not null"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type RoomMember struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID             uuid.UUID `json:"room_id" gorm:"not null;uniqueIndex:idx_room_member"`
	UserID             uuid.UUID `json:"user_id" gorm:"not null;uniqueIndex:idx_room_member"`
	Role               string    `json:"role" gorm:"default:member"`
	PositionInFormation int      `json:"position_in_formation" gorm:"default:0"`
	JoinedAt           time.Time `json:"joined_at"`
	IsOnline           bool      `json:"is_online" gorm:"default:false"`
}

type RoomRoleHistory struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID    uuid.UUID  `json:"room_id" gorm:"not null"`
	UserID    uuid.UUID  `json:"user_id" gorm:"not null"`
	OldRole   string     `json:"old_role,omitempty"`
	NewRole   string     `json:"new_role" gorm:"not null"`
	ChangedBy uuid.UUID  `json:"changed_by" gorm:"not null"`
	Reason    string     `json:"reason,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateRoomRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=100"`
	Description     string `json:"description,omitempty"`
	StartDate       string `json:"start_date,omitempty"`
	EndDate         string `json:"end_date,omitempty"`
	StartLocation   string `json:"start_location,omitempty"`
	EndLocation     string `json:"end_location,omitempty"`
	MaxMembers      int    `json:"max_members,omitempty"`
	IsPublic        bool   `json:"is_public,omitempty"`
	TouringType     string `json:"touring_type,omitempty"`
	Difficulty      string `json:"difficulty,omitempty"`
}

type RoomResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	StartDate     *string   `json:"start_date,omitempty"`
	EndDate       *string   `json:"end_date,omitempty"`
	StartLocation string    `json:"start_location,omitempty"`
	EndLocation   string    `json:"end_location,omitempty"`
	DistanceKm    float64   `json:"distance_km,omitempty"`
	MemberCount   int       `json:"member_count"`
	IsPublic      bool      `json:"is_public"`
	TouringType   string    `json:"touring_type"`
	Difficulty    string    `json:"difficulty"`
	CreatedBy     uuid.UUID `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}
