package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// W3C Trace Context constants
const (
	TraceParentHeader = "traceparent"
	TraceStateHeader  = "tracestate"
	TraceIDKey        = "trace_id"
	SpanIDKey         = "span_id"
)

// W3C traceparent format: version-traceid-spanid-flags
var traceParentRegex = regexp.MustCompile(`^([0-9a-f]{2})-([0-9a-f]{32})-([0-9a-f]{16})-([0-9a-f]{2})$`)

// TraceContext creates a middleware that handles W3C Trace Context
func TraceContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		var traceID, spanID string
		
		// Try to extract existing traceparent header
		traceparent := c.GetHeader(TraceParentHeader)
		if traceparent != "" {
			if matches := traceParentRegex.FindStringSubmatch(traceparent); len(matches) == 5 {
				traceID = matches[2]
				// Generate new span ID for this request (child span)
				spanID = generateSpanID()
			}
		}
		
		// Generate new trace if no valid traceparent found
		if traceID == "" {
			traceID = generateTraceID()
			spanID = generateSpanID()
		}
		
		// Create new traceparent header for response
		newTraceParent := fmt.Sprintf("00-%s-%s-01", traceID, spanID)
		
		// Add to Gin context
		c.Set(TraceIDKey, traceID)
		c.Set(SpanIDKey, spanID)
		
		// Add to request context for downstream services
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		ctx = context.WithValue(ctx, SpanIDKey, spanID)
		c.Request = c.Request.WithContext(ctx)
		
		// Add to response headers
		c.Header(TraceParentHeader, newTraceParent)
		
		c.Next()
	}
}

// RequestLogger creates a middleware that logs HTTP requests with W3C trace context
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Extract trace ID from context
		traceID, _ := c.Get(TraceIDKey)
		if traceID == nil {
			traceID = "unknown"
		}
		
		// Calculate duration
		duration := time.Since(start)
		
		// Log structured request details
		logger.InfoContext(c.Request.Context(), "HTTP request completed",
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"response_size", c.Writer.Size(),
		)
	}
}

// generateTraceID creates a new W3C compliant trace ID (16 bytes = 32 hex chars)
func generateTraceID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to a simple trace ID if crypto/rand fails
		return strings.Repeat("0", 32)
	}
	return hex.EncodeToString(bytes)
}

// generateSpanID creates a new W3C compliant span ID (8 bytes = 16 hex chars)
func generateSpanID() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to a simple span ID if crypto/rand fails
		return strings.Repeat("0", 16)
	}
	return hex.EncodeToString(bytes)
}