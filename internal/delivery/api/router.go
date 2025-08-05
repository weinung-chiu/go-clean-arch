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
		v1.POST("/articles", createArticle(app))
		v1.PUT("/articles/:id", updateArticle(app))
		v1.POST("/articles/:id/publish", publishArticle(app))
	}
}

// CreateArticleRequest defines the request body for creating an article
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}

// UpdateArticleRequest defines the request body for updating an article
type UpdateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// listPublishedArticles returns a handler for listing published articles
func listPublishedArticles(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		publishedOnly := true
		filter := usecase.BlogFilter{
			PublishedOnly: &publishedOnly,
		}

		articles, err := app.ListArticles(c.Request.Context(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"articles": articles})
	}
}

// createArticle returns a handler for creating a new article
func createArticle(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateArticleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		article, err := app.CreateArticle(c.Request.Context(), req.Title, req.Content, req.AuthorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"article": article})
	}
}

// updateArticle returns a handler for updating an existing article
func updateArticle(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		articleID := c.Param("id")
		if articleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
			return
		}

		var req UpdateArticleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		article, err := app.UpdateArticle(c.Request.Context(), articleID, req.Title, req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"article": article})
	}
}

// publishArticle returns a handler for publishing an existing article
func publishArticle(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		articleID := c.Param("id")
		if articleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
			return
		}

		article, err := app.PublishArticle(c.Request.Context(), articleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish article"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"article": article})
	}
}
