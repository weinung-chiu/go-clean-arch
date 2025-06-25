package api

import (
	"github.com/gin-gonic/gin"
	"go-clean-arch/internal/usecase"
)

func RegisterRoutes(r *gin.Engine, app *usecase.Application) {
	r.GET("/home", func(c *gin.Context) {
		c.File("internal/delivery/api/qna.html")
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")
	sessions := v1.Group("/sessions")
	sessions.GET(":id/ws", HandlerWebSocketSession(app))
	sessions.GET("/", HandlerListSessions(app))
	sessions.POST("/", HandlerNewSession(app))
	sessions.GET(":id", HandlerGetSession(app))
	sessions.POST(":id/questions", HandlerSubmitQuestion(app))

	questions := v1.Group("/questions")
	questions.POST(":id/upvote", HandlerUpvoteQuestion(app))
}
