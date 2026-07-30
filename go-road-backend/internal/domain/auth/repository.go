package auth

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error

	CreateDeviceToken(ctx context.Context, token *DeviceToken) error
	DeleteDeviceToken(ctx context.Context, id uuid.UUID) error
	ListDeviceTokens(ctx context.Context, userID uuid.UUID) ([]DeviceToken, error)

	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
