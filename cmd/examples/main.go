package main

import (
	"context"
	"errors"
	"fmt"
	"go-clean-arch/internal/application"
	"go-clean-arch/internal/common"
	"go-clean-arch/internal/router"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var traceIDKey = "traceIDKey"

func main() {
	rootLogger := slog.New(common.NewTraceHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	//rootLogger := slog.New(common.NewTraceHandler(slog.NewTextHandler(os.Stdout, nil)))
	rootLogger.Info("initial app with log/slog", "my_key", "my_value")

	rootCtx := context.Background()

	// logger 在 constructor 顯式帶入
	params := application.NewExampleAppParams{Logger: rootLogger}
	app := application.NewExampleApplication(params)

	// Create gin router
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()

	// Register all handlers
	router.SetupExampleHandlers(ginRouter, app, rootLogger)

	// Build HTTP server
	port := os.Getenv("PORT")
	httpAddr := fmt.Sprintf("0.0.0.0:%s", port)
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
		}
	}()

	<-rootCtx.Done()
}
