package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "go-road-backend/internal/domain/auth"
)

type authRepository struct {
	db *Database
}

func NewAuthRepository(db *Database) domain.Repository {
	return &authRepository{db: db}
}

func (r *authRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *authRepository) CreateDeviceToken(ctx context.Context, token *domain.DeviceToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *authRepository) DeleteDeviceToken(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.DeviceToken{}, "id = ?", id).Error
}

func (r *authRepository) ListDeviceTokens(ctx context.Context, userID uuid.UUID) ([]domain.DeviceToken, error) {
	var tokens []domain.DeviceToken
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true", userID).Find(&tokens).Error
	return tokens, err
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *authRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ? AND is_revoked = false AND expires_at > NOW()", tokenHash).First(&token).Error
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}
	return &token, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_revoked": true,
		"revoked_at": gorm.Expr("NOW()"),
	}).Error
}

func (r *authRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Where("user_id = ? AND is_revoked = false", userID).Updates(map[string]interface{}{
		"is_revoked": true,
		"revoked_at": gorm.Expr("NOW()"),
	}).Error
}
