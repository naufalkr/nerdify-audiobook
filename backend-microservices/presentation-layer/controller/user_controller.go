package controller

import (
	"bytes"
	"io"
	"log"
	"microservice/user/data-layer/config"
	"microservice/user/domain-layer/middleware"
	"microservice/user/domain-layer/service"
	"microservice/user/helpers/dto"
	"microservice/user/helpers/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles HTTP requests related to users
type UserController struct {
	userService       *service.UserService
	tokenMaker        utils.TokenMaker
	cloudinaryService *config.CloudinaryService
}

// NewUserController creates a new instance of UserController
func NewUserController(userService *service.UserService, tokenMaker utils.TokenMaker, cloudinaryService *config.CloudinaryService) *UserController {
	return &UserController{
		userService:       userService,
		tokenMaker:        tokenMaker,
		cloudinaryService: cloudinaryService,
	}
}

// SetupRoutes registers all the routes for the user controller
func (c *UserController) SetupRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		// Public routes (no authentication required)
		auth.POST("/register", c.Register)
		auth.POST("/login", c.Login)
		auth.POST("/verify-email", c.VerifyEmail)
		auth.POST("/resend-verification-email", c.ResendVerificationEmailWithRateLimiting) // Updated to use rate-limited version
		auth.POST("/forgot-password", c.ForgotPassword)
		auth.POST("/reset-password", c.ResetPassword)
		auth.POST("/refresh-token", c.RefreshToken)
		auth.GET("/verify-email-link", c.VerifyEmailByLink)
		auth.POST("/verify-email-link", c.VerifyEmailByLink)
	}

	// Protected routes (authentication required)
	user := router.Group("/api/users")
	user.Use(middleware.AuthMiddleware(c.tokenMaker))
	{
		user.POST("/logout", c.Logout)
		user.GET("/profile", c.GetUserProfile)
		user.GET("/:id/profile", middleware.RoleCheckMiddleware([]string{"SUPERADMIN"}), c.GetUserProfileById)
		user.DELETE("/:id", middleware.RoleCheckMiddleware([]string{"SUPERADMIN"}), c.DeleteUserByID)
		user.POST("/update-profile", c.UpdateProfile)
		user.POST("/upload-profile-image", c.UploadProfileImage) // Changed from /:id to use token
		user.GET("/list", c.ListUsers)

		// User role management route
		user.PUT("/:id/change-role", middleware.RoleCheckMiddleware([]string{"SUPERADMIN"}), c.ChangeUserRole)
		// ... other user routes
	}
}

// Register handles user registration requests
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.userService.Register(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrInvalidEmail, dto.ErrUserAlreadyExist:
			statusCode = http.StatusBadRequest
		case dto.ErrCreateUserFailed:
			statusCode = http.StatusInternalServerError
		case dto.ErrSendEmailFailed:
			// Continue with registration even if email sending fails
			ctx.JSON(http.StatusCreated, gin.H{
				"user":    resp,
				"warning": "Email verification could not be sent. Please request a new verification email."})
			return
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": dto.MESSAGE_REGISTER_SUCCESS,
		"data":    resp,
	})
}

// Login handles user login requests
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that at least one of email or username is provided
	if req.Email == "" && req.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Either email or username is required"})
		return
	}

	resp, err := c.userService.Login(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrUserNotFound, dto.ErrPasswordNotMatch:
			statusCode = http.StatusUnauthorized
		case dto.ErrInvalidEmail:
			statusCode = http.StatusBadRequest
		case dto.ErrUserNotVerified:
			// For unverified users, return 200 but with a special message
			ctx.JSON(http.StatusOK, gin.H{
				"error": "User email not verified",
				"user":  resp,
			})
			return
		case dto.ErrLoginFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_LOGIN_SUCCESS,
		"data":    resp,
	})
}

// VerifyEmail handles email verification requests using OTP
func (c *UserController) VerifyEmail(ctx *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.userService.VerifyEmail(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrOTPNotMatch, dto.ErrOTPExpired:
			statusCode = http.StatusBadRequest
		case dto.ErrTooManyAttempts:
			statusCode = http.StatusTooManyRequests
		case dto.ErrUpdateUserFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_VERIFY_EMAIL_SUCCESS,
		"data":    resp,
	})
}

// VerifyEmailByLink handles email verification requests using the link method
func (c *UserController) VerifyEmailByLink(ctx *gin.Context) {
	var req dto.VerifyEmailLinkRequest

	// Get query parameters from the URL for GET requests
	if ctx.Request.Method == "GET" {
		req.Email = ctx.Query("email")
		req.Token = ctx.Query("token")
	} else {
		// Handle POST requests with JSON body
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Validate required fields
	if req.Email == "" || req.Token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and token are required"})
		return
	}

	resp, err := c.userService.VerifyEmailByLink(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Verifikasi email berhasil",
		"data":    resp,
	})
}

