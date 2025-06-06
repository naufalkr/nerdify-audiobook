package controller

import (
	"log"
	"microservice/user/data-layer/config"
	"microservice/user/domain-layer/service"
	"microservice/user/helpers/dto"
	"microservice/user/helpers/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantController handles HTTP requests related to tenants
type TenantController struct {
	tenantService     service.TenantService
	tokenMaker        utils.TokenMaker
	cloudinaryService *config.CloudinaryService
}

// NewTenantController creates a new instance of TenantController with dependency injection
func NewTenantController(tenantService service.TenantService, tokenMaker utils.TokenMaker, cloudinaryService *config.CloudinaryService) *TenantController {
	return &TenantController{
		tenantService:     tenantService,
		tokenMaker:        tokenMaker,
		cloudinaryService: cloudinaryService,
	}
}

// UpdateProfile handles tenant profile updates
func (c *TenantController) UpdateProfile(ctx *gin.Context) {
	// Get tenant ID from context
	tenantID, exists := ctx.Get("tenant_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse request
	var req dto.TenantProfileUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update profile
	if err := c.tenantService.UpdateTenantProfile(ctx.Request.Context(), tenantID.(uuid.UUID), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ListTenants handles listing all tenants (superadmin only)
func (c *TenantController) ListTenants(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	search := ctx.Query("search")
	isActiveStr := ctx.Query("is_active")
	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	req := dto.TenantListRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		IsActive: isActive,
	}

	// Get tenants
	tenants, total, err := c.tenantService.ListTenants(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": tenants,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// DirectInviteUserToTenant invites a user to a tenant directly with user ID (SuperAdmin only)
func (c *TenantController) DirectInviteUserToTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	userID := ctx.Param("userID")
	superadminID, _ := ctx.Get("userID") // Set by auth middleware

	// Call service method (need to implement this method in service layer)
	err := c.tenantService.DirectInviteUserToTenant(ctx.Request.Context(), tenantID, userID, superadminID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User invited successfully by SuperAdmin"})
}

// GetTenantDetailsByID gets detailed tenant information (for Admin & SuperAdmin)
func (c *TenantController) GetTenantDetailsByID(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	adminID, _ := ctx.Get("userID") // Set by auth middleware

	// Call service method
	tenant, err := c.tenantService.GetTenantDetailsByID(ctx.Request.Context(), tenantID, adminID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// GetUserCurrentTenant gets the current tenant for a user using token
func (c *TenantController) GetUserCurrentTenant(ctx *gin.Context) {
	userID, _ := ctx.Get("userID") // Set by auth middleware

	// Call service method
	tenant, err := c.tenantService.GetUserCurrentTenant(ctx.Request.Context(), userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// The following are admin functions that were moved from TenantAdminController

// PromoteToAdmin promotes a user to ADMIN role within a tenant
func (c *TenantController) PromoteToAdmin(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	userID := ctx.Param("userID")
	promotedBy, _ := ctx.Get("userID") // Set by auth middleware

	// Implement the service call
	err := c.tenantService.PromoteToAdmin(ctx.Request.Context(), tenantID, userID, promotedBy.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}

// DemoteFromAdmin demotes a user from ADMIN role to regular user
func (c *TenantController) DemoteFromAdmin(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	userID := ctx.Param("userID")
	demotedBy, _ := ctx.Get("userID") // Set by auth middleware

	err := c.tenantService.DemoteFromAdmin(ctx.Request.Context(), tenantID, userID, demotedBy.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User demoted from admin successfully"})
}

// UpdateTenant handles tenant update by an admin
func (c *TenantController) UpdateTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	userID, _ := ctx.Get("userID") // Set by auth middleware

	var req dto.TenantUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service method
	tenant, err := c.tenantService.UpdateTenantByAdmin(ctx.Request.Context(), tenantID, userID.(string), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// InviteUserToTenant invites a user to a tenant by an admin
func (c *TenantController) InviteUserToTenant(ctx *gin.Context) {
	var req dto.InviteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminUserID, _ := ctx.Get("userID") // Set by auth middleware

	// Call service method
	err := c.tenantService.InviteUserByAdmin(ctx.Request.Context(), adminUserID.(string), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User invited successfully"})
}

// RemoveUserFromTenant removes a user from a tenant by an admin
func (c *TenantController) RemoveUserFromTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	userID := ctx.Param("userID")
	adminUserID, _ := ctx.Get("userID") // Extract admin user ID from context

	// Validate parameters
	if tenantID == "" || userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tenantID and userID are required"})
		return
	}

	// Call service method
	err := c.tenantService.RemoveUserFromTenant(ctx.Request.Context(), tenantID, userID, adminUserID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User removed from tenant successfully"})
}

// GetTenantUsers gets all users in a tenant with pagination
func (c *TenantController) GetTenantUsers(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")

	// Parse pagination parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Call service method
	users, total, err := c.tenantService.GetTenantUsers(ctx.Request.Context(), tenantID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// UpdateTenantContact updates the contact email of a tenant
func (c *TenantController) UpdateTenantContact(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")

	var req struct {
		ContactEmail string `json:"contactEmail" binding:"required,email"`
		ContactPhone string `json:"contactPhone"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service method
	err := c.tenantService.UpdateTenantContact(ctx.Request.Context(), tenantID, req.ContactEmail, req.ContactPhone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant contact updated successfully"})
}

// UpdateTenantLogo updates the logo/profile picture of a tenant
func (c *TenantController) UpdateTenantLogo(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")

	// Get the uploaded file
	file, err := ctx.FormFile("logo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	log.Printf("Uploading tenant logo: %s, size: %d bytes, content-type: %s", file.Filename, file.Size, file.Header.Get("Content-Type"))

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
	result, err := c.cloudinaryService.UploadFile(ctx.Request.Context(), file, "tenants", transformations)
	if err != nil {
		log.Printf("Error uploading file to Cloudinary: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Call service method
	err = c.tenantService.UpdateTenantLogo(ctx.Request.Context(), tenantID, result.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tenant logo updated successfully",
		"logoUrl": result.URL,
	})
}

// UpdateTenantSubscription updates the subscription details of a tenant
func (c *TenantController) UpdateTenantSubscription(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")

	var req struct {
		SubscriptionPlan      string `json:"SubscriptionPlan" binding:"required,oneof=Basic Premium Enterprise"`
		MaxUsers              int    `json:"max_users"`
		SubscriptionStartDate string `json:"subscriptionStartDate"`
		SubscriptionEndDate   string `json:"subscriptionEndDate"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set subscription dates automatically if not provided
	now := time.Now()
	startDate := now.Format("2006-01-02")
	endDate := now.AddDate(1, 0, 0).Format("2006-01-02") // 1 year from now

	// Use provided dates if they exist
	if req.SubscriptionStartDate != "" {
		startDate = req.SubscriptionStartDate
	}
	if req.SubscriptionEndDate != "" {
		endDate = req.SubscriptionEndDate
	}

	// Set default max users by plan if not specified
	maxUsers := req.MaxUsers
	if maxUsers <= 0 {
		switch req.SubscriptionPlan {
		case "Basic":
			maxUsers = 50
		case "Premium":
			maxUsers = 100
		case "Enterprise":
			maxUsers = 500
		}
	}

	// Call service method
	err := c.tenantService.UpdateTenantSubscription(
		ctx.Request.Context(),
		tenantID,
		req.SubscriptionPlan,
		startDate,
		endDate,
		maxUsers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant subscription updated successfully"})
}

// CreateTenant creates a new tenant
func (c *TenantController) CreateTenant(ctx *gin.Context) {
	// Parse request
	var req dto.TenantCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service method
	tenant, err := c.tenantService.CreateTenant(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": tenant})
}

// GetAllTenants retrieves all tenants (for superadmin)
func (c *TenantController) GetAllTenants(ctx *gin.Context) {
	// Call service method
	tenants, err := c.tenantService.GetAllTenants(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenants})
}

// DeleteTenant soft-deletes a tenant
func (c *TenantController) DeleteTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")

	// Call service method
	err := c.tenantService.DeleteTenant(ctx.Request.Context(), tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}

// GetTenantByID retrieves tenant details by ID
func (c *TenantController) GetTenantByID(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	userID, _ := ctx.Get("userID")

	// Call service method
	tenant, err := c.tenantService.GetTenantByID(ctx.Request.Context(), tenantID, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// GetCurrentTenantDetails retrieves detailed information about the current tenant
func (c *TenantController) GetCurrentTenantDetails(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	userID, _ := ctx.Get("userID")

	// Call service method
	tenant, err := c.tenantService.GetCurrentTenantDetails(ctx.Request.Context(), tenantID, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// UpdateSubscription updates the tenant's subscription
func (c *TenantController) UpdateSubscription(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")

	var req dto.SubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service method
	tenant, err := c.tenantService.UpdateSubscription(ctx.Request.Context(), tenantID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// GetUserTenants retrieves tenants that the current user belongs to
func (c *TenantController) GetUserTenants(ctx *gin.Context) {
	// Get current user ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call service method
	tenants, err := c.tenantService.GetUserTenants(ctx.Request.Context(), userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenants})
}

// UpdateCurrentTenant handles tenant update by an admin using their token to identify the tenant
func (c *TenantController) UpdateCurrentTenant(ctx *gin.Context) {
	userID, _ := ctx.Get("userID") // Set by auth middleware

	var req dto.TenantUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service method that will determine the tenant from the admin's token
	tenant, err := c.tenantService.UpdateCurrentTenant(ctx.Request.Context(), userID.(string), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// GetCurrentTenantUsers gets all users in the admin's current tenant with pagination
func (c *TenantController) GetCurrentTenantUsers(ctx *gin.Context) {
	userID, _ := ctx.Get("userID") // Set by auth middleware

	// Parse pagination parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Call service method that determines tenant from admin's token
	users, total, err := c.tenantService.GetCurrentTenantUsers(ctx.Request.Context(), userID.(string), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// InviteUserToCurrentTenant invites a user to the admin's current tenant
func (c *TenantController) InviteUserToCurrentTenant(ctx *gin.Context) {
	var req dto.InviteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminUserID, _ := ctx.Get("userID") // Set by auth middleware

	// Call service method that determines tenant from admin's token
	err := c.tenantService.InviteUserToCurrentTenant(ctx.Request.Context(), adminUserID.(string), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User invited successfully"})
}

// RemoveUserFromCurrentTenant removes a user from the admin's current tenant
func (c *TenantController) RemoveUserFromCurrentTenant(ctx *gin.Context) {
	userID := ctx.Param("userID")
	adminUserID, _ := ctx.Get("userID") // Extract admin user ID from context

	// Validate parameters
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	// Call service method that determines tenant from admin's token
	err := c.tenantService.RemoveUserFromCurrentTenant(ctx.Request.Context(), userID, adminUserID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User removed from tenant successfully"})
}

// UpdateCurrentSubscription updates the subscription for the admin's current tenant
func (c *TenantController) UpdateCurrentSubscription(ctx *gin.Context) {
	userID, _ := ctx.Get("userID") // Set by auth middleware

	var req dto.SubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service method that determines tenant from admin's token
	tenant, err := c.tenantService.UpdateCurrentSubscription(ctx.Request.Context(), userID.(string), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tenant})
}

// UpdateCurrentTenantContact updates the contact information for the admin's current tenant
func (c *TenantController) UpdateCurrentTenantContact(ctx *gin.Context) {
	userID, _ := ctx.Get("userID") // Set by auth middleware

	var req struct {
		ContactEmail string `json:"contactEmail" binding:"required,email"`
		ContactPhone string `json:"contactPhone"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service method that determines tenant from admin's token
	err := c.tenantService.UpdateCurrentTenantContact(ctx.Request.Context(), userID.(string), req.ContactEmail, req.ContactPhone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant contact updated successfully"})
}

// UpdateCurrentTenantLogo updates the logo for the admin's current tenant
func (c *TenantController) UpdateCurrentTenantLogo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID") // Set by auth middleware
	if !exists || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the uploaded file
	file, err := ctx.FormFile("logo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	log.Printf("Uploading tenant logo: %s, size: %d bytes, content-type: %s", file.Filename, file.Size, file.Header.Get("Content-Type"))

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
	result, err := c.cloudinaryService.UploadFile(ctx.Request.Context(), file, "tenants", transformations)
	if err != nil {
		log.Printf("Error uploading file to Cloudinary: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Update tenant logo with Cloudinary image URL
	err = c.tenantService.UpdateCurrentTenantLogo(ctx.Request.Context(), userID.(string), result.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tenant logo updated successfully",
		"logoUrl": result.URL,
	})
}

// UpdateTenant handles tenant update by an admin
