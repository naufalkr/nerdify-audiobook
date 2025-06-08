package seed

import (
	"log"

	"gorm.io/gorm"
)

// RunAllSeeders runs all database seeders in the correct order
func RunAllSeeders(db *gorm.DB) error {
	log.Println("Starting database seeding process...")

	// Define the seeding order (important due to foreign key dependencies)
	seeders := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"Authors", AuthorSeeder},
		{"Genres", GenreSeeder},
		{"Readers", ReaderSeeder},
		{"Users", UserSeeder},
		{"Audiobooks", AudiobookSeeder},
		{"Audiobook-Genre Relationships", AudiobookGenreSeeder},
		{"Tracks", TrackSeeder},
		{"Analytics", AnalyticsSeeder},
	}

	// Run each seeder
	for _, seeder := range seeders {
		log.Printf("Running %s seeder...", seeder.name)

		if err := seeder.fn(db); err != nil {
			log.Printf("Error running %s seeder: %v", seeder.name, err)
			return err
		}

		log.Printf("%s seeder completed successfully", seeder.name)
	}

	log.Println("All database seeders completed successfully!")
	return nil
}

// RunSpecificSeeder runs a specific seeder by name
func RunSpecificSeeder(db *gorm.DB, seederName string) error {
	seeders := map[string]func(*gorm.DB) error{
		"authors":          AuthorSeeder,
		"genres":           GenreSeeder,
		"readers":          ReaderSeeder,
		"users":            UserSeeder,
		"audiobooks":       AudiobookSeeder,
		"audiobook-genres": AudiobookGenreSeeder,
		"tracks":           TrackSeeder,
		"analytics":        AnalyticsSeeder,
	}

	if seeder, exists := seeders[seederName]; exists {
		log.Printf("Running %s seeder...", seederName)
		return seeder(db)
	}

	log.Printf("Seeder '%s' not found", seederName)
	return nil
}

// ClearAllData clears all seeded data (useful for testing)
func ClearAllData(db *gorm.DB) error {
	log.Println("Clearing all seeded data...")

	// Clear in reverse order to avoid foreign key constraints
	tables := []string{
		"analytics",
		"tracks",
		"audiobook_genres",
		"audiobooks",
		"users",
		"readers",
		"genres",
		"authors",
	}

	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			log.Printf("Error clearing table %s: %v", table, err)
			return err
		}
		log.Printf("Cleared table: %s", table)
	}

	log.Println("All data cleared successfully!")
	return nil
}

// GetSeedingStats returns statistics about seeded data
func GetSeedingStats(db *gorm.DB) map[string]int64 {
	stats := make(map[string]int64)

	tables := map[string]interface{}{
		"authors":          &struct{}{},
		"genres":           &struct{}{},
		"readers":          &struct{}{},
		"users":            &struct{}{},
		"audiobooks":       &struct{}{},
		"tracks":           &struct{}{},
		"analytics":        &struct{}{},
		"audiobook_genres": &struct{}{},
	}

	for tableName := range tables {
		var count int64
		db.Table(tableName).Count(&count)
		stats[tableName] = count
	}

	return stats
}