// ResendVerificationEmail handles resending verification email requests
func (c *UserController) ResendVerificationEmail(ctx *gin.Context) {
	var req dto.SendVerificationEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.userService.ResendVerificationEmail(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrUserAlreadyVerified:
			statusCode = http.StatusBadRequest
		case dto.ErrSendEmailFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_RESEND_VERIFICATION_EMAIL_SUCCESS,
	})
}

// ResendVerificationEmailWithRateLimiting handles resending verification email requests with rate limiting
func (c *UserController) ResendVerificationEmailWithRateLimiting(ctx *gin.Context) {
	var req dto.SendVerificationEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.userService.ResendVerificationEmailWithRateLimiting(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch {
		case err == dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case err == dto.ErrUserAlreadyVerified:
			statusCode = http.StatusBadRequest
		case err == dto.ErrSendEmailFailed:
			statusCode = http.StatusInternalServerError
		case err == dto.ErrTooManyResendAttempts:
			statusCode = http.StatusTooManyRequests
		case strings.Contains(err.Error(), dto.ErrInCooldownPeriod.Error()):
			// Special case for cooldown error with minutes information
			statusCode = http.StatusTooManyRequests
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_RESEND_VERIFICATION_EMAIL_SUCCESS,
	})
}

// ForgotPassword handles forgot password requests
func (c *UserController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.userService.ForgotPassword(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrUserNotFound:
			// For security, don't reveal that user doesn't exist
			ctx.JSON(http.StatusOK, gin.H{
				"message": dto.MESSAGE_FORGOT_PASSWORD_SUCCESS,
			})
			return
		case dto.ErrSendEmailFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_FORGOT_PASSWORD_SUCCESS,
	})
}

// ResetPassword handles password reset requests
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.userService.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrTokenInvalid:
			statusCode = http.StatusBadRequest
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrUpdateUserFailed:
			statusCode = http.StatusInternalServerError
		case dto.ErrSamePassword:
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_RESET_PASSWORD_SUCCESS,
	})
}

// RefreshToken handles token refresh requests
func (c *UserController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.userService.RefreshToken(ctx.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrTokenInvalid:
			statusCode = http.StatusUnauthorized
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrLoginFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_REFRESH_TOKEN_SUCCESS,
		"data":    resp,
	})
}

// Logout handles user logout requests
func (c *UserController) Logout(ctx *gin.Context) {
	// Ambil userID dari konteks yang diatur oleh AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := c.userService.Logout(ctx.Request.Context(), userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_LOGOUT_SUCCESS,
	})
}

// DeleteUserByID handles user deletion requests
func (c *UserController) DeleteUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Log the request
	log.Printf("DeleteUserByID: Request received to delete user with ID: %s", userID)

	// Validate the UUID format
	_, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("DeleteUserByID: Invalid user ID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = c.userService.DeleteUserByID(ctx.Request.Context(), userID)
	if err != nil {
		log.Printf("DeleteUserByID: Error during user deletion: %v", err)

		// Determine appropriate status code based on error
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" || strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	log.Printf("DeleteUserByID: Successfully deleted user with ID: %s", userID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_DELETE_USER_SUCCESS,
	})
}

// HardDeleteUserByID handles permanent user deletion requests
func (c *UserController) HardDeleteUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Log the request
	log.Printf("HardDeleteUserByID: Request received to permanently delete user with ID: %s", userID)

	// Validate the UUID format
	_, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("HardDeleteUserByID: Invalid user ID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = c.userService.HardDeleteUserByID(ctx.Request.Context(), userID)
	if err != nil {
		log.Printf("HardDeleteUserByID: Error during permanent user deletion: %v", err)

		// Determine appropriate status code based on error
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" || strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HardDeleteUserByID: Successfully permanently deleted user with ID: %s", userID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus pengguna secara permanen",
	})
}

// GetUserProfile retrieves the profile of the currently authenticated user
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	// Get userID from the context set by AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user profile using the token instead of passing the ID directly
	profile, err := c.userService.GetUserProfileFromToken(ctx.Request.Context(), userID.(string))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == dto.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    profile,
	})
}

// GetUserProfileById retrieves a user profile by ID (only for SUPERADMIN)
func (c *UserController) GetUserProfileById(ctx *gin.Context) {
	// The CheckRoleMiddleware already verified this is a SUPERADMIN
	targetUserID := ctx.Param("id")
	if targetUserID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get the admin's ID from the token context instead of expecting it as a parameter
	adminID, exists := ctx.Get("userID")
	if !exists || adminID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - invalid token"})
		return
	}

	// NOTE: GetUserProfileById method doesn't exist in UserService, need to be implemented
	// Using GetUserProfileFromToken as temporary workaround
	profile, err := c.userService.GetUserProfileFromToken(ctx.Request.Context(), targetUserID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == dto.ErrUserNotFound {
			statusCode = http.StatusNotFound
		} else if err == dto.ErrUnauthorized {
			statusCode = http.StatusForbidden
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    profile,
	})
}

