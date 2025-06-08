package route

import (
	"catalog-service/domain_layer/middleware"
	"catalog-service/domain_layer/service"
	"catalog-service/presentation_layer/controller"

	"github.com/gin-gonic/gin"
)

// UserRoutes sets up all user-related routes
func UserRoutes(router *gin.RouterGroup, userController *controller.UserController, userManagementService *service.UserManagementService) {
	users := router.Group("/users")

	// All user management routes are protected (SuperAdmin only)
	users.Use(middleware.RequireSuperAdminWithAPIValidationMiddleware(userManagementService))
	{
		users.POST("", userController.CreateUser)
		users.GET("", userController.GetAllUsers)
		users.GET("/search", userController.GetUserByEmail)
		users.GET("/:id", userController.GetUserByID)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
		users.GET("/role/:role", userController.GetUsersByRole)
	}
}
