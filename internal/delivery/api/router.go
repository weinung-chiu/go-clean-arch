package api

import (
	"go-clean-arch/internal/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, app *usecase.Application) {
	r.GET("/home", func(c *gin.Context) {
		c.File("internal/delivery/api/qna.html")
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Session routes
	sessions := v1.Group("/sessions")
	sessions.GET(":id/ws", HandlerWebSocketSession(app))
	sessions.GET("/", HandlerListSessions(app))
	sessions.POST("/", HandlerNewSession(app))
	sessions.GET(":id", HandlerGetSession(app))

	// Auth routes (no auth required)
	auth := v1.Group("/auth")
	auth.POST("/register", HandlerRegister(app))
	auth.POST("/login", HandlerLogin(app))

	// Protected routes (require authentication)
	protected := v1.Group("/")
	protected.Use(AuthMiddleware(app))
	{
		// Protected session routes
		protected.POST("sessions/:id/questions", HandlerSubmitQuestion(app))

		// Protected question routes
		protected.POST("questions/:id/upvote", HandlerUpvoteQuestion(app))

		// Protected auth routes
		protected.GET("auth/profile", HandlerGetProfile(app))
	}
}
