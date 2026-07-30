package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/social"
	"go-road-backend/internal/repository/redis"
)

type socialService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewSocialService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &socialService{repo: repo, cache: cache, logger: logger}
}

func (s *socialService) CreatePost(ctx context.Context, caption string, photos []string, authorID uuid.UUID) (*domain.Post, error) {
	post := &domain.Post{
		AuthorID:  authorID,
		Caption:   caption,
		Photos:    photos,
		IsPublic:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	s.invalidateFeedCache(ctx, authorID)

	return post, nil
}

func (s *socialService) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	post, err := s.repo.FindPostByID(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *socialService) DeletePost(ctx context.Context, id, userID uuid.UUID) error {
	post, err := s.repo.FindPostByID(ctx, id)
	if err != nil {
		return errors.New("post not found")
	}
	if post.AuthorID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.SoftDeletePost(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	s.invalidateFeedCache(ctx, userID)
	return nil
}

func (s *socialService) LikePost(ctx context.Context, postID, userID uuid.UUID) error {
	post, err := s.repo.FindPostByID(ctx, postID)
	if err != nil {
		return errors.New("post not found")
	}

	block, _ := s.repo.GetBlock(ctx, post.AuthorID, userID)
	if block != nil {
		return errors.New("blocked from interacting")
	}

	like := &domain.PostLike{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateLike(ctx, like); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	post.LikesCount++
	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return fmt.Errorf("failed to update like count: %w", err)
	}

	cacheKey := fmt.Sprintf("cache:feed:user:%s:page:", post.AuthorID.String())
	s.cache.Delete(ctx, cacheKey)

	return nil
}

func (s *socialService) UnlikePost(ctx context.Context, postID, userID uuid.UUID) error {
	post, err := s.repo.FindPostByID(ctx, postID)
	if err != nil {
		return errors.New("post not found")
	}

	if err := s.repo.DeleteLike(ctx, postID, userID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	if post.LikesCount > 0 {
		post.LikesCount--
		s.repo.UpdatePost(ctx, post)
	}

	return nil
}

func (s *socialService) CreateComment(ctx context.Context, postID, authorID uuid.UUID, content string, replyToID *uuid.UUID) (*domain.PostComment, error) {
	comment := &domain.PostComment{
		PostID:    postID,
		AuthorID:  authorID,
		Content:   content,
		ReplyToID: replyToID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	post, err := s.repo.FindPostByID(ctx, postID)
	if err == nil {
		post.CommentsCount++
		s.repo.UpdatePost(ctx, post)
	}

	return comment, nil
}

func (s *socialService) DeleteComment(ctx context.Context, id, userID uuid.UUID) error {
	comment, err := s.repo.FindCommentByID(ctx, id)
	if err != nil {
		return errors.New("comment not found")
	}
	if comment.AuthorID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteComment(ctx, id)
}

func (s *socialService) ListComments(ctx context.Context, postID uuid.UUID, cursor string, limit int) ([]domain.PostComment, string, bool, error) {
	return s.repo.ListComments(ctx, postID, cursor, limit)
}

func (s *socialService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	block, _ := s.repo.GetBlock(ctx, followingID, followerID)
	if block != nil {
		return errors.New("blocked from following")
	}

	follow := &domain.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}
	return s.repo.CreateFollow(ctx, follow)
}

func (s *socialService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.repo.DeleteFollow(ctx, followerID, followingID)
}

func (s *socialService) GetFollowers(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Follow, string, bool, error) {
	return s.repo.ListFollowers(ctx, userID, cursor, limit)
}

func (s *socialService) GetFollowing(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Follow, string, bool, error) {
	return s.repo.ListFollowing(ctx, userID, cursor, limit)
}

func (s *socialService) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	if blockerID == blockedID {
		return errors.New("cannot block yourself")
	}

	block := &domain.Block{
		BlockerID: blockerID,
		BlockedID: blockedID,
	}

	if err := s.repo.CreateBlock(ctx, block); err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	s.repo.DeleteFollow(ctx, blockerID, blockedID)
	s.repo.DeleteFollow(ctx, blockedID, blockerID)

	return nil
}

func (s *socialService) GetFeed(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.Post, string, bool, error) {
	cacheKey := fmt.Sprintf("cache:feed:user:%s:page:%s", userID.String(), cursor)
	var cached struct {
		Posts   []domain.Post `json:"posts"`
		Cursor  string        `json:"cursor"`
		HasMore bool          `json:"has_more"`
	}
	if cursor == "" {
		if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
			return cached.Posts, cached.Cursor, cached.HasMore, nil
		}
	}

	posts, nextCursor, hasMore, err := s.repo.ListFeed(ctx, userID, cursor, limit)
	if err != nil {
		return nil, "", false, err
	}

	if cursor == "" {
		s.cache.SetJSON(ctx, cacheKey, map[string]interface{}{
			"posts":    posts,
			"cursor":   nextCursor,
			"has_more": hasMore,
		}, 5*time.Minute)
	}

	return posts, nextCursor, hasMore, nil
}

func (s *socialService) GetExploreFeed(ctx context.Context, cursor string, limit int) ([]domain.Post, string, bool, error) {
	cacheKey := fmt.Sprintf("cache:feed:explore:page:%s", cursor)
	var cached struct {
		Posts   []domain.Post `json:"posts"`
		Cursor  string        `json:"cursor"`
		HasMore bool          `json:"has_more"`
	}
	if cursor == "" {
		if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
			return cached.Posts, cached.Cursor, cached.HasMore, nil
		}
	}

	posts, nextCursor, hasMore, err := s.repo.ListExplorePosts(ctx, cursor, limit)
	if err != nil {
		return nil, "", false, err
	}

	if cursor == "" {
		s.cache.SetJSON(ctx, cacheKey, map[string]interface{}{
			"posts":    posts,
			"cursor":   nextCursor,
			"has_more": hasMore,
		}, 5*time.Minute)
	}

	return posts, nextCursor, hasMore, nil
}

func (s *socialService) ReportContent(ctx context.Context, reporterID uuid.UUID, reportedType string, reportedID uuid.UUID, reason, description string) error {
	report := &domain.Report{
		ReporterID:  reporterID,
		ReportedType: reportedType,
		ReportedID:  reportedID,
		Reason:      reason,
		Description: description,
		Status:      "pending",
	}
	return s.repo.CreateReport(ctx, report)
}

func (s *socialService) invalidateFeedCache(ctx context.Context, userID uuid.UUID) {
	data, _ := json.Marshal(map[string]string{"user_id": userID.String(), "action": "invalidate_feed"})
	s.cache.Publish(ctx, "social:feed:invalidate", string(data))
}
