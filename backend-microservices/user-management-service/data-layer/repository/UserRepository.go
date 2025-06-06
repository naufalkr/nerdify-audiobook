package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"microservice/user/data-layer/entity"
	"microservice/user/helpers/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(ctx context.Context, tx *gorm.DB) error
	RollbackTx(ctx context.Context, tx *gorm.DB)

	CreateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error)
	FindUserByUsername(ctx context.Context, tx *gorm.DB, username string) (entity.User, error)
	FindUserById(ctx context.Context, tx *gorm.DB, id string) (entity.User, error)
	Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	GetUserRole(ctx context.Context, tx *gorm.DB, roleID uuid.UUID) (entity.Role, error)
	FindAllUser(ctx context.Context, offset int, limit int, tx *gorm.DB) ([]entity.User, int64, error)
	FindAllUsersByKeyword(ctx context.Context, offset int, limit int, keyword string, tx *gorm.DB) ([]entity.User, int64, error)
	SoftDeleteUserByID(ctx context.Context, id string, tx *gorm.DB) error
	HardDeleteUserByID(ctx context.Context, id string, tx *gorm.DB) error
	ExistsByEmailOrUsername(ctx context.Context, tx *gorm.DB, email, username string) (bool, error)

	FindAll(ctx context.Context, page, limit int) ([]entity.User, int, error)
	ForceRefreshUserByID(ctx context.Context, id string) (entity.User, error)
	FindUserByIdWithLock(ctx context.Context, tx *gorm.DB, id string) (entity.User, error)
	UpdateUserRoleDirectSQL(ctx context.Context, tx *gorm.DB, userID string, roleID string) error
	ExecuteRawUpdateQuery(ctx context.Context, userID string, roleID string) error
	GetUserRoleInTenant(ctx context.Context, userID uuid.UUID, tenantID uuid.UUID) (entity.Role, error)
	FindUsersByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	return r.db.WithContext(ctx).Begin(), nil
}

func (r *userRepository) CommitTx(ctx context.Context, tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *userRepository) RollbackTx(ctx context.Context, tx *gorm.DB) {
	tx.Rollback()
}

func (r *userRepository) CreateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	exec := r.db
	if tx != nil {
		exec = tx
	}
	// Ensure OTPCreatedAt and VerificationTokenCreatedAt are nil if not set
	if user.OTPCreatedAt.IsZero() {
		user.OTPCreatedAt = nil
	}
	if user.VerificationTokenCreatedAt.IsZero() {
		user.VerificationTokenCreatedAt = nil
	}
	err := exec.WithContext(ctx).Create(&user).Error
	return user, err
}

func (r *userRepository) FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, error) {
	if email == "" {
		return entity.User{}, fmt.Errorf("email cannot be empty")
	}

	var user entity.User
	exec := r.db
	if tx != nil {
		exec = tx
	}
	err := exec.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found with email: %s", email)
			return entity.User{}, fmt.Errorf("user not found: %w", dto.ErrUserNotFound)
		}
		log.Printf("Error finding user by email %s: %v", email, err)
		return entity.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindUserByUsername(ctx context.Context, tx *gorm.DB, username string) (entity.User, error) {
	var user entity.User
	exec := r.db
	if tx != nil {
		exec = tx
	}
	err := exec.WithContext(ctx).
		Preload("Role").
		Where("user_name = ?", username).
		First(&user).Error
	return user, err
}

func (r *userRepository) FindUserById(ctx context.Context, tx *gorm.DB, id string) (entity.User, error) {
	var user entity.User
	exec := r.db
	if tx != nil {
		exec = tx
	}
	err := exec.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error
	return user, err
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	log.Printf("DEBUG: Repository - updating user: ID=%s, Username=%s", user.ID.String(), user.UserName)
	exec := r.db
	if tx != nil {
		exec = tx
	}
	result := exec.WithContext(ctx).Save(&user)
	if result.Error != nil {
		log.Printf("DEBUG: Repository - error updating user: %v", result.Error)
		return user, result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("DEBUG: Repository - warning: no rows affected when updating user")
	} else {
		log.Printf("DEBUG: Repository - successfully updated user with %d rows affected", result.RowsAffected)
	}
	return user, result.Error
}

func (r *userRepository) GetUserRole(ctx context.Context, tx *gorm.DB, roleID uuid.UUID) (entity.Role, error) {
	var role entity.Role
	exec := r.db
	if tx != nil {
		exec = tx
	}
	err := exec.WithContext(ctx).First(&role, "id = ?", roleID).Error
	return role, err
}

func (r *userRepository) FindAllUser(ctx context.Context, offset int, limit int, tx *gorm.DB) ([]entity.User, int64, error) {
	exec := r.db
	if tx != nil {
		exec = tx
	}
	var users []entity.User
	var count int64
	err := exec.WithContext(ctx).
		Model(&entity.User{}).
		Count(&count).
		Limit(limit).
		Offset(offset).
		Preload("Role").
		Find(&users).Error
	return users, count, err
}

