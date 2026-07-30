package social

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	CreatePost(ctx context.Context, caption string, photos []string, authorID uuid.UUID) (*Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*Post, error)
	DeletePost(ctx context.Context, id, userID uuid.UUID) error

	LikePost(ctx context.Context, postID, userID uuid.UUID) error
	UnlikePost(ctx context.Context, postID, userID uuid.UUID) error

	CreateComment(ctx context.Context, postID, authorID uuid.UUID, content string, replyToID *uuid.UUID) (*PostComment, error)
	DeleteComment(ctx context.Context, id, userID uuid.UUID) error
	ListComments(ctx context.Context, postID uuid.UUID, cursor string, limit int) ([]PostComment, string, bool, error)

	FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	GetFollowers(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Follow, string, bool, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Follow, string, bool, error)

	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error

	GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]Post, string, bool, error)
	GetExploreFeed(ctx context.Context, cursor string, limit int) ([]Post, string, bool, error)

	ReportContent(ctx context.Context, reporterID uuid.UUID, reportedType string, reportedID uuid.UUID, reason, description string) error
}
