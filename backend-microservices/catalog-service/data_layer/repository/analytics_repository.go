package repository

import (
	"catalog-service/data_layer/entity"
	"time"

	"gorm.io/gorm"
)

// AnalyticsRepositoryInterface defines the contract for analytics repository
type AnalyticsRepositoryInterface interface {
	Create(analytics *entity.Analytics) error
	GetByID(id uint) (*entity.Analytics, error)
	GetAll(offset, limit int) ([]entity.Analytics, int64, error)
	Delete(id uint) error
	GetByUserID(userID string, offset, limit int) ([]entity.Analytics, int64, error)
	GetByAudiobookID(audiobookID uint, offset, limit int) ([]entity.Analytics, int64, error)
	GetByEventType(eventType string, offset, limit int) ([]entity.Analytics, int64, error)
	GetByDateRange(startDate, endDate time.Time, offset, limit int) ([]entity.Analytics, int64, error)
	GetAnalyticsSummary(audiobookID uint) (map[string]int64, error)
	DeleteByUserID(userID string) error
	DeleteByAudiobookID(audiobookID uint) error
}

// AnalyticsRepository implements AnalyticsRepositoryInterface
type AnalyticsRepository struct {
	db *gorm.DB
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepositoryInterface {
	return &AnalyticsRepository{db: db}
}

// Create creates a new analytics record
func (r *AnalyticsRepository) Create(analytics *entity.Analytics) error {
	return r.db.Create(analytics).Error
}

// GetByID retrieves an analytics record by ID
func (r *AnalyticsRepository) GetByID(id uint) (*entity.Analytics, error) {
	var analytics entity.Analytics
	err := r.db.Preload("Audiobook").Preload("User").First(&analytics, id).Error
	if err != nil {
		return nil, err
	}
	return &analytics, nil
}

// GetAll retrieves all analytics records with pagination
func (r *AnalyticsRepository) GetAll(offset, limit int) ([]entity.Analytics, int64, error) {
	var analytics []entity.Analytics
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Analytics{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Preload("Audiobook").Preload("User").Offset(offset).Limit(limit).Find(&analytics).Error; err != nil {
		return nil, 0, err
	}

	return analytics, total, nil
}

// Delete deletes an analytics record by ID
func (r *AnalyticsRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Analytics{}, id).Error
}

// GetByUserID retrieves analytics records by user ID
func (r *AnalyticsRepository) GetByUserID(userID string, offset, limit int) ([]entity.Analytics, int64, error) {
	var analytics []entity.Analytics
	var total int64

	dbQuery := r.db.Model(&entity.Analytics{}).Where("user_id = ?", userID).Preload("Audiobook").Preload("User")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&analytics).Error; err != nil {
		return nil, 0, err
	}

	return analytics, total, nil
}

// GetByAudiobookID retrieves analytics records by audiobook ID
func (r *AnalyticsRepository) GetByAudiobookID(audiobookID uint, offset, limit int) ([]entity.Analytics, int64, error) {
	var analytics []entity.Analytics
	var total int64

	dbQuery := r.db.Model(&entity.Analytics{}).Where("audiobook_id = ?", audiobookID).Preload("Audiobook").Preload("User")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&analytics).Error; err != nil {
		return nil, 0, err
	}

	return analytics, total, nil
}

// GetByEventType retrieves analytics records by event type
func (r *AnalyticsRepository) GetByEventType(eventType string, offset, limit int) ([]entity.Analytics, int64, error) {
	var analytics []entity.Analytics
	var total int64

	dbQuery := r.db.Model(&entity.Analytics{}).Where("event_type = ?", eventType).Preload("Audiobook").Preload("User")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&analytics).Error; err != nil {
		return nil, 0, err
	}

	return analytics, total, nil
}

// GetByDateRange retrieves analytics records within a date range
func (r *AnalyticsRepository) GetByDateRange(startDate, endDate time.Time, offset, limit int) ([]entity.Analytics, int64, error) {
	var analytics []entity.Analytics
	var total int64

	dbQuery := r.db.Model(&entity.Analytics{}).
		Where("event_timestamp BETWEEN ? AND ?", startDate, endDate).
		Preload("Audiobook").Preload("User")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&analytics).Error; err != nil {
		return nil, 0, err
	}

	return analytics, total, nil
}

// GetAnalyticsSummary retrieves analytics summary for an audiobook
func (r *AnalyticsRepository) GetAnalyticsSummary(audiobookID uint) (map[string]int64, error) {
	summary := make(map[string]int64)

	var results []struct {
		EventType string
		Count     int64
	}

	err := r.db.Model(&entity.Analytics{}).
		Select("event_type, COUNT(*) as count").
		Where("audiobook_id = ?", audiobookID).
		Group("event_type").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		summary[result.EventType] = result.Count
	}

	return summary, nil
}

// DeleteByUserID deletes all analytics records for a user
func (r *AnalyticsRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.Analytics{}).Error
}

// DeleteByAudiobookID deletes all analytics records for an audiobook
func (r *AnalyticsRepository) DeleteByAudiobookID(audiobookID uint) error {
	return r.db.Where("audiobook_id = ?", audiobookID).Delete(&entity.Analytics{}).Error
}
