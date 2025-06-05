// filepath: /home/alvn/Documents/playground/kp/clean_architecture/user_management/domain-layer/service/tenant_service.go
package service

import (
	"context"
	"fmt"
	"log"
	"microservice/user/data-layer/entity"
	"microservice/user/data-layer/repository"
	"microservice/user/helpers/dto"
	"microservice/user/helpers/utils"

	"strings"
	"time"

	"github.com/google/uuid"
)

// TenantService defines the interface for tenant-related operations
type TenantService interface {
	CreateTenant(ctx context.Context, req dto.TenantCreateRequest) (dto.TenantResponse, error)
	GetAllTenants(ctx context.Context) ([]dto.TenantResponse, error)
	GetTenantByID(ctx context.Context, tenantID string, userID string) (dto.TenantResponse, error)
	UpdateTenant(ctx context.Context, tenantID string, req dto.TenantUpdateRequest) (dto.TenantResponse, error)
	DeleteTenant(ctx context.Context, tenantID string) error
	UpdateSubscription(ctx context.Context, tenantID string, req dto.SubscriptionRequest) (dto.TenantResponse, error)
	InviteUserToTenant(ctx context.Context, req dto.InviteUserRequest) error
	GetUserTenants(ctx context.Context, userID string) ([]dto.TenantResponse, error)
	GetCurrentTenantDetails(ctx context.Context, tenantID string, userID string) (dto.TenantResponse, error)

	// Admin management functions
	PromoteToAdmin(ctx context.Context, tenantID, userID, promotedBy string) error
	DemoteFromAdmin(ctx context.Context, tenantID, userID, demotedBy string) error

	// Admin tenant management functions
	UpdateTenantByAdmin(ctx context.Context, tenantID, adminUserID string, req dto.TenantUpdateRequest) (dto.TenantResponse, error)
	InviteUserByAdmin(ctx context.Context, adminUserID string, req dto.InviteUserRequest) error
	RemoveUserFromTenant(ctx context.Context, tenantID, userID, adminUserID string) error
	GetTenantUsers(ctx context.Context, tenantID string, page, limit int) ([]dto.UserResponse, int, error)

	// New methods for enhanced tenant detail management
	DirectInviteUserToTenant(ctx context.Context, tenantID, userID string, superadminID string) error
	GetTenantDetailsByID(ctx context.Context, tenantID, adminID string) (dto.TenantResponse, error)
	GetUserCurrentTenant(ctx context.Context, userID string) (dto.TenantResponse, error)

	// Profile and contact management
	UpdateTenantContact(ctx context.Context, tenantID, contactEmail, contactPhone string) error
	UpdateTenantLogo(ctx context.Context, tenantID, logoURL string) error
	UpdateTenantSubscription(ctx context.Context, tenantID, plan, startDate, endDate string, maxUsers int) error
	UpdateTenantProfile(ctx context.Context, tenantID uuid.UUID, req dto.TenantProfileUpdateRequest) error

	// Listing and filtering
	ListTenants(ctx context.Context, req dto.TenantListRequest) ([]dto.TenantListResponse, int64, error)

	// Token-based admin tenant management functions
	UpdateCurrentTenant(ctx context.Context, adminUserID string, req dto.TenantUpdateRequest) (dto.TenantResponse, error)
	GetCurrentTenantUsers(ctx context.Context, adminUserID string, page, limit int) ([]dto.UserResponse, int, error)
	InviteUserToCurrentTenant(ctx context.Context, adminUserID string, req dto.InviteUserRequest) error
	RemoveUserFromCurrentTenant(ctx context.Context, userID, adminUserID string) error
	UpdateCurrentSubscription(ctx context.Context, adminUserID string, req dto.SubscriptionRequest) (dto.TenantResponse, error)
	UpdateCurrentTenantContact(ctx context.Context, adminUserID string, contactEmail, contactPhone string) error
	UpdateCurrentTenantLogo(ctx context.Context, adminUserID string, logoURL string) error

	// New method to remove a user from all tenants they are part of
	RemoveUserFromAllTenants(ctx context.Context, userID string) error

	// New method to remove superadmin from all tenants
	RemoveSuperadminFromAllTenants(ctx context.Context) error

	// External API methods
	GetUserByID(ctx context.Context, userID string) (dto.ExternalUserInfoResponse, error)
}

// tenantService implements the TenantService interface
type tenantService struct {
	repo        repository.TenantRepository
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	emailSender utils.EmailSender
}

// NewTenantService creates a new instance of the tenant service
func NewTenantService(repo repository.TenantRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository, emailSender utils.EmailSender) TenantService {
	return &tenantService{
		repo:        repo,
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		emailSender: emailSender,
	}
}

// CreateTenant creates a new tenant
func (s *tenantService) CreateTenant(ctx context.Context, req dto.TenantCreateRequest) (dto.TenantResponse, error) {
	tenant := entity.Tenant{
		ID:                    uuid.New(),
		Name:                  req.Name,
		Description:           req.Description,
		LogoURL:               req.LogoURL,
		ContactEmail:          req.ContactEmail,
		ContactPhone:          req.ContactPhone,
		MaxUsers:              15,          // Default max users for new tenants
		SubscriptionPlan:      "",          // No subscription by default
		SubscriptionStartDate: time.Time{}, // No subscription by default
		SubscriptionEndDate:   time.Time{}, // No subscription by default
		IsActive:              true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	createdTenant, err := s.repo.Create(ctx, nil, tenant)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrCreateTenantFailed
	}

	return mapTenantToResponse(createdTenant), nil
}

