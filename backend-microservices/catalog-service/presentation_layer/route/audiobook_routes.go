package route

import (
	"catalog-service/domain_layer/middleware"
	"catalog-service/domain_layer/service"
	"catalog-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// AudiobookRoutes sets up all audiobook-related routes
func AudiobookRoutes(router *gin.RouterGroup, audiobookController *controller.AudiobookController, userManagementService *service.UserManagementService) {
	audiobooks := router.Group("/audiobooks")
	{
		// Public routes (no authentication required)
		audiobooks.GET("", audiobookController.GetAllAudiobooks)
		audiobooks.GET("/search", audiobookController.SearchAudiobooks)
		audiobooks.GET("/:id", audiobookController.GetAudiobookByID)

		// Protected routes (SuperAdmin only)
		adminRoutes := audiobooks.Group("")
		adminRoutes.Use(middleware.RequireSuperAdminWithAPIValidationMiddleware(userManagementService))
		{
			adminRoutes.POST("", audiobookController.CreateAudiobook)
			adminRoutes.PUT("/:id", audiobookController.UpdateAudiobook)
			adminRoutes.DELETE("/:id", audiobookController.DeleteAudiobook)

			// Genre management
			adminRoutes.POST("/:id/genres", audiobookController.AddGenresToAudiobook)
			adminRoutes.DELETE("/:id/genres", audiobookController.RemoveGenresFromAudiobook)
		}
	}
}
