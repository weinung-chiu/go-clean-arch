package router

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	application "go-clean-arch/internal/application"
)

func SetupDevHandlers(router *gin.Engine, app *application.DevelopApplication) {
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(
		requestid.New(
			requestid.WithGenerator(func() string { return uuid.New().String() }),
			//requestid.WithCustomHeaderStrKey("your-customer-key"),
		),
	)
	//ginRouter.Use(LoggerMiddleware(ctx))

	r := router.Group("/api")
	v1 := r.Group("/v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}
