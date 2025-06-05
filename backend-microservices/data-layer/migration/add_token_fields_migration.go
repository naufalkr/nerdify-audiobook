package migration

import (
	"log"
	model "microservice/user/data-layer/entity"

	"gorm.io/gorm"
)

// AddTokenFieldsToUsers adds token-related columns to the users table to support token-based authentication
func AddTokenFieldsToUsers(db *gorm.DB) error {
	log.Println("Adding token fields to users table...")

	// Check if auth_token column already exists
	hasAuthToken := db.Migrator().HasColumn(&model.User{}, "auth_token")
	if !hasAuthToken {
		if err := db.Exec("ALTER TABLE users ADD COLUMN auth_token VARCHAR(500)").Error; err != nil {
			log.Printf("Error adding auth_token column: %v", err)
			return err
		}
		log.Println("Added auth_token column to users table")
	}

	// Check if token_expiry column already exists
	hasTokenExpiry := db.Migrator().HasColumn(&model.User{}, "token_expiry")
	if !hasTokenExpiry {
		if err := db.Exec("ALTER TABLE users ADD COLUMN token_expiry TIMESTAMP").Error; err != nil {
			log.Printf("Error adding token_expiry column: %v", err)
			return err
		}
		log.Println("Added token_expiry column to users table")
	}

	// Check if refresh_token column already exists
	hasRefreshToken := db.Migrator().HasColumn(&model.User{}, "refresh_token")
	if !hasRefreshToken {
		if err := db.Exec("ALTER TABLE users ADD COLUMN refresh_token VARCHAR(500)").Error; err != nil {
			log.Printf("Error adding refresh_token column: %v", err)
			return err
		}
		log.Println("Added refresh_token column to users table")
	}

	// Check if token_created_at column already exists
	hasTokenCreatedAt := db.Migrator().HasColumn(&model.User{}, "token_created_at")
	if !hasTokenCreatedAt {
		if err := db.Exec("ALTER TABLE users ADD COLUMN token_created_at TIMESTAMP").Error; err != nil {
			log.Printf("Error adding token_created_at column: %v", err)
			return err
		}
		log.Println("Added token_created_at column to users table")
	}

	log.Println("Successfully added token fields to users table")
	return nil
}
