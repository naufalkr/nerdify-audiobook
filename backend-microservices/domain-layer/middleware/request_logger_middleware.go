package middleware

import (
	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware logs the request body for debugging
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log POST/PUT/PATCH requests (which have bodies)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// Read the body
			bodyBytes, _ := io.ReadAll(c.Request.Body)

			// Log it
			log.Printf("[DEBUG] Request to %s %s: %s", c.Request.Method, c.Request.URL.Path, string(bodyBytes))

			// Put the body back
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Continue processing the request
		c.Next()
	}
}
