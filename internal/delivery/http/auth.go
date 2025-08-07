package http

import (
	"go-clean-arch/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// Authorization header key
	AuthorizationHeader = "Authorization"

	// Bearer token prefix
	BearerPrefix = "Bearer "

	// Context key for authenticated user
	UserKey = "user"
)

// JWTAuthMiddleware creates a middleware that validates JWT tokens using Application
func JWTAuthMiddleware(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		// Validate token using AuthService
		user, err := app.AuthService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set authenticated user in context
		c.Set(UserKey, user)
		c.Next()
	}
}

// OptionalJWTAuthMiddleware creates a middleware that optionally validates JWT tokens
// If token is present and valid, sets user in context. If not present, continues without authentication.
func OptionalJWTAuthMiddleware(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			// Invalid format, continue without authentication
			c.Next()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			// No token after Bearer, continue without authentication
			c.Next()
			return
		}

		// Validate token using AuthService
		user, err := app.AuthService.ValidateToken(tokenString)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Set authenticated user in context
		c.Set(UserKey, user)
		c.Next()
	}
}
