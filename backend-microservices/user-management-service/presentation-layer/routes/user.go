package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// UserRoutes sets up all user-related routes including authentication and admin functions
func UserRoutes(router *gin.Engine, controller *controller.UserController, roleController *controller.RoleController, tokenMaker utils.TokenMaker) {
	// Public auth routes (no authentication required)
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", controller.Register)
		authRoutes.POST("/login", controller.Login)
		authRoutes.POST("/verify-email", controller.VerifyEmail)
		authRoutes.POST("/resend-verification-email", controller.ResendVerificationEmail)
		authRoutes.POST("/forgot-password", controller.ForgotPassword)
		authRoutes.POST("/reset-password", controller.ResetPassword)
		authRoutes.POST("/refresh-token", controller.RefreshToken)
	}

	// Protected user routes (authentication required)
	userRoutes := router.Group("/api/users")
	userRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	{
		userRoutes.POST("/logout", controller.Logout)
		userRoutes.GET("/profile", controller.GetUserProfile)
		userRoutes.PUT("/profile", controller.UpdateProfile)
		userRoutes.PATCH("/profile", controller.UpdateProfile)
		userRoutes.PUT("/profile ", controller.UpdateProfile)      // Support URL with trailing space
		userRoutes.PATCH("/profile ", controller.UpdateProfile)    // Support URL with trailing space + PATCH method
		userRoutes.DELETE("/profile", controller.DeleteOwnAccount) // Route for users to delete their own account
		userRoutes.POST("/profile/upload-image", controller.UploadProfileImage)

		// Email update with verification
		userRoutes.PUT("/email", controller.UpdateEmail)
		userRoutes.POST("/email/verify", controller.VerifyEmailUpdate)
	}

	// SuperAdmin only routes for role management
	adminRoutes := router.Group("/api/admin/users")
	adminRoutes.Use(middleware.AuthMiddleware(tokenMaker))
	adminRoutes.Use(middleware.RoleMiddleware("SUPERADMIN"))
	{
		// Get list of all users with pagination
		// GET /api/admin/users?page=1&limit=10
		adminRoutes.GET("", controller.ListUsers)

		// Create new user by admin
		// POST /api/admin/users?verified=true
		adminRoutes.POST("", controller.CreateUserByAdmin)

		// Get user profile by ID
		// GET /api/admin/users/:id
		adminRoutes.GET("/:id", controller.GetUserProfileById)
		// Update user data (including email and username) by ID
		// PUT /api/admin/users/:id
		adminRoutes.PUT("/:id", middleware.RequestLoggerMiddleware(), controller.UpdateUserData)
		// Change user role (SUPERADMIN only)
		// Support both PUT and POST methods for better compatibility
		adminRoutes.PUT("/:id/role", controller.ChangeUserRole)

		// Verify user email by ID (SUPERADMIN only)
		adminRoutes.POST("/:id/verify-email", controller.VerifyEmailById)

		// Soft Delete user by ID
		adminRoutes.DELETE("/:id", controller.DeleteUserByID)

		// Hard Delete user by ID
		adminRoutes.DELETE("/:id/permanent", controller.HardDeleteUserByID)
	}
}
