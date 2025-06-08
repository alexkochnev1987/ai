package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"jsonplaceholder-api/internal/models"

	"gorm.io/gorm"
)

// AuthRepositoryInterface defines the contract for auth repository
type AuthRepositoryInterface interface {
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID int) error
	CleanupExpiredTokens(ctx context.Context) error
}

// AuthRepository implements AuthRepositoryInterface
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository
func NewAuthRepository(db *gorm.DB) AuthRepositoryInterface {
	return &AuthRepository{db: db}
}

// CreateRefreshToken creates a new refresh token in the database
func (r *AuthRepository) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

// GetRefreshToken retrieves a refresh token by token string
func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&refreshToken).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("refresh token not found or expired: %w", err)
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	
	return &refreshToken, nil
}

// RevokeRefreshToken revokes a specific refresh token
func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	result := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true)
	
	if result.Error != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("refresh token not found")
	}
	
	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a specific user
func (r *AuthRepository) RevokeAllUserTokens(ctx context.Context, userID int) error {
	if err := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error; err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}
	
	return nil
}

// CleanupExpiredTokens removes expired refresh tokens from the database
func (r *AuthRepository) CleanupExpiredTokens(ctx context.Context) error {
	if err := r.db.WithContext(ctx).
		Where("expires_at < ? OR is_revoked = ?", time.Now(), true).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	
	return nil
} 