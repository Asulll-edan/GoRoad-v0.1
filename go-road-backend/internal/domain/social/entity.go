package social

import (
	"time"

	"github.com/google/uuid"
)

type Badge struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Code        string                 `json:"code" gorm:"uniqueIndex;not null"`
	Name        string                 `json:"name" gorm:"not null"`
	Description string                 `json:"description,omitempty"`
	IconURL     string                 `json:"icon_url,omitempty"`
	Category    string                 `json:"category" gorm:"not null"`
	Tier        string                 `json:"tier" gorm:"not null"`
	Criteria    map[string]interface{} `json:"criteria" gorm:"type:jsonb;not null"`
	IsHidden    bool                   `json:"is_hidden" gorm:"default:false"`
	CreatedAt   time.Time              `json:"created_at"`
}

type UserBadge struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"not null;uniqueIndex:idx_user_badge"`
	BadgeID   uuid.UUID  `json:"badge_id" gorm:"not null;uniqueIndex:idx_user_badge"`
	AwardedAt time.Time  `json:"awarded_at"`
	TouringID *uuid.UUID `json:"touring_id,omitempty"`
}

type Post struct {
	ID            uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID        *uuid.UUID             `json:"room_id,omitempty"`
	AuthorID      uuid.UUID              `json:"author_id" gorm:"not null"`
	Caption       string                 `json:"caption,omitempty"`
	Photos        []string               `json:"photos,omitempty" gorm:"type:text[]"`
	StatsSnapshot map[string]interface{} `json:"stats_snapshot,omitempty" gorm:"type:jsonb;default:'{}'"`
	RouteSnapshot map[string]interface{} `json:"route_snapshot,omitempty" gorm:"type:jsonb;default:'{}'"`
	IsPublic      bool                   `json:"is_public" gorm:"default:true"`
	LikesCount    int                    `json:"likes_count" gorm:"default:0"`
	CommentsCount int                    `json:"comments_count" gorm:"default:0"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	DeletedAt     *time.Time             `json:"deleted_at,omitempty"`
}

type PostLike struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PostID    uuid.UUID `json:"post_id" gorm:"not null;uniqueIndex:idx_post_like"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;uniqueIndex:idx_post_like"`
	CreatedAt time.Time `json:"created_at"`
}

type PostComment struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PostID    uuid.UUID  `json:"post_id" gorm:"not null"`
	AuthorID  uuid.UUID  `json:"author_id" gorm:"not null"`
	Content   string     `json:"content" gorm:"not null"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Follow struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FollowerID  uuid.UUID `json:"follower_id" gorm:"not null;uniqueIndex:idx_follow"`
	FollowingID uuid.UUID `json:"following_id" gorm:"not null;uniqueIndex:idx_follow"`
	CreatedAt   time.Time `json:"created_at"`
}

type Report struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ReporterID   uuid.UUID  `json:"reporter_id" gorm:"not null"`
	ReportedType string     `json:"reported_type" gorm:"not null"`
	ReportedID   uuid.UUID  `json:"reported_id" gorm:"not null"`
	Reason       string     `json:"reason" gorm:"not null"`
	Description  string     `json:"description,omitempty"`
	Status       string     `json:"status" gorm:"not null;default:pending"`
	ReviewedBy   *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Notification struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID              `json:"user_id" gorm:"not null;index"`
	Type      string                 `json:"type" gorm:"not null"`
	Title     string                 `json:"title" gorm:"not null"`
	Body      string                 `json:"body,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty" gorm:"type:jsonb;default:'{}'"`
	IsRead    bool                   `json:"is_read" gorm:"default:false"`
	ReadAt    *time.Time             `json:"read_at,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}
