package models

import (
	"time"
)

// Standard API Response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Paginated response
type PaginatedResponse struct {
	APIResponse
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// Helper functions for creating responses
func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

func NewErrorResponse(message, error string) APIResponse {
	return APIResponse{
		Success:   false,
		Message:   message,
		Error:     error,
		Timestamp: time.Now(),
	}
} 