// GetAllTenants gets all tenants (for SuperAdmin only)
func (s *tenantService) GetAllTenants(ctx context.Context) ([]dto.TenantResponse, error) {
	tenants, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.TenantResponse
	for _, tenant := range tenants {
		response = append(response, mapTenantToResponse(tenant))
	}

	return response, nil
}

// GetTenantByID gets a tenant by ID
func (s *tenantService) GetTenantByID(ctx context.Context, tenantID string, userID string) (dto.TenantResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	// Check if user is part of this tenant
	isInTenant, err := s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, err
	}

	if !isInTenant {
		return dto.TenantResponse{}, dto.ErrUserNotInTenant
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	return mapTenantToResponse(tenant), nil
}

// UpdateTenant updates a tenant
func (s *tenantService) UpdateTenant(ctx context.Context, tenantID string, req dto.TenantUpdateRequest) (dto.TenantResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Description != "" {
		tenant.Description = req.Description
	}
	if req.LogoURL != "" {
		tenant.LogoURL = req.LogoURL
	}
	if req.ContactEmail != "" {
		tenant.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != "" {
		tenant.ContactPhone = req.ContactPhone
	}

	tenant.UpdatedAt = time.Now()

	updatedTenant, err := s.repo.Update(ctx, nil, tenant)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUpdateTenantFailed
	}

	return mapTenantToResponse(updatedTenant), nil
}

// DeleteTenant performs a soft delete on a tenant
func (s *tenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	tenant.IsActive = false
	tenant.UpdatedAt = time.Now()

	_, err = s.repo.Update(ctx, nil, tenant)
	if err != nil {
		return dto.ErrUpdateTenantFailed
	}

	return nil
}

// UpdateSubscription updates a tenant's subscription
func (s *tenantService) UpdateSubscription(ctx context.Context, tenantID string, req dto.SubscriptionRequest) (dto.TenantResponse, error) {
	// Parse tenant ID
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	// Check if user making the request is superadmin
	userIDCtx := ctx.Value("userID")
	if userIDCtx == nil {
		return dto.TenantResponse{}, dto.ErrUnauthorized
	}

	userID, ok := userIDCtx.(string)
	if !ok {
		return dto.TenantResponse{}, dto.ErrUnauthorized
	}

	// Get user role
	user, err := s.userRepo.FindUserById(ctx, nil, userID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrRoleNotFound
	}

	// Only superadmin can update subscription
	if role.Name != "SUPERADMIN" {
		return dto.TenantResponse{}, dto.ErrUnauthorized
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	// Set subscription dates
	now := time.Now()
	startDate := now
	endDate := now.AddDate(1, 0, 0) // 1 year from now

	// Use provided dates if they exist
	if req.SubscriptionStartDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", req.SubscriptionStartDate); err == nil {
			startDate = parsedStartDate
		}
	}

	if req.SubscriptionEndDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", req.SubscriptionEndDate); err == nil {
			endDate = parsedEndDate
		}
	}

	// Set the subscription plan
	tenant.SubscriptionPlan = string(req.SubscriptionPlan)

	// Set max users based on plan if not explicitly provided
	if req.MaxUsers <= 0 {
		switch req.SubscriptionPlan {
		case dto.PlanBasic:
			tenant.MaxUsers = 50
		case dto.PlanPremium:
			tenant.MaxUsers = 100
		case dto.PlanEnterprise:
			tenant.MaxUsers = 500
		}
	} else {
		tenant.MaxUsers = req.MaxUsers
	}

	tenant.SubscriptionStartDate = startDate
	tenant.SubscriptionEndDate = endDate
	tenant.UpdatedAt = time.Now()

	updatedTenant, err := s.repo.Update(ctx, nil, tenant)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUpdateTenantFailed
	}

	return mapTenantToResponse(updatedTenant), nil
}

// InviteUserToTenant invites a user to join a tenant
func (s *tenantService) InviteUserToTenant(ctx context.Context, req dto.InviteUserRequest) error {
	tenantUUID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	// Check if tenant exists
	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	// Check if user exists
	user, err := s.userRepo.FindUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Get user's role
	// Check if RoleID is nil before dereferencing
	if user.RoleID == nil {
		return dto.ErrRoleNotFound
	}
	role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		return dto.ErrRoleNotFound
	}

	// Superadmin cannot be invited to any tenant
	if role.Name == "SUPERADMIN" {
		return dto.ErrSuperadminCannotJoinTenant
	}

	// Check if user is already in any tenant
	existingTenants, err := s.repo.GetTenantsByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	// If user already has a tenant and is not a superadmin, they cannot join another tenant
	if len(existingTenants) > 0 && role.Name != "SUPERADMIN" {
		return dto.ErrUserAlreadyInTenant
	}

	// Check if tenant has reached max users
	count, err := s.repo.CountTenantUsers(ctx, tenantUUID)
	if err != nil {
		return err
	}

	if count >= tenant.MaxUsers {
		return dto.ErrMaxUserLimitReached
	}

	// Create user-tenant relationship
	userTenant := entity.UserTenant{
		ID:        uuid.New(),
		UserID:    user.ID,
		TenantID:  tenant.ID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.AddUserToTenant(ctx, nil, userTenant)
	if err != nil {
		return err
	}

	// Send invitation email
	emailBody := "Anda telah diundang untuk bergabung dengan tenant " + tenant.Name
	err = s.emailSender.Send(user.Email, "Undangan Tenant", emailBody)
	if err != nil {
		// Log error but don't fail the operation
		// Consider implementing a job queue for retry mechanism
	}

	return nil
}

