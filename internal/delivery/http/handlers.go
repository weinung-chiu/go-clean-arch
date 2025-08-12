package http

import (
	"go-clean-arch/internal/entity"
	"go-clean-arch/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateArticleRequest defines the request body for creating an article
type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UpdateArticleRequest defines the request body for updating an article
type UpdateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// HealthCheck returns a handler for health check endpoint
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// ListPublishedArticles returns a handler for listing published articles
func ListPublishedArticles(app *usecase.Application) gin.HandlerFunc {
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

// CreateArticle returns a handler for creating a new article
func CreateArticle(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateArticleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		authedUser, ok := user.(*entity.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
			return
		}
		authorID := authedUser.ID
		article, err := app.CreateArticle(c.Request.Context(), req.Title, req.Content, authorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"article": article})
	}
}

// UpdateArticle returns a handler for updating an existing article
func UpdateArticle(app *usecase.Application) gin.HandlerFunc {
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

		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		authedUser, ok := user.(*entity.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
			return
		}
		authorID := authedUser.ID
		article, err := app.UpdateArticle(c.Request.Context(), articleID, req.Title, req.Content, authorID)
		if err != nil {
			if err.Error() == "unauthorized: user can only update their own articles" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: you can only update your own articles"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"article": article})
	}
}

// PublishArticle returns a handler for publishing an existing article
func PublishArticle(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		articleID := c.Param("id")
		if articleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
			return
		}

		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		authedUser, ok := user.(*entity.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
			return
		}
		authorID := authedUser.ID
		article, err := app.PublishArticle(c.Request.Context(), articleID, authorID)
		if err != nil {
			if err.Error() == "unauthorized: user can only publish their own articles" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: you can only publish your own articles"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish article"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"article": article})
	}
}

// DeleteArticle returns a handler for deleting an existing article
func DeleteArticle(app *usecase.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		articleID := c.Param("id")
		if articleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
			return
		}

		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		authedUser, ok := user.(*entity.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
			return
		}
		authorID := authedUser.ID
		err := app.DeleteArticle(c.Request.Context(), articleID, authorID)
		if err != nil {
			if err.Error() == "unauthorized: user can only delete their own articles" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: you can only delete your own articles"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
	}
}