// UpdateProfile handles user profile update requests
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// Get userID from the context set by AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Printf("DEBUG: Updating profile for user ID: %s", userID)

	// Log the raw request body for debugging
	requestBody, _ := io.ReadAll(ctx.Request.Body)
	log.Printf("DEBUG: Raw request body: %s", string(requestBody))

	// Restore the request body for further processing
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	var req dto.UserProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("DEBUG: JSON binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("DEBUG: Profile update request: username=%s, email=%s, fullName=%s",
		req.UserName, req.Email, req.FullName)

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		log.Printf("DEBUG: UUID parsing error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = c.userService.UpdateUserProfile(ctx.Request.Context(), userUUID, req)
	if err != nil {
		log.Printf("DEBUG: Update profile error: %v", err)
		statusCode := http.StatusInternalServerError
		if err == dto.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

// UploadProfileImage handles user profile image upload requests
func (c *UserController) UploadProfileImage(ctx *gin.Context) {
	// Get userID from the token context instead of URL parameter
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file, err := ctx.FormFile("profile_image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	log.Printf("Uploading file: %s, size: %d bytes, content-type: %s", file.Filename, file.Size, file.Header.Get("Content-Type"))

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images are allowed"})
		return
	}

	// Define image transformations for Cloudinary
	transformations := &config.ImageTransformation{
		Width:   500,    // Resize to 500px width
		Height:  500,    // Resize to 500px height
		Crop:    "fill", // Fill mode for cropping
		Quality: 80,     // 80% quality
		Format:  "auto", // Auto format (webp if supported)
	}

	// Upload directly to Cloudinary
	result, err := c.cloudinaryService.UploadFile(ctx.Request.Context(), file, "profiles", transformations)
	if err != nil {
		log.Printf("Error uploading file to Cloudinary: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Update user profile with Cloudinary image URL
	updatedUser, err := c.userService.UpdateProfileImage(ctx.Request.Context(), userID.(string), result.URL)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == dto.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile image uploaded successfully",
		"data":    updatedUser,
	})
}

func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

// ListUsers handles listing all users (superadmin only)
func (c *UserController) ListUsers(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	search := ctx.Query("search")
	role := ctx.Query("role")
	status := ctx.Query("status")

	req := dto.UserListRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Role:     role,
		Status:   status,
	}

	// Get users
	users, total, err := c.userService.ListUsers(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ChangeUserRole handles requests to change a user's role
func (c *UserController) ChangeUserRole(ctx *gin.Context) {
	// Get the target user ID from the URL parameter
	userID := ctx.Param("id")
	// Get the ID of the user making the change (must be SUPERADMIN, already verified by middleware)
	changerID, _ := ctx.Get("userID")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format. Role is required",
		})
		return
	}

	// Call the service to change the role
	updatedUser, err := c.userService.ChangeUserRole(ctx.Request.Context(), userID, req.Role, changerID.(string))
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch err {
		case dto.ErrUserNotFound, dto.ErrRoleNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrUnauthorized:
			statusCode = http.StatusForbidden
		case dto.ErrUpdateUserFailed, dto.ErrUpdateRoleFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"data":    updatedUser,
	})
}

// UpdateUserData handles admin requests to update user data including username and email
func (c *UserController) UpdateUserData(ctx *gin.Context) {
	// Get the target user ID from the URL parameter
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get the ID of the admin making the change (must be SUPERADMIN, already verified by middleware)
	adminID, exists := ctx.Get("userID")
	if !exists || adminID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - invalid token"})
		return
	}

	// Parse request body
	var req dto.UserDataUpdateRequest

	// Log the raw request body to debug
	bodyBytes, _ := ctx.GetRawData()
	log.Printf("UpdateUserData raw request body: %s", string(bodyBytes))

	// We need to set the body back since GetRawData consumes it
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the parsed DTO to debug
	log.Printf("UpdateUserData parsed DTO: username=%s, email=%s, fullName=%s",
		req.Username, req.Email, req.FullName)

	// Call the service to update the user data
	updatedProfile, err := c.userService.UpdateUserData(ctx.Request.Context(), userID, adminID.(string), req)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch err {
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrUnauthorized:
			statusCode = http.StatusForbidden
		case dto.ErrEmailAlreadyInUse, dto.ErrUsernameAlreadyInUse:
			statusCode = http.StatusConflict
		case dto.ErrUpdateUserFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_UPDATE_USER_DATA_SUCCESS,
		"data":    updatedProfile,
	})
}

