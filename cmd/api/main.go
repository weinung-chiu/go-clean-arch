package main

import (
	"context"
	"errors"
	"fmt"
	"go-clean-arch/internal/adapter"
	"go-clean-arch/internal/config"
	httpdelivery "go-clean-arch/internal/delivery/http"
	"go-clean-arch/internal/platform/logger"
	"go-clean-arch/internal/usecase"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	AppName  = "go-clean-arch"
	AppBuild = "dev"
)

func main() {
	// Load application configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Set up logging
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		slog.Default().Error("Failed to parse log level", "error", err)
		os.Exit(1)
	}

	rootLogger := slog.New(logger.NewSimpleHandler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
	))
	rootLogger = rootLogger.With("service", AppName, "build", AppBuild)

	// Initialize application dependencies
	blogRepo, err := adapter.NewPostgresBlogRepository(cfg.DatabaseDSN)
	if err != nil {
		rootLogger.Error("Failed to create blog repository", "error", err)
		os.Exit(1)
	}

	app, err := usecase.NewApplication(usecase.NewApplicationParams{
		Logger:   rootLogger,
		BlogRepo: blogRepo,
	})
	if err != nil {
		rootLogger.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	// Set up HTTP server
	gin.SetMode(gin.ReleaseMode)
	ginRouter := httpdelivery.SetupRouter(app, rootLogger)

	// Configure server
	httpAddr := fmt.Sprintf(":%d", cfg.ApiPort)
	server := &http.Server{
		Addr:           httpAddr,
		Handler:        ginRouter,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		rootLogger.InfoContext(ctx, "Starting HTTP server", "addr", httpAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			rootLogger.ErrorContext(ctx, "HTTP server failed", "error", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	select {
	case sig := <-sigChan:
		rootLogger.InfoContext(ctx, "Received shutdown signal", "signal", sig)
	case <-ctx.Done():
		rootLogger.InfoContext(ctx, "Context cancelled")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	rootLogger.InfoContext(shutdownCtx, "Shutting down HTTP server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		rootLogger.ErrorContext(shutdownCtx, "HTTP server shutdown failed", "error", err)
		os.Exit(1)
	}

	rootLogger.InfoContext(shutdownCtx, "Server shutdown complete")
}
