package main

import (
	"content-management-service/data_layer/migration"
	"content-management-service/helpers/config"
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Command line flags
	var (
		migrateOnly  = flag.Bool("migrate", false, "Run only database migrations")
		seedOnly     = flag.Bool("seed", false, "Run only database seeding")
		seedSpecific = flag.String("seed-specific", "", "Run specific seeder (authors, genres, readers, users, audiobooks, tracks, analytics)")
		clearData    = flag.Bool("clear", false, "Clear all seeded data")
		showStats    = flag.Bool("stats", false, "Show seeding statistics")
		help         = flag.Bool("help", false, "Show help information")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Load database configuration
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch {
	case *clearData:
		fmt.Println("Clearing all seeded data...")
		if err := migration.ClearSeededData(db); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		fmt.Println("✅ Data cleared successfully!")

	case *showStats:
		fmt.Println("Database Seeding Statistics:")
		fmt.Println("============================")
		stats := migration.GetSeedingStatistics(db)
		for table, count := range stats {
			fmt.Printf("%-20s: %d records\n", table, count)
		}

	case *migrateOnly:
		fmt.Println("Running database migrations...")
		if err := migration.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✅ Migration completed successfully!")

	case *seedOnly:
		fmt.Println("Running database seeders...")
		if err := migration.SeedDatabase(db); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		fmt.Println("✅ Seeding completed successfully!")

	case *seedSpecific != "":
		fmt.Printf("Running %s seeder...\n", *seedSpecific)
		if err := migration.SeedSpecific(db, *seedSpecific); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		fmt.Printf("✅ %s seeder completed successfully!\n", *seedSpecific)

	default:
		// Run both migration and seeding
		fmt.Println("Running database migration and seeding...")
		if err := migration.AutoMigrateAndSeed(db); err != nil {
			log.Fatalf("Migration and seeding failed: %v", err)
		}
		fmt.Println("✅ Migration and seeding completed successfully!")

		// Show final statistics
		fmt.Println("\nFinal Database Statistics:")
		fmt.Println("==========================")
		stats := migration.GetSeedingStatistics(db)
		for table, count := range stats {
			fmt.Printf("%-20s: %d records\n", table, count)
		}
	}
}

func showHelp() {
	fmt.Println("Database Migration and Seeding Tool")
	fmt.Println("===================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seeder/main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -migrate           Run only database migrations")
	fmt.Println("  -seed              Run only database seeding")
	fmt.Println("  -seed-specific     Run specific seeder:")
	fmt.Println("                     authors, genres, readers, users,")
	fmt.Println("                     audiobooks, tracks, analytics")
	fmt.Println("  -clear             Clear all seeded data")
	fmt.Println("  -stats             Show seeding statistics")
	fmt.Println("  -help              Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seeder/main.go                    # Run migration and seeding")
	fmt.Println("  go run cmd/seeder/main.go -migrate           # Run only migrations")
	fmt.Println("  go run cmd/seeder/main.go -seed              # Run only seeding")
	fmt.Println("  go run cmd/seeder/main.go -seed-specific authors  # Seed only authors")
	fmt.Println("  go run cmd/seeder/main.go -clear             # Clear all data")
	fmt.Println("  go run cmd/seeder/main.go -stats             # Show statistics")
}
