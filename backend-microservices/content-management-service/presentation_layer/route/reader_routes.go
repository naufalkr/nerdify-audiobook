package route

import (
	"content-management-service/domain_layer/middleware"
	"content-management-service/domain_layer/service"
	"content-management-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// ReaderRoutes sets up all reader-related routes
func ReaderRoutes(router *gin.RouterGroup, readerController *controller.ReaderController, userManagementService *service.UserManagementService) {
	readers := router.Group("/readers")
	{
		// Public routes (no authentication required)
		readers.GET("", readerController.GetAllReaders)
		readers.GET("/search", readerController.SearchReaders)
		readers.GET("/:id", readerController.GetReaderByID)

		// Protected routes (SuperAdmin only)
		adminRoutes := readers.Group("")
		adminRoutes.Use(middleware.RequireSuperAdminWithAPIValidationMiddleware(userManagementService))
		{
			adminRoutes.POST("", readerController.CreateReader)
			adminRoutes.PUT("/:id", readerController.UpdateReader)
			adminRoutes.DELETE("/:id", readerController.DeleteReader)
		}
	}
}