// GetUserTenants gets all tenants for a specific user
func (s *tenantService) GetUserTenants(ctx context.Context, userID string) ([]dto.TenantResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, dto.ErrUserNotFound
	}

	tenants, err := s.repo.GetTenantsByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	var response []dto.TenantResponse
	for _, tenant := range tenants {
		response = append(response, mapTenantToResponse(tenant))
	}

	return response, nil
}

// PromoteToAdmin promotes a user to ADMIN role within a tenant
func (s *tenantService) PromoteToAdmin(ctx context.Context, tenantID, userID, promotedBy string) error {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Parse but use the promotedBy value directly where needed
	_, err = uuid.Parse(promotedBy)
	if err != nil {
		return dto.ErrUserNotFound
	}

	isInTenant, err := s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
	if err != nil || !isInTenant {
		return dto.ErrUserNotInTenant
	}

	user, err := s.userRepo.FindUserById(ctx, nil, userID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Check if user is already a superadmin
	currentRole, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		return dto.ErrRoleNotFound
	}

	if currentRole.Name == "SUPERADMIN" {
		return dto.ErrSuperadminCannotJoinTenant
	}

	adminRole, err := s.roleRepo.FindByName(ctx, "ADMIN")
	if err != nil {
		return dto.ErrRoleNotFound
	}

	// Create a pointer to the UUID
	roleIDPtr := &adminRole.ID
	user.RoleID = roleIDPtr
	user.UpdatedAt = time.Now()

	_, err = s.userRepo.Update(ctx, nil, user)
	if err != nil {
		return dto.ErrUpdateUserFailed
	}

	return nil
}

// DemoteFromAdmin demotes a user from ADMIN role to regular user
func (s *tenantService) DemoteFromAdmin(ctx context.Context, tenantID, userID, demotedBy string) error {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Parse but use the demotedBy value directly where needed
	_, err = uuid.Parse(demotedBy)
	if err != nil {
		return dto.ErrUserNotFound
	}

	isInTenant, err := s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
	if err != nil || !isInTenant {
		return dto.ErrUserNotInTenant
	}

	user, err := s.userRepo.FindUserById(ctx, nil, userID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Check if user is a superadmin
	currentRole, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		return dto.ErrRoleNotFound
	}

	if currentRole.Name == "SUPERADMIN" {
		return dto.ErrSuperadminCannotJoinTenant
	}

	userRole, err := s.roleRepo.FindByName(ctx, "USER")
	if err != nil {
		return dto.ErrRoleNotFound
	}

	// Create a pointer to the UUID
	roleIDPtr := &userRole.ID
	user.RoleID = roleIDPtr
	user.UpdatedAt = time.Now()

	_, err = s.userRepo.Update(ctx, nil, user)
	if err != nil {
		return dto.ErrUpdateUserFailed
	}

	return nil
}

