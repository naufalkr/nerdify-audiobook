package migration

import (
	"log"

	"gorm.io/gorm"
)

// AddPendingEmailMigration adds the pending_email column to the users table
func AddPendingEmailMigration(db *gorm.DB) error {
	log.Println("Running migration: Add pending_email column to users table")

	// Check if the column already exists
	hasColumn := db.Migrator().HasColumn("users", "pending_email")
	if !hasColumn {
		err := db.Exec("ALTER TABLE users ADD COLUMN pending_email VARCHAR(100)").Error
		if err != nil {
			log.Printf("Error adding pending_email column: %v", err)
			return err
		}
		log.Println("Successfully added pending_email column to users table")
	} else {
		log.Println("Column pending_email already exists, skipping")
	}

	return nil
}
