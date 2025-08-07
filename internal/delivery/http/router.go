package http

import (
	"go-clean-arch/internal/usecase"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the HTTP router with middleware chain and routes
func SetupRouter(app *usecase.Application, logger *slog.Logger) *gin.Engine {
	// Create router
	router := gin.New()

	// Health check endpoint (no logging middleware)
	router.GET("/health", HealthCheck())

	// Setup middleware chain for API routes only
	apiGroup := router.Group("")
	apiGroup.Use(TraceContext())        // W3C Trace Context (first for tracing)
	apiGroup.Use(RequestLogger(logger)) // Request logging with trace IDs
	apiGroup.Use(gin.Recovery())        // Panic recovery (last safety net)

	// API v1 routes with middleware
	v1 := apiGroup.Group("/api/v1")
	{
		// Auth endpoints (public)
		v1.POST("/auth/register", Register(app))
		v1.POST("/auth/login", Login(app))

		// Public articles endpoints
		v1.GET("/articles", ListPublishedArticles(app))

		// Protected articles endpoints
		protected := v1.Group("")
		protected.Use(JWTAuthMiddleware(app))
		{
			protected.POST("/articles", CreateArticle(app))
			protected.PUT("/articles/:id", UpdateArticle(app))
			protected.POST("/articles/:id/publish", PublishArticle(app))
			protected.DELETE("/articles/:id", DeleteArticle(app))
		}
	}

	return router
}
