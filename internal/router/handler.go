package router

import (
	"context"
	application "go-clean-arch/internal/application"
	"go-clean-arch/internal/common"
	"log/slog"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupDevHandlers(
	router *gin.Engine,
	app *application.DevelopApplication,
	logger *slog.Logger,
) {
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(TraceIDMiddleware("X-Trace-ID"))
	router.Use(AccessLogMiddleware(logger.With("service", "DUMMY_SERVICE_NAME")))

	r := router.Group("/api")
	v1 := r.Group("/v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}

// TraceIDMiddleware 從 request 中取得已有的 trace ID ，如果不存在就建立新的。取得之後會放入 Context 內供後續使用
func TraceIDMiddleware(field string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(field)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request.Context(), common.TraceIDCtxKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// AccessLogMiddleware 實作 Access Log
// 見 https://vgjira.atlassian.net/wiki/spaces/GTAT/pages/190809401/DataDog+GT+Log
// TODO: 尚未完全符合格式、trace-ID, time 有 duplicate
func AccessLogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logger = logger.With(
			slog.Time("timestamp", start),
			slog.String("trace_id", common.GetTraceID(c.Request.Context())),
		)
		groupLogger := logger.WithGroup("log")
		groupLogger = groupLogger.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
		)

		c.Next()

		groupLogger = groupLogger.With(
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int("status_code", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)

		if c.Writer.Status() >= 300 {
			groupLogger.ErrorContext(c.Request.Context(), "request failed")
		} else {
			groupLogger.InfoContext(c.Request.Context(), "request complete")
		}
	}
}
