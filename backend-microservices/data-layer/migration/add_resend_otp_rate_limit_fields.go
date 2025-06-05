package migration

import (
	"log"

	"gorm.io/gorm"
)

// AddResendOTPRateLimitFields adds fields to the users table to track OTP resend rate limiting
func AddResendOTPRateLimitFields(db *gorm.DB) error {
	log.Println("Adding OTP resend rate limit fields to users table...")

	// Check if resend_count column already exists
	hasResendCount := db.Migrator().HasColumn("users", "resend_count")
	if !hasResendCount {
		if err := db.Exec("ALTER TABLE users ADD COLUMN resend_count INT DEFAULT 0").Error; err != nil {
			log.Printf("Error adding resend_count column: %v", err)
			return err
		}
		log.Println("Added resend_count column to users table")
	}

	// Check if last_resend_at column already exists
	hasLastResendAt := db.Migrator().HasColumn("users", "last_resend_at")
	if !hasLastResendAt {
		if err := db.Exec("ALTER TABLE users ADD COLUMN last_resend_at TIMESTAMP").Error; err != nil {
			log.Printf("Error adding last_resend_at column: %v", err)
			return err
		}
		log.Println("Added last_resend_at column to users table")
	}

	// Check if cooldown_started_at column already exists
	hasCooldownStartedAt := db.Migrator().HasColumn("users", "cooldown_started_at")
	if !hasCooldownStartedAt {
		if err := db.Exec("ALTER TABLE users ADD COLUMN cooldown_started_at TIMESTAMP").Error; err != nil {
			log.Printf("Error adding cooldown_started_at column: %v", err)
			return err
		}
		log.Println("Added cooldown_started_at column to users table")
	}

	log.Println("Successfully added OTP resend rate limit fields to users table")
	return nil
}
