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

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Articles endpoints
		v1.GET("/articles", listPublishedArticles(app))
	}
}

// listPublishedArticles returns a handler for listing published articles
func listPublishedArticles(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		publishedOnly := true
		filter := usecase.BlogFilter{
			PublishedOnly: &publishedOnly,
		}

		articles, err := app.ListPosts(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"articles": articles})
	}
}
