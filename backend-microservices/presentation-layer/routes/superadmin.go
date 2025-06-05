package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SuperAdminRoutes sets up routes specifically for SuperAdmin role
// These routes provide system-wide management capabilities
func SuperAdminRoutes(router *gin.Engine, userController *controller.UserController, tenantController *controller.TenantController, tenantAPIController *controller.TenantAPIController, roleController *controller.RoleController, auditController *controller.AuditController, tokenMaker utils.TokenMaker) {
	// SuperAdmin routes group - requires SuperAdmin role only
	superAdminRoutes := router.Group("/api/v1/superadmin")
	superAdminRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	superAdminRoutes.Use(middleware.RoleMiddleware("SuperAdmin"))
	{
		// ===== SYSTEM-WIDE USER MANAGEMENT =====
		superAdminRoutes.GET("/users", userController.ListUsers)
		superAdminRoutes.POST("/users", userController.CreateUserByAdmin)
		superAdminRoutes.GET("/users/:id", userController.GetUserProfileById)
		superAdminRoutes.PUT("/users/:id", userController.UpdateUserData)
		superAdminRoutes.DELETE("/users/:id", userController.DeleteUserByID)
		superAdminRoutes.DELETE("/users/:id/permanent", userController.HardDeleteUserByID)
		superAdminRoutes.POST("/users/:id/verify-email", userController.VerifyEmailById)
		superAdminRoutes.PUT("/users/:id/role", userController.ChangeUserRole)

		// ===== SYSTEM-WIDE TENANT MANAGEMENT =====
		superAdminRoutes.GET("/tenants", tenantController.GetAllTenants)
		superAdminRoutes.POST("/tenants", tenantController.CreateTenant)
		superAdminRoutes.GET("/tenants/:tenantID", tenantController.GetCurrentTenantDetails)
		superAdminRoutes.PUT("/tenants/:tenantID", tenantController.UpdateTenant)
		superAdminRoutes.DELETE("/tenants/:tenantID", tenantController.DeleteTenant)

		// ===== TENANT-USER RELATIONSHIP MANAGEMENT =====
		superAdminRoutes.GET("/tenants/:tenantID/users", tenantAPIController.GetTenantUsersByTenantID)
		superAdminRoutes.POST("/tenants/:tenantID/users/:userID/invite", tenantController.DirectInviteUserToTenant)
		superAdminRoutes.DELETE("/tenants/:tenantID/users/:userID", tenantController.RemoveUserFromTenant)
		superAdminRoutes.POST("/tenants/:tenantID/users/:userID/promote", tenantController.PromoteToAdmin)
		superAdminRoutes.POST("/tenants/:tenantID/users/:userID/demote", tenantController.DemoteFromAdmin)

		// ===== SUBSCRIPTION & BILLING MANAGEMENT =====
		superAdminRoutes.GET("/tenants/:tenantID/subscription", tenantController.GetCurrentTenantDetails) // Returns subscription info
		superAdminRoutes.POST("/tenants/:tenantID/subscription", tenantController.UpdateTenantSubscription)

		// ===== TENANT CUSTOMIZATION =====
		superAdminRoutes.PATCH("/tenants/:tenantID/contact", tenantController.UpdateTenantContact)
		superAdminRoutes.PATCH("/tenants/:tenantID/logo", tenantController.UpdateTenantLogo)

		// ===== ROLE MANAGEMENT =====
		superAdminRoutes.GET("/roles", roleController.ListAllRoles)
		superAdminRoutes.POST("/roles", roleController.CreateRole)
		superAdminRoutes.GET("/roles/:id", roleController.GetRoleByID)
		superAdminRoutes.GET("/roles/name/:name", roleController.GetRoleByName)
		superAdminRoutes.PUT("/roles/:id", roleController.UpdateRole)
		superAdminRoutes.GET("/roles/system", roleController.GetSystemRoles)
		superAdminRoutes.POST("/roles/seed", roleController.SeedDefaultRoles)

		// ===== SYSTEM AUDIT & MONITORING =====
		superAdminRoutes.GET("/audit-logs", auditController.GetAuditLogs)
		superAdminRoutes.GET("/audit-logs/users/:userID", auditController.GetAuditLogs)     // User-specific audit logs
		superAdminRoutes.GET("/audit-logs/tenants/:tenantID", auditController.GetAuditLogs) // Tenant-specific audit logs

		// ===== SYSTEM STATISTICS & ANALYTICS =====
		superAdminRoutes.GET("/stats/overview", func(c *gin.Context) {
			// System-wide overview statistics
			c.JSON(200, gin.H{"message": "System overview statistics - to be implemented"})
		})
		superAdminRoutes.GET("/stats/tenants", func(c *gin.Context) {
			// Tenant statistics across the system
			c.JSON(200, gin.H{"message": "Tenant statistics - to be implemented"})
		})
		superAdminRoutes.GET("/stats/users", func(c *gin.Context) {
			// User statistics across the system
			c.JSON(200, gin.H{"message": "User statistics - to be implemented"})
		})

		// ===== SYSTEM MAINTENANCE =====
		superAdminRoutes.POST("/system/maintenance/enable", func(c *gin.Context) {
			// Enable maintenance mode
			c.JSON(200, gin.H{"message": "Maintenance mode enabled - to be implemented"})
		})
		superAdminRoutes.POST("/system/maintenance/disable", func(c *gin.Context) {
			// Disable maintenance mode
			c.JSON(200, gin.H{"message": "Maintenance mode disabled - to be implemented"})
		})

		// ===== EXTERNAL API MANAGEMENT =====
		// Routes for managing API keys and external service access
		superAdminRoutes.GET("/api-keys", func(c *gin.Context) {
			// List API keys for external services
			c.JSON(200, gin.H{"message": "API key management - to be implemented"})
		})
		superAdminRoutes.POST("/api-keys", func(c *gin.Context) {
			// Create new API key
			c.JSON(200, gin.H{"message": "API key creation - to be implemented"})
		})
		superAdminRoutes.DELETE("/api-keys/:keyID", func(c *gin.Context) {
			// Revoke API key
			c.JSON(200, gin.H{"message": "API key revocation - to be implemented"})
		})
	}
}
