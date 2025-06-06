package controller

import (
	"microservice/user/domain-layer/service"
	"microservice/user/helpers/dto"
	"microservice/user/helpers/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantAPIController handles HTTP requests related to external tenant API access and user tenant context
type TenantAPIController struct {
	tenantService     service.TenantService
	userTenantService service.UserTenantContextService
	tokenMaker        utils.TokenMaker
}

// NewTenantAPIController creates a new instance of TenantAPIController
func NewTenantAPIController(tenantService service.TenantService, userTenantService service.UserTenantContextService, tokenMaker utils.TokenMaker) *TenantAPIController {
	return &TenantAPIController{
		tenantService:     tenantService,
		userTenantService: userTenantService,
		tokenMaker:        tokenMaker,
	}
}

// ===== TENANT MANAGEMENT APIs =====

// ValidateTenantAccess checks if a tenant exists and if the accessing service has rights to access
func (c *TenantAPIController) ValidateTenantAccess(ctx *gin.Context) {
	tenantID := ctx.Param("id")
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
		return
	}

	// Get service identity from API key (in a real implementation, validate this API key)
	apiKey := ctx.GetHeader("X-API-Key")
	if apiKey == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
		return
	}

	// Simplified for example - in real world, validate the API key and check permissions
	// For now we'll just check if the tenant exists

	// Get tenant info
	userID := ctx.GetString("userID") // From JWT auth, if applicable
	tenant, err := c.tenantService.GetTenantByID(ctx, tenantID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	// Return only necessary information for the external service
	ctx.JSON(http.StatusOK, gin.H{
		"tenantId":   tenant.ID,
		"tenantName": tenant.Name,
		"isActive":   tenant.IsActive,
	})
}

// ListTenants provides a list of tenants for external service consumption
// This can be used by the alat service to select valid tenants for assignment
func (c *TenantAPIController) ListTenants(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Create a filter request
	isActive := true
	filter := dto.TenantListRequest{
		Page:     page,
		PageSize: limit,
		IsActive: &isActive, // Only return active tenants
	}

	// Get tenants
	tenants, total, err := c.tenantService.ListTenants(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenants"})
		return
	}

	// Return simplified tenant data for external service
	ctx.JSON(http.StatusOK, gin.H{
		"tenants": tenants,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetTenantById retrieves a specific tenant by ID for external service consumption
func (c *TenantAPIController) GetTenantById(ctx *gin.Context) {
	tenantID := ctx.Param("id")
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
		return
	}

	// Get tenant info (use empty string for userID to bypass user-specific permissions)
	tenant, err := c.tenantService.GetTenantByID(ctx, tenantID, "")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	// Return tenant data
	ctx.JSON(http.StatusOK, gin.H{
		"tenant": tenant,
	})
}

// ===== AUTHENTICATION & AUTHORIZATION APIs =====

// ValidateJWTToken validates a JWT token from another microservice
func (c *TenantAPIController) ValidateJWTToken(ctx *gin.Context) {
	var request dto.TokenValidationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		// Try to get token from header as fallback
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusBadRequest, dto.TokenValidationResponse{
				IsValid: false,
				Error:   "Token is required in request body or Authorization header",
			})
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		request.Token = tokenString
	}

	// Validate the token using the token maker
	claims, err := c.tokenMaker.ParseAccessToken(request.Token)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.TokenValidationResponse{
			IsValid: false,
			Error:   "Invalid or expired token",
		})
		return
	}

	// Get user details from the service
	user, err := c.tenantService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.TokenValidationResponse{
			IsValid: false,
			Error:   "User not found",
		})
		return
	}

	// Build user info response
	userInfo := &dto.ExternalUserInfoResponse{
		UserID:   claims.UserID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RoleID:   user.RoleID,
		Status:   user.Status,
		IsActive: user.Status == "active",
	}

	// Return successful validation response
	ctx.JSON(http.StatusOK, dto.TokenValidationResponse{
		IsValid:  true,
		UserInfo: userInfo,
	})
}

// GetUserInfoFromToken extracts user information from Authorization header
func (c *TenantAPIController) GetUserInfoFromToken(ctx *gin.Context) {
	// Get token from Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Extract token (remove "Bearer " prefix)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// Parse and validate the token
	claims, err := c.tokenMaker.ParseAccessToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Return user information
	ctx.JSON(http.StatusOK, gin.H{
		"userID":   claims.UserID,
		"email":    claims.Email,
		"userRole": claims.RoleName,
	})
}

