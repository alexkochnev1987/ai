package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jsonplaceholder-api/internal/models"
	"jsonplaceholder-api/internal/services"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, models.NewErrorResponse("User already exists", ""))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to register user", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse("User registered successfully", response))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid credentials", ""))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to login", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Login successful", response))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.APIResponse{data=models.TokenPair}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err.Error() == "invalid refresh token" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid refresh token", ""))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to refresh token", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Token refreshed successfully", tokens))
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token and logout user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	err := h.authService.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to logout", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("Logout successful", nil))
}

// Me godoc
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("User not authenticated", ""))
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Invalid user ID", ""))
		return
	}

	userResponse, err := h.authService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get user", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("User retrieved successfully", userResponse))
} 