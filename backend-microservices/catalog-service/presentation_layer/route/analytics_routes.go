package route

import (
	"catalog-service/domain_layer/middleware"
	"catalog-service/domain_layer/service"
	"catalog-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// AnalyticsRoutes sets up all analytics-related routes
func AnalyticsRoutes(router *gin.RouterGroup, analyticsController *controller.AnalyticsController, userManagementService *service.UserManagementService) {
	analytics := router.Group("/analytics")

	// All analytics routes are protected (SuperAdmin only)
	analytics.Use(middleware.RequireSuperAdminWithAPIValidationMiddleware(userManagementService))
	{
		analytics.POST("", analyticsController.CreateAnalyticsEvent)
		analytics.GET("/:id", analyticsController.GetAnalyticsByID)
		analytics.DELETE("/:id", analyticsController.DeleteAnalytics)

		// Query routes
		analytics.GET("/date-range", analyticsController.GetAnalyticsByDateRange)
		analytics.GET("/user/:user_id", analyticsController.GetAnalyticsByUser)
		analytics.GET("/audiobook/:audiobook_id", analyticsController.GetAnalyticsByAudiobook)
		analytics.GET("/event/:event_type", analyticsController.GetAnalyticsByEventType)
		analytics.GET("/summary", analyticsController.GetAnalyticsSummary)
	}
}
