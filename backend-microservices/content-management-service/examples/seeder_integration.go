package main

import (
	"content-management-service/data_layer/migration"
	"content-management-service/helpers/config"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Example: Integration of seeders into main application
func exampleIntegration() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Check for seeding flag
	shouldSeed := flag.Bool("seed", false, "Run database seeding after migration")
	flag.Parse()

	// Initialize database
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Check if this is a fresh database or development environment
	isDevelopment := os.Getenv("GO_ENV") == "development"
	isProduction := os.Getenv("GO_ENV") == "production"

	if *shouldSeed || isDevelopment {
		// Run migration and seeding for development
		log.Println("Running migration and seeding for development environment...")
		if err := migration.AutoMigrateAndSeed(db); err != nil {
			log.Fatalf("Failed to migrate and seed database: %v", err)
		}

		// Show seeding statistics
		stats := migration.GetSeedingStatistics(db)
		log.Println("Database seeding completed. Statistics:")
		for table, count := range stats {
			log.Printf("  %s: %d records", table, count)
		}
	} else if isProduction {
		// Run only migration for production
		log.Println("Running migration for production environment...")
		if err := migration.AutoMigrate(db); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
	} else {
		// Default: run migration only
		log.Println("Running database migration...")
		if err := migration.AutoMigrate(db); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
	}

	// Continue with rest of application setup...
	log.Println("Application setup completed!")
}

// Example: Conditional seeding based on environment
func conditionalSeeding() {
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Always run migrations
	if err := migration.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Check if database is empty (no authors exist)
	stats := migration.GetSeedingStatistics(db)
	if stats["authors"] == 0 {
		log.Println("Database appears to be empty. Running seeders...")
		if err := migration.SeedDatabase(db); err != nil {
			log.Printf("Warning: Seeding failed: %v", err)
		}
	} else {
		log.Println("Database already contains data. Skipping seeding.")
	}
}

// Example: Selective seeding
func selectiveSeeding() {
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := migration.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Seed basic data if not exists
	stats := migration.GetSeedingStatistics(db)

	if stats["authors"] == 0 {
		log.Println("Seeding authors...")
		migration.SeedSpecific(db, "authors")
	}

	if stats["genres"] == 0 {
		log.Println("Seeding genres...")
		migration.SeedSpecific(db, "genres")
	}

	if stats["readers"] == 0 {
		log.Println("Seeding readers...")
		migration.SeedSpecific(db, "readers")
	}

	// Only seed test data in development
	if os.Getenv("GO_ENV") == "development" {
		if stats["audiobooks"] == 0 {
			log.Println("Seeding development data...")
			migration.SeedSpecific(db, "users")
			migration.SeedSpecific(db, "audiobooks")
			migration.SeedSpecific(db, "audiobook-genres")
			migration.SeedSpecific(db, "tracks")
			migration.SeedSpecific(db, "analytics")
		}
	}
}

func main() {
	log.Println("This is an example integration file.")
	log.Println("Choose one of the example functions to see different seeding approaches:")
	log.Println("1. exampleIntegration() - Full integration with flags")
	log.Println("2. conditionalSeeding() - Seed only if database is empty")
	log.Println("3. selectiveSeeding() - Selective seeding based on environment")

	// Uncomment one of these to try:
	// exampleIntegration()
	// conditionalSeeding()
	// selectiveSeeding()
}
