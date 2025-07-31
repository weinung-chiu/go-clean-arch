package api

import (
	"go-clean-arch/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all HTTP routes for the API.
func RegisterRoutes(r *gin.Engine, app *usecase.Application) {
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
