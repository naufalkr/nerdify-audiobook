package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserTenantContextRequest untuk request set current tenant
type UserTenantContextRequest struct {
	TenantID uuid.UUID `json:"tenantId" validate:"required"`
}

// UserTenantContextResponse untuk response current tenant
type UserTenantContextResponse struct {
	UserID     uuid.UUID `json:"userId"`
	TenantID   uuid.UUID `json:"tenantId"`
	TenantName string    `json:"tenantName"`
	UserRole   string    `json:"userRole"`
	IsActive   bool      `json:"isActive"`
	JoinedAt   time.Time `json:"joinedAt"`
}

// UserTenantsListResponse untuk response list tenant user
type UserTenantsListResponse struct {
	UserID  uuid.UUID                   `json:"userId"`
	Tenants []UserTenantContextResponse `json:"tenants"`
	Current *UserTenantContextResponse  `json:"current,omitempty"`
	Total   int                         `json:"total"`
}

// TenantUsersResponse untuk response users dalam tenant
type TenantUsersResponse struct {
	TenantID uuid.UUID        `json:"tenantId"`
	Users    []TenantUserInfo `json:"users"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

// TenantUserInfo informasi user dalam tenant
type TenantUserInfo struct {
	UserID   uuid.UUID `json:"userId"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
	Role     string    `json:"role"`
	IsActive bool      `json:"isActive"`
	JoinedAt time.Time `json:"joinedAt"`
}

// UserTenantAccessValidationRequest untuk validasi akses user ke tenant
type UserTenantAccessValidationRequest struct {
	UserID   uuid.UUID `json:"userId" validate:"required"`
	TenantID uuid.UUID `json:"tenantId" validate:"required"`
}

// UserTenantAccessValidationResponse untuk response validasi akses
type UserTenantAccessValidationResponse struct {
	UserID     uuid.UUID `json:"userId"`
	TenantID   uuid.UUID `json:"tenantId"`
	HasAccess  bool      `json:"hasAccess"`
	UserRole   string    `json:"userRole,omitempty"`
	TenantName string    `json:"tenantName,omitempty"`
}

// SwitchTenantRequest untuk request switch tenant
type SwitchTenantRequest struct {
	TenantID uuid.UUID `json:"tenantId" validate:"required"`
}

// SwitchTenantResponse untuk response switch tenant
type SwitchTenantResponse struct {
	Success     bool      `json:"success"`
	NewTenantID uuid.UUID `json:"newTenantId"`
	TenantName  string    `json:"tenantName"`
	UserRole    string    `json:"userRole"`
	Message     string    `json:"message"`
}
