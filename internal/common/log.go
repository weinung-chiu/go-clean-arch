package common

import (
	"context"
	"log/slog"
)

type traceIDKeyType struct{} // 不導出 type，避免其他套件自建 key
// TraceIDCtxKey 是導出的全域變數，可在整個 codebase 中作為 context key 使用
var TraceIDCtxKey = traceIDKeyType{}

func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(TraceIDCtxKey).(string); ok {
		return v
	}
	return ""
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
	if traceID, ok := ctx.Value(TraceIDCtxKey).(string); ok {
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
