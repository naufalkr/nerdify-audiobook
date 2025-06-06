package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"microservice/user/data-layer/entity"
)

// UserTenantRepository interface untuk mengelola user-tenant relationships
type UserTenantRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, tx *gorm.DB, userTenant entity.UserTenant) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.UserTenant, error)
	Update(ctx context.Context, tx *gorm.DB, userTenant entity.UserTenant) (entity.UserTenant, error)
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	// User-Tenant specific operations
	GetUserTenantByUserIDAndTenantID(userID, tenantID string) (entity.UserTenant, error)
	GetUserTenants(userID string) ([]entity.UserTenant, error)
	GetUserCurrentTenant(userID string) (entity.UserTenant, error)
	UpdateUserCurrentTenant(userID, tenantID string) error

	// Tenant users operations with pagination
	GetTenantUsersWithPagination(tenantID string, page, limit int) ([]entity.UserTenant, int, error)

	// Validation operations
	IsUserInTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error)
	HasUserAccessToTenant(userID, tenantID string) (bool, error)

	// Batch operations
	GetUserTenantsByTenantID(tenantID string) ([]entity.UserTenant, error)
	DeactivateUserFromTenant(ctx context.Context, tx *gorm.DB, userID, tenantID uuid.UUID) error
	ActivateUserInTenant(ctx context.Context, tx *gorm.DB, userID, tenantID uuid.UUID) error
}

type userTenantRepository struct {
	db *gorm.DB
}

func NewUserTenantRepository(db *gorm.DB) UserTenantRepository {
	return &userTenantRepository{db: db}
}

func (r *userTenantRepository) Create(ctx context.Context, tx *gorm.DB, userTenant entity.UserTenant) error {
	exec := r.db
	if tx != nil {
		exec = tx
	}
	return exec.WithContext(ctx).Create(&userTenant).Error
}

func (r *userTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (entity.UserTenant, error) {
	var userTenant entity.UserTenant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&userTenant).Error
	return userTenant, err
}

func (r *userTenantRepository) Update(ctx context.Context, tx *gorm.DB, userTenant entity.UserTenant) (entity.UserTenant, error) {
	exec := r.db
	if tx != nil {
		exec = tx
	}

	userTenant.UpdatedAt = time.Now()
	err := exec.WithContext(ctx).Save(&userTenant).Error
	return userTenant, err
}

func (r *userTenantRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	exec := r.db
	if tx != nil {
		exec = tx
	}
	return exec.WithContext(ctx).Delete(&entity.UserTenant{}, "id = ?", id).Error
}

func (r *userTenantRepository) GetUserTenantByUserIDAndTenantID(userID, tenantID string) (entity.UserTenant, error) {
	var userTenant entity.UserTenant
	err := r.db.Where("user_id = ? AND tenant_id = ? AND is_active = ?", userID, tenantID, true).First(&userTenant).Error
	return userTenant, err
}

func (r *userTenantRepository) GetUserTenants(userID string) ([]entity.UserTenant, error) {
	var userTenants []entity.UserTenant
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&userTenants).Error
	return userTenants, err
}

func (r *userTenantRepository) GetUserCurrentTenant(userID string) (entity.UserTenant, error) {
	var userTenant entity.UserTenant
	// Untuk saat ini, ambil tenant pertama yang aktif
	// TODO: Implementasi field 'is_current' atau mechanism untuk menentukan current tenant
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&userTenant).Error
	return userTenant, err
}

func (r *userTenantRepository) UpdateUserCurrentTenant(userID, tenantID string) error {
	// Start a transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First, set all user's tenant relationships to not current (jika ada field is_current)
	// Untuk saat ini, kita skip ini karena belum ada field is_current di entity

	// Verify the user has access to the requested tenant
	var userTenant entity.UserTenant
	err := tx.Where("user_id = ? AND tenant_id = ? AND is_active = ?", userID, tenantID, true).First(&userTenant).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the tenant's updated_at to mark it as recently accessed
	userTenant.UpdatedAt = time.Now()
	err = tx.Save(&userTenant).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *userTenantRepository) GetTenantUsersWithPagination(tenantID string, page, limit int) ([]entity.UserTenant, int, error) {
	var userTenants []entity.UserTenant
	var count int64

	// Count total records
	err := r.db.Model(&entity.UserTenant{}).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated results
	err = r.db.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&userTenants).Error
	if err != nil {
		return nil, 0, err
	}

	return userTenants, int(count), nil
}

func (r *userTenantRepository) IsUserInTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.UserTenant{}).
		Where("user_id = ? AND tenant_id = ? AND is_active = ?", userID, tenantID, true).
		Count(&count).Error
	return count > 0, err
}

func (r *userTenantRepository) HasUserAccessToTenant(userID, tenantID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.UserTenant{}).
		Where("user_id = ? AND tenant_id = ? AND is_active = ?", userID, tenantID, true).
		Count(&count).Error
	return count > 0, err
}

func (r *userTenantRepository) GetUserTenantsByTenantID(tenantID string) ([]entity.UserTenant, error) {
	var userTenants []entity.UserTenant
	err := r.db.Where("tenant_id = ? AND is_active = ?", tenantID, true).Find(&userTenants).Error
	return userTenants, err
}

func (r *userTenantRepository) DeactivateUserFromTenant(ctx context.Context, tx *gorm.DB, userID, tenantID uuid.UUID) error {
	exec := r.db
	if tx != nil {
		exec = tx
	}

	return exec.WithContext(ctx).Model(&entity.UserTenant{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Update("is_active", false).Error
}

func (r *userTenantRepository) ActivateUserInTenant(ctx context.Context, tx *gorm.DB, userID, tenantID uuid.UUID) error {
	exec := r.db
	if tx != nil {
		exec = tx
	}

	return exec.WithContext(ctx).Model(&entity.UserTenant{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Update("is_active", true).Error
}