// VerifyEmailById handles admin requests to verify a user's email directly by user ID
func (c *UserController) VerifyEmailById(ctx *gin.Context) {
	// Get the user ID from the URL parameter
	userID := ctx.Param("id")

	// Validate the UUID format
	_, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Call the service to verify the user's email
	resp, err := c.userService.VerifyEmailById(ctx.Request.Context(), userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrUpdateUserFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email pengguna berhasil diverifikasi oleh admin",
		"data":    resp,
	})
}

// UpdateEmail handles user email update requests
func (c *UserController) UpdateEmail(ctx *gin.Context) {
	// Get userID from context set by AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.EmailUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = c.userService.UpdateUserEmail(ctx.Request.Context(), userUUID, req.Email)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrUserAlreadyExist:
			statusCode = http.StatusBadRequest
		case dto.ErrSendEmailFailed:
			statusCode = http.StatusInternalServerError
		case dto.ErrUpdateUserFailed:
			statusCode = http.StatusInternalServerError
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Verification code sent to your new email. Please verify to complete the update.",
	})
}

// VerifyEmailUpdate handles verification for email updates
func (c *UserController) VerifyEmailUpdate(ctx *gin.Context) {
	// Get userID from context set by AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.EmailUpdateVerificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = c.userService.VerifyEmailUpdate(ctx.Request.Context(), userUUID, req.OTP)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case dto.ErrOTPNotMatch, dto.ErrOTPExpired:
			statusCode = http.StatusBadRequest
		case dto.ErrTooManyAttempts:
			statusCode = http.StatusTooManyRequests
		case dto.ErrUpdateUserFailed:
			statusCode = http.StatusInternalServerError
		default:
			// For custom errors not in the predefined list
			statusCode = http.StatusBadRequest
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email updated successfully",
	})
}

// UpdateFullName handles user full name update
func (c *UserController) UpdateFullName(ctx *gin.Context) {
	// Get userID from context set by AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		FullName string `json:"full_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Create a DTO with only the fullName field
	updateRequest := dto.UserProfileUpdateRequest{
		FullName: req.FullName,
	}

	err = c.userService.UpdateUserProfile(ctx.Request.Context(), userUUID, updateRequest)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == dto.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Full name updated successfully",
	})
}

// UpdateAddress handles user address update
func (c *UserController) UpdateAddress(ctx *gin.Context) {
	// Get userID from context set by AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Alamat    string  `json:"alamat" binding:"required"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Create a DTO with only the address fields
	updateRequest := dto.UserProfileUpdateRequest{
		Alamat:    req.Alamat,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	err = c.userService.UpdateUserProfile(ctx.Request.Context(), userUUID, updateRequest)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == dto.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Address updated successfully",
	})
}

// CreateUserByAdmin handles admin requests to create new users
func (c *UserController) CreateUserByAdmin(ctx *gin.Context) {
	// Parse request body
	var req dto.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get isVerified flag from query param (default to false)
	isVerifiedStr := ctx.DefaultQuery("verified", "false")
	isVerified := isVerifiedStr == "true"

	// Get admin ID from context for audit purposes
	adminID, exists := ctx.Get("userID")
	if !exists || adminID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Log action
	log.Printf("CreateUserByAdmin: Admin %s is creating a new user with email: %s", adminID, req.Email)

	// Call service to create user
	resp, err := c.userService.CreateUserByAdmin(ctx.Request.Context(), req, isVerified)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case dto.ErrInvalidEmail, dto.ErrUserAlreadyExist:
			statusCode = http.StatusBadRequest
		case dto.ErrCreateUserFailed:
			statusCode = http.StatusInternalServerError
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat oleh admin",
		"data":    resp,
	})
}

// DeleteOwnAccount handles requests to delete the authenticated user's own account
func (c *UserController) DeleteOwnAccount(ctx *gin.Context) {
	// Get userID from context set by AuthMiddleware
	userID, exists := ctx.Get("userID")
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Log the request
	log.Printf("DeleteOwnAccount: Request received to delete user's own account with ID: %s", userID)

	// Parse userID to ensure it's a valid UUID
	_, err := uuid.Parse(userID.(string))
	if err != nil {
		log.Printf("DeleteOwnAccount: Invalid user ID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Use the existing service method to delete the user
	err = c.userService.DeleteUserByID(ctx.Request.Context(), userID.(string))
	if err != nil {
		log.Printf("DeleteOwnAccount: Error during user deletion: %v", err)

		// Determine appropriate status code based on error
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" || strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	log.Printf("DeleteOwnAccount: Successfully deleted user's own account with ID: %s", userID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Akun berhasil dihapus",
	})
}
