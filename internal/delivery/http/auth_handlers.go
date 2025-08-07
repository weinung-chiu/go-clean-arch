package http

import (
	"go-clean-arch/internal/entity"
	"go-clean-arch/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token     string      `json:"token"`
	User      UserProfile `json:"user"`
	ExpiresAt string      `json:"expires_at"`
}

// UserProfile represents the user profile in API responses
type UserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// Login handles user authentication
func Login(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		// Authenticate user
		result, err := app.AuthService.LoginWithPassword(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return authentication result
		c.JSON(http.StatusOK, AuthResponse{
			Token: result.Token,
			User: UserProfile{
				ID:       result.User.ID,
				Username: result.User.Username,
				Name:     result.User.Name,
			},
			ExpiresAt: result.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}

// Register handles user registration
func Register(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		// Register user
		result, err := app.AuthService.RegisterWithPassword(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return authentication result
		c.JSON(http.StatusCreated, AuthResponse{
			Token: result.Token,
			User: UserProfile{
				ID:       result.User.ID,
				Username: result.User.Username,
				Name:     result.User.Name,
			},
			ExpiresAt: result.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}

// GetAuthenticatedUser extracts the authenticated user from Gin context
func GetAuthenticatedUser(c *gin.Context) (*entity.User, bool) {
	user, exists := c.Get(UserKey)
	if !exists {
		return nil, false
	}

	if authenticatedUser, ok := user.(*entity.User); ok {
		return authenticatedUser, true
	}

	return nil, false
}
