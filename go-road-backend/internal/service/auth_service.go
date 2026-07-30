package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"go-road-backend/internal/config"
	domain "go-road-backend/internal/domain/auth"
	"go-road-backend/internal/repository/redis"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already registered")
	ErrUsernameAlreadyExists = errors.New("username already taken")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidRefreshToken   = errors.New("invalid or expired refresh token")
)

type authService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	cfg    *config.Config
	logger *zap.Logger
}

func NewAuthService(repo domain.Repository, cache redis.CacheRepository, cfg *config.Config, logger *zap.Logger) domain.Service {
	return &authService{
		repo:   repo,
		cache:  cache,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	existing, _ := s.repo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	existing, _ = s.repo.FindByUsername(ctx, req.Username)
	if existing != nil {
		return nil, ErrUsernameAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:        req.Email,
		Username:     req.Username,
		FullName:     req.FullName,
		PasswordHash: string(hash),
		Phone:        req.Phone,
		RidingSkill:  "beginner",
		RidingRole:   "member",
		IsActive:     true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	storedToken, err := s.repo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if storedToken.IsRevoked || storedToken.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.repo.FindByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	s.repo.RevokeRefreshToken(ctx, storedToken.ID)

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	cacheKey := fmt.Sprintf("cache:user:%s:public_profile", userID.String())

	var cached domain.UserResponse
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	resp := toUserResponse(user)
	s.cache.SetJSON(ctx, cacheKey, resp, 15*time.Minute)

	return resp, nil
}

func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*domain.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if v, ok := updates["full_name"]; ok {
		user.FullName = v.(string)
	}
	if v, ok := updates["bio"]; ok {
		user.Bio = v.(string)
	}
	if v, ok := updates["riding_skill"]; ok {
		user.RidingSkill = v.(string)
	}
	if v, ok := updates["riding_role"]; ok {
		user.RidingRole = v.(string)
	}
	if v, ok := updates["phone"]; ok {
		user.Phone = v.(string)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	cacheKeys := []string{
		fmt.Sprintf("cache:user:%s:public_profile", userID.String()),
		fmt.Sprintf("cache:user:%s", userID.String()),
	}
	s.cache.Delete(ctx, cacheKeys...)

	return toUserResponse(user), nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hash)
	return s.repo.Update(ctx, user)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil
	}

	_ = user
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	return nil
}

func (s *authService) RegisterDevice(ctx context.Context, userID uuid.UUID, token, platform string) error {
	device := &domain.DeviceToken{
		UserID:   userID,
		Token:    token,
		Platform: platform,
		IsActive: true,
	}
	return s.repo.CreateDeviceToken(ctx, device)
}

func (s *authService) UnregisterDevice(ctx context.Context, deviceID uuid.UUID) error {
	return s.repo.DeleteDeviceToken(ctx, deviceID)
}

func (s *authService) generateAuthResponse(ctx context.Context, user *domain.User) (*domain.AuthResponse, error) {
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *toUserResponse(user),
	}, nil
}

func (s *authService) generateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"email":    user.Email,
		"username": user.Username,
		"role":     "user",
		"exp":      time.Now().Add(time.Duration(s.cfg.JWTExpiryHour) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *authService) generateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := uuid.NewRandom(); err != nil {
		return "", err
	}

	token := uuid.New().String() + uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.repo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", err
	}

	return token, nil
}

func toUserResponse(user *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		FullName:    user.FullName,
		PhotoURL:    user.PhotoURL,
		RidingSkill: user.RidingSkill,
		RidingRole:  user.RidingRole,
		IsVerified:  user.IsVerified,
	}
}
