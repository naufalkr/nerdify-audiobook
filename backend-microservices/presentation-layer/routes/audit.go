package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupAuditRoutes registers all routes for audit functionality
func SetupAuditRoutes(r *gin.Engine, auditController *controller.AuditController, tokenMaker utils.TokenMaker) {
	// Admin only routes - Superadmins can view all logs
	adminAudit := r.Group("/api/admin/audit-logs")
	adminAudit.Use(middleware.AuthMiddleware(tokenMaker), middleware.RoleMiddleware("SUPERADMIN"))
	{
		// View all audit logs with filters
		adminAudit.GET("", auditController.GetAuditLogs)
		adminAudit.GET("/:id", auditController.GetAuditLogByID)
		adminAudit.GET("/export", auditController.ExportAuditLogs)
		adminAudit.GET("/statistics", auditController.GetAuditStatistics)
	}

	// User-specific audit logs - Any authenticated user can view their own logs
	userAudit := r.Group("/api/audit-logs")
	userAudit.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// View user's own audit logs
		userAudit.GET("/my-logs", auditController.GetUserAuditLogs)
	}
}
