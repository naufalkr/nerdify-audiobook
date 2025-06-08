package repository

import (
	"catalog-service/data_layer/entity"

	"gorm.io/gorm"
)

// AuthorRepositoryInterface defines the contract for author repository
type AuthorRepositoryInterface interface {
	Create(author *entity.Author) error
	GetByID(id uint) (*entity.Author, error)
	GetAll(offset, limit int) ([]entity.Author, int64, error)
	Update(author *entity.Author) error
	Delete(id uint) error
	SearchByName(query string, offset, limit int) ([]entity.Author, int64, error)
	ExistsByName(name string) (bool, error)
}

// AuthorRepository implements AuthorRepositoryInterface
type AuthorRepository struct {
	db *gorm.DB
}

// NewAuthorRepository creates a new author repository
func NewAuthorRepository(db *gorm.DB) AuthorRepositoryInterface {
	return &AuthorRepository{db: db}
}

// Create creates a new author
func (r *AuthorRepository) Create(author *entity.Author) error {
	return r.db.Create(author).Error
}

// GetByID retrieves an author by ID
func (r *AuthorRepository) GetByID(id uint) (*entity.Author, error) {
	var author entity.Author
	err := r.db.First(&author, id).Error
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// GetAll retrieves all authors with pagination
func (r *AuthorRepository) GetAll(offset, limit int) ([]entity.Author, int64, error) {
	var authors []entity.Author
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Author{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&authors).Error; err != nil {
		return nil, 0, err
	}

	return authors, total, nil
}

// Update updates an existing author
func (r *AuthorRepository) Update(author *entity.Author) error {
	return r.db.Save(author).Error
}

// Delete deletes an author by ID
func (r *AuthorRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Author{}, id).Error
}

// SearchByName searches authors by name
func (r *AuthorRepository) SearchByName(query string, offset, limit int) ([]entity.Author, int64, error) {
	var authors []entity.Author
	var total int64

	dbQuery := r.db.Model(&entity.Author{})
	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query+"%")
	}

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&authors).Error; err != nil {
		return nil, 0, err
	}

	return authors, total, nil
}

// ExistsByName checks if an author exists by name
func (r *AuthorRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Author{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
