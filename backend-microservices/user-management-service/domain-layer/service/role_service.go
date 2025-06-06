package service

import (
	"context"
	"time"

	"microservice/user/data-layer/entity"
	"microservice/user/data-layer/repository"

	"github.com/google/uuid"
)

type RoleService interface {
	CreateRole(ctx context.Context, name, description string, isSystem bool) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	ListAll(ctx context.Context) ([]entity.Role, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	SeedDefaultRoles(ctx context.Context) error
	UpdateRole(ctx context.Context, id uuid.UUID, name, description string) (*entity.Role, error)
	CountRoles(ctx context.Context) (int64, error)
	SearchRoles(ctx context.Context, query string, page, pageSize int) ([]entity.Role, int64, error)
	GetSystemRoles(ctx context.Context) ([]entity.Role, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	BulkDeleteRoles(ctx context.Context, ids []uuid.UUID) (int, error)
}

type roleService struct {
	repo repository.RoleRepository
}

// NewRoleService creates a new role service
func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{
		repo: repo,
	}
}

// CreateRole creates a new role
func (s *roleService) CreateRole(ctx context.Context, name, description string, isSystem bool) (*entity.Role, error) {
	role := &entity.Role{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		IsSystem:    isSystem,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.Create(ctx, role)
}

// FindByName finds a role by name
func (s *roleService) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	return s.repo.FindByName(ctx, name)
}

// FindByID finds a role by ID
func (s *roleService) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	return s.repo.FindByID(ctx, id)
}

// ListAll lists all roles
func (s *roleService) ListAll(ctx context.Context) ([]entity.Role, error) {
	return s.repo.ListAll(ctx)
}

// DeleteByID deletes a role by ID
func (s *roleService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteByID(ctx, id)
}

// SeedDefaultRoles seeds default roles
func (s *roleService) SeedDefaultRoles(ctx context.Context) error {
	return s.repo.SeedDefaultRoles(ctx)
}

// UpdateRole updates an existing role
func (s *roleService) UpdateRole(ctx context.Context, id uuid.UUID, name, description string) (*entity.Role, error) {
	// Find the role first to check if it exists
	existingRole, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the fields
	existingRole.Name = name
	existingRole.Description = description
	existingRole.UpdatedAt = time.Now()

	// Save the changes
	return s.repo.Update(ctx, existingRole)
}

// CountRoles returns the total number of roles
func (s *roleService) CountRoles(ctx context.Context) (int64, error) {
	return s.repo.CountRoles(ctx)
}

// SearchRoles searches roles based on a query with pagination
func (s *roleService) SearchRoles(ctx context.Context, query string, page, pageSize int) ([]entity.Role, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size
	} else if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}

	return s.repo.SearchRoles(ctx, query, page, pageSize)
}

// GetSystemRoles retrieves all system-defined roles
func (s *roleService) GetSystemRoles(ctx context.Context) ([]entity.Role, error) {
	return s.repo.GetSystemRoles(ctx)
}

// ExistsByName checks if a role exists by name
func (s *roleService) ExistsByName(ctx context.Context, name string) (bool, error) {
	role, err := s.repo.FindByName(ctx, name)
	if err != nil {
		// If the error is "role not found", return false with no error
		if err.Error() == "role not found" {
			return false, nil
		}
		// Otherwise, return the error
		return false, err
	}
	// Role was found, return true
	return role != nil, nil
}

// BulkDeleteRoles deletes multiple roles by their IDs
func (s *roleService) BulkDeleteRoles(ctx context.Context, ids []uuid.UUID) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDeleteRoles(ctx, ids)
}
