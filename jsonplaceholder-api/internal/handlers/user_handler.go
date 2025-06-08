package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"jsonplaceholder-api/internal/models"
	"jsonplaceholder-api/internal/services"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get list of all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.UserResponse}
// @Failure 500 {object} models.APIResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.GetUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get users", err.Error()))
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	response := models.PaginatedResponse{
		APIResponse: models.APIResponse{
			Success: true,
			Message: "Users retrieved successfully",
			Data:    users,
		},
		Pagination: models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid user ID", "ID must be a positive integer"))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("User not found", ""))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to get user", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("User retrieved successfully", user))
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateUserRequest true "User data"
// @Success 201 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, models.NewErrorResponse("User already exists", ""))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to create user", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse("User created successfully", user))
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "User data"
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid user ID", "ID must be a positive integer"))
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("User not found", ""))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to update user", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("User updated successfully", user))
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid user ID", "ID must be a positive integer"))
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("User not found", ""))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to delete user", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("User deleted successfully", nil))
} 