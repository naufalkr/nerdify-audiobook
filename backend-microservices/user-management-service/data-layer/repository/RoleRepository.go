package repository

import (
	"context"
	"errors"
	"microservice/user/data-layer/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepository defines the contract for role data access
type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) (*entity.Role, error)
	ListAll(ctx context.Context) ([]entity.Role, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	SeedDefaultRoles(ctx context.Context) error
	Update(ctx context.Context, role *entity.Role) (*entity.Role, error)
	CountRoles(ctx context.Context) (int64, error)
	SearchRoles(ctx context.Context, query string, page, pageSize int) ([]entity.Role, int64, error)
	GetSystemRoles(ctx context.Context) ([]entity.Role, error)
	BulkDeleteRoles(ctx context.Context, ids []uuid.UUID) (int, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("role not found")
	}
	return &role, err
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("role not found")
	}
	return &role, err
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) ListAll(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.WithContext(ctx).Order("created_at asc").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	role, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return errors.New("cannot delete a system-defined role")
	}
	return r.db.WithContext(ctx).Delete(&entity.Role{}, "id = ?", id).Error
}

func (r *roleRepository) SeedDefaultRoles(ctx context.Context) error {
	defaultRoles := []entity.Role{
		{
			ID:          uuid.New(),
			Name:        "SUPERADMIN",
			Description: "Super administrator with all access",
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "ADMIN",
			Description: "Tenant admin with local control",
			IsSystem:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "USER",
			Description: "Regular user",
			IsSystem:    true,
		},
	}

	for _, role := range defaultRoles {
		var existing entity.Role
		err := r.db.WithContext(ctx).Where("name = ?", role.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.WithContext(ctx).Create(&role).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	// Check if the role exists
	var existingRole entity.Role
	err := r.db.WithContext(ctx).First(&existingRole, "id = ?", role.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("role not found")
	}
	if err != nil {
		return nil, err
	}

	// If it's a system role, don't allow changing the name
	if existingRole.IsSystem && existingRole.Name != role.Name {
		return nil, errors.New("cannot change the name of a system role")
	}

	// Update the role
	role.UpdatedAt = time.Now()
	err = r.db.WithContext(ctx).Save(role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) CountRoles(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Role{}).Count(&count).Error
	return count, err
}

func (r *roleRepository) SearchRoles(ctx context.Context, query string, page, pageSize int) ([]entity.Role, int64, error) {
	var roles []entity.Role
	var totalCount int64

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Build query with search conditions
	db := r.db.WithContext(ctx).Model(&entity.Role{})

	// Add search filter if query is provided
	if query != "" {
		searchQuery := "%" + query + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", searchQuery, searchQuery)
	}

	// Count total matching records
	err := db.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Execute query with pagination
	err = db.Order("created_at asc").
		Limit(pageSize).
		Offset(offset).
		Find(&roles).Error

	return roles, totalCount, err
}

func (r *roleRepository) GetSystemRoles(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.WithContext(ctx).
		Where("is_system = ?", true).
		Order("created_at asc").
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) BulkDeleteRoles(ctx context.Context, ids []uuid.UUID) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// Convert UUID slice to string slice for SQL IN clause
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	// Find all system roles in the provided IDs
	var systemRoles []entity.Role
	err := r.db.WithContext(ctx).
		Where("id IN ? AND is_system = ?", ids, true).
		Find(&systemRoles).Error
	if err != nil {
		return 0, err
	}

	// If any system roles were found, return an error
	if len(systemRoles) > 0 {
		return 0, errors.New("cannot delete system-defined roles")
	}

	// Delete non-system roles
	result := r.db.WithContext(ctx).
		Where("id IN ? AND is_system = ?", ids, false).
		Delete(&entity.Role{})
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}
