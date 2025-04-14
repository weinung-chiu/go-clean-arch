package main

import (
	"context"
	"go-clean-arch/internal/common"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

var traceIDKey = "traceIDKey"

func main() {
	//rootLogger := slog.New(NewTraceHandler(slog.NewJSONHandler(os.Stdout, nil)))
	rootLogger := slog.New(common.NewTraceHandler(slog.NewTextHandler(os.Stdout, nil)))
	rootLogger.Info("initial app with log/slog", "my_key", "my_value")

	rootCtx := context.Background()
	ctx := middleware(rootCtx)

	svc := NewService(rootLogger)
	svc.DoSomething(ctx)
}

func middleware(ctx context.Context) context.Context {
	traceID := uuid.New().String()
	return context.WithValue(ctx, traceIDKey, traceID)
}

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{logger.With("component", "service")}
}

func (s *Service) DoSomething(ctx context.Context) {
	s.logger.InfoContext(ctx, "do something...", "key_in_svc", "value_in_svc")
}
