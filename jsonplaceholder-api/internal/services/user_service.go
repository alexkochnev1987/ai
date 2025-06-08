package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"jsonplaceholder-api/internal/models"
	"jsonplaceholder-api/internal/repositories"
)

type UserServiceInterface interface {
	GetUsers(ctx context.Context, page, limit int) ([]*models.UserResponse, int64, error)
	GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error)
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUsers(ctx context.Context, page, limit int) ([]*models.UserResponse, int64, error) {
	offset := (page - 1) * limit
	users, total, err := s.userRepo.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	// Convert to response format
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return userResponses, total, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Check if user already exists by email
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
		Address:  req.Address,
		Company:  req.Company,
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ToResponse(), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Username != nil {
		// Check if new username already exists (excluding current user)
		existingUser, err := s.userRepo.GetByUsername(ctx, *req.Username)
		if err == nil && existingUser.ID != id {
			return nil, fmt.Errorf("username already exists")
		}
		user.Username = *req.Username
	}
	if req.Email != nil {
		// Check if new email already exists (excluding current user)
		existingUser, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil && existingUser.ID != id {
			return nil, fmt.Errorf("email already exists")
		}
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Website != nil {
		user.Website = *req.Website
	}
	if req.Address != nil {
		user.Address = *req.Address
	}
	if req.Company != nil {
		user.Company = *req.Company
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user.ToResponse(), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

// Helper methods

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
} 