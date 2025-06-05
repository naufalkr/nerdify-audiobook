package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// UserTenantContextRoutes sets up routes for user-tenant context management
func UserTenantContextRoutes(router *gin.Engine, controller *controller.TenantAPIController, tokenMaker utils.TokenMaker) {
	// ================================= USER TENANT CONTEXT ROUTES =================================

	// All routes require authentication
	apiRoutes := router.Group("/api/v1/user-tenant")
	apiRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// Current tenant management
		apiRoutes.GET("/current", controller.GetCurrentTenant)
		apiRoutes.PUT("/current", controller.SetCurrentTenant)

		// User tenants management
		apiRoutes.GET("/tenants", controller.GetUserTenants)
		apiRoutes.POST("/switch", controller.SwitchTenant)

		// Access validation
		apiRoutes.POST("/validate-access", controller.ValidateUserTenantAccessContext)

		// Admin routes - get users in current tenant
		apiRoutes.GET("/users", controller.GetTenantUsersContext)
	}

	// SuperAdmin routes - access to any tenant
	superAdminRoutes := router.Group("/api/v1/user-tenant")
	superAdminRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	superAdminRoutes.Use(middleware.RoleMiddleware("SuperAdmin"))
	{
		// Get users by specific tenant ID
		superAdminRoutes.GET("/tenants/:tenantId/users", controller.GetTenantUsersByTenantID)
	}
}