func (r *userRepository) FindAllUsersByKeyword(ctx context.Context, offset int, limit int, keyword string, tx *gorm.DB) ([]entity.User, int64, error) {
	exec := r.db
	if tx != nil {
		exec = tx
	}
	var users []entity.User
	var count int64
	err := exec.WithContext(ctx).
		Model(&entity.User{}).
		Where("user_name ILIKE ? OR email ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Count(&count).
		Limit(limit).
		Offset(offset).
		Preload("Role").
		Find(&users).Error
	return users, count, err
}

func (r *userRepository) SoftDeleteUserByID(ctx context.Context, id string, tx *gorm.DB) error {
	exec := r.db
	if tx != nil {
		exec = tx
	}

	// Add debug logging
	log.Printf("SoftDeleteUserByID: Attempting to soft delete user with ID: %s", id)

	// Parse the ID to ensure it's a valid UUID
	userUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("SoftDeleteUserByID: Invalid user ID format: %v", err)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// First, check if the user exists
	var user entity.User
	if err := exec.WithContext(ctx).Where("id = ?", userUUID).First(&user).Error; err != nil {
		log.Printf("SoftDeleteUserByID: User not found: %v", err)
		return err
	}

	// Execute the soft delete with result capture
	result := exec.WithContext(ctx).Where("id = ?", userUUID).Delete(&entity.User{})
	if result.Error != nil {
		log.Printf("SoftDeleteUserByID: Error performing soft delete: %v", result.Error)
		return result.Error
	}

	// Log the number of affected rows to validate the operation
	log.Printf("SoftDeleteUserByID: Soft delete complete. Rows affected: %d", result.RowsAffected)

	// Verify the delete was successful by trying to fetch the record (with Unscoped)
	var count int64
	exec.WithContext(ctx).Model(&entity.User{}).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", userUUID).Count(&count)
	log.Printf("SoftDeleteUserByID: Verification - Found %d soft-deleted records with this ID", count)

	return nil
}

func (r *userRepository) HardDeleteUserByID(ctx context.Context, id string, tx *gorm.DB) error {
	exec := r.db
	if tx != nil {
		exec = tx
	}

	// Add debug logging
	log.Printf("HardDeleteUserByID: Attempting to permanently delete user with ID: %s", id)

	// Parse the ID to ensure it's a valid UUID
	userUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("HardDeleteUserByID: Invalid user ID format: %v", err)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// First, check if the user exists
	var user entity.User
	if err := exec.WithContext(ctx).Unscoped().Where("id = ?", userUUID).First(&user).Error; err != nil {
		log.Printf("HardDeleteUserByID: User not found: %v", err)
		return err
	}

	// Execute the hard delete with result capture
	result := exec.WithContext(ctx).Unscoped().Where("id = ?", userUUID).Delete(&entity.User{})
	if result.Error != nil {
		log.Printf("HardDeleteUserByID: Error performing hard delete: %v", result.Error)
		return result.Error
	}

	// Log the number of affected rows to validate the operation
	log.Printf("HardDeleteUserByID: Hard delete complete. Rows affected: %d", result.RowsAffected)

	// Verify the delete was successful by trying to fetch the record (with Unscoped)
	var count int64
	exec.WithContext(ctx).Model(&entity.User{}).Unscoped().Where("id = ?", userUUID).Count(&count)
	if count > 0 {
		log.Printf("HardDeleteUserByID: Warning - Record still exists after hard delete!")
	} else {
		log.Printf("HardDeleteUserByID: Verification - Record successfully deleted from database")
	}

	return nil
}

func (r *userRepository) ExistsByEmailOrUsername(ctx context.Context, tx *gorm.DB, email, username string) (bool, error) {
	exec := r.db
	if tx != nil {
		exec = tx
	}

	var count int64
	err := exec.WithContext(ctx).
		Model(&entity.User{}).
		Where("(email = ? OR user_name = ?) AND deleted_at IS NULL", email, username).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepository) FindAll(ctx context.Context, page, limit int) ([]entity.User, int, error) {
	var users []entity.User
	var count int64

	// Hitung total users
	if err := r.db.Model(&entity.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get users with pagination
	if err := r.db.Preload("Role").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(count), nil
}

// Add a new method to force refresh user data
func (r *userRepository) ForceRefreshUserByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error
	return user, err
}

// FindUserByIdWithLock gets user by ID with an exclusive row lock for update
func (r *userRepository) FindUserByIdWithLock(ctx context.Context, tx *gorm.DB, id string) (entity.User, error) {
	var user entity.User
	exec := r.db
	if tx != nil {
		exec = tx
	}

	err := exec.WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error

	return user, err
}

// UpdateUserRoleDirectSQL performs a direct SQL update to change a user's role
func (r *userRepository) UpdateUserRoleDirectSQL(ctx context.Context, tx *gorm.DB, userID string, roleID string) error {
	exec := r.db
	if tx != nil {
		exec = tx
	}

	// Get the correct table name from GORM's naming strategy
	tableName := r.db.NamingStrategy.TableName("users")
	log.Printf("DEBUG: Using table name: %s for role update", tableName)

	// Debug: Check current state before update
	var currentUser struct {
		ID     string
		RoleID string
	}
	if err := exec.WithContext(ctx).Table(tableName).Where("id = ?", userID).
		Select("id, role_id").Scan(&currentUser).Error; err != nil {
		log.Printf("DEBUG: Error checking current user state: %v", err)
	} else {
		log.Printf("DEBUG: Current user state - ID: %s, RoleID: %s", currentUser.ID, currentUser.RoleID)
	}

	// Debug: Check if role exists
	var roleExists bool
	if err := exec.WithContext(ctx).Table("roles").Where("id = ?", roleID).
		Select("1").Scan(&roleExists).Error; err != nil {
		log.Printf("DEBUG: Error checking role existence: %v", err)
	} else {
		log.Printf("DEBUG: Role ID %s exists: %v", roleID, roleExists)
	}

	// Debug the query we're about to execute
	updateSQL := fmt.Sprintf("UPDATE %s SET role_id = '%s', updated_at = '%s' WHERE id = '%s'",
		tableName, roleID, time.Now().Format("2006-01-02 15:04:05"), userID)
	log.Printf("DEBUG: SQL Update Query: %s", updateSQL)

	// Execute the update with explicit table name
	result := exec.WithContext(ctx).
		Debug().
		Table(tableName).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"role_id":    roleID,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Printf("DEBUG: SQL Update Error: %v", result.Error)
		return result.Error
	}

	log.Printf("DEBUG: SQL Update Result - RowsAffected: %d, Error: %v", result.RowsAffected, result.Error)

	// Verify the update immediately
	var updatedUser struct {
		ID     string
		RoleID string
	}
	verifyErr := exec.WithContext(ctx).
		Table(tableName).
		Where("id = ?", userID).
		Select("id, role_id").
		Scan(&updatedUser).Error

	if verifyErr != nil {
		log.Printf("DEBUG: Failed to verify role update: %v", verifyErr)
	} else {
		log.Printf("DEBUG: Verification after update - User ID: %s, New RoleID: %s", updatedUser.ID, updatedUser.RoleID)
		if updatedUser.RoleID != roleID {
			log.Printf("DEBUG: ⚠️ WARNING: Role verification failed! Expected %s but found %s", roleID, updatedUser.RoleID)
		} else {
			log.Printf("DEBUG: ✓ Role successfully updated to %s", roleID)
		}
	}

	// Additional verification: Check if the role is properly linked
	var roleName string
	if err := exec.WithContext(ctx).Table("roles").
		Where("id = ?", updatedUser.RoleID).
		Select("name").
		Scan(&roleName).Error; err != nil {
		log.Printf("DEBUG: Error verifying role name: %v", err)
	} else {
		log.Printf("DEBUG: Verified role name: %s", roleName)
	}

	return nil
}

