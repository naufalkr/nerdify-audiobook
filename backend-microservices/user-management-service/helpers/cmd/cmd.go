package cmd

import (
	"log"
	"microservice/user/data-layer/config"
	"microservice/user/data-layer/migration"
	"microservice/user/helpers/seeder"
)

// FlagOptions contains the command-line flag values
type FlagOptions struct {
	Migrate      bool
	PartialIndex bool
	Seed         bool
	SeedUsers    bool
	Reset        bool
	ResetAll     bool
}

// RunMigrationCommand runs appropriate migration commands based on provided flags
func RunMigrationCommand(opts FlagOptions) {
	// Setup database connection using the same method as migrations
	log.Println("Connecting to database...")
	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	// Handle reset-all flag
	if opts.ResetAll {
		log.Println("Completely resetting database...")
		if err := migration.ResetDatabaseCompletely(db); err != nil {
			log.Fatalf("Failed to reset database completely: %v", err)
		}
		// After reset-all, consider migration done
		opts.Migrate = false
	} else if opts.Reset {
		// Handle regular reset
		log.Println("Resetting database...")
		if err := migration.ResetDatabase(db); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
		// After reset, consider migration done
		opts.Migrate = false
	}

	// Run migrations if requested
	if opts.Migrate {
		log.Println("Running migrations...")
		if err := migration.RunMigration(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Run seeders if requested
	if opts.Seed {
		log.Println("Running all seeders...")
		if err := seeder.RunAllSeeders(db); err != nil {
			log.Fatalf("Failed to run seeders: %v", err)
		}
	}

	// Only seed verified users if requested
	if opts.SeedUsers {
		log.Println("Creating verified users...")
		if err := seeder.SeedVerifiedUser(db); err != nil {
			log.Fatalf("Failed to create verified users: %v", err)
		}
	}

	// Run partial unique index migration only if requested
	if opts.PartialIndex {
		log.Println("Running partial unique index migration...")
		if err := migration.PartialUniqueIndexMigration(db); err != nil {
			log.Fatalf("Failed to run partial unique index migration: %v", err)
		}
		log.Println("Partial unique index migration completed successfully")
	}

	// If no flags specified, print usage
	if !opts.Migrate && !opts.Seed && !opts.SeedUsers && !opts.Reset && !opts.ResetAll && !opts.PartialIndex {
		log.Println("No action specified. Use one of these flags:")
		log.Println("  -migrate: Run database migrations")
		log.Println("  -partial-index: Run only the partial unique index migration")
		log.Println("  -seed: Run database seeders")
		log.Println("  -seed-users: Run only the verified users seeder")
		log.Println("  -reset: Reset database (drop all tables and re-migrate)")
		log.Println("  -reset-all: Completely reset database (drops everything including columns)")
	} else {
		log.Println("Database operations completed successfully")
	}
}
