package migration

import (
	"log"

	"gorm.io/gorm"
)

// DropTenantIDFromUsers removes the tenant_id column from the users table
func DropTenantIDFromUsers(db *gorm.DB) error {
	log.Println("Dropping tenant_id column from users table...")

	// Check if the tenant_id column exists in the users table
	if db.Migrator().HasColumn(&User{}, "TenantID") {
		// Drop the tenant_id column
		if err := db.Migrator().DropColumn(&User{}, "tenant_id"); err != nil {
			log.Printf("Failed to drop tenant_id column: %v", err)
			return err
		}
		log.Println("Successfully dropped tenant_id column from users table")
	} else {
		log.Println("tenant_id column does not exist in users table, skipping")
	}

	return nil
}
