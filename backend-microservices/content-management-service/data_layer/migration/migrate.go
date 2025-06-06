package migration

import (
	"content-management-service/data_layer/entity"
	"content-management-service/data_layer/migration/seed"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate runs database migration for all entities
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(
		&entity.Author{},
		&entity.Reader{},
		&entity.Genre{},
		&entity.User{},
		&entity.Audiobook{},
		&entity.Track{},
		&entity.Analytics{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

// AutoMigrateAndSeed runs database migration and seeds data
func AutoMigrateAndSeed(db *gorm.DB) error {
	// Run migrations first
	if err := AutoMigrate(db); err != nil {
		return err
	}

	// Run seeders
	if err := seed.RunAllSeeders(db); err != nil {
		log.Printf("Seeding failed: %v", err)
		return err
	}

	return nil
}

// SeedDatabase runs only the seeders without migration
func SeedDatabase(db *gorm.DB) error {
	return seed.RunAllSeeders(db)
}

// SeedSpecific runs a specific seeder
func SeedSpecific(db *gorm.DB, seederName string) error {
	return seed.RunSpecificSeeder(db, seederName)
}

// ClearSeededData clears all seeded data
func ClearSeededData(db *gorm.DB) error {
	return seed.ClearAllData(db)
}

// GetSeedingStatistics returns seeding statistics
func GetSeedingStatistics(db *gorm.DB) map[string]int64 {
	return seed.GetSeedingStats(db)
}
