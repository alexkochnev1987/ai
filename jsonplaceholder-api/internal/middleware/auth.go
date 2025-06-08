package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"jsonplaceholder-api/internal/models"
)

// JWTAuth middleware validates JWT tokens
func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header required", ""))
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid authorization header format", "Expected 'Bearer <token>'"))
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token", err.Error()))
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
			// Set user information in context
			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("username", claims.Username)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token claims", ""))
			c.Abort()
			return
		}
	}
} 