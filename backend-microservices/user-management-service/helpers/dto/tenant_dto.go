package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type TenantProfileUpdateRequest struct {
	Name         string                `json:"name" binding:"required"`
	Description  string                `json:"description"`
	ContactEmail string                `json:"contact_email" binding:"required,email"`
	ContactPhone string                `json:"contact_phone"`
	Logo         *multipart.FileHeader `form:"logo"`
}

type TenantListResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	LogoURL          string    `json:"logo_url"`
	ContactEmail     string    `json:"contact_email"`
	ContactPhone     string    `json:"contact_phone"`
	SubscriptionPlan string    `json:"subscription_plan"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        string    `json:"created_at"`
}

type TenantListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}
