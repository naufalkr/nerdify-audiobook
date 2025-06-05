package seeder

import (
	"log"

	"gorm.io/gorm"
)

// RunAllSeeders runs all seeders
func RunAllSeeders(db *gorm.DB) error {
	log.Println("Running all seeders...")

	// First ensure roles exist
	if err := SeedDefaultRoles(db); err != nil {
		log.Printf("Failed to seed default roles: %v", err)
		return err
	}

	// Then create superadmin user
	if err := SeedSuperAdminUser(db); err != nil {
		log.Printf("Failed to seed superadmin user: %v", err)
		return err
	}

	// Create verified regular user
	if err := SeedVerifiedUser(db); err != nil {
		log.Printf("Failed to seed verified regular user: %v", err)
		return err
	}

	log.Println("All seeders completed successfully")
	return nil
}
