package logger

import (
	"context"
	"log/slog"
)

type SimpleHandler struct {
	base slog.Handler
}

func NewSimpleHandler(base slog.Handler) *SimpleHandler {
	return &SimpleHandler{base: base}
}

func (h *SimpleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

func (h *SimpleHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.base.Handle(ctx, record)
}

func (h *SimpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.base.WithAttrs(attrs)
}

func (h *SimpleHandler) WithGroup(name string) slog.Handler {
	return h.base.WithGroup(name)
}
