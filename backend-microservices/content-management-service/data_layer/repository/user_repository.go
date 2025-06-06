package repository

import (
	"content-management-service/data_layer/entity"

	"gorm.io/gorm"
)

// UserRepositoryInterface defines the contract for user repository
type UserRepositoryInterface interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetAll(offset, limit int) ([]entity.User, int64, error)
	Update(user *entity.User) error
	Delete(id string) error
	ExistsByID(id string) (bool, error)
	GetByRole(role string, offset, limit int) ([]entity.User, int64, error)
}

// UserRepository implements UserRepositoryInterface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users with pagination
func (r *UserRepository) GetAll(offset, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	// Count total records
	if err := r.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}

// ExistsByID checks if a user exists by ID
func (r *UserRepository) ExistsByID(id string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetByRole retrieves users by role with pagination
func (r *UserRepository) GetByRole(role string, offset, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	dbQuery := r.db.Model(&entity.User{}).Where("role = ?", role)

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := dbQuery.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
