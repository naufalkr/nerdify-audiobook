package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	SuperAdminRole = "superadmin"
)

// RoleMiddleware verifies that the user has one of the required roles
// Usage: RoleMiddleware("Admin") or RoleMiddleware("Admin", "SuperAdmin")
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the user role from the context
		userRole, exists := ctx.Get("userRole")
		if !exists {
			log.Printf("RoleMiddleware: userRole not found in context")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			return
		}

		// Debug log the roles for troubleshooting
		log.Printf("RoleMiddleware: Required roles: %v, User role: %v", requiredRoles, userRole)

		// Convert to string and compare (case-insensitive)
		userRoleStr, ok := userRole.(string)
		if !ok {
			log.Printf("RoleMiddleware: userRole is not a string: %T", userRole)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role format"})
			return
		}

		// Check if user has any of the required roles (case insensitive comparison)
		hasPermission := false
		for _, requiredRole := range requiredRoles {
			if strings.EqualFold(userRoleStr, requiredRole) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			log.Printf("RoleMiddleware: Access denied. Required: %v, Actual: %s", requiredRoles, userRoleStr)
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		ctx.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Get user role from context
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user has superadmin role
		if strings.ToLower(role.(string)) != SuperAdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Get user role from context
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, r := range roles {
			if strings.ToLower(role.(string)) == strings.ToLower(r) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
