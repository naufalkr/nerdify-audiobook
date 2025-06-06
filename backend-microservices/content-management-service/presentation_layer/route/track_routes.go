package route

import (
	"content-management-service/domain_layer/middleware"
	"content-management-service/domain_layer/service"
	"content-management-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// TrackRoutes sets up all track-related routes
func TrackRoutes(router *gin.RouterGroup, trackController *controller.TrackController, userManagementService *service.UserManagementService) {
	tracks := router.Group("/tracks")
	{
		// Public routes (no authentication required)
		tracks.GET("", trackController.GetAllTracks)
		tracks.GET("/:id", trackController.GetTrackByID)
		tracks.GET("/audiobook/:audiobook_id", trackController.GetTracksByAudiobook)

		// Protected routes (SuperAdmin only)
		adminRoutes := tracks.Group("")
		adminRoutes.Use(middleware.RequireSuperAdminWithAPIValidationMiddleware(userManagementService))
		{
			adminRoutes.POST("", trackController.CreateTrack)
			adminRoutes.PUT("/:id", trackController.UpdateTrack)
			adminRoutes.DELETE("/:id", trackController.DeleteTrack)
			adminRoutes.PUT("/audiobook/:audiobook_id/order", trackController.UpdateTrackOrder)
		}
	}
}
