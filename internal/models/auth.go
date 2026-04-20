package models

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         User      `json:"user"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Name     string         `json:"name" validate:"required,min=2,max=100"`
	Username string         `json:"username" validate:"required,min=2,max=50"`
	Email    string         `json:"email" validate:"required,email"`
	Password string         `json:"password" validate:"required,min=8"`
	Address  AddressRequest `json:"address"`
	Phone    string         `json:"phone" validate:"required"`
	Website  string         `json:"website"`
	Company  CompanyRequest `json:"company"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        int       `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" db:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" db:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	IsRevoked bool      `json:"is_revoked" db:"is_revoked" gorm:"default:false"`
	
	// Relations
	User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TableName returns the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
} 