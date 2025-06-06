package repository

import (
	"content-management-service/data_layer/entity"

	"gorm.io/gorm"
)

// ReaderRepositoryInterface defines the contract for reader repository
type ReaderRepositoryInterface interface {
	Create(reader *entity.Reader) error
	GetByID(id uint) (*entity.Reader, error)
	GetAll(offset, limit int) ([]entity.Reader, int64, error)
	Update(reader *entity.Reader) error
	Delete(id uint) error
	SearchByName(query string, offset, limit int) ([]entity.Reader, int64, error)
	ExistsByName(name string) (bool, error)
}

// ReaderRepository implements ReaderRepositoryInterface
type ReaderRepository struct {
	db *gorm.DB
}

// NewReaderRepository creates a new reader repository
func NewReaderRepository(db *gorm.DB) ReaderRepositoryInterface {
	return &ReaderRepository{db: db}
}

// Create creates a new reader
func (r *ReaderRepository) Create(reader *entity.Reader) error {
	return r.db.Create(reader).Error
}

// GetByID retrieves a reader by ID
func (r *ReaderRepository) GetByID(id uint) (*entity.Reader, error) {
	var reader entity.Reader
	err := r.db.First(&reader, id).Error
	if err != nil {
		return nil, err
	}
	return &reader, nil
}

// GetAll retrieves all readers with pagination
func (r *ReaderRepository) GetAll(offset, limit int) ([]entity.Reader, int64, error) {
	var readers []entity.Reader
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Reader{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&readers).Error; err != nil {
		return nil, 0, err
	}

	return readers, total, nil
}

// Update updates an existing reader
func (r *ReaderRepository) Update(reader *entity.Reader) error {
	return r.db.Save(reader).Error
}

// Delete deletes a reader by ID
func (r *ReaderRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Reader{}, id).Error
}

// SearchByName searches readers by name
func (r *ReaderRepository) SearchByName(query string, offset, limit int) ([]entity.Reader, int64, error) {
	var readers []entity.Reader
	var total int64

	dbQuery := r.db.Model(&entity.Reader{})
	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query+"%")
	}

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&readers).Error; err != nil {
		return nil, 0, err
	}

	return readers, total, nil
}

// ExistsByName checks if a reader exists by name
func (r *ReaderRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Reader{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
