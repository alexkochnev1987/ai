package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"jsonplaceholder-api/internal/config"
	"jsonplaceholder-api/internal/models"
	"jsonplaceholder-api/internal/repositories"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	GetUserByID(ctx context.Context, userID uint) (*models.UserResponse, error)
}

type AuthService struct {
	userRepo repositories.UserRepositoryInterface
	authRepo repositories.AuthRepositoryInterface
	config   *config.Config
}

func NewAuthService(
	userRepo repositories.UserRepositoryInterface,
	authRepo repositories.AuthRepositoryInterface,
	config *config.Config,
) AuthServiceInterface {
	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
		config:   config,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	// Check if username already exists
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check if username exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Website:  req.Website,
		Password: hashedPassword,
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Save refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpiry),
	}
	if err := s.authRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is disabled")
	}

	// Verify password
	if !s.checkPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Save refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpiry),
	}
	if err := s.authRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	// Get refresh token from database
	storedToken, err := s.authRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Generate new tokens
	tokens, err := s.generateTokens(&storedToken.User)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Delete old refresh token
	if err := s.authRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	// Save new refresh token
	newRefreshToken := &models.RefreshToken{
		UserID:    storedToken.UserID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpiry),
	}
	if err := s.authRepo.CreateRefreshToken(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return tokens, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.authRepo.DeleteRefreshToken(ctx, refreshToken)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

// Helper methods

func (s *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) generateTokens(user *models.User) (*models.TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := &models.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "jsonplaceholder-api",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.SecretKey))
}

func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
} 