// ExecuteRawUpdateQuery executes a raw SQL update query for role changes
func (r *userRepository) ExecuteRawUpdateQuery(ctx context.Context, userID string, roleID string) error {
	// Get table name from GORM configuration
	tableName := r.db.NamingStrategy.TableName("users")
	log.Printf("ExecuteRawUpdateQuery using table: %s", tableName)

	// Format raw SQL query
	updateStmt := fmt.Sprintf(`
		UPDATE %s 
		SET role_id = '%s', updated_at = '%s' 
		WHERE id = '%s'`,
		tableName, roleID, time.Now().Format("2006-01-02 15:04:05"), userID)

	log.Printf("Executing raw SQL: %s", updateStmt)

	// Execute as raw SQL to bypass any GORM caching
	result := r.db.WithContext(ctx).Exec(updateStmt)
	if result.Error != nil {
		log.Printf("Raw SQL error: %v", result.Error)
		return result.Error
	}

	log.Printf("Raw SQL result: affected rows = %d", result.RowsAffected)

	// Verify the update with a separate query
	var currentRole string
	verifyStmt := fmt.Sprintf("SELECT role_id FROM %s WHERE id = '%s'", tableName, userID)
	err := r.db.WithContext(ctx).Raw(verifyStmt).Scan(&currentRole).Error

	if err != nil {
		log.Printf("Error verifying role update: %v", err)
	} else {
		log.Printf("Verification: User %s now has role_id = %s", userID, currentRole)
		if currentRole != roleID {
			log.Printf("⚠️ WARNING: Role verification failed! Expected %s but found %s", roleID, currentRole)
		} else {
			log.Printf("✓ Role successfully updated to %s", roleID)
		}
	}

	return nil
}

// GetUserRoleInTenant gets the role of a user in a specific tenant
func (r *userRepository) GetUserRoleInTenant(ctx context.Context, userID uuid.UUID, tenantID uuid.UUID) (entity.Role, error) {

	// 1. Check if user is part of the tenant
	var userTenant entity.UserTenant
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		First(&userTenant).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Role{}, fmt.Errorf("user is not part of this tenant")
		}
		return entity.Role{}, err
	}

	// 2. Get the user's role
	var user entity.User
	err = r.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Role{}, fmt.Errorf("user not found")
		}
		return entity.Role{}, err
	}

	return user.Role, nil
}

// FindUsersByRoleID finds all users with a specific role ID
func (r *userRepository) FindUsersByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
