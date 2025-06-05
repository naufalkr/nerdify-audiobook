package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RoleCheckMiddleware verifies that the user has at least one of the required roles
func RoleCheckMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the user role from the context
		userRole, exists := ctx.Get("userRole")
		if !exists {
			log.Printf("RoleCheckMiddleware: userRole not found in context")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			return
		}

		// Debug log the roles for troubleshooting
		log.Printf("RoleCheckMiddleware: Allowed roles: %v, User role: %v", allowedRoles, userRole)

		// Convert to string
		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("RoleCheckMiddleware: userRole is not a string: %T", userRole)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role format"})
			return
		}

		// Check if user has at least one of the allowed roles
		for _, role := range allowedRoles {
			if strings.EqualFold(userRoleStr, role) {
				// User has one of the allowed roles
				ctx.Next()
				return
			}
		}

		// No matching role found
		log.Printf("RoleCheckMiddleware: Access denied. Allowed: %v, Actual: %s", allowedRoles, userRoleStr)
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	}
}
