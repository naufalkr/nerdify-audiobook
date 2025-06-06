package migration

import (
	"log"
	model "microservice/user/data-layer/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Temporary User struct for legacy migration cleanup
type User struct {
	TenantID uuid.UUID `gorm:"type:char(36)"`
}

func (User) TableName() string {
	return "users"
}

// FixRoleIDColumnType ensures the role_id column in users table is of type uuid
func FixRoleIDColumnType(db *gorm.DB) error {
	log.Println("Fixing role_id column type in users table...")

	// Check if role_id column exists and is not of type uuid
	if db.Migrator().HasColumn(&model.User{}, "role_id") {
		// Get the current column type
		var columnType string
		err := db.Raw(`
			SELECT data_type 
			FROM information_schema.columns 
			WHERE table_name = 'users' AND column_name = 'role_id'
		`).Scan(&columnType).Error

		if err != nil {
			log.Printf("Error checking role_id column type: %v", err)
			return err
		}

		// If the column is not of type uuid, alter it
		if columnType != "uuid" {
			log.Printf("Converting role_id column from %s to uuid type", columnType)
			err = db.Exec(`
				ALTER TABLE users 
				ALTER COLUMN role_id TYPE uuid USING role_id::uuid
			`).Error

			if err != nil {
				log.Printf("Error converting role_id column type: %v", err)
				return err
			}
			log.Println("Successfully converted role_id column to uuid type")
		} else {
			log.Println("role_id column is already of type uuid")
		}
	} else {
		log.Println("role_id column does not exist in users table")
	}

	return nil
}

// RunMigration handles the database migration for all models
func RunMigration(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Run migrations for all models in correct order
	err := db.AutoMigrate(
		&model.Role{},       // Roles migrated first
		&model.Tenant{},     // Tenants
		&model.User{},       // Users
		&model.UserTenant{}, // User-tenant relation
		&model.AuditLog{},   // Audit log
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	// Fix role_id column type if needed
	if err := FixRoleIDColumnType(db); err != nil {
		log.Printf("Warning: Failed to fix role_id column type: %v", err)
	}

	// Try to drop tenant_id from users if still exists (legacy cleanup)
	if db.Migrator().HasColumn(&User{}, "tenant_id") {
		err := db.Migrator().DropColumn(&User{}, "tenant_id")
		if err != nil {
			log.Printf("Warning: Failed to drop legacy tenant_id column: %v", err)
		} else {
			log.Println("Dropped legacy tenant_id column from users")
		}
	}

	// Add token fields to users table for token-based authentication
	if err := AddTokenFieldsToUsers(db); err != nil {
		log.Printf("Warning: Failed to add token fields to users: %v", err)
	}

	// Add pending_email field to users table
	if err := AddPendingEmailMigration(db); err != nil {
		log.Printf("Warning: Failed to add pending_email field to users: %v", err)
	}

	// Add fields for OTP resend rate limiting
	if err := AddResendOTPRateLimitFields(db); err != nil {
		log.Printf("Warning: Failed to add resend OTP rate limit fields to users: %v", err)
	}

	log.Println("Migration completed successfully")
	return nil
}

// SeedDefaultRoles creates system roles if not present
func SeedDefaultRoles(db *gorm.DB) error {
	log.Println("Seeding default roles...")

	defaultRoles := []struct {
		Name        string
		Description string
	}{
		{"SUPERADMIN", "Super administrator with full access"},
		{"ADMIN", "Tenant administrator"},
		{"USER", "Regular user"},
	}

	for _, role := range defaultRoles {
		var existing model.Role
		result := db.Where("name = ?", role.Name).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			newRole := model.Role{
				ID:          uuid.New(),
				Name:        role.Name,
				Description: role.Description,
				IsSystem:    true,
			}

			if err := db.Create(&newRole).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
				return err
			}

			log.Printf("Created role: %s", role.Name)
		}
	}

	log.Println("Default roles seeded successfully")
	return nil
}

// ResetDatabase drops all tables and recreates them
func ResetDatabase(db *gorm.DB) error {
	log.Println("Resetting database (dropping all tables)...")

	err := db.Migrator().DropTable(
		&model.AuditLog{},
		&model.UserTenant{},
		&model.User{},
		&model.Tenant{},
		&model.Role{},
	)

	if err != nil {
		log.Printf("Failed to drop tables: %v", err)
		return err
	}

	log.Println("All tables dropped")
	return RunMigration(db)
}

// TruncateAllTables empties all rows from tables without dropping them
func TruncateAllTables(db *gorm.DB) error {
	log.Println("Truncating all tables...")

	tables := []string{"audit_logs", "user_tenants", "users", "tenants", "roles"}
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			log.Printf("Failed to truncate table %s: %v", table, err)
		}
	}

	log.Println("All tables truncated successfully")
	return nil
}
