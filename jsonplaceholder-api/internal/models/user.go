package models

import (
	"time"
)

// User model for database (matches JSONPlaceholder structure)
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"not null;size:100"`
	Username string `json:"username" gorm:"not null;unique;size:50"`
	Email    string `json:"email" gorm:"not null;unique;size:100"`
	Phone    string `json:"phone" gorm:"size:50"`
	Website  string `json:"website" gorm:"size:100"`

	// Authentication fields
	Password  string `json:"-" gorm:"not null;size:255"` // Hidden from JSON
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Embedded structs for JSONPlaceholder compatibility
	Address Address `json:"address" gorm:"embedded;embeddedPrefix:address_"`
	Company Company `json:"company" gorm:"embedded;embeddedPrefix:company_"`
}

// Address embedded struct
type Address struct {
	Street  string `json:"street" gorm:"size:100"`
	Suite   string `json:"suite" gorm:"size:50"`
	City    string `json:"city" gorm:"size:50"`
	Zipcode string `json:"zipcode" gorm:"size:20"`
	Geo     Geo    `json:"geo" gorm:"embedded;embeddedPrefix:geo_"`
}

// Geo embedded struct
type Geo struct {
	Lat string `json:"lat" gorm:"size:20"`
	Lng string `json:"lng" gorm:"size:20"`
}

// Company embedded struct
type Company struct {
	Name        string `json:"name" gorm:"size:100"`
	CatchPhrase string `json:"catchPhrase" gorm:"size:200"`
	BS          string `json:"bs" gorm:"size:200"`
}

// Request/Response models
type UserResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	Address  Address `json:"address"`
	Company  Company `json:"company"`
}

type CreateUserRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    string  `json:"phone,omitempty"`
	Website  string  `json:"website,omitempty"`
	Password string  `json:"password" binding:"required,min=6"`
	Address  Address `json:"address"`
	Company  Company `json:"company"`
}

type UpdateUserRequest struct {
	Name     *string  `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Username *string  `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    *string  `json:"email,omitempty" binding:"omitempty,email"`
	Phone    *string  `json:"phone,omitempty"`
	Website  *string  `json:"website,omitempty"`
	Address  *Address `json:"address,omitempty"`
	Company  *Company `json:"company,omitempty"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Website:  u.Website,
		Address:  u.Address,
		Company:  u.Company,
	}
} 