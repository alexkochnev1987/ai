package services

import (
	"context"
	"fmt"

	"jsonplaceholder-api/internal/models"
	"jsonplaceholder-api/internal/repositories"
)

// UserServiceInterface defines the contract for user service
type UserServiceInterface interface {
	GetAll(ctx context.Context, pagination models.PaginationParams) (*models.PaginatedResponse, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	Update(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id int) error
}

// UserService implements UserServiceInterface
type UserService struct {
	userRepo repositories.UserRepositoryInterface
	authService AuthServiceInterface
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repositories.UserRepositoryInterface,
	authService AuthServiceInterface,
) UserServiceInterface {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
	}
}

// GetAll retrieves all users with pagination
func (s *UserService) GetAll(ctx context.Context, pagination models.PaginationParams) (*models.PaginatedResponse, error) {
	users, total, err := s.userRepo.GetAll(ctx, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	meta := pagination.CalculateMeta(total)

	return &models.PaginatedResponse{
		Data:       users,
		Pagination: meta,
	}, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
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
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
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

	return user, nil
}

// Update updates an existing user
func (s *UserService) Update(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Prepare updates map
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Username != nil {
		// Check if username already exists for a different user
		if exists, err := s.userRepo.UsernameExists(ctx, *req.Username); err != nil {
			return nil, fmt.Errorf("failed to check username existence: %w", err)
		} else if exists {
			// Check if it's not the same user
			existingUser, err := s.userRepo.GetByUsername(ctx, *req.Username)
			if err == nil && existingUser.ID != id {
				return nil, fmt.Errorf("username already exists")
			}
		}
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		// Check if email already exists for a different user
		if exists, err := s.userRepo.EmailExists(ctx, *req.Email); err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		} else if exists {
			// Check if it's not the same user
			existingUser, err := s.userRepo.GetByEmail(ctx, *req.Email)
			if err == nil && existingUser.ID != id {
				return nil, fmt.Errorf("email already exists")
			}
		}
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}

	// Handle address updates
	if req.Address != nil {
		updates["address_street"] = req.Address.Street
		updates["address_suite"] = req.Address.Suite
		updates["address_city"] = req.Address.City
		updates["address_zipcode"] = req.Address.Zipcode
		updates["address_geo_lat"] = req.Address.Geo.Lat
		updates["address_geo_lng"] = req.Address.Geo.Lng
	}

	// Handle company updates
	if req.Company != nil {
		updates["company_name"] = req.Company.Name
		updates["company_catch_phrase"] = req.Company.CatchPhrase
		updates["company_bs"] = req.Company.BS
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	user, err := s.userRepo.Update(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete deletes a user by ID
func (s *UserService) Delete(ctx context.Context, id int) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
} 