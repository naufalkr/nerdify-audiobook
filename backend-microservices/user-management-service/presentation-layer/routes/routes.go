package routes

import (
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes is the main function to setup all routes for the application
// It delegates to specific route setup functions from other files
func SetupRoutes(r *gin.Engine, userController *controller.UserController, roleController *controller.RoleController, tenantController *controller.TenantController, tenantAPIController *controller.TenantAPIController, auditController *controller.AuditController, tokenMaker utils.TokenMaker) {
	// Setup user and authentication routes
	UserRoutes(r, userController, roleController, tokenMaker)

	// Setup authentication-specific routes (like verify email links)
	AuthRoutes(r, userController, tokenMaker)

	// Setup role management routes (legacy - keeping for backward compatibility)
	RoleRoutes(r, roleController)

	// Setup tenant routes (legacy - keeping for backward compatibility)
	TenantRoutes(r, tenantController, tokenMaker)

	// Setup user-tenant context routes
	UserTenantContextRoutes(r, tenantAPIController, tokenMaker)

	// Setup external API routes for service communication
	ExternalAPIRoutes(r, tenantAPIController, tokenMaker)

	// ===== NEW ROLE-BASED ROUTE ORGANIZATION =====
	// Setup dedicated Admin routes
	AdminRoutes(r, userController, tenantController, tenantAPIController, auditController, tokenMaker)

	// Setup dedicated SuperAdmin routes
	SuperAdminRoutes(r, userController, tenantController, tenantAPIController, roleController, auditController, tokenMaker)

	// Setup audit routes
	SetupAuditRoutes(r, auditController, tokenMaker)
}
