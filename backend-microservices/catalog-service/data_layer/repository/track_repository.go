package repository

import (
	"catalog-service/data_layer/entity"

	"gorm.io/gorm"
)

// TrackRepositoryInterface defines the contract for track repository
type TrackRepositoryInterface interface {
	Create(track *entity.Track) error
	GetByID(id uint) (*entity.Track, error)
	GetByIDWithRelations(id uint) (*entity.Track, error)
	GetAll(offset, limit int) ([]entity.Track, int64, error)
	Update(track *entity.Track) error
	Delete(id uint) error
	GetByAudiobookID(audiobookID uint) ([]entity.Track, error)
	SearchByTitle(query string, offset, limit int) ([]entity.Track, int64, error)
	DeleteByAudiobookID(audiobookID uint) error
}

// TrackRepository implements TrackRepositoryInterface
type TrackRepository struct {
	db *gorm.DB
}

// NewTrackRepository creates a new track repository
func NewTrackRepository(db *gorm.DB) TrackRepositoryInterface {
	return &TrackRepository{db: db}
}

// Create creates a new track
func (r *TrackRepository) Create(track *entity.Track) error {
	return r.db.Create(track).Error
}

// GetByID retrieves a track by ID
func (r *TrackRepository) GetByID(id uint) (*entity.Track, error) {
	var track entity.Track
	err := r.db.First(&track, id).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}

// GetByIDWithRelations retrieves a track by ID with relations
func (r *TrackRepository) GetByIDWithRelations(id uint) (*entity.Track, error) {
	var track entity.Track
	err := r.db.Preload("Audiobook").First(&track, id).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}

// GetAll retrieves all tracks with pagination
func (r *TrackRepository) GetAll(offset, limit int) ([]entity.Track, int64, error) {
	var tracks []entity.Track
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Track{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&tracks).Error; err != nil {
		return nil, 0, err
	}

	return tracks, total, nil
}

// Update updates an existing track
func (r *TrackRepository) Update(track *entity.Track) error {
	return r.db.Save(track).Error
}

// Delete deletes a track by ID
func (r *TrackRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Track{}, id).Error
}

// GetByAudiobookID retrieves all tracks for a specific audiobook
func (r *TrackRepository) GetByAudiobookID(audiobookID uint) ([]entity.Track, error) {
	var tracks []entity.Track
	err := r.db.Where("audiobook_id = ?", audiobookID).Order("id ASC").Find(&tracks).Error
	return tracks, err
}

// SearchByTitle searches tracks by title
func (r *TrackRepository) SearchByTitle(query string, offset, limit int) ([]entity.Track, int64, error) {
	var tracks []entity.Track
	var total int64

	dbQuery := r.db.Model(&entity.Track{}).Preload("Audiobook")
	if query != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+query+"%")
	}

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&tracks).Error; err != nil {
		return nil, 0, err
	}

	return tracks, total, nil
}

// DeleteByAudiobookID deletes all tracks for a specific audiobook
func (r *TrackRepository) DeleteByAudiobookID(audiobookID uint) error {
	return r.db.Where("audiobook_id = ?", audiobookID).Delete(&entity.Track{}).Error
}
