package seed

import (
	"catalog-service/data_layer/entity"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// UserSeeder seeds the users table
func UserSeeder(db *gorm.DB) error {
	log.Println("Seeding users...")

	users := []entity.User{
		{ID: "auth0|admin123456789", Role: "admin"},
		{ID: "auth0|superadmin987654321", Role: "superadmin"},
		{ID: "auth0|user001234567890", Role: "user"},
		{ID: "auth0|user002345678901", Role: "user"},
		{ID: "auth0|user003456789012", Role: "user"},
		{ID: "auth0|user004567890123", Role: "user"},
		{ID: "auth0|user005678901234", Role: "user"},
		{ID: "auth0|user006789012345", Role: "user"},
		{ID: "auth0|user007890123456", Role: "user"},
		{ID: "auth0|user008901234567", Role: "user"},
		{ID: "auth0|user009012345678", Role: "user"},
		{ID: "auth0|user010123456789", Role: "user"},
		{ID: "auth0|moderator111111111", Role: "moderator"},
		{ID: "auth0|moderator222222222", Role: "moderator"},
		{ID: "auth0|moderator333333333", Role: "moderator"},
		{ID: "google-oauth2|111111111111111111111", Role: "user"},
		{ID: "google-oauth2|222222222222222222222", Role: "user"},
		{ID: "google-oauth2|333333333333333333333", Role: "user"},
		{ID: "google-oauth2|444444444444444444444", Role: "user"},
		{ID: "google-oauth2|555555555555555555555", Role: "user"},
		{ID: "facebook|123456789012345", Role: "user"},
		{ID: "facebook|234567890123456", Role: "user"},
		{ID: "facebook|345678901234567", Role: "user"},
		{ID: "github|12345678", Role: "user"},
		{ID: "github|23456789", Role: "user"},
		{ID: "github|34567890", Role: "user"},
		{ID: "linkedin|abcd1234", Role: "user"},
		{ID: "linkedin|efgh5678", Role: "user"},
		{ID: "twitter|wxyz9876", Role: "user"},
		{ID: "twitter|stuv5432", Role: "user"},
	}

	// Generate additional users with UUIDs
	for i := 1; i <= 50; i++ {
		userID := fmt.Sprintf("uuid|user-%03d-uuid-%d", i, 1000000+i)
		role := "user"
		if i%15 == 0 {
			role = "moderator"
		} else if i%25 == 0 {
			role = "admin"
		}
		users = append(users, entity.User{
			ID:   userID,
			Role: role,
		})
	}

	// Check if users already exist to avoid duplicates
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already seeded, skipping...")
		return nil
	}

	// Create users in batches
	batchSize := 15
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		if err := db.Create(users[i:end]).Error; err != nil {
			log.Printf("Error seeding users batch %d-%d: %v", i, end-1, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d users", len(users))
	return nil
}