// UpdateTenantByAdmin updates a tenant by an admin user
func (s *tenantService) UpdateTenantByAdmin(ctx context.Context, tenantID, adminUserID string, req dto.TenantUpdateRequest) (dto.TenantResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	// Get admin user and role
	admin, err := s.userRepo.FindUserById(ctx, nil, adminUserID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	// Check admin role
	adminRole, err := s.roleRepo.FindByID(ctx, *admin.RoleID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrRoleNotFound
	}

	// Check if admin is SuperAdmin
	isSuperAdmin := adminRole.Name == "SUPERADMIN"

	// If not SuperAdmin, check if admin is part of the tenant
	if !isSuperAdmin {
		isInTenant, err := s.repo.IsUserInTenant(ctx, adminUUID, tenantUUID)
		if err != nil || !isInTenant {
			return dto.TenantResponse{}, dto.ErrUserNotInTenant
		}
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Description != "" {
		tenant.Description = req.Description
	}
	if req.LogoURL != "" {
		tenant.LogoURL = req.LogoURL
	}
	if req.ContactEmail != "" {
		tenant.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != "" {
		tenant.ContactPhone = req.ContactPhone
	}

	tenant.UpdatedAt = time.Now()

	updatedTenant, err := s.repo.Update(ctx, nil, tenant)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUpdateTenantFailed
	}

	return mapTenantToResponse(updatedTenant), nil
}

// InviteUserByAdmin allows an admin to invite a user to a tenant
func (s *tenantService) InviteUserByAdmin(ctx context.Context, adminUserID string, req dto.InviteUserRequest) error {
	tenantUUID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Get admin user and role
	admin, err := s.userRepo.FindUserById(ctx, nil, adminUserID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Check admin role
	adminRole, err := s.roleRepo.FindByID(ctx, *admin.RoleID)
	if err != nil {
		return dto.ErrRoleNotFound
	}

	// Check if admin is SuperAdmin
	isSuperAdmin := adminRole.Name == "SUPERADMIN"

	// If not SuperAdmin, check if admin is part of the tenant
	if !isSuperAdmin {
		isInTenant, err := s.repo.IsUserInTenant(ctx, adminUUID, tenantUUID)
		if err != nil || !isInTenant {
			return dto.ErrUserNotInTenant
		}
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	user, err := s.userRepo.FindUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrUserNotFound
	}

	count, err := s.repo.CountTenantUsers(ctx, tenantUUID)
	if err != nil {
		return err
	}

	if count >= tenant.MaxUsers {
		return dto.ErrMaxUserLimitReached
	}

	userTenant := entity.UserTenant{
		ID:        uuid.New(),
		UserID:    user.ID,
		TenantID:  tenant.ID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.AddUserToTenant(ctx, nil, userTenant)
	if err != nil {
		return err
	}

	emailBody := "Anda telah diundang untuk bergabung dengan tenant " + tenant.Name
	err = s.emailSender.Send(user.Email, "Undangan Tenant", emailBody)
	if err != nil {
		// Log error but continue
	}

	return nil
}

// RemoveUserFromTenant allows an admin to remove a user from a tenant
func (s *tenantService) RemoveUserFromTenant(ctx context.Context, tenantID, userID, adminUserID string) error {
	log.Printf("Starting RemoveUserFromTenant: tenantID=%s, userID=%s, adminUserID=%s", tenantID, userID, adminUserID)

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		log.Printf("Error parsing tenantID: %v", err)
		return dto.ErrTenantNotFound
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Error parsing userID: %v", err)
		return dto.ErrUserNotFound
	}

	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		log.Printf("Error parsing adminUserID: %v", err)
		return dto.ErrUserNotFound
	}

	// Check if the admin is a SUPERADMIN
	admin, err := s.userRepo.FindUserById(ctx, nil, adminUserID)
	if err != nil {
		log.Printf("Error finding admin user: %v", err)
		return dto.ErrUserNotFound
	}

	adminRole, err := s.roleRepo.FindByID(ctx, *admin.RoleID)
	if err != nil {
		log.Printf("Error finding admin role: %v", err)
		return dto.ErrRoleNotFound
	}

	isSuperAdmin := adminRole.Name == "SUPERADMIN"

	// Declare isInTenant variable
	var isInTenant bool

	// Check if admin is part of the tenant (skip for SUPERADMIN)
	if !isSuperAdmin {
		isInTenant, err = s.repo.IsUserInTenant(ctx, adminUUID, tenantUUID)
		if err != nil {
			log.Printf("Error checking if admin is in tenant: %v", err)
			return err
		}
		if !isInTenant {
			log.Printf("Admin is not part of the tenant")
			return dto.ErrUserNotInTenant
		}
	}

	// Check if the user to be removed is part of the tenant
	isInTenant, err = s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
	if err != nil {
		log.Printf("Error checking if user is in tenant: %v", err)
		return err
	}
	if !isInTenant {
		log.Printf("User is not part of the tenant")
		return dto.ErrUserNotInTenant
	}

	// Prevent admin from removing themselves
	if userID == adminUserID {
		log.Printf("Admin cannot remove themselves")
		return dto.ErrCannotRemoveSelf
	}

	// Check if the user is a SUPERADMIN
	user, err := s.userRepo.FindUserById(ctx, nil, userID)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return dto.ErrUserNotFound
	}

	userRole, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		log.Printf("Error finding user role: %v", err)
		return dto.ErrRoleNotFound
	}

	if userRole.Name == "SUPERADMIN" {
		log.Printf("Cannot remove a SUPERADMIN")
		return dto.ErrCannotRemoveSuperadmin
	}

	// Perform the removal
	log.Printf("Attempting to remove user from tenant")
	err = s.repo.RemoveUserFromTenant(ctx, nil, userUUID, tenantUUID)
	if err != nil {
		log.Printf("Error removing user from tenant: %v", err)
		return err
	}

	log.Printf("Successfully removed user from tenant")
	return nil
}

// GetTenantUsers gets a paginated list of users in a tenant
func (s *tenantService) GetTenantUsers(ctx context.Context, tenantID string, page, limit int) ([]dto.UserResponse, int, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, 0, dto.ErrTenantNotFound
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, totalCount, err := s.repo.GetUsersByTenantID(ctx, tenantUUID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:         user.ID.String(),
			Username:   user.UserName,
			Email:      user.Email,
			UserRole:   user.Role.Name,
			IsVerified: user.IsVerified,
		})
	}

	return userResponses, totalCount, nil
}

