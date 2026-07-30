package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/social"
)

type socialRepository struct {
	db *Database
}

func NewSocialRepository(db *Database) domain.Repository {
	return &socialRepository{db: db}
}

func (r *socialRepository) CreatePost(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *socialRepository) FindPostByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	var post domain.Post
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *socialRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *socialRepository) SoftDeletePost(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Post{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *socialRepository) ListFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Post, string, bool, error) {
	var posts []domain.Post
	db := r.db.WithContext(ctx).
		Joins("JOIN follows ON follows.following_id = posts.author_id").
		Where("follows.follower_id = ? AND posts.deleted_at IS NULL", userID)

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				CreatedAt time.Time `json:"created_at"`
				ID        string    `json:"id"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(posts.created_at, posts.id) < (?, ?)", c.CreatedAt, c.ID)
			}
		}
	}

	db = db.Order("posts.created_at DESC, posts.id DESC").Limit(limit + 1)
	if err := db.Find(&posts).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}

	var nextCursor string
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		cursorData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return posts, nextCursor, hasMore, nil
}

func (r *socialRepository) ListExplorePosts(ctx context.Context, cursor string, limit int) ([]domain.Post, string, bool, error) {
	var posts []domain.Post
	db := r.db.WithContext(ctx).Where("is_public = true AND deleted_at IS NULL")

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				CreatedAt time.Time `json:"created_at"`
				ID        string    `json:"id"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
			}
		}
	}

	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&posts).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}

	var nextCursor string
	if len(posts) > 0 {
		last := posts[len(posts)-1]
		cursorData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return posts, nextCursor, hasMore, nil
}

func (r *socialRepository) CreateLike(ctx context.Context, like *domain.PostLike) error {
	return r.db.WithContext(ctx).Create(like).Error
}

func (r *socialRepository) DeleteLike(ctx context.Context, postID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).Delete(&domain.PostLike{}).Error
}

func (r *socialRepository) CreateComment(ctx context.Context, comment *domain.PostComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *socialRepository) FindCommentByID(ctx context.Context, id uuid.UUID) (*domain.PostComment, error) {
	var comment domain.PostComment
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *socialRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.PostComment{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *socialRepository) ListComments(ctx context.Context, postID uuid.UUID, cursor string, limit int) ([]domain.PostComment, string, bool, error) {
	var comments []domain.PostComment
	db := r.db.WithContext(ctx).Where("post_id = ? AND deleted_at IS NULL", postID)

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				CreatedAt time.Time `json:"created_at"`
				ID        string    `json:"id"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
			}
		}
	}

	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&comments).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(comments) > limit
	if hasMore {
		comments = comments[:limit]
	}

	var nextCursor string
	if len(comments) > 0 {
		last := comments[len(comments)-1]
		cursorData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return comments, nextCursor, hasMore, nil
}

func (r *socialRepository) CreateFollow(ctx context.Context, follow *domain.Follow) error {
	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *socialRepository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&domain.Follow{}).Error
}

func (r *socialRepository) GetFollow(ctx context.Context, followerID, followingID uuid.UUID) (*domain.Follow, error) {
	var follow domain.Follow
	err := r.db.WithContext(ctx).Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&follow).Error
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

func (r *socialRepository) ListFollowers(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Follow, string, bool, error) {
	var follows []domain.Follow
	db := r.db.WithContext(ctx).Where("following_id = ?", userID)

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				CreatedAt time.Time `json:"created_at"`
				ID        string    `json:"id"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
			}
		}
	}

	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&follows).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(follows) > limit
	if hasMore {
		follows = follows[:limit]
	}

	var nextCursor string
	if len(follows) > 0 {
		last := follows[len(follows)-1]
		cursorData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return follows, nextCursor, hasMore, nil
}

func (r *socialRepository) ListFollowing(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Follow, string, bool, error) {
	var follows []domain.Follow
	db := r.db.WithContext(ctx).Where("follower_id = ?", userID)

	if cursor != "" {
		data, err := base64.URLEncoding.DecodeString(cursor)
		if err == nil {
			var c struct {
				CreatedAt time.Time `json:"created_at"`
				ID        string    `json:"id"`
			}
			if err := json.Unmarshal(data, &c); err == nil {
				db = db.Where("(created_at, id) < (?, ?)", c.CreatedAt, c.ID)
			}
		}
	}

	db = db.Order("created_at DESC, id DESC").Limit(limit + 1)
	if err := db.Find(&follows).Error; err != nil {
		return nil, "", false, err
	}

	hasMore := len(follows) > limit
	if hasMore {
		follows = follows[:limit]
	}

	var nextCursor string
	if len(follows) > 0 {
		last := follows[len(follows)-1]
		cursorData, _ := json.Marshal(struct {
			CreatedAt time.Time `json:"created_at"`
			ID        string    `json:"id"`
		}{CreatedAt: last.CreatedAt, ID: last.ID.String()})
		nextCursor = base64.URLEncoding.EncodeToString(cursorData)
	}

	return follows, nextCursor, hasMore, nil
}

func (r *socialRepository) CreateBlock(ctx context.Context, block *domain.Block) error {
	return r.db.WithContext(ctx).Create(block).Error
}

func (r *socialRepository) GetBlock(ctx context.Context, blockerID, blockedID uuid.UUID) (*domain.Block, error) {
	var block domain.Block
	err := r.db.WithContext(ctx).Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *socialRepository) CreateReport(ctx context.Context, report *domain.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *socialRepository) UpdateReport(ctx context.Context, id uuid.UUID, status string, reviewedBy uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Report{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"reviewed_by":  reviewedBy,
		"reviewed_at":  gorm.Expr("NOW()"),
	}).Error
}
