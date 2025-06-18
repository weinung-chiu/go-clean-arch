package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-clean-arch/internal/adapter"
	"go-clean-arch/internal/config"
	"go-clean-arch/internal/delivery/api"
	"go-clean-arch/internal/platform/logger"
	"go-clean-arch/internal/usecase"
	"log/slog"
	"os"
)

var (
	AppName  = "live-QA"
	AppBuild = "no_build"
)

func main() {
	// TODO: load from .env file only if local development
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load("./configs/.env")
		if err != nil {
			slog.Default().Error(fmt.Sprintf("failed to load .env file: %v", err))
		}
	}
	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error(fmt.Sprintf("failed to load config: %v", err))
		os.Exit(1)
	}

	var logLevel slog.Level
	err = logLevel.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		slog.Default().Error(fmt.Sprintf("failed to unmarshal log level: %v", err))
		os.Exit(1)
	}

	rootLogger := slog.New(logger.NewSimpleHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
	rootLogger = rootLogger.With("service", AppName, "build", AppBuild)

	memSessionRepo := adapter.NewMemorySessionRepo()
	memQuestionRepo := adapter.NewMemoryQuestionRepo()
	memEventBroker := adapter.NewInMemorySessionEventBroker()

	app, err := usecase.NewApplication(usecase.NewApplicationParams{
		Logger:       rootLogger,
		SessionRepo:  memSessionRepo,
		QuestionRepo: memQuestionRepo,
		EventBus:     memEventBroker,
	})
	if err != nil {
		rootLogger.Error("Failed to create application", "error", err)
		os.Exit(0)
	}

	rootCtx := context.Background()

	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()

	api.RegisterRoutes(ginRouter, app)

	// Build HTTP server
	httpAddr := fmt.Sprintf("0.0.0.0:%d", cfg.ApiPort)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: ginRouter,
	}

	// Run the server in a goroutine
	go func() {
		rootLogger.InfoContext(rootCtx, fmt.Sprintf("HTTP server is on http://%s", httpAddr))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			rootLogger.ErrorContext(rootCtx, "failed to start HTTP server", "error", err)
			os.Exit(1)
		}
	}()

	// for testing purposes, should do it in WebSocket handler
	time.Sleep(10 * time.Second)
	go func() {
		events, err := app.SubscribeToLastestSessionEvents(rootCtx)
		if err != nil {
			rootLogger.ErrorContext(rootCtx, "failed to subscribe to session events", "error", err)
			os.Exit(1)
		}
		for event := range events {
			switch event.Type {
			case usecase.SessionEventQuestionSubmitted:
				rootLogger.InfoContext(rootCtx, "New question submitted", "session_id", event.SessionID, "question", event.Payload)
			default:
				rootLogger.WarnContext(rootCtx, "Unknown event type", "type", event.Type)
			}
		}
	}()

	<-rootCtx.Done()

}
