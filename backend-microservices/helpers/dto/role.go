package dto

import (
	"time"

	"github.com/google/uuid"
)

// RoleRequest represents the incoming request structure for creating or updating a role
type RoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description"`
	IsSystem    bool   `json:"isSystem"`
}

// RoleResponse represents the outgoing response structure for role data
type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"isSystem"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// RoleListResponse represents paginated response for roles
type RoleListResponse struct {
	Roles      []RoleResponse `json:"roles"`
	TotalCount int64          `json:"totalCount"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
}

// BulkDeleteRequest represents a request to delete multiple roles
type BulkDeleteRequest struct {
	IDs []uuid.UUID `json:"ids" binding:"required"`
}

// BulkDeleteResponse represents a response after bulk deletion
type BulkDeleteResponse struct {
	DeletedCount int `json:"deletedCount"`
}
