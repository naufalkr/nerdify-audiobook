package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware verifies API keys for external service access
func APIKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if this is the SuperAdmin validation endpoint - we'll bypass API key check for this endpoint
		if ctx.Request.URL.Path == "/api/external/auth/validate-superadmin" && ctx.Request.Method == "GET" {
			// For SuperAdmin validation, we only need the JWT token, which is checked in the handler
			ctx.Next()
			return
		}

		// Get the API key from header
		apiKey := ctx.GetHeader("X-API-Key")

		// If no API key is provided
		if apiKey == "" {
			log.Println("APIKeyMiddleware: Missing API key header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			return
		}

		// In a production environment, you would validate the API key against stored keys
		// This could be from database, environment variables, or a key management service
		validAPIKeys := strings.Split(os.Getenv("VALID_API_KEYS"), ",")

		// For testing, if no valid keys are configured, accept a default test key
		if len(validAPIKeys) == 0 || (len(validAPIKeys) == 1 && validAPIKeys[0] == "") {
			if apiKey == "alat-service-api-key" {
				ctx.Next()
				return
			}
		} else {
			// Check if the provided API key is valid
			for _, key := range validAPIKeys {
				if apiKey == key {
					ctx.Next()
					return
				}
			}
		}

		// If we reach here, the API key is invalid
		log.Println("APIKeyMiddleware: Invalid API key")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
	}
}
