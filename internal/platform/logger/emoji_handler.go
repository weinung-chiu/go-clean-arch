package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

// EmojiHandler demonstrates the purpose of internal/platform/ by providing
// a custom logging format that's both human-readable and AI-parseable.
// This shows how platform-specific concerns (logging format, observability tools)
// are abstracted away from business logic while providing enhanced functionality.
//
// Designed for local development only - outputs single-line emoji-formatted logs to stdout.
//
// Example output: 15:04:05 [INFO] 🔍 TRACE(abc123...) 🏢 API 📖 GET(/articles) ✅ 200 ⏱️ 45ms
//
// Benefits:
// - Human-friendly: Visual emojis make logs easy to scan
// - AI-friendly: Semantic structure enables automated log analysis
// - Platform abstraction: Business logic doesn't know about log formatting
// - Clean Architecture: Infrastructure concerns isolated in platform layer
type EmojiHandler struct{}

// NewEmojiHandler creates a new AI-friendly logging handler for local development
func NewEmojiHandler() *EmojiHandler {
	return &EmojiHandler{}
}

func (h *EmojiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Enable all levels for local development
	return level >= slog.LevelDebug
}

func (h *EmojiHandler) Handle(ctx context.Context, record slog.Record) error {
	// Create AI-friendly formatted message
	formatted := h.formatForAI(record)

	// Output directly as single line with timestamp and level
	levelStr := record.Level.String()
	timeStr := record.Time.Format("15:04:05")

	// Write complete formatted line directly to stdout
	_, err := fmt.Fprintf(os.Stdout, "%s [%s] %s\n", timeStr, levelStr, formatted)
	return err
}

func (h *EmojiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, return same handler (attrs will be in the record)
	return h
}

func (h *EmojiHandler) WithGroup(name string) slog.Handler {
	// For simplicity, return same handler
	return h
}

// formatForAI converts log records into AI and human-friendly format with semantic emojis
func (h *EmojiHandler) formatForAI(record slog.Record) string {
	var parts []string
	attrs := make(map[string]interface{})

	// Extract attributes from record
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	})

	// Add trace context if available
	if traceID, ok := attrs["trace_id"].(string); ok && traceID != "" {
		parts = append(parts, fmt.Sprintf("🔍 TRACE(%s)", traceID[:8]+"..."))
	}

	// Add service context
	if service, ok := attrs["service"].(string); ok {
		parts = append(parts, fmt.Sprintf("🏢 %s", strings.ToUpper(service)))
	}

	// Add HTTP method and path if available
	if method, ok := attrs["method"].(string); ok {
		if path, pathOk := attrs["path"].(string); pathOk {
			emoji := getMethodEmoji(method)
			parts = append(parts, fmt.Sprintf("%s %s(%s)", emoji, method, path))
		}
	}

	// Add status with appropriate emoji
	if status, ok := attrs["status"].(int); ok {
		emoji := getStatusEmoji(status)
		parts = append(parts, fmt.Sprintf("%s %d", emoji, status))
	}

	// Add duration if available
	if duration, ok := attrs["duration"]; ok {
		if d, ok := duration.(time.Duration); ok {
			parts = append(parts, fmt.Sprintf("⏱️ %v", d.Round(time.Millisecond)))
		}
	}

	// Add user context if available
	if userID, ok := attrs["user_id"].(string); ok {
		parts = append(parts, fmt.Sprintf("👤 %s", userID))
	}

	// Add client IP if available
	if clientIP, ok := attrs["client_ip"].(string); ok {
		parts = append(parts, fmt.Sprintf("🌍 %s", clientIP))
	}

	// Add error information if available
	if errorValue, ok := attrs["error"]; ok {
		parts = append(parts, fmt.Sprintf("💥 ERROR(%v)", errorValue))
	}

	// Combine with original message
	formatted := strings.Join(parts, " ")
	if record.Message != "" {
		if formatted != "" {
			formatted = fmt.Sprintf("%s | %s", formatted, record.Message)
		} else {
			formatted = record.Message
		}
	}

	return formatted
}

// getMethodEmoji returns appropriate emoji for HTTP methods
func getMethodEmoji(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "📖"
	case "POST":
		return "📝"
	case "PUT":
		return "✏️"
	case "DELETE":
		return "🗑️"
	case "PATCH":
		return "🔧"
	default:
		return "🌐"
	}
}

// getStatusEmoji returns appropriate emoji for HTTP status codes
func getStatusEmoji(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "✅"
	case status >= 300 && status < 400:
		return "↩️"
	case status >= 400 && status < 500:
		return "❌"
	case status >= 500:
		return "💥"
	default:
		return "❓"
	}
}
