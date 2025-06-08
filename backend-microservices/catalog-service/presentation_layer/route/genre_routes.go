package route

import (
	"catalog-service/domain_layer/middleware"
	"catalog-service/domain_layer/service"
	"catalog-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// GenreRoutes sets up all genre-related routes
func GenreRoutes(router *gin.RouterGroup, genreController *controller.GenreController, userManagementService *service.UserManagementService) {
	genres := router.Group("/genres")
	{
		// Public routes (no authentication required)
		genres.GET("", genreController.GetAllGenres)
		genres.GET("/:id", genreController.GetGenreByID)
		genres.POST("/batch", genreController.GetGenresByIDs)

		// Protected routes (SuperAdmin only)
		adminRoutes := genres.Group("")
		adminRoutes.Use(middleware.RequireSuperAdminWithAPIValidationMiddleware(userManagementService))
		{
			adminRoutes.POST("", genreController.CreateGenre)
			adminRoutes.PUT("/:id", genreController.UpdateGenre)
			adminRoutes.DELETE("/:id", genreController.DeleteGenre)
		}
	}
}
