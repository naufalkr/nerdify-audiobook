package routes

import (
	"log"
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// TenantRoutes mengatur routes untuk fitur tenant
func TenantRoutes(router *gin.Engine, controller *controller.TenantController, tokenMaker utils.TokenMaker) {
	// ================================= TENANT ROUTES =================================

	// ================================= USER ROUTES =================================
	setupUserRoutes(router, controller, tokenMaker)

	// ================================= ADMIN ROUTES =================================
	setupAdminRoutes(router, controller, tokenMaker)

	// ================================= SUPERADMIN ROUTES =================================
	setupSuperAdminRoutes(router, controller, tokenMaker)
}

// setupUserRoutes mengatur routes untuk user reguler dengan token
func setupUserRoutes(router *gin.Engine, controller *controller.TenantController, tokenMaker utils.TokenMaker) {
	// User routes untuk mendapatkan tenants yang diikuti user
	userRoutes := router.Group("/api/tenants")
	userRoutes.Use(middleware.AuthMiddleware(tokenMaker)) // Selalu menggunakan token
	{
		// Mendapatkan tenant yang diikuti oleh user yang sedang login
		userRoutes.GET("/user-tenants", controller.GetUserTenants)

		// Mendapatkan detail tenant yang sedang diakses user
		userRoutes.GET("/detail/tenant", controller.GetUserCurrentTenant)

		// Mendapatkan daftar user dalam tenant
		userRoutes.GET("/users", controller.GetCurrentTenantUsers)
	}
}

// setupAdminRoutes mengatur routes untuk admin dengan token
func setupAdminRoutes(router *gin.Engine, controller *controller.TenantController, tokenMaker utils.TokenMaker) {
	// Admin routes untuk mengelola tenant dimana admin berada (menggunakan token)
	adminRoutes := router.Group("/api/admin/tenants")
	adminRoutes.Use(middleware.AuthMiddleware(tokenMaker)) // Selalu menggunakan token
	adminRoutes.Use(middleware.RoleCheckMiddleware([]string{"ADMIN"}))
	{
		// Admin mendapatkan detail tenant dimana dia berada (tanpa perlu ID)
		adminRoutes.GET("", controller.GetUserCurrentTenant)

		// Admin mengupdate tenant dimana dia berada
		adminRoutes.PUT("", controller.UpdateCurrentTenant)

		// Admin melihat daftar user di tenant
		adminRoutes.GET("/users", controller.GetCurrentTenantUsers)

		// Admin mengundang user ke tenant
		adminRoutes.POST("/invite", controller.InviteUserToCurrentTenant)

		// Admin menghapus user dari tenant
		adminRoutes.DELETE("/users/:userID", controller.RemoveUserFromCurrentTenant)

		// Admin mengupdate subscription tenant
		adminRoutes.POST("/subscription", controller.UpdateCurrentSubscription)

		// Admin mengupdate contact email tenant
		adminRoutes.PATCH("/contact", controller.UpdateCurrentTenantContact)

		// Admin mengupdate logo tenant
		adminRoutes.PATCH("/logo", controller.UpdateCurrentTenantLogo)
	}

}

// setupSuperAdminRoutes mengatur routes untuk SuperAdmin dengan ID dan token
func setupSuperAdminRoutes(router *gin.Engine, controller *controller.TenantController, tokenMaker utils.TokenMaker) {
	// Satu router group untuk semua route SuperAdmin
	superAdminRoutes := router.Group("/api/superadmin")
	superAdminRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	superAdminRoutes.Use(middleware.RoleMiddleware("SUPERADMIN"))
	{
		// CRUD operasi tenant
		superAdminRoutes.POST("/tenants", controller.CreateTenant)
		superAdminRoutes.GET("/tenants", controller.GetAllTenants)
		superAdminRoutes.DELETE("/tenants/:tenantID", controller.DeleteTenant)
		superAdminRoutes.PUT("/tenants/:tenantID", controller.UpdateTenant)

		// Get detailed information about a tenant
		superAdminRoutes.GET("/tenants/:tenantID/details", controller.GetCurrentTenantDetails)

		// Subscription management
		superAdminRoutes.POST("/tenants/:tenantID/subscription", controller.UpdateTenantSubscription)

		// User invitation and management
		superAdminRoutes.POST("/tenants/invite", controller.InviteUserToTenant)
		superAdminRoutes.POST("/tenants/:tenantID/users/:userID/promote", controller.PromoteToAdmin)
		superAdminRoutes.POST("/tenants/:tenantID/users/:userID/demote", controller.DemoteFromAdmin)
		superAdminRoutes.POST("/tenants/:tenantID/users/:userID/invite", controller.DirectInviteUserToTenant)
		superAdminRoutes.GET("/tenants/:tenantID/users", controller.GetTenantUsers)
		superAdminRoutes.DELETE("/tenants/:tenantID/users/:userID", func(c *gin.Context) {
			tenantID := c.Param("tenantID")
			userID := c.Param("userID")
			// Debug log to verify parameters
			log.Printf("Attempting to remove user %s from tenant %s by SuperAdmin", userID, tenantID)
			controller.RemoveUserFromTenant(c)
		})

		// Tenant customization
		superAdminRoutes.PATCH("/tenants/:tenantID/contact", controller.UpdateTenantContact)
		superAdminRoutes.PATCH("/tenants/:tenantID/logo", controller.UpdateTenantLogo)
	}
}
