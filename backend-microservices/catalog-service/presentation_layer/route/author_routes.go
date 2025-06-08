package route

import (
	"catalog-service/domain_layer/middleware"
	"catalog-service/domain_layer/service"
	"catalog-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// AuthorRoutes sets up routes for author endpoints
func AuthorRoutes(router *gin.RouterGroup, authorController *controller.AuthorController, userManagementService *service.UserManagementService) {
	authors := router.Group("/authors")
	{
		// Public routes (no authentication required)
		authors.GET("", authorController.GetAllAuthors)
		authors.GET("/search", authorController.SearchAuthors)
		authors.GET("/:id", authorController.GetAuthor)

		// Protected routes (SuperAdmin only)
		adminRoutes := authors.Group("")
		adminRoutes.Use(middleware.RequireSuperAdminWithAPIValidationMiddleware(userManagementService))
		{
			adminRoutes.POST("", authorController.CreateAuthor)
			adminRoutes.PUT("/:id", authorController.UpdateAuthor)
			adminRoutes.DELETE("/:id", authorController.DeleteAuthor)
		}
	}
}
