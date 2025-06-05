package repository

import (
	"context"
	"log"
	"microservice/user/data-layer/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TenantRepository adalah interface untuk mengakses data tenant
type TenantRepository interface {
	Create(ctx context.Context, tx *gorm.DB, tenant entity.Tenant) (entity.Tenant, error)
	FindAll(ctx context.Context) ([]entity.Tenant, error)
	FindByID(ctx context.Context, id uuid.UUID) (entity.Tenant, error)
	Update(ctx context.Context, tx *gorm.DB, tenant entity.Tenant) (entity.Tenant, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GetTenantsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Tenant, error)
	AddUserToTenant(ctx context.Context, tx *gorm.DB, userTenant entity.UserTenant) error
	RemoveUserFromTenant(ctx context.Context, tx *gorm.DB, userID, tenantID uuid.UUID) error
	CountTenantUsers(ctx context.Context, tenantID uuid.UUID) (int, error)
	IsUserInTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error)
	GetUsersByTenantID(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.User, int, error)
}

// tenantRepository adalah implementasi TenantRepository
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository membuat instance baru TenantRepository
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db}
}

// Create membuat tenant baru
func (r *tenantRepository) Create(ctx context.Context, tx *gorm.DB, tenant entity.Tenant) (entity.Tenant, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.Create(&tenant).Error; err != nil {
		return entity.Tenant{}, err
	}

	return tenant, nil
}

// FindAll mengembalikan semua tenant
func (r *tenantRepository) FindAll(ctx context.Context) ([]entity.Tenant, error) {
	var tenants []entity.Tenant
	if err := r.db.Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

// FindByID mengembalikan tenant dengan ID tertentu
func (r *tenantRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.Tenant, error) {
	var tenant entity.Tenant
	if err := r.db.Where("id = ?", id).First(&tenant).Error; err != nil {
		return entity.Tenant{}, err
	}
	return tenant, nil
}

// Update memperbarui tenant
func (r *tenantRepository) Update(ctx context.Context, tx *gorm.DB, tenant entity.Tenant) (entity.Tenant, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.Save(&tenant).Error; err != nil {
		return entity.Tenant{}, err
	}
	return tenant, nil
}

// SoftDelete melakukan soft delete pada tenant
func (r *tenantRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(&entity.Tenant{}, "id = ?", id).Error
}

// GetTenantsByUserID mengembalikan tenant yang diikuti oleh user
func (r *tenantRepository) GetTenantsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Tenant, error) {
	var tenants []entity.Tenant
	if err := r.db.Joins("JOIN user_tenants ON tenants.id = user_tenants.tenant_id").
		Where("user_tenants.user_id = ? AND user_tenants.is_active = ?", userID, true).
		Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

// AddUserToTenant menambahkan user ke tenant
func (r *tenantRepository) AddUserToTenant(ctx context.Context, tx *gorm.DB, userTenant entity.UserTenant) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(&userTenant).Error
}

// RemoveUserFromTenant menghapus user dari tenant
func (r *tenantRepository) RemoveUserFromTenant(ctx context.Context, tx *gorm.DB, userID, tenantID uuid.UUID) error {
	log.Printf("Starting RemoveUserFromTenant in repository: userID=%s, tenantID=%s", userID, tenantID)

	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.Model(&entity.UserTenant{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Update("is_active", false).
		Error
	if err != nil {
		log.Printf("Error performing soft delete in repository: %v", err)
		return err
	}

	log.Printf("Successfully performed soft delete in repository")
	return nil
}

// CountTenantUsers menghitung jumlah user di tenant
func (r *tenantRepository) CountTenantUsers(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.Model(&entity.UserTenant{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// IsUserInTenant memeriksa apakah user merupakan anggota tenant
func (r *tenantRepository) IsUserInTenant(ctx context.Context, userID, tenantID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.UserTenant{}).
		Where("user_id = ? AND tenant_id = ? AND is_active = ?", userID, tenantID, true).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUsersByTenantID mengembalikan daftar user di tenant dengan pagination
func (r *tenantRepository) GetUsersByTenantID(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.User, int, error) {
	var users []entity.User
	var count int64

	// Hitung total users
	if err := r.db.Model(&entity.User{}).
		Joins("JOIN user_tenants ON users.id = user_tenants.user_id").
		Where("user_tenants.tenant_id = ? AND user_tenants.is_active = ?", tenantID, true).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get users with pagination
	if err := r.db.Preload("Role").
		Joins("JOIN user_tenants ON users.id = user_tenants.user_id").
		Where("user_tenants.tenant_id = ? AND user_tenants.is_active = ?", tenantID, true).
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(count), nil
}
