package repository

import (
	"content-management-service/data_layer/entity"

	"gorm.io/gorm"
)

// GenreRepositoryInterface defines the contract for genre repository
type GenreRepositoryInterface interface {
	Create(genre *entity.Genre) error
	GetByID(id uint) (*entity.Genre, error)
	GetAll(offset, limit int) ([]entity.Genre, int64, error)
	Update(genre *entity.Genre) error
	Delete(id uint) error
	SearchByName(query string, offset, limit int) ([]entity.Genre, int64, error)
	ExistsByName(name string) (bool, error)
	GetByIDs(ids []uint) ([]entity.Genre, error)
}

// GenreRepository implements GenreRepositoryInterface
type GenreRepository struct {
	db *gorm.DB
}

// NewGenreRepository creates a new genre repository
func NewGenreRepository(db *gorm.DB) GenreRepositoryInterface {
	return &GenreRepository{db: db}
}

// Create creates a new genre
func (r *GenreRepository) Create(genre *entity.Genre) error {
	return r.db.Create(genre).Error
}

// GetByID retrieves a genre by ID
func (r *GenreRepository) GetByID(id uint) (*entity.Genre, error) {
	var genre entity.Genre
	err := r.db.First(&genre, id).Error
	if err != nil {
		return nil, err
	}
	return &genre, nil
}

// GetAll retrieves all genres with pagination
func (r *GenreRepository) GetAll(offset, limit int) ([]entity.Genre, int64, error) {
	var genres []entity.Genre
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Genre{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&genres).Error; err != nil {
		return nil, 0, err
	}

	return genres, total, nil
}

// Update updates an existing genre
func (r *GenreRepository) Update(genre *entity.Genre) error {
	return r.db.Save(genre).Error
}

// Delete deletes a genre by ID
func (r *GenreRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Genre{}, id).Error
}

// SearchByName searches genres by name
func (r *GenreRepository) SearchByName(query string, offset, limit int) ([]entity.Genre, int64, error) {
	var genres []entity.Genre
	var total int64

	dbQuery := r.db.Model(&entity.Genre{})
	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query+"%")
	}

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&genres).Error; err != nil {
		return nil, 0, err
	}

	return genres, total, nil
}

// ExistsByName checks if a genre exists by name
func (r *GenreRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Genre{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// GetByIDs retrieves genres by multiple IDs
func (r *GenreRepository) GetByIDs(ids []uint) ([]entity.Genre, error) {
	var genres []entity.Genre
	err := r.db.Where("id IN ?", ids).Find(&genres).Error
	return genres, err
}
