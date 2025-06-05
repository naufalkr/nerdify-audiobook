package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// AuthRoutes sets up all authentication-related routes including verify email link
func AuthRoutes(router *gin.Engine, controller *controller.UserController, tokenMaker utils.TokenMaker) {
	// Email verification routes
	verifyRoutes := router.Group("/api/auth")
	{
		// Add both GET and POST endpoints for link verification
		verifyRoutes.GET("/verify-email-link", controller.VerifyEmailByLink)
		verifyRoutes.POST("/verify-email-link", controller.VerifyEmailByLink)
	}

	// Protected user profile routes
	profileRoutes := router.Group("/api/profile")
	profileRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	{
		profileRoutes.GET("/me", controller.GetUserProfile)
		profileRoutes.GET("/user/:id", middleware.RoleCheckMiddleware([]string{"SUPERADMIN"}), controller.GetUserProfileById)
	}
}
