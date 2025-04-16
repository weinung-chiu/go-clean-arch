package router

import "C"
import (
	"context"
	"go-clean-arch/internal/application"
	"go-clean-arch/internal/common"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupExampleHandlers(
	router *gin.Engine,
	app *application.ExampleApplication,
	logger *slog.Logger,
) {
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(TraceIDMiddleware("X-Trace-ID"))
	router.Use(AccessLogMiddleware(logger.With("service", "DUMMY_SERVICE_NAME")))

	r := router.Group("/api")
	v1 := r.Group("/v1")
	examples := v1.Group("/examples")

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	examples.GET("/logging", loggingExampleHandler(app))
	examples.GET("/random", randomExampleHandler(app))
	examples.GET("/clock", clockExampleHandler(app))
}

type RandomExample struct {
	MockRNG   int
	PseudoRNG int
	CryptoRNG int
}

func loggingExampleHandler(app *application.ExampleApplication) gin.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		err := app.DemoLog(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
		}
	}
}

func randomExampleHandler(app *application.ExampleApplication) gin.HandlerFunc {
	type response struct {
		Message string `json:"message"`
		Results struct {
			MockRNG   int `json:"mock"`
			PseudoRNG int `json:"pseudo"`
			CryptoRNG int `json:"crypto"`
		} `json:"result"`
	}

	return func(c *gin.Context) {
		result, _ := app.DemoRandom(c.Request.Context())

		c.JSON(http.StatusOK, response{
			Message: "hello world",
			Results: struct {
				MockRNG   int `json:"mock"`
				PseudoRNG int `json:"pseudo"`
				CryptoRNG int `json:"crypto"`
			}{MockRNG: result.MockRNG, PseudoRNG: result.PseudoRNG, CryptoRNG: result.CryptoRNG},
		})
	}
}

func clockExampleHandler(app *application.ExampleApplication) gin.HandlerFunc {
	type response struct {
		Message string `json:"message"`
		Results struct {
			Before string `json:"before"`
			After  string `json:"after"`
			Delta  string `json:"delta"`
		} `json:"result"`
	}

	return func(c *gin.Context) {
		result, _ := app.DemoClock(c.Request.Context())

		c.JSON(http.StatusOK, response{
			Message: "example of mock clock, 此例模擬一個假的時間起點(before)及假的處理時間(delta)，結束於 after",
			Results: struct {
				Before string `json:"before"`
				After  string `json:"after"`
				Delta  string `json:"delta"`
			}{
				Before: result.Before.Format(time.RFC3339Nano),
				After:  result.After.Format(time.RFC3339Nano),
				Delta:  result.Delta.String(),
			},
		})
	}
}

// TraceIDMiddleware 從 request 中取得已有的 trace ID ，如果不存在就建立新的。取得之後會放入 Context 內供後續使用
func TraceIDMiddleware(field string) gin.HandlerFunc {
	// TODO: 用 "github.com/gin-contrib/requestid" 處理，加入 response Header
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
		slog.Debug("Trace-ID", "id", common.GetTraceID(c.Request.Context()))
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
