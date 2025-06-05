package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// AdminRoutes sets up routes specifically for Admin role (not SuperAdmin)
// These routes allow admins to manage their tenant and users within their tenant
func AdminRoutes(router *gin.Engine, userController *controller.UserController, tenantController *controller.TenantController, tenantAPIController *controller.TenantAPIController, auditController *controller.AuditController, tokenMaker utils.TokenMaker) {
	// Admin routes group - requires Admin or SuperAdmin role
	adminRoutes := router.Group("/api/v1/admin")
	adminRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	adminRoutes.Use(middleware.RoleMiddleware("Admin", "SuperAdmin")) // Allow both Admin and SuperAdmin
	{
		// ===== TENANT MANAGEMENT =====
		// Admin can view and update their own tenant
		adminRoutes.GET("/tenant", tenantController.GetUserCurrentTenant)
		adminRoutes.PUT("/tenant", tenantController.UpdateCurrentTenant)
		adminRoutes.PATCH("/tenant/contact", tenantController.UpdateCurrentTenantContact)
		adminRoutes.PATCH("/tenant/logo", tenantController.UpdateCurrentTenantLogo)

		// ===== USER MANAGEMENT IN TENANT =====
		// Admin can manage users in their tenant
		adminRoutes.GET("/users", tenantAPIController.GetTenantUsersContext)
		adminRoutes.POST("/users/invite", tenantController.InviteUserToCurrentTenant)
		adminRoutes.DELETE("/users/:userID", tenantController.RemoveUserFromCurrentTenant)

		// Admin can view user profiles in their tenant
		adminRoutes.GET("/users/:userID/profile", userController.GetUserProfileById)

		// ===== USER TENANT CONTEXT MANAGEMENT =====
		// Admin can validate user access to tenant
		adminRoutes.POST("/validate-user-access", tenantAPIController.ValidateUserTenantAccessContext)

		// ===== SUBSCRIPTION MANAGEMENT =====
		// Admin can view and update subscription (limited scope)
		adminRoutes.GET("/subscription", tenantController.GetUserCurrentTenant) // Returns subscription info
		adminRoutes.POST("/subscription", tenantController.UpdateCurrentSubscription)

		// ===== AUDIT & MONITORING =====
		// Admin can view audit logs for their tenant
		adminRoutes.GET("/audit-logs", auditController.GetAuditLogs)

		// ===== TENANT STATS & ANALYTICS =====
		// Admin can view tenant statistics
		adminRoutes.GET("/stats/users", func(c *gin.Context) {
			// This could be implemented to show user count, active users, etc.
			c.JSON(200, gin.H{"message": "User statistics endpoint - to be implemented"})
		})
		adminRoutes.GET("/stats/usage", func(c *gin.Context) {
			// This could show tenant usage statistics
			c.JSON(200, gin.H{"message": "Usage statistics endpoint - to be implemented"})
		})
	}
}