// UpdateTenantContact updates the contact information for a tenant
func (s *tenantService) UpdateTenantContact(ctx context.Context, tenantID, contactEmail, contactPhone string) error {
	// Parse tenant ID
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	// Check if user making the request is admin or superadmin
	userIDCtx := ctx.Value("userID")
	if userIDCtx != nil {
		userID, ok := userIDCtx.(string)
		if ok {
			// Get user role
			user, err := s.userRepo.FindUserById(ctx, nil, userID)
			if err == nil {
				role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
				if err == nil {
					// SuperAdmin can update any tenant without being a member
					if role.Name == "SUPERADMIN" {
						// Continue with tenant update
					} else if role.Name == "ADMIN" {
						// Admin must be a member of the tenant
						userUUID, _ := uuid.Parse(userID)
						isInTenant, err := s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
						if err != nil || !isInTenant {
							return dto.ErrUserNotInTenant
						}
					}
				}
			}
		}
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	tenant.ContactEmail = contactEmail
	if contactPhone != "" {
		tenant.ContactPhone = contactPhone
	}

	tenant.UpdatedAt = time.Now()

	_, err = s.repo.Update(ctx, nil, tenant)
	if err != nil {
		return dto.ErrUpdateTenantFailed
	}

	return nil
}

// UpdateTenantLogo updates the logo URL for a tenant
func (s *tenantService) UpdateTenantLogo(ctx context.Context, tenantID, logoURL string) error {
	// Parse tenant ID
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	// Check if user making the request is admin or superadmin
	userIDCtx := ctx.Value("userID")
	if userIDCtx != nil {
		userID, ok := userIDCtx.(string)
		if ok {
			// Get user role
			user, err := s.userRepo.FindUserById(ctx, nil, userID)
			if err == nil {
				role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
				if err == nil {
					// SuperAdmin can update any tenant without being a member
					if role.Name == "SUPERADMIN" {
						// Continue with tenant update
					} else if role.Name == "ADMIN" {
						// Admin must be a member of the tenant
						userUUID, _ := uuid.Parse(userID)
						isInTenant, err := s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
						if err != nil || !isInTenant {
							return dto.ErrUserNotInTenant
						}
					}
				}
			}
		}
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	tenant.LogoURL = logoURL
	tenant.UpdatedAt = time.Now()

	_, err = s.repo.Update(ctx, nil, tenant)
	if err != nil {
		return dto.ErrUpdateTenantFailed
	}

	return nil
}

// UpdateTenantSubscription updates subscription details for a tenant
func (s *tenantService) UpdateTenantSubscription(ctx context.Context, tenantID, plan, startDate, endDate string, maxUsers int) error {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.ErrTenantNotFound
	}

	// Parse dates
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return dto.ErrInvalidDate
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return dto.ErrInvalidDate
	}

	// Validate subscription plan
	subscriptionPlan := dto.SubscriptionPlan(plan)
	if subscriptionPlan != dto.PlanBasic && subscriptionPlan != dto.PlanPremium && subscriptionPlan != dto.PlanEnterprise {
		return dto.ErrInvalidSubscriptionPlan
	}

	tenant.SubscriptionPlan = plan
	tenant.SubscriptionStartDate = startDateParsed
	tenant.SubscriptionEndDate = endDateParsed

	if maxUsers > 0 {
		tenant.MaxUsers = maxUsers
	}

	tenant.UpdatedAt = time.Now()

	_, err = s.repo.Update(ctx, nil, tenant)
	if err != nil {
		return dto.ErrUpdateTenantFailed
	}

	return nil
}

// UpdateTenantProfile updates a tenant's profile
func (s *tenantService) UpdateTenantProfile(ctx context.Context, tenantID uuid.UUID, req dto.TenantProfileUpdateRequest) error {
	// Get tenant
	tenant, err := s.repo.FindByID(ctx, tenantID)
	if err != nil {
		return err
	}

	// Update basic info
	tenant.Name = req.Name
	tenant.Description = req.Description
	tenant.ContactEmail = req.ContactEmail
	tenant.ContactPhone = req.ContactPhone

	// Handle logo upload if provided
	if req.Logo != nil {
		// Delete old logo if exists
		if tenant.LogoURL != "" {
			oldPath := strings.TrimPrefix(tenant.LogoURL, "/uploads/")
			utils.DeleteFile(oldPath)
		}

		// Upload new logo
		filePath, err := utils.UploadFile(req.Logo, "logos")
		if err != nil {
			return fmt.Errorf("failed to upload logo: %v", err)
		}

		tenant.LogoURL = utils.GetFileURL(filePath)
	}

	// Update tenant
	_, err = s.repo.Update(ctx, nil, tenant)
	return err
}

