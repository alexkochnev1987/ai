package routes

import (
	"jsonplaceholder-api/internal/config"
	"jsonplaceholder-api/internal/handlers"
	"jsonplaceholder-api/internal/middleware"
	"jsonplaceholder-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Apply global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Security())
	r.Use(middleware.RequestID())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.HealthResponse{
			Status:  "ok",
			Message: "API is running",
			Version: cfg.App.Version,
		})
	})

	// API version 1
	v1 := r.Group("/api/v1")
	{
		// API info endpoint
		v1.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.NewSuccessResponse("JSONPlaceholder API v1", map[string]interface{}{
				"version": cfg.App.Version,
				"endpoints": map[string]interface{}{
					"authentication": map[string]string{
						"POST /api/v1/auth/register": "Register new user",
						"POST /api/v1/auth/login":    "Login user",
						"POST /api/v1/auth/refresh":  "Refresh access token",
						"POST /api/v1/auth/logout":   "Logout user",
						"GET /api/v1/auth/me":        "Get current user (requires auth)",
					},
					"users": map[string]string{
						"GET /api/v1/users":       "Get all users (with pagination)",
						"GET /api/v1/users/{id}":  "Get user by ID",
						"POST /api/v1/users":      "Create new user (requires auth)",
						"PUT /api/v1/users/{id}":  "Update user (requires auth)",
						"DELETE /api/v1/users/{id}": "Delete user (requires auth)",
					},
				},
				"documentation": "See README.md for detailed API usage examples",
			}))
		})
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
			
			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.JWTAuth(cfg.JWT.SecretKey))
			{
				authProtected.GET("/me", authHandler.Me)
			}
		}

		// Users routes
		users := v1.Group("/users")
		{
			// Public routes (read-only like JSONPlaceholder)
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)

			// Protected routes (write operations require authentication)
			usersProtected := users.Group("")
			usersProtected.Use(middleware.JWTAuth(cfg.JWT.SecretKey))
			{
				usersProtected.POST("", userHandler.CreateUser)
				usersProtected.PUT("/:id", userHandler.UpdateUser)
				usersProtected.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("Route not found", "The requested endpoint does not exist"))
	})

	return r
} 