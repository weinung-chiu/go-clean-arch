package main

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"os"
)

type traceIDKeyType struct{}

var traceIDKey = traceIDKeyType{}

func main() {
	//rootLogger := slog.New(NewTraceHandler(slog.NewJSONHandler(os.Stdout, nil)))
	rootLogger := slog.New(NewTraceHandler(slog.NewTextHandler(os.Stdout, nil)))
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

// TraceHandler wraps a slog.Handler and injects trace_id from context into logs
type TraceHandler struct {
	handler slog.Handler
}

func NewTraceHandler(h slog.Handler) *TraceHandler {
	return &TraceHandler{handler: h}
}

func (h *TraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	// Inject trace ID if present
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.handler.Handle(ctx, r)
}

func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{handler: h.handler.WithGroup(name)}
}