// ListTenants retrieves a paginated list of tenants
func (s *tenantService) ListTenants(ctx context.Context, req dto.TenantListRequest) ([]dto.TenantListResponse, int64, error) {
	// Get all tenants
	tenants, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Apply filters
	filteredTenants := make([]entity.Tenant, 0)
	for _, tenant := range tenants {
		if req.Search != "" && !strings.Contains(strings.ToLower(tenant.Name), strings.ToLower(req.Search)) {
			continue
		}
		if req.IsActive != nil && tenant.IsActive != *req.IsActive {
			continue
		}
		filteredTenants = append(filteredTenants, tenant)
	}

	// Calculate pagination
	total := int64(len(filteredTenants))
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if end > len(filteredTenants) {
		end = len(filteredTenants)
	}
	if start > len(filteredTenants) {
		start = len(filteredTenants)
	}

	// Convert to response DTO
	response := make([]dto.TenantListResponse, 0)
	for _, tenant := range filteredTenants[start:end] {
		response = append(response, dto.TenantListResponse{
			ID:               tenant.ID,
			Name:             tenant.Name,
			Description:      tenant.Description,
			LogoURL:          tenant.LogoURL,
			ContactEmail:     tenant.ContactEmail,
			ContactPhone:     tenant.ContactPhone,
			SubscriptionPlan: tenant.SubscriptionPlan,
			IsActive:         tenant.IsActive,
			CreatedAt:        tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, total, nil
}

// GetCurrentTenantDetails gets the details of a tenant for a user based on token
func (s *tenantService) GetCurrentTenantDetails(ctx context.Context, tenantID string, userID string) (dto.TenantResponse, error) {
	// Parse tenantID to UUID
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	// Check if user is part of this tenant
	isInTenant, err := s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, err
	}

	if !isInTenant {
		return dto.TenantResponse{}, dto.ErrUserNotInTenant
	}

	// Get tenant details
	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	// Return the tenant details
	return mapTenantToResponse(tenant), nil
}

// DirectInviteUserToTenant mengundang user langsung ke tenant (hanya bisa dilakukan oleh SUPERADMIN)
func (s *tenantService) DirectInviteUserToTenant(ctx context.Context, tenantID, userID string, superadminID string) error {
	log.Printf("DirectInviteUserToTenant called with tenantID: %s, userID: %s", tenantID, userID)

	// Parse IDs
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		log.Printf("Error parsing tenant ID: %v", err)
		return dto.ErrInvalidTenantID
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		return dto.ErrInvalidUserID
	}

	// Verify tenant exists
	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		log.Printf("Error finding tenant: %v", err)
		return err
	}

	// Verify user exists and is not a SUPERADMIN
	user, err := s.userRepo.FindUserById(ctx, nil, userID)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return err
	}

	userRole, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		log.Printf("Error finding user role: %v", err)
		return err
	}

	if userRole.Name == "SUPERADMIN" {
		log.Printf("Cannot invite SUPERADMIN to tenant")
		return dto.ErrCannotInviteSuperadmin
	}

	// Check if user is already in this specific tenant
	isInTenant, err := s.repo.IsUserInTenant(ctx, userUUID, tenantUUID)
	if err != nil {
		log.Printf("Error checking if user is in tenant: %v", err)
		return err
	}
	if isInTenant {
		log.Printf("User is already in this tenant")
		return dto.ErrUserInSameTenant
	}

	// Check if user is in any other tenant
	existingTenants, err := s.repo.GetTenantsByUserID(ctx, userUUID)
	if err != nil {
		log.Printf("Error checking user's existing tenants: %v", err)
		return err
	}

	// If user is already in any tenant, they cannot be invited to another tenant
	if len(existingTenants) > 0 {
		log.Printf("User is already in another tenant")
		return dto.ErrUserInOtherTenant
	}

	// Check if tenant has reached its user limit
	userCount, err := s.repo.CountTenantUsers(ctx, tenantUUID)
	if err != nil {
		log.Printf("Error counting users in tenant: %v", err)
		return err
	}

	if userCount >= tenant.MaxUsers {
		log.Printf("Tenant has reached maximum user limit: %d/%d", userCount, tenant.MaxUsers)
		return dto.ErrTenantMaxUsersReached
	}

	// Create user-tenant relationship
	userTenant := entity.UserTenant{
		ID:        uuid.New(),
		UserID:    userUUID,
		TenantID:  tenantUUID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.AddUserToTenant(ctx, nil, userTenant); err != nil {
		log.Printf("Error creating user-tenant relationship: %v", err)
		return err
	}

	log.Printf("Successfully invited user %s to tenant %s", userID, tenantID)
	return nil
}

// GetTenantDetailsByID gets detailed tenant information for admin and superadmin
func (s *tenantService) GetTenantDetailsByID(ctx context.Context, tenantID, adminID string) (dto.TenantResponse, error) {
	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	// Check if admin exists
	admin, err := s.userRepo.FindUserById(ctx, nil, adminID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	// Get admin role
	role, err := s.roleRepo.FindByID(ctx, *admin.RoleID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrRoleNotFound
	}

	// Check if user is ADMIN or SUPERADMIN
	isAdmin := role.Name == "ADMIN"
	isSuperAdmin := role.Name == "SUPERADMIN"

	if !isAdmin && !isSuperAdmin {
		return dto.TenantResponse{}, dto.ErrUnauthorized
	}

	// For admin, check if they are in the tenant
	if isAdmin {
		isInTenant, err := s.repo.IsUserInTenant(ctx, adminUUID, tenantUUID)
		if err != nil {
			return dto.TenantResponse{}, err
		}

		if !isInTenant {
			return dto.TenantResponse{}, dto.ErrUserNotInTenant
		}
	}

	// Get tenant details
	tenant, err := s.repo.FindByID(ctx, tenantUUID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrTenantNotFound
	}

	// Get additional tenant metadata only for admin/superadmin
	userCount, _ := s.repo.CountTenantUsers(ctx, tenantUUID)

	// Create enhanced response with admin-specific information
	response := mapTenantToResponse(tenant)
	response.UserCount = userCount
	response.CanEdit = true

	return response, nil
}

// GetUserCurrentTenant gets the current tenant for a user using token
func (s *tenantService) GetUserCurrentTenant(ctx context.Context, userID string) (dto.TenantResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrUserNotFound
	}

	// Get user's tenants
	tenants, err := s.repo.GetTenantsByUserID(ctx, userUUID)
	if err != nil {
		return dto.TenantResponse{}, err
	}

	if len(tenants) == 0 {
		return dto.TenantResponse{}, dto.ErrUserNotInTenant
	}

	// Return first tenant (could be enhanced to track and return most recently used tenant)
	return mapTenantToResponse(tenants[0]), nil
}

// UpdateCurrentTenant updates the tenant for an admin using their token
func (s *tenantService) UpdateCurrentTenant(ctx context.Context, adminUserID string, req dto.TenantUpdateRequest) (dto.TenantResponse, error) {
	// Validate admin user ID
	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrInvalidID
	}

	// Get admin's tenant
	tenants, err := s.repo.GetTenantsByUserID(ctx, adminUUID)
	if err != nil || len(tenants) == 0 {
		return dto.TenantResponse{}, dto.ErrUserNotInTenant
	}

	// Get the tenant that the admin belongs to
	tenantID := tenants[0].ID.String()

	// Check if user is an admin in this tenant
	role, err := s.userRepo.GetUserRoleInTenant(ctx, adminUUID, tenants[0].ID)
	if err != nil {
		return dto.TenantResponse{}, err
	}

	if role.Name != "ADMIN" && role.Name != "SUPERADMIN" {
		return dto.TenantResponse{}, dto.ErrUnauthorized
	}

	// Use the existing UpdateTenantByAdmin method
	return s.UpdateTenantByAdmin(ctx, tenantID, adminUserID, req)
}

// GetCurrentTenantUsers gets all users in the admin's current tenant with pagination
func (s *tenantService) GetCurrentTenantUsers(ctx context.Context, adminUserID string, page, limit int) ([]dto.UserResponse, int, error) {
	// Validate admin user ID
	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return nil, 0, dto.ErrInvalidID
	}

	// Get admin's tenant
	tenants, err := s.repo.GetTenantsByUserID(ctx, adminUUID)
	if err != nil || len(tenants) == 0 {
		return nil, 0, dto.ErrUserNotInTenant
	}

	tenantID := tenants[0].ID.String()

	// Check if user is an admin in this tenant
	role, err := s.userRepo.GetUserRoleInTenant(ctx, adminUUID, tenants[0].ID)
	if err != nil {
		return nil, 0, err
	}

	if role.Name != "ADMIN" && role.Name != "SUPERADMIN" {
		return nil, 0, dto.ErrUnauthorized
	}

	// Use existing GetTenantUsers method
	return s.GetTenantUsers(ctx, tenantID, page, limit)
}

// InviteUserToCurrentTenant invites a user to the admin's current tenant
func (s *tenantService) InviteUserToCurrentTenant(ctx context.Context, adminUserID string, req dto.InviteUserRequest) error {
	// Validate admin user ID
	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.ErrInvalidID
	}

	// Get admin's tenant
	tenants, err := s.repo.GetTenantsByUserID(ctx, adminUUID)
	if err != nil || len(tenants) == 0 {
		return dto.ErrUserNotInTenant
	}

	// Set the tenant ID in the request
	req.TenantID = tenants[0].ID.String()

	// Check if user is an admin in this tenant
	role, err := s.userRepo.GetUserRoleInTenant(ctx, adminUUID, tenants[0].ID)
	if err != nil {
		return err
	}

	if role.Name != "ADMIN" && role.Name != "SUPERADMIN" {
		return dto.ErrUnauthorized
	}

	// Use existing InviteUserByAdmin method
	return s.InviteUserByAdmin(ctx, adminUserID, req)
}

// RemoveUserFromCurrentTenant removes a user from the admin's current tenant
func (s *tenantService) RemoveUserFromCurrentTenant(ctx context.Context, userID, adminUserID string) error {
	// Validate admin user ID
	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.ErrInvalidID
	}

	// Get admin's tenant
	tenants, err := s.repo.GetTenantsByUserID(ctx, adminUUID)
	if err != nil || len(tenants) == 0 {
		return dto.ErrUserNotInTenant
	}

	tenantID := tenants[0].ID.String()

	// Check if user is an admin in this tenant
	role, err := s.userRepo.GetUserRoleInTenant(ctx, adminUUID, tenants[0].ID)
	if err != nil {
		return err
	}

	if role.Name != "ADMIN" && role.Name != "SUPERADMIN" {
		return dto.ErrUnauthorized
	}

	// Use existing RemoveUserFromTenant method
	return s.RemoveUserFromTenant(ctx, tenantID, userID, adminUserID)
}

// UpdateCurrentSubscription updates the subscription for the admin's current tenant
func (s *tenantService) UpdateCurrentSubscription(ctx context.Context, adminUserID string, req dto.SubscriptionRequest) (dto.TenantResponse, error) {
	// Validate admin user ID
	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.TenantResponse{}, dto.ErrInvalidID
	}

	// Get admin's tenant
	tenants, err := s.repo.GetTenantsByUserID(ctx, adminUUID)
	if err != nil || len(tenants) == 0 {
		return dto.TenantResponse{}, dto.ErrUserNotInTenant
	}

	tenantID := tenants[0].ID.String()

	// Check if user is an admin in this tenant
	role, err := s.userRepo.GetUserRoleInTenant(ctx, adminUUID, tenants[0].ID)
	if err != nil {
		return dto.TenantResponse{}, err
	}

	if role.Name != "ADMIN" && role.Name != "SUPERADMIN" {
		return dto.TenantResponse{}, dto.ErrUnauthorized
	}

	// Use existing UpdateSubscription method
	return s.UpdateSubscription(ctx, tenantID, req)
}

// UpdateCurrentTenantContact updates the contact information for the admin's current tenant
func (s *tenantService) UpdateCurrentTenantContact(ctx context.Context, adminUserID string, contactEmail, contactPhone string) error {
	// Validate admin user ID
	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.ErrInvalidID
	}

	// Get admin's tenant
	tenants, err := s.repo.GetTenantsByUserID(ctx, adminUUID)
	if err != nil || len(tenants) == 0 {
		return dto.ErrUserNotInTenant
	}

	tenantID := tenants[0].ID.String()

	// Check if user is an admin in this tenant
	role, err := s.userRepo.GetUserRoleInTenant(ctx, adminUUID, tenants[0].ID)
	if err != nil {
		return err
	}

	if role.Name != "ADMIN" && role.Name != "SUPERADMIN" {
		return dto.ErrUnauthorized
	}

	// Use existing UpdateTenantContact method
	return s.UpdateTenantContact(ctx, tenantID, contactEmail, contactPhone)
}

