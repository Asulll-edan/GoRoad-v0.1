package social

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreatePost(ctx context.Context, post *Post) error
	FindPostByID(ctx context.Context, id uuid.UUID) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	SoftDeletePost(ctx context.Context, id uuid.UUID) error
	ListFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Post, string, bool, error)
	ListExplorePosts(ctx context.Context, cursor string, limit int) ([]Post, string, bool, error)

	CreateLike(ctx context.Context, like *PostLike) error
	DeleteLike(ctx context.Context, postID, userID uuid.UUID) error

	CreateComment(ctx context.Context, comment *PostComment) error
	FindCommentByID(ctx context.Context, id uuid.UUID) (*PostComment, error)
	DeleteComment(ctx context.Context, id uuid.UUID) error
	ListComments(ctx context.Context, postID uuid.UUID, cursor string, limit int) ([]PostComment, string, bool, error)

	CreateFollow(ctx context.Context, follow *Follow) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	GetFollow(ctx context.Context, followerID, followingID uuid.UUID) (*Follow, error)
	ListFollowers(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Follow, string, bool, error)
	ListFollowing(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Follow, string, bool, error)

	CreateBlock(ctx context.Context, block *Block) error
	GetBlock(ctx context.Context, blockerID, blockedID uuid.UUID) (*Block, error)

	CreateReport(ctx context.Context, report *Report) error
	UpdateReport(ctx context.Context, id uuid.UUID, status string, reviewedBy uuid.UUID) error
}

type Block struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	BlockerID uuid.UUID `json:"blocker_id" gorm:"not null"`
	BlockedID uuid.UUID `json:"blocked_id" gorm:"not null"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
