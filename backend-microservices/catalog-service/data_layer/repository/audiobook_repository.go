package repository

import (
	"catalog-service/data_layer/entity"

	"gorm.io/gorm"
)

// AudiobookRepositoryInterface defines the contract for audiobook repository
type AudiobookRepositoryInterface interface {
	Create(audiobook *entity.Audiobook) error
	GetByID(id uint) (*entity.Audiobook, error)
	GetByIDWithRelations(id uint) (*entity.Audiobook, error)
	GetAll(offset, limit int) ([]entity.Audiobook, int64, error)
	GetAllWithRelations(offset, limit int) ([]entity.Audiobook, int64, error)
	Update(audiobook *entity.Audiobook) error
	Delete(id uint) error
	SearchByTitle(query string, offset, limit int) ([]entity.Audiobook, int64, error)
	GetByAuthorID(authorID uint, offset, limit int) ([]entity.Audiobook, int64, error)
	GetByReaderID(readerID uint, offset, limit int) ([]entity.Audiobook, int64, error)
	GetByGenreID(genreID uint, offset, limit int) ([]entity.Audiobook, int64, error)
	AssignGenres(audiobookID uint, genreIDs []uint) error
	RemoveGenres(audiobookID uint, genreIDs []uint) error
	RemoveAllGenres(audiobookID uint) error
}

// AudiobookRepository implements AudiobookRepositoryInterface
type AudiobookRepository struct {
	db *gorm.DB
}

// NewAudiobookRepository creates a new audiobook repository
func NewAudiobookRepository(db *gorm.DB) AudiobookRepositoryInterface {
	return &AudiobookRepository{db: db}
}

// Create creates a new audiobook
func (r *AudiobookRepository) Create(audiobook *entity.Audiobook) error {
	return r.db.Create(audiobook).Error
}

// GetByID retrieves an audiobook by ID
func (r *AudiobookRepository) GetByID(id uint) (*entity.Audiobook, error) {
	var audiobook entity.Audiobook
	err := r.db.First(&audiobook, id).Error
	if err != nil {
		return nil, err
	}
	return &audiobook, nil
}

// GetByIDWithRelations retrieves an audiobook by ID with all relations
func (r *AudiobookRepository) GetByIDWithRelations(id uint) (*entity.Audiobook, error) {
	var audiobook entity.Audiobook
	err := r.db.Preload("Author").Preload("Reader").Preload("Genres").Preload("Tracks").First(&audiobook, id).Error
	if err != nil {
		return nil, err
	}
	return &audiobook, nil
}

// GetAll retrieves all audiobooks with pagination
func (r *AudiobookRepository) GetAll(offset, limit int) ([]entity.Audiobook, int64, error) {
	var audiobooks []entity.Audiobook
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Audiobook{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&audiobooks).Error; err != nil {
		return nil, 0, err
	}

	return audiobooks, total, nil
}

// GetAllWithRelations retrieves all audiobooks with relations and pagination
func (r *AudiobookRepository) GetAllWithRelations(offset, limit int) ([]entity.Audiobook, int64, error) {
	var audiobooks []entity.Audiobook
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Audiobook{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with relations
	if err := r.db.Preload("Author").Preload("Reader").Preload("Genres").Offset(offset).Limit(limit).Find(&audiobooks).Error; err != nil {
		return nil, 0, err
	}

	return audiobooks, total, nil
}

// Update updates an existing audiobook
func (r *AudiobookRepository) Update(audiobook *entity.Audiobook) error {
	return r.db.Save(audiobook).Error
}

// Delete deletes an audiobook by ID
func (r *AudiobookRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Audiobook{}, id).Error
}

// SearchByTitle searches audiobooks by title
func (r *AudiobookRepository) SearchByTitle(query string, offset, limit int) ([]entity.Audiobook, int64, error) {
	var audiobooks []entity.Audiobook
	var total int64

	dbQuery := r.db.Model(&entity.Audiobook{}).Preload("Author").Preload("Reader").Preload("Genres")
	if query != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+query+"%")
	}

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&audiobooks).Error; err != nil {
		return nil, 0, err
	}

	return audiobooks, total, nil
}

// GetByAuthorID retrieves audiobooks by author ID
func (r *AudiobookRepository) GetByAuthorID(authorID uint, offset, limit int) ([]entity.Audiobook, int64, error) {
	var audiobooks []entity.Audiobook
	var total int64

	dbQuery := r.db.Model(&entity.Audiobook{}).Where("author_id = ?", authorID).Preload("Author").Preload("Reader").Preload("Genres")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&audiobooks).Error; err != nil {
		return nil, 0, err
	}

	return audiobooks, total, nil
}

// GetByReaderID retrieves audiobooks by reader ID
func (r *AudiobookRepository) GetByReaderID(readerID uint, offset, limit int) ([]entity.Audiobook, int64, error) {
	var audiobooks []entity.Audiobook
	var total int64

	dbQuery := r.db.Model(&entity.Audiobook{}).Where("reader_id = ?", readerID).Preload("Author").Preload("Reader").Preload("Genres")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&audiobooks).Error; err != nil {
		return nil, 0, err
	}

	return audiobooks, total, nil
}

// GetByGenreID retrieves audiobooks by genre ID
func (r *AudiobookRepository) GetByGenreID(genreID uint, offset, limit int) ([]entity.Audiobook, int64, error) {
	var audiobooks []entity.Audiobook
	var total int64

	dbQuery := r.db.Model(&entity.Audiobook{}).
		Joins("JOIN audiobook_genres ON audiobooks.id = audiobook_genres.audiobook_id").
		Where("audiobook_genres.genre_id = ?", genreID).
		Preload("Author").Preload("Reader").Preload("Genres")

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&audiobooks).Error; err != nil {
		return nil, 0, err
	}

	return audiobooks, total, nil
}

// AssignGenres assigns genres to an audiobook
func (r *AudiobookRepository) AssignGenres(audiobookID uint, genreIDs []uint) error {
	var audiobook entity.Audiobook
	if err := r.db.First(&audiobook, audiobookID).Error; err != nil {
		return err
	}

	var genres []entity.Genre
	if err := r.db.Where("id IN ?", genreIDs).Find(&genres).Error; err != nil {
		return err
	}

	return r.db.Model(&audiobook).Association("Genres").Append(genres)
}

// RemoveGenres removes genres from an audiobook
func (r *AudiobookRepository) RemoveGenres(audiobookID uint, genreIDs []uint) error {
	var audiobook entity.Audiobook
	if err := r.db.First(&audiobook, audiobookID).Error; err != nil {
		return err
	}

	var genres []entity.Genre
	if err := r.db.Where("id IN ?", genreIDs).Find(&genres).Error; err != nil {
		return err
	}

	return r.db.Model(&audiobook).Association("Genres").Delete(genres)
}

// RemoveAllGenres removes all genres from an audiobook
func (r *AudiobookRepository) RemoveAllGenres(audiobookID uint) error {
    var audiobook entity.Audiobook
    if err := r.db.First(&audiobook, audiobookID).Error; err != nil {
        return err
    }

    // Clear all genre associations
    return r.db.Model(&audiobook).Association("Genres").Clear()
}
