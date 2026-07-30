package auth

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error

	RegisterDevice(ctx context.Context, userID uuid.UUID, token, platform string) error
	UnregisterDevice(ctx context.Context, deviceID uuid.UUID) error
}
