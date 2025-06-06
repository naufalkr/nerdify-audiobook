package migration

import (
	"log"

	"gorm.io/gorm"
)

// PartialUniqueIndexMigration modifies existing unique indexes to be partial indexes
// that only consider non-deleted records
func PartialUniqueIndexMigration(db *gorm.DB) error {
	log.Println("Running PartialUniqueIndexMigration - Start")

	// Drop existing unique indexes
	if err := db.Exec("DROP INDEX IF EXISTS idx_users_email").Error; err != nil {
		log.Printf("Error dropping idx_users_email: %v", err)
		return err
	}

	if err := db.Exec("DROP INDEX IF EXISTS idx_users_user_name").Error; err != nil {
		log.Printf("Error dropping idx_users_user_name: %v", err)
		return err
	}

	// Create partial unique indexes that exclude soft-deleted records
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active ON users (email) WHERE deleted_at IS NULL").Error; err != nil {
		log.Printf("Error creating idx_users_email_active: %v", err)
		return err
	}

	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_user_name_active ON users (user_name) WHERE deleted_at IS NULL").Error; err != nil {
		log.Printf("Error creating idx_users_user_name_active: %v", err)
		return err
	}

	log.Println("Running PartialUniqueIndexMigration - Complete")
	return nil
}