// ValidateUserPermissions validates if user has specific permissions
func (c *TenantAPIController) ValidateUserPermissions(ctx *gin.Context) {
	var req struct {
		Token        string   `json:"token" binding:"required"`
		TenantID     string   `json:"tenantId"`
		RequiredRole string   `json:"requiredRole"`
		Permissions  []string `json:"permissions"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Parse and validate the token
	claims, err := c.tokenMaker.ParseAccessToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Invalid or expired token",
		})
		return
	}

	// Check role permission if required
	hasRolePermission := true
	if req.RequiredRole != "" {
		hasRolePermission = claims.RoleName == req.RequiredRole ||
			claims.RoleName == "SUPERADMIN" // SUPERADMIN has all permissions
	}

	// Check tenant access if tenantID is provided
	hasTenantAccess := true
	if req.TenantID != "" {
		// Here you would typically check if user belongs to the tenant
		// For now, we'll assume the check is done via the tenant service
		_, err := c.tenantService.GetTenantByID(ctx, req.TenantID, claims.UserID)
		if err != nil {
			hasTenantAccess = false
		}
	}

	// Return permission validation result
	ctx.JSON(http.StatusOK, gin.H{
		"valid":             hasRolePermission && hasTenantAccess,
		"hasRolePermission": hasRolePermission,
		"hasTenantAccess":   hasTenantAccess,
		"userID":            claims.UserID,
		"userRole":          claims.RoleName,
	})
}

// ValidateIsSuperAdmin validates if a user has SuperAdmin role
// This is a specialized endpoint for inter-service communication
func (c *TenantAPIController) ValidateIsSuperAdmin(ctx *gin.Context) {
	// Get token from Authorization header
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Authorization header is required",
		})
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Parse and validate the token
	claims, err := c.tokenMaker.ParseAccessToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"valid": false,
			"error": "Invalid or expired token",
		})
		return
	}

	// Check if user has SuperAdmin role
	isSuperAdmin := false
	if claims.RoleName != "" {
		isSuperAdmin = claims.RoleName == "SUPERADMIN" || claims.RoleName == "SuperAdmin"
	}

	// Return validation result
	ctx.JSON(http.StatusOK, gin.H{
		"valid":        isSuperAdmin,
		"userID":       claims.UserID,
		"userRole":     claims.RoleName,
		"isSuperAdmin": isSuperAdmin,
	})
}

// ===== USER TENANT CONTEXT APIs =====

// GetCurrentTenant retrieves the current tenant for the authenticated user
// @Summary Get current tenant
// @Description Get the current tenant context for the authenticated user
// @Tags User Tenant Context
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.UserTenantContextResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user-tenant/current [get]
func (c *TenantAPIController) GetCurrentTenant(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	result, err := c.userTenantService.GetUserCurrentTenant(ctx.Request.Context(), userUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current tenant", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Current tenant retrieved successfully", "data": result})
}

// SetCurrentTenant sets the current tenant for the authenticated user
// @Summary Set current tenant
// @Description Set the current tenant context for the authenticated user
// @Tags User Tenant Context
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.UserTenantContextRequest true "Tenant context request"
// @Success 200 {object} response.Response{data=dto.UserTenantContextResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user-tenant/current [put]
func (c *TenantAPIController) SetCurrentTenant(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req dto.UserTenantContextRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate tenant ID
	tenantID, err := uuid.Parse(req.TenantID.String())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	result, err := c.userTenantService.SetUserCurrentTenant(ctx.Request.Context(), userUUID, tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set current tenant", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Current tenant set successfully", "data": result})
}

// GetUserTenants retrieves all tenants for the authenticated user
// @Summary Get user tenants
// @Description Get all tenants that the authenticated user has access to
// @Tags User Tenant Context
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.UserTenantsListResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user-tenant/tenants [get]
func (c *TenantAPIController) GetUserTenants(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	result, err := c.userTenantService.GetUserTenants(ctx.Request.Context(), userUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tenants", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User tenants retrieved successfully", "data": result})
}

// SwitchTenant switches the current tenant for the authenticated user
// @Summary Switch tenant
// @Description Switch to a different tenant for the authenticated user
// @Tags User Tenant Context
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.SwitchTenantRequest true "Switch tenant request"
// @Success 200 {object} response.Response{data=dto.SwitchTenantResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user-tenant/switch [post]
func (c *TenantAPIController) SwitchTenant(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req dto.SwitchTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate tenant ID
	tenantID, err := uuid.Parse(req.TenantID.String())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	result, err := c.userTenantService.SwitchUserTenant(ctx.Request.Context(), userUUID, tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to switch tenant", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant switched successfully", "data": result})
}

// ValidateUserTenantAccessContext validates if user has access to a specific tenant
// @Summary Validate user tenant access
// @Description Validate if the authenticated user has access to a specific tenant
// @Tags User Tenant Context
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.UserTenantAccessValidationRequest true "Access validation request"
// @Success 200 {object} response.Response{data=dto.UserTenantAccessValidationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user-tenant/validate-access [post]
func (c *TenantAPIController) ValidateUserTenantAccessContext(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req dto.UserTenantAccessValidationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate tenant ID
	tenantID, err := uuid.Parse(req.TenantID.String())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID", "details": err.Error()})
		return
	}

	result, err := c.userTenantService.ValidateUserTenantAccess(ctx.Request.Context(), userUUID, tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate access", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Access validation completed", "data": result})
}

// GetTenantUsersContext retrieves users in a tenant (admin function)
// @Summary Get tenant users
// @Description Get all users in the current tenant (admin only)
// @Tags User Tenant Context
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=dto.TenantUsersResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user-tenant/users [get]
func (c *TenantAPIController) GetTenantUsersContext(ctx *gin.Context) {
	// Validate user is admin
	role, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role.(string) != "Admin" && role.(string) != "SuperAdmin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin role required"})
		return
	}

	// Get current tenant ID from user's tenant context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get user's current tenant
	currentTenant, err := c.userTenantService.GetUserCurrentTenant(ctx.Request.Context(), userUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current tenant", "details": err.Error()})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	result, err := c.userTenantService.GetTenantUsers(ctx.Request.Context(), currentTenant.TenantID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant users", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant users retrieved successfully", "data": result})
}

// GetTenantUsersByTenantID retrieves users in a specific tenant (super admin function)
// @Summary Get users by tenant ID
// @Description Get all users in a specific tenant (super admin only)
// @Tags User Tenant Context
// @Security ApiKeyAuth
// @Produce json
// @Param tenantId path string true "Tenant ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=dto.TenantUsersResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user-tenant/tenants/{tenantId}/users [get]
func (c *TenantAPIController) GetTenantUsersByTenantID(ctx *gin.Context) {
	// Validate user is super admin
	role, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if role.(string) != "SuperAdmin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied. SuperAdmin role required"})
		return
	}

	// Get tenant ID from path
	tenantIDStr := ctx.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	result, err := c.userTenantService.GetTenantUsers(ctx.Request.Context(), tenantID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant users", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant users retrieved successfully", "data": result})
}

// ===== BUSINESS LOGIC APIs =====

// GetTenantSubscription returns tenant subscription information
func (c *TenantAPIController) GetTenantSubscription(ctx *gin.Context) {
	tenantID := ctx.Param("id")
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
		return
	}

	// Get tenant info
	tenant, err := c.tenantService.GetTenantByID(ctx, tenantID, "")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	// Return subscription information
	ctx.JSON(http.StatusOK, gin.H{
		"tenantID":              tenant.ID,
		"subscriptionPlan":      tenant.SubscriptionPlan,
		"subscriptionStartDate": tenant.SubscriptionStartDate,
		"subscriptionEndDate":   tenant.SubscriptionEndDate,
		"isActive":              tenant.IsActive,
	})
}

// GetTenantLimits returns tenant usage limits based on subscription
func (c *TenantAPIController) GetTenantLimits(ctx *gin.Context) {
	tenantID := ctx.Param("id")
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
		return
	}

	// Get tenant info
	tenant, err := c.tenantService.GetTenantByID(ctx, tenantID, "")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	// Calculate limits based on subscription plan
	limits := map[string]interface{}{
		"maxUsers":         tenant.MaxUsers,
		"maxAssets":        getMaxAssetsForPlan(string(tenant.SubscriptionPlan)),
		"maxRentals":       getMaxRentalsForPlan(string(tenant.SubscriptionPlan)),
		"subscriptionPlan": tenant.SubscriptionPlan,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tenantID": tenant.ID,
		"limits":   limits,
	})
}

// GetTenantUsers returns users belonging to a tenant (External API version)
func (c *TenantAPIController) GetTenantUsers(ctx *gin.Context) {
	tenantID := ctx.Param("id")
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
		return
	}

	// Get tenant users with pagination (default page=1, limit=100)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 100
	}

	users, total, err := c.tenantService.GetTenantUsers(ctx, tenantID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant users"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tenantID": tenantID,
		"users":    users,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// ValidateUserTenantAccess validates if a user has access to a specific tenant (External API version)
func (c *TenantAPIController) ValidateUserTenantAccess(ctx *gin.Context) {
	tenantID := ctx.Param("id")
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
		return
	}

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Check if user has access to tenant - simplified validation
	// Try to get tenant by ID with the user ID to check access
	_, err := c.tenantService.GetTenantByID(ctx, tenantID, req.UserID)
	hasAccess := err == nil

	ctx.JSON(http.StatusOK, gin.H{
		"userID":    req.UserID,
		"tenantID":  tenantID,
		"hasAccess": hasAccess,
	})
}

// GetUserTenantsExternal returns all tenants a user belongs to (External API version)
func (c *TenantAPIController) GetUserTenantsExternal(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get user tenants
	tenants, err := c.tenantService.GetUserTenants(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tenants"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"userID":  userID,
		"tenants": tenants,
	})
}

// Helper functions for subscription limits
func getMaxAssetsForPlan(plan string) int {
	switch plan {
	case "basic":
		return 10
	case "premium":
		return 50
	case "enterprise":
		return 200
	default:
		return 5
	}
}

func getMaxRentalsForPlan(plan string) int {
	switch plan {
	case "basic":
		return 5
	case "premium":
		return 25
	case "enterprise":
		return 100
	default:
		return 2
	}
}
