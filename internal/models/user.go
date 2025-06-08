package models

import (
	"time"
)

// User represents a user in the system matching JSONPlaceholder structure
type User struct {
	ID       int     `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Name     string  `json:"name" db:"name" validate:"required,min=2,max=100"`
	Username string  `json:"username" db:"username" validate:"required,min=2,max=50" gorm:"uniqueIndex"`
	Email    string  `json:"email" db:"email" validate:"required,email" gorm:"uniqueIndex"`
	Address  Address `json:"address" db:"address" gorm:"embedded;embeddedPrefix:address_"`
	Phone    string  `json:"phone" db:"phone" validate:"required"`
	Website  string  `json:"website" db:"website"`
	Company  Company `json:"company" db:"company" gorm:"embedded;embeddedPrefix:company_"`
	
	// Authentication fields (not exposed in JSON)
	PasswordHash string    `json:"-" db:"password_hash" gorm:"not null"`
	CreatedAt    time.Time `json:"-" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"-" db:"updated_at" gorm:"autoUpdateTime"`
}

// Address represents the nested address structure
type Address struct {
	Street  string `json:"street" db:"street"`
	Suite   string `json:"suite" db:"suite"`
	City    string `json:"city" db:"city"`
	Zipcode string `json:"zipcode" db:"zipcode"`
	Geo     Geo    `json:"geo" db:"geo" gorm:"embedded;embeddedPrefix:geo_"`
}

// Geo represents geographical coordinates
type Geo struct {
	Lat string `json:"lat" db:"lat"`
	Lng string `json:"lng" db:"lng"`
}

// Company represents the company information
type Company struct {
	Name        string `json:"name" db:"name"`
	CatchPhrase string `json:"catchPhrase" db:"catch_phrase"`
	BS          string `json:"bs" db:"bs"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name     string            `json:"name" validate:"required,min=2,max=100"`
	Username string            `json:"username" validate:"required,min=2,max=50"`
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=8"`
	Address  AddressRequest    `json:"address"`
	Phone    string            `json:"phone" validate:"required"`
	Website  string            `json:"website"`
	Company  CompanyRequest    `json:"company"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name     *string            `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Username *string            `json:"username,omitempty" validate:"omitempty,min=2,max=50"`
	Email    *string            `json:"email,omitempty" validate:"omitempty,email"`
	Address  *AddressRequest    `json:"address,omitempty"`
	Phone    *string            `json:"phone,omitempty"`
	Website  *string            `json:"website,omitempty"`
	Company  *CompanyRequest    `json:"company,omitempty"`
}

// AddressRequest represents the address in request payloads
type AddressRequest struct {
	Street  string     `json:"street"`
	Suite   string     `json:"suite"`
	City    string     `json:"city"`
	Zipcode string     `json:"zipcode"`
	Geo     GeoRequest `json:"geo"`
}

// GeoRequest represents geographical coordinates in requests
type GeoRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// CompanyRequest represents company information in requests
type CompanyRequest struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchPhrase"`
	BS          string `json:"bs"`
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
} 