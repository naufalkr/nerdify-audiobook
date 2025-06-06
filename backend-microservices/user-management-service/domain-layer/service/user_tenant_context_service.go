package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"microservice/user/data-layer/repository"
	"microservice/user/helpers/dto"
)

type UserTenantContextService interface {
	// Mendapatkan current tenant user
	GetUserCurrentTenant(ctx context.Context, userID uuid.UUID) (*dto.UserTenantContextResponse, error)

	// Set current tenant user
	SetUserCurrentTenant(ctx context.Context, userID, tenantID uuid.UUID) (*dto.UserTenantContextResponse, error)

	// Mendapatkan semua tenant yang dimiliki user
	GetUserTenants(ctx context.Context, userID uuid.UUID) (*dto.UserTenantsListResponse, error)

	// Switch tenant user
	SwitchUserTenant(ctx context.Context, userID, tenantID uuid.UUID) (*dto.SwitchTenantResponse, error)

	// Validasi akses user ke tenant
	ValidateUserTenantAccess(ctx context.Context, userID, tenantID uuid.UUID) (*dto.UserTenantAccessValidationResponse, error)

	// Mendapatkan users dalam tenant (untuk admin)
	GetTenantUsers(ctx context.Context, tenantID uuid.UUID, page, limit int) (*dto.TenantUsersResponse, error)
}

type userTenantContextService struct {
	userTenantRepo repository.UserTenantRepository
	userRepo       repository.UserRepository
	tenantRepo     repository.TenantRepository
	db             *gorm.DB
}

func NewUserTenantContextService(
	userTenantRepo repository.UserTenantRepository,
	userRepo repository.UserRepository,
	tenantRepo repository.TenantRepository,
	db *gorm.DB,
) UserTenantContextService {
	return &userTenantContextService{
		userTenantRepo: userTenantRepo,
		userRepo:       userRepo,
		tenantRepo:     tenantRepo,
		db:             db,
	}
}

func (s *userTenantContextService) GetUserCurrentTenant(ctx context.Context, userID uuid.UUID) (*dto.UserTenantContextResponse, error) {
	// Cek user exists
	user, err := s.userRepo.FindUserById(ctx, nil, userID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get current tenant
	userTenant, err := s.userTenantRepo.GetUserCurrentTenant(userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get current tenant: %w", err)
	}

	// Get tenant details
	tenant, err := s.tenantRepo.FindByID(ctx, userTenant.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	// Get user role in tenant (from user's role)
	userRole := "User" // default
	if user.RoleID != nil {
		role, err := s.userRepo.GetUserRole(ctx, nil, *user.RoleID)
		if err == nil {
			userRole = role.Name
		}
	}

	return &dto.UserTenantContextResponse{
		UserID:     userTenant.UserID,
		TenantID:   userTenant.TenantID,
		TenantName: tenant.Name,
		UserRole:   userRole,
		IsActive:   userTenant.IsActive,
		JoinedAt:   userTenant.CreatedAt,
	}, nil
}

func (s *userTenantContextService) SetUserCurrentTenant(ctx context.Context, userID, tenantID uuid.UUID) (*dto.UserTenantContextResponse, error) {
	// Validasi user exists
	user, err := s.userRepo.FindUserById(ctx, nil, userID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validasi tenant exists
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	// Validasi user memiliki akses ke tenant
	userTenant, err := s.userTenantRepo.GetUserTenantByUserIDAndTenantID(userID.String(), tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("user does not have access to this tenant: %w", err)
	}

	if !userTenant.IsActive {
		return nil, errors.New("user access to this tenant is inactive")
	}

	// Update current tenant (set semua tenant user menjadi tidak current, lalu set yang dipilih menjadi current)
	err = s.userTenantRepo.UpdateUserCurrentTenant(userID.String(), tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to update current tenant: %w", err)
	}

	// Get user role
	userRole := "User"
	if user.RoleID != nil {
		role, err := s.userRepo.GetUserRole(ctx, nil, *user.RoleID)
		if err == nil {
			userRole = role.Name
		}
	}

	return &dto.UserTenantContextResponse{
		UserID:     userID,
		TenantID:   tenantID,
		TenantName: tenant.Name,
		UserRole:   userRole,
		IsActive:   userTenant.IsActive,
		JoinedAt:   userTenant.CreatedAt,
	}, nil
}

func (s *userTenantContextService) GetUserTenants(ctx context.Context, userID uuid.UUID) (*dto.UserTenantsListResponse, error) {
	// Validasi user exists
	user, err := s.userRepo.FindUserById(ctx, nil, userID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get all user tenants
	userTenants, err := s.userTenantRepo.GetUserTenants(userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user tenants: %w", err)
	}

	var tenants []dto.UserTenantContextResponse
	var currentTenant *dto.UserTenantContextResponse

	for _, ut := range userTenants {
		// Get tenant details
		tenant, err := s.tenantRepo.FindByID(ctx, ut.TenantID)
		if err != nil {
			continue // skip jika tenant tidak ditemukan
		}

		// Get user role
		userRole := "User"
		if user.RoleID != nil {
			role, err := s.userRepo.GetUserRole(ctx, nil, *user.RoleID)
			if err == nil {
				userRole = role.Name
			}
		}

		tenantInfo := dto.UserTenantContextResponse{
			UserID:     ut.UserID,
			TenantID:   ut.TenantID,
			TenantName: tenant.Name,
			UserRole:   userRole,
			IsActive:   ut.IsActive,
			JoinedAt:   ut.CreatedAt,
		}

		tenants = append(tenants, tenantInfo)

		// Check if this is current tenant - untuk saat ini, kita ambil yang pertama sebagai current
		// TODO: Implementasi field is_current di masa depan
		if currentTenant == nil {
			currentTenant = &tenantInfo
		}
	}

	return &dto.UserTenantsListResponse{
		UserID:  userID,
		Tenants: tenants,
		Current: currentTenant,
		Total:   len(tenants),
	}, nil
}

func (s *userTenantContextService) SwitchUserTenant(ctx context.Context, userID, tenantID uuid.UUID) (*dto.SwitchTenantResponse, error) {
	// Set current tenant
	result, err := s.SetUserCurrentTenant(ctx, userID, tenantID)
	if err != nil {
		return &dto.SwitchTenantResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to switch tenant: %v", err),
		}, err
	}

	return &dto.SwitchTenantResponse{
		Success:     true,
		NewTenantID: result.TenantID,
		TenantName:  result.TenantName,
		UserRole:    result.UserRole,
		Message:     "Successfully switched to tenant: " + result.TenantName,
	}, nil
}

func (s *userTenantContextService) ValidateUserTenantAccess(ctx context.Context, userID, tenantID uuid.UUID) (*dto.UserTenantAccessValidationResponse, error) {
	// Get user
	user, err := s.userRepo.FindUserById(ctx, nil, userID.String())
	if err != nil {
		return &dto.UserTenantAccessValidationResponse{
			UserID:    userID,
			TenantID:  tenantID,
			HasAccess: false,
		}, nil
	}

	// Get tenant
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return &dto.UserTenantAccessValidationResponse{
			UserID:    userID,
			TenantID:  tenantID,
			HasAccess: false,
		}, nil
	}

	// Check user-tenant relationship
	userTenant, err := s.userTenantRepo.GetUserTenantByUserIDAndTenantID(userID.String(), tenantID.String())
	if err != nil {
		return &dto.UserTenantAccessValidationResponse{
			UserID:     userID,
			TenantID:   tenantID,
			HasAccess:  false,
			TenantName: tenant.Name,
		}, nil
	}

	// Get user role
	userRole := "User"
	if user.RoleID != nil {
		role, err := s.userRepo.GetUserRole(ctx, nil, *user.RoleID)
		if err == nil {
			userRole = role.Name
		}
	}

	return &dto.UserTenantAccessValidationResponse{
		UserID:     userID,
		TenantID:   tenantID,
		HasAccess:  userTenant.IsActive,
		UserRole:   userRole,
		TenantName: tenant.Name,
	}, nil
}

func (s *userTenantContextService) GetTenantUsers(ctx context.Context, tenantID uuid.UUID, page, limit int) (*dto.TenantUsersResponse, error) {
	// Validasi tenant exists
	_, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	// Set default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Get tenant users with pagination
	userTenants, total, err := s.userTenantRepo.GetTenantUsersWithPagination(tenantID.String(), page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant users: %w", err)
	}

	var users []dto.TenantUserInfo
	for _, ut := range userTenants {
		// Get user details
		user, err := s.userRepo.FindUserById(ctx, nil, ut.UserID.String())
		if err != nil {
			continue // skip jika user tidak ditemukan
		}

		// Get user role
		userRole := "User"
		if user.RoleID != nil {
			role, err := s.userRepo.GetUserRole(ctx, nil, *user.RoleID)
			if err == nil {
				userRole = role.Name
			}
		}

		userInfo := dto.TenantUserInfo{
			UserID:   ut.UserID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     userRole,
			IsActive: ut.IsActive,
			JoinedAt: ut.CreatedAt,
		}

		users = append(users, userInfo)
	}

	return &dto.TenantUsersResponse{
		TenantID: tenantID,
		Users:    users,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}
