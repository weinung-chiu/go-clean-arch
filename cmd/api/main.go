package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go-clean-arch/internal/adapter"
	"go-clean-arch/internal/config"
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

	// --- DEMO: in-memory repositories ---
	memSessionRepo := adapter.NewMemorySessionRepo()
	memQuestionRepo := adapter.NewMemoryQuestionRepo()

	app, err := usecase.NewApplication(usecase.NewApplicationParams{
		Logger:       rootLogger,
		SessionRepo:  memSessionRepo,
		QuestionRepo: memQuestionRepo,
	})
	if err != nil {
		rootLogger.Error("Failed to create application", "error", err)
		os.Exit(0)
	}

	ctx := context.Background()

	// 1. Create a new session
	session, err := app.NewSession(ctx, "My Awesome Conference Talk")
	if err != nil {
		rootLogger.Error("Failed to create session", "error", err)
		return
	}
	rootLogger.Debug("Created session", "session", session)

	// 2. List all sessions
	sessions, err := app.ListSessions(ctx)
	if err != nil {
		rootLogger.Error("Failed to list sessions", "error", err)
		return
	}
	rootLogger.Debug("Sessions", "sessions", sessions)

	// 3. Submit a question
	question, err := app.SubmitQuestion(ctx, session.ID, "What was the biggest challenge?", "Alice")
	if err != nil {
		rootLogger.Error("Failed to submit question", "error", err)
		return
	}
	rootLogger.Debug("Submitted question", "question", question)

	// 4. Get session with questions
	sess, questions, err := app.GetSession(ctx, session.ID)
	if err != nil {
		rootLogger.Error("Failed to get session with questions", "error", err)
		return
	}
	rootLogger.Debug("Session with questions", "session", sess, "questions", questions)
}
