package dto

import (
	"time"
)

// TenantResponse adalah DTO untuk response tenant
type TenantResponse struct {
	ID                    string           `json:"id"`
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	LogoURL               string           `json:"logo_url"`
	ContactEmail          string           `json:"contact_email"`
	ContactPhone          string           `json:"contact_phone"`
	MaxUsers              int              `json:"max_users"`
	CurrentUsers          int              `json:"current_users,omitempty"`
	SubscriptionPlan      SubscriptionPlan `json:"subscription_plan"`
	SubscriptionStartDate time.Time        `json:"subscription_start_date"`
	SubscriptionEndDate   time.Time        `json:"subscription_end_date"`
	IsActive              bool             `json:"is_active"`
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`
	UserCount             int              `json:"user_count,omitempty"` // Added for admin/superadmin view
	CanEdit               bool             `json:"can_edit,omitempty"`   // Indicates if user can edit the tenant
}

// TenantCreateRequest adalah DTO untuk membuat tenant baru
type TenantCreateRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	LogoURL      string `json:"logo_url"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
}

// TenantUpdateRequest adalah DTO untuk update tenant
type TenantUpdateRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	LogoURL      string `json:"logo_url"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
}

// SubscriptionPlan represents the available subscription plans
type SubscriptionPlan string

const (
	PlanBasic      SubscriptionPlan = "Basic"
	PlanPremium    SubscriptionPlan = "Premium"
	PlanEnterprise SubscriptionPlan = "Enterprise"
)

// SubscriptionRequest adalah DTO untuk update subscription tenant
type SubscriptionRequest struct {
	SubscriptionPlan SubscriptionPlan `json:"subscriptionPlan" binding:"required,oneof=Basic Premium Enterprise"`
	MaxUsers         int              `json:"maxUsers"`
	// These fields are optional and will be auto-filled if not provided
	SubscriptionStartDate string `json:"subscriptionStartDate,omitempty"`
	SubscriptionEndDate   string `json:"subscriptionEndDate,omitempty"`
}

// InviteUserRequest adalah DTO untuk invite user ke tenant
type InviteUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	TenantID string `json:"tenant_id" binding:"omitempty,uuid"`
}

// AppError merupakan struktur untuk error yang lebih deskriptif
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error mengimplementasikan interface error
func (e AppError) Error() string {
	return e.Message
}

// NewError membuat error baru dengan code dan message
func NewError(code, message string) error {
	return AppError{
		Code:    code,
		Message: message,
	}
}

// Custom errors untuk tenant operations
var (
	ErrTenantNotFound             = NewError("TENANT_NOT_FOUND", "Tenant tidak ditemukan")
	ErrCreateTenantFailed         = NewError("CREATE_TENANT_FAILED", "Gagal membuat tenant")
	ErrUpdateTenantFailed         = NewError("UPDATE_TENANT_FAILED", "Gagal mengupdate tenant")
	ErrDeleteTenantFailed         = NewError("DELETE_TENANT_FAILED", "Gagal menghapus tenant")
	ErrUserNotInTenant            = NewError("USER_NOT_IN_TENANT", "User tidak berada dalam tenant ini")
	ErrInvalidSubscriptionPlan    = NewError("INVALID_SUBSCRIPTION_PLAN", "Subscription plan tidak valid")
	ErrMaxUserLimitReached        = NewError("MAX_USER_LIMIT_REACHED", "Tenant telah mencapai batas maksimum user")
	ErrCannotRemoveSelf           = NewError("CANNOT_REMOVE_SELF", "Tidak dapat menghapus diri sendiri dari tenant")
	ErrCannotRemoveSuperadmin     = NewError("CANNOT_REMOVE_SUPERADMIN", "Tidak dapat menghapus superadmin dari tenant")
	ErrInvalidDate                = NewError("INVALID_DATE", "Format tanggal tidak valid")
	ErrUserAlreadyInTenant        = NewError("USER_ALREADY_IN_TENANT", "User sudah berada dalam tenant ini")
	ErrSuperadminCannotJoinTenant = NewError("SUPERADMIN_CANNOT_JOIN_TENANT", "Superadmin tidak dapat bergabung dengan tenant manapun")
)
