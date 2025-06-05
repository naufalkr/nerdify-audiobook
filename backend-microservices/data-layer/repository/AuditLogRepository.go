package repository

import (
	"context"
	"microservice/user/data-layer/entity"

	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(ctx context.Context, tx *gorm.DB, log entity.AuditLog) error
	FindByID(ctx context.Context, id string) (entity.AuditLog, error)
	FindByEntityTypeAndID(ctx context.Context, entityType string, entityID string) ([]entity.AuditLog, error)
	FindByUserID(ctx context.Context, userID string) ([]entity.AuditLog, error)
	FindByTenantID(ctx context.Context, tenantID string) ([]entity.AuditLog, error)
	FindAll(ctx context.Context, limit, offset int) ([]entity.AuditLog, int, error)
	GetDB() *gorm.DB
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, tx *gorm.DB, log entity.AuditLog) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&log).Error
	}
	return r.db.WithContext(ctx).Create(&log).Error
}

func (r *auditLogRepository) FindByID(ctx context.Context, id string) (entity.AuditLog, error) {
	var log entity.AuditLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&log).Error
	return log, err
}

func (r *auditLogRepository) FindByEntityTypeAndID(ctx context.Context, entityType string, entityID string) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByUserID(ctx context.Context, userID string) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByTenantID(ctx context.Context, tenantID string) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.AuditLog, int, error) {
	var logs []entity.AuditLog
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&entity.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, int(total), err
}

func (r *auditLogRepository) GetDB() *gorm.DB {
	return r.db
}
