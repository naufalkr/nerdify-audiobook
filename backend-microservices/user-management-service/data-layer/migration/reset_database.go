package migration

import (
	"log"
	model "microservice/user/data-layer/entity"

	"gorm.io/gorm"
)

// ResetDatabaseCompletely drops and recreates all tables with updated schema
func ResetDatabaseCompletely(db *gorm.DB) error {
	log.Println("Completely resetting database (dropping all tables and recreating with updated schema)...")

	// Drop all tables in reverse order of dependencies
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

	log.Println("All tables dropped successfully")

	// Re-run migrations with updated schema
	return RunMigration(db)
}
