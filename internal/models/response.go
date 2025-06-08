package models

import "time"

// APIResponse represents a standardized API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError represents an error response
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{}      `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page         int   `json:"page"`
	PerPage      int   `json:"per_page"`
	Total        int64 `json:"total"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Page    int `json:"page" form:"page" validate:"min=1"`
	PerPage int `json:"per_page" form:"per_page" validate:"min=1,max=100"`
}

// DefaultPagination returns default pagination parameters
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:    1,
		PerPage: 10,
	}
}

// Offset calculates the offset for database queries
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// CalculateMeta calculates pagination metadata
func (p PaginationParams) CalculateMeta(total int64) PaginationMeta {
	totalPages := int((total + int64(p.PerPage) - 1) / int64(p.PerPage))
	
	return PaginationMeta{
		Page:        p.Page,
		PerPage:     p.PerPage,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     p.Page < totalPages,
		HasPrevious: p.Page > 1,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
} 