package seeder

import (
	"log"
	model "microservice/user/data-layer/entity"
	"microservice/user/helpers/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserData represents data for a verified user to be created
type UserData struct {
	Username string
	Email    string
	FullName string
	Address  string
	Password string
}

// SeedVerifiedUser creates multiple regular users with verified status
func SeedVerifiedUser(db *gorm.DB) error {
	log.Println("Seeding multiple verified regular users...")

	// Find or create default tenant
	var defaultTenant model.Tenant
	if err := db.Where("name = ?", "Default Tenant").First(&defaultTenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Default tenant does not exist, creating new one...")
			// Create default tenant
			defaultTenant = model.Tenant{
				ID:                    uuid.New(),
				Name:                  "Default Tenant",
				Description:           "Default system tenant",
				LogoURL:               "",
				ContactEmail:          "admin@example.com",
				ContactPhone:          "",
				MaxUsers:              50,
				SubscriptionPlan:      "Enterprise",
				SubscriptionStartDate: time.Now(),
				SubscriptionEndDate:   time.Now().AddDate(1, 0, 0), // 1 year subscription
				IsActive:              true,
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			}
			if err := db.Create(&defaultTenant).Error; err != nil {
				log.Printf("Failed to create default tenant: %v", err)
				return err
			}
			log.Println("Created default tenant")
		} else {
			log.Printf("Error finding default tenant: %v", err)
			return err
		}
	} else {
		log.Println("Default tenant already exists, using existing one")
	}

	// Find USER role
	var userRole model.Role
	if err := db.Where("name = ?", "USER").First(&userRole).Error; err != nil {
		log.Printf("USER role not found: %v", err)
		return err
	}

	// Define users to create
	usersToCreate := []UserData{
		{
			Username: "user1",
			Email:    "user1@example.com",
			FullName: "Regular User One",
			Address:  "123 First Street, City",
			Password: "password123",
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
			FullName: "Regular User Two",
			Address:  "456 Second Avenue, City",
			Password: "password123",
		},
		{
			Username: "user3",
			Email:    "user3@example.com",
			FullName: "Regular User Three",
			Address:  "789 Third Road, City",
			Password: "password123",
		},
		{
			Username: "user4",
			Email:    "user4@example.com",
			FullName: "Regular User Four",
			Address:  "101 Fourth Boulevard, City",
			Password: "password123",
		},
		{
			Username: "user5",
			Email:    "user5@example.com",
			FullName: "Regular User Five",
			Address:  "202 Fifth Lane, City",
			Password: "password123",
		},
	}

	// Create each user in the database
	for _, userData := range usersToCreate {
		// Check if this user already exists
		var existingUser model.User
		result := db.Where("user_name = ?", userData.Username).First(&existingUser)
		if result.Error == nil {
			log.Printf("User %s already exists, skipping", userData.Username)
			continue
		} else if result.Error != gorm.ErrRecordNotFound {
			log.Printf("Error checking for existing user %s: %v", userData.Username, result.Error)
			return result.Error
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(userData.Password)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", userData.Username, err)
			return err
		}

		// Create the user
		user := model.User{
			ID:              uuid.New(),
			RoleID:          &userRole.ID,
			UserName:        userData.Username,
			Email:           userData.Email,
			Password:        hashedPassword,
			FullName:        userData.FullName,
			Alamat:          userData.Address,
			ProfileImageURL: "",
			IsVerified:      true, // User is already verified
			Status:          "active",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		log.Printf("Creating verified user: %s", userData.Username)
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", userData.Username, err)
			return err
		}

		// Create a relation between user and tenant
		userTenant := model.UserTenant{
			ID:        uuid.New(),
			UserID:    user.ID,
			TenantID:  defaultTenant.ID,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		log.Printf("Creating user-tenant relationship for user %s", userData.Username)
		if err := db.Create(&userTenant).Error; err != nil {
			log.Printf("Failed to create user-tenant relationship for user %s: %v", userData.Username, err)
			return err
		}

		log.Printf("Successfully created verified user: %s", userData.Username)
	}

	log.Printf("Successfully created %d verified users", len(usersToCreate))
	return nil
}
