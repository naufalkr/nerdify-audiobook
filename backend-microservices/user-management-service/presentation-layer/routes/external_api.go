package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// ExternalAPIRoutes sets up routes for external service API access
func ExternalAPIRoutes(r *gin.Engine, tenantController *controller.TenantAPIController, tokenMaker utils.TokenMaker) {
	// The external API route group allows services like "asset_management" to access tenant data
	api := r.Group("/api/external")

	// Add API Key middleware for all external routes
	api.Use(middleware.APIKeyMiddleware())

	// ===== AUTHENTICATION & AUTHORIZATION APIs =====
	// Token validation and user info for other microservices
	api.POST("/auth/validate-token", tenantController.ValidateJWTToken)
	api.POST("/auth/user-info", tenantController.GetUserInfoFromToken)
	api.POST("/auth/validate-user-permissions", tenantController.ValidateUserPermissions)
	api.GET("/auth/validate-superadmin", tenantController.ValidateIsSuperAdmin) // This endpoint bypasses API key validation, requires only JWT token

	// ===== TENANT MANAGEMENT APIs =====
	// Basic tenant operations for external service consumption
	api.GET("/tenants", tenantController.ListTenants)
	api.GET("/tenants/:id", tenantController.GetTenantById)
	api.GET("/tenants/:id/validate", tenantController.ValidateTenantAccess)

	// ===== BUSINESS LOGIC APIs =====
	// Tenant subscription and limits for business validation
	api.GET("/tenants/:id/subscription", tenantController.GetTenantSubscription)
	api.GET("/tenants/:id/limits", tenantController.GetTenantLimits)
	api.GET("/tenants/:id/users", tenantController.GetTenantUsers)
	api.POST("/tenants/:id/validate-user-access", tenantController.ValidateUserTenantAccess)
	api.GET("/users/:userId/tenants", tenantController.GetUserTenants)
}
