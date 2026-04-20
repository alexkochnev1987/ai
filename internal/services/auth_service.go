package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"jsonplaceholder-api/internal/config"
	"jsonplaceholder-api/internal/models"
	"jsonplaceholder-api/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthServiceInterface defines the contract for auth service
type AuthServiceInterface interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID int) error
	ValidateToken(tokenString string) (*models.JWTClaims, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

// AuthService implements AuthServiceInterface
type AuthService struct {
	config         *config.Config
	userRepo       repositories.UserRepositoryInterface
	authRepo       repositories.AuthRepositoryInterface
}

// NewAuthService creates a new auth service
func NewAuthService(
	config *config.Config,
	userRepo repositories.UserRepositoryInterface,
	authRepo repositories.AuthRepositoryInterface,
) AuthServiceInterface {
	return &AuthService{
		config:   config,
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Check if email already exists
	if exists, err := s.userRepo.EmailExists(ctx, req.Email); err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	if exists, err := s.userRepo.UsernameExists(ctx, req.Username); err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash the password
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Address: models.Address{
			Street:  req.Address.Street,
			Suite:   req.Address.Suite,
			City:    req.Address.City,
			Zipcode: req.Address.Zipcode,
			Geo: models.Geo{
				Lat: req.Address.Geo.Lat,
				Lng: req.Address.Geo.Lng,
			},
		},
		Phone:   req.Phone,
		Website: req.Website,
		Company: models.Company{
			Name:        req.Company.Name,
			CatchPhrase: req.Company.CatchPhrase,
			BS:          req.Company.BS,
		},
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Save refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: tokenPair.ExpiresAt,
	}

	if err := s.authRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.LoginResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         *user,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := s.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Save refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: tokenPair.ExpiresAt,
	}

	if err := s.authRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.LoginResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User:         *user,
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.TokenPair, error) {
	// Get refresh token from database
	refreshToken, err := s.authRepo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Generate new tokens
	tokenPair, err := s.generateTokenPair(&refreshToken.User)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Revoke old refresh token
	if err := s.authRepo.RevokeRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Save new refresh token
	newRefreshToken := &models.RefreshToken{
		UserID:    refreshToken.UserID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: tokenPair.ExpiresAt,
	}

	if err := s.authRepo.CreateRefreshToken(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return tokenPair, nil
}

// Logout revokes a specific refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.authRepo.RevokeRefreshToken(ctx, refreshToken)
}

// LogoutAll revokes all refresh tokens for a user
func (s *AuthService) LogoutAll(ctx context.Context, userID int) error {
	return s.authRepo.RevokeAllUserTokens(ctx, userID)
}

// ValidateToken validates and parses a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with a plain password
func (s *AuthService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// generateTokenPair generates an access token and refresh token pair
func (s *AuthService) generateTokenPair(user *models.User) (*models.TokenPair, error) {
	// Generate access token
	accessClaims := &models.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Issuer,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshTokenString, err := s.generateRandomToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshExpiresAt := time.Now().Add(s.config.JWT.RefreshTokenDuration)

	return &models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    refreshExpiresAt,
	}, nil
}

// generateRandomToken generates a random token for refresh tokens
func (s *AuthService) generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
} 