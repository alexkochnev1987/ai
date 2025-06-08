package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"jsonplaceholder-api/internal/config"
	"jsonplaceholder-api/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CORSMiddleware configures CORS settings
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	config := cors.DefaultConfig()
	
	if cfg.IsDevelopment() {
		config.AllowAllOrigins = true
	} else {
		// In production, specify allowed origins
		config.AllowOrigins = []string{"https://yourdomain.com"}
	}
	
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return cors.New(config)
}

// LoggerMiddleware creates a custom logger middleware
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return log.Printf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		return ""
	})
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INTERNAL_ERROR",
					Message: "Internal server error",
					Details: gin.H{"error": err},
				},
				Timestamp: time.Now(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INTERNAL_ERROR",
					Message: "Internal server error",
				},
				Timestamp: time.Now(),
			})
		}
	})
}

// ValidationMiddleware validates request payload using struct tags
func ValidationMiddleware() *validator.Validate {
	return validator.New()
}

// PaginationMiddleware extracts and validates pagination parameters
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagination := models.DefaultPagination()

		if pageStr := c.Query("page"); pageStr != "" {
			if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
				pagination.Page = page
			}
		}

		if perPageStr := c.Query("per_page"); perPageStr != "" {
			if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 && perPage <= 100 {
				pagination.PerPage = perPage
			}
		}

		c.Set("pagination", pagination)
		c.Next()
	}
}

// GetPaginationFromContext extracts pagination from gin context
func GetPaginationFromContext(c *gin.Context) models.PaginationParams {
	if pagination, exists := c.Get("pagination"); exists {
		if p, ok := pagination.(models.PaginationParams); ok {
			return p
		}
	}
	return models.DefaultPagination()
}

// RequestIDMiddleware adds a request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

// ValidateJSON validates JSON request body
func ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			var validationErrors []models.ValidationError
			
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrs {
					validationErrors = append(validationErrors, models.ValidationError{
						Field:   e.Field(),
						Tag:     e.Tag(),
						Value:   e.Param(),
						Message: getValidationMessage(e),
					})
				}
			}

			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "VALIDATION_ERROR",
					Message: "Invalid request data",
					Details: validationErrors,
				},
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// getValidationMessage returns a human-readable validation error message
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + e.Param() + " characters long"
	case "max":
		return "Must be at most " + e.Param() + " characters long"
	default:
		return "Invalid value"
	}
} 