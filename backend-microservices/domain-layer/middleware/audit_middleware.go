package middleware

import (
	"bytes"
	"io"
	"microservice/user/domain-layer/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditMiddlewareFunc creates a middleware that logs all HTTP requests
func AuditMiddlewareFunc(auditService service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for certain paths if needed
		if shouldSkipAudit(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Read and restore the request body
		var requestBody []byte
		if c.Request.Body != nil && c.Request.Method != "GET" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Capture the response
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// Process the request
		c.Next()

		// Extract user ID from context
		userID, _ := c.Get("userID")
		userIDStr, ok := userID.(string)
		if !ok {
			userIDStr = ""
		}

		// Create audit log entry
		auditService.LogActivity(c.Request.Context(), &service.LogActivityRequest{
			UserID:     userIDStr,
			EntityID:   "http-" + uuid.New().String(),
			EntityType: "HTTP Request",
			Action:     c.Request.Method + " " + c.Request.URL.Path,
			Detail:     c.Request.URL.RawQuery,
			OldValue:   string(requestBody),
			NewValue:   responseWriter.body.String(),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		})
	}
}

// shouldSkipAudit determines if audit logging should be skipped for a path
func shouldSkipAudit(path string) bool {
	// Skip health checks, static files, etc.
	return strings.HasPrefix(path, "/health") ||
		strings.HasPrefix(path, "/static") ||
		path == "/favicon.ico" ||
		path == "/"
}

// responseBodyWriter is a custom response writer that captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AuditMiddleware is a struct-based middleware for audit logging
type AuditMiddleware struct {
	auditService service.AuditService
}

// NewAuditMiddleware creates a new audit middleware
func NewAuditMiddleware(auditService service.AuditService) *AuditMiddleware {
	return &AuditMiddleware{
		auditService: auditService,
	}
}

// AuditLog returns a gin.HandlerFunc that logs all requests
func (m *AuditMiddleware) AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for certain paths if needed
		if shouldSkipAudit(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Read and restore the request body
		var requestBody []byte
		if c.Request.Body != nil && c.Request.Method != "GET" {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Capture the response
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// Process the request
		c.Next()

		// Extract user ID from context
		userID, _ := c.Get("userID")
		userIDStr, ok := userID.(string)
		if !ok {
			userIDStr = ""
		}

		// Create audit log entry
		m.auditService.LogActivity(c.Request.Context(), &service.LogActivityRequest{
			UserID:     userIDStr,
			EntityID:   "http-" + uuid.New().String(),
			EntityType: "HTTP Request",
			Action:     c.Request.Method + " " + c.Request.URL.Path,
			Detail:     c.Request.URL.RawQuery,
			OldValue:   string(requestBody),
			NewValue:   responseWriter.body.String(),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		})
	}
}
