package middleware

import (
	"log"
	"microservice/user/helpers/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies JWT token and extracts user claims
func AuthMiddleware(tokenMaker utils.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("AuthMiddleware: Missing Authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
				"code":  "MISSING_AUTH_HEADER",
			})
			return
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("AuthMiddleware: Invalid authorization format - missing Bearer prefix")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
				"code":  "INVALID_AUTH_FORMAT",
			})
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("AuthMiddleware: Validating token at %v", time.Now())

		// Verify and parse the token
		claims, err := tokenMaker.ParseAccessToken(tokenString)
		if err != nil {
			log.Printf("AuthMiddleware: Token validation failed: %v", err)

			// Check specific error types
			switch {
			case strings.Contains(err.Error(), "token expired"):
				log.Printf("AuthMiddleware: Token expired at %v", time.Now())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token has expired",
					"code":  "TOKEN_EXPIRED",
				})
			case strings.Contains(err.Error(), "invalid token signature"):
				log.Printf("AuthMiddleware: Invalid token signature")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token signature",
					"code":  "INVALID_TOKEN_SIGNATURE",
				})
			default:
				log.Printf("AuthMiddleware: Invalid token - %v", err)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid or expired token",
					"code":  "INVALID_TOKEN",
				})
			}
			return
		}

		// Store user information in the context for later use
		ctx.Set("userID", claims.UserID)
		ctx.Set("userRole", claims.RoleName)

		log.Printf("AuthMiddleware: Auth successful for user: %s with role: %s at %v",
			claims.UserID, claims.RoleName, time.Now())
		ctx.Next()
	}
}
