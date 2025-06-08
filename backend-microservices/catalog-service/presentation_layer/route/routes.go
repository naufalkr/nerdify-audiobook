package route

import (
	"catalog-service/domain_layer/service"
	"catalog-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	authorController *controller.AuthorController,
	readerController *controller.ReaderController,
	genreController *controller.GenreController,
	audiobookController *controller.AudiobookController,
	trackController *controller.TrackController,
	userController *controller.UserController,
	analyticsController *controller.AnalyticsController,
	userManagementService *service.UserManagementService,
) {
	// API versioning
	api := router.Group("/api/v1")

	// Health check endpoint
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "catalog-service",
		})
	})

	// Setup route groups
	AuthorRoutes(api, authorController, userManagementService)
	ReaderRoutes(api, readerController, userManagementService)
	GenreRoutes(api, genreController, userManagementService)
	AudiobookRoutes(api, audiobookController, userManagementService)
	TrackRoutes(api, trackController, userManagementService)
	UserRoutes(api, userController, userManagementService)
	AnalyticsRoutes(api, analyticsController, userManagementService)
}