// UpdateCurrentTenantLogo updates the logo for the admin's current tenant
func (s *tenantService) UpdateCurrentTenantLogo(ctx context.Context, adminUserID string, logoURL string) error {
	// Validate admin user ID
	adminUUID, err := uuid.Parse(adminUserID)
	if err != nil {
		return dto.ErrInvalidID
	}

	// Get admin's tenant
	tenants, err := s.repo.GetTenantsByUserID(ctx, adminUUID)
	if err != nil || len(tenants) == 0 {
		return dto.ErrUserNotInTenant
	}

	tenantID := tenants[0].ID.String()

	// Check if user is an admin in this tenant
	role, err := s.userRepo.GetUserRoleInTenant(ctx, adminUUID, tenants[0].ID)
	if err != nil {
		return err
	}

	if role.Name != "ADMIN" && role.Name != "SUPERADMIN" {
		return dto.ErrUnauthorized
	}

	// Use existing UpdateTenantLogo method
	return s.UpdateTenantLogo(ctx, tenantID, logoURL)
}

// RemoveUserFromAllTenants removes a user from all tenants they are part of
func (s *tenantService) RemoveUserFromAllTenants(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Get all tenants the user is part of
	tenants, err := s.repo.GetTenantsByUserID(ctx, userUUID)
	if err != nil {
		return err
	}

	// Remove user from each tenant
	for _, tenant := range tenants {
		err = s.repo.RemoveUserFromTenant(ctx, nil, userUUID, tenant.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveSuperadminFromAllTenants removes all superadmin users from all tenants
func (s *tenantService) RemoveSuperadminFromAllTenants(ctx context.Context) error {
	log.Printf("Starting RemoveSuperadminFromAllTenants")

	// Get superadmin role
	superadminRole, err := s.roleRepo.FindByName(ctx, "SUPERADMIN")
	if err != nil {
		log.Printf("Error finding superadmin role: %v", err)
		return dto.ErrRoleNotFound
	}

	// Get all users with superadmin role
	superadmins, err := s.userRepo.FindUsersByRoleID(ctx, superadminRole.ID)
	if err != nil {
		log.Printf("Error finding superadmin users: %v", err)
		return err
	}

	// Remove each superadmin from all tenants
	for _, superadmin := range superadmins {
		// Get all tenants for this superadmin
		tenants, err := s.repo.GetTenantsByUserID(ctx, superadmin.ID)
		if err != nil {
			log.Printf("Error getting tenants for superadmin %s: %v", superadmin.ID, err)
			continue
		}

		// Remove superadmin from each tenant
		for _, tenant := range tenants {
			err = s.repo.RemoveUserFromTenant(ctx, nil, superadmin.ID, tenant.ID)
			if err != nil {
				log.Printf("Error removing superadmin %s from tenant %s: %v", superadmin.ID, tenant.ID, err)
				continue
			}
			log.Printf("Successfully removed superadmin %s from tenant %s", superadmin.ID, tenant.ID)
		}
	}

	log.Printf("Completed RemoveSuperadminFromAllTenants")
	return nil
}

// GetUserByID retrieves user information by user ID for external API usage
func (s *tenantService) GetUserByID(ctx context.Context, userID string) (dto.ExternalUserInfoResponse, error) {
	// Validate user ID format
	_, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %s", userID)
		return dto.ExternalUserInfoResponse{}, fmt.Errorf("invalid user ID format")
	}

	// Get user by ID
	user, err := s.userRepo.FindUserById(ctx, nil, userID)
	if err != nil {
		log.Printf("User not found: %s", userID)
		return dto.ExternalUserInfoResponse{}, fmt.Errorf("user not found")
	}

	role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		log.Printf("Role not found for user %s: %v", userID, err)
		return dto.ExternalUserInfoResponse{}, fmt.Errorf("role not found")
	}

	// Return user info for external API
	return dto.ExternalUserInfoResponse{
		UserID:   user.ID.String(),
		Username: user.UserName,
		Email:    user.Email,
		Role:     role.Name,
		RoleID:   user.RoleID.String(),
		Status:   user.Status,
		IsActive: user.Status == "active",
	}, nil
}

// mapTenantToResponse maps a Tenant entity to a TenantResponse DTO
func mapTenantToResponse(tenant entity.Tenant) dto.TenantResponse {
	return dto.TenantResponse{
		ID:                    tenant.ID.String(),
		Name:                  tenant.Name,
		Description:           tenant.Description,
		LogoURL:               tenant.LogoURL,
		ContactEmail:          tenant.ContactEmail,
		ContactPhone:          tenant.ContactPhone,
		SubscriptionPlan:      dto.SubscriptionPlan(tenant.SubscriptionPlan),
		SubscriptionStartDate: tenant.SubscriptionStartDate,
		SubscriptionEndDate:   tenant.SubscriptionEndDate,
		MaxUsers:              tenant.MaxUsers,
		IsActive:              tenant.IsActive,
		CreatedAt:             tenant.CreatedAt,
		UpdatedAt:             tenant.UpdatedAt,
	}
}
