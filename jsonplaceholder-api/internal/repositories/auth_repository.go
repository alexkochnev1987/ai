package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"jsonplaceholder-api/internal/models"
)

type AuthRepositoryInterface interface {
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteExpiredTokens(ctx context.Context) error
	DeleteUserTokens(ctx context.Context, userID uint) error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepositoryInterface {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid refresh token")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return &refreshToken, nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	result := r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete refresh token: %w", result.Error)
	}
	return nil
}

func (r *AuthRepository) DeleteExpiredTokens(ctx context.Context) error {
	result := r.db.WithContext(ctx).Where("expires_at <= ?", time.Now()).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", result.Error)
	}
	return nil
}

func (r *AuthRepository) DeleteUserTokens(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user tokens: %w", result.Error)
	}
	return nil
} 