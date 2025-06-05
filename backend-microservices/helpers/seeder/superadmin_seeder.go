package seeder

import (
	"log"
	model "microservice/user/data-layer/entity"
	"microservice/user/helpers/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedSuperAdminUser creates a verified SUPERADMIN user if it doesn't exist
func SeedSuperAdminUser(db *gorm.DB) error {
	log.Println("Seeding verified SUPERADMIN user...")

	// Check if superadmin already exists
	var existingUser model.User
	result := db.Where("user_name = ?", "superadmin").First(&existingUser)
	if result.Error == nil {
		log.Println("SUPERADMIN user already exists, skipping")
		return nil
	} else if result.Error != gorm.ErrRecordNotFound {
		// Only log as error if it's not just "record not found"
		log.Printf("Error checking for existing superadmin: %v", result.Error)
		return result.Error
	}
	log.Println("SUPERADMIN user does not exist, creating new one...")

	// Find SUPERADMIN role
	var superadminRole model.Role
	if err := db.Where("name = ?", "SUPERADMIN").First(&superadminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("SUPERADMIN role not found, creating default roles first...")
			// Seed default roles
			if err := SeedDefaultRoles(db); err != nil {
				log.Printf("Failed to seed default roles: %v", err)
				return err
			}

			// Try again to find the role
			if err := db.Where("name = ?", "SUPERADMIN").First(&superadminRole).Error; err != nil {
				log.Printf("SUPERADMIN role still not found even after seeding: %v", err)
				return err
			}
			log.Println("Created SUPERADMIN role")
		} else {
			log.Printf("Error finding SUPERADMIN role: %v", err)
			return err
		}
	} else {
		log.Println("Found existing SUPERADMIN role")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword("superadmin123")
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}

	// Create the superadmin user
	superadminUser := model.User{
		ID:              uuid.New(),
		RoleID:          &superadminRole.ID,
		UserName:        "superadmin",
		Email:           "superadmin@example.com",
		Password:        hashedPassword,
		FullName:        "Super Administrator",
		Alamat:          "",
		ProfileImageURL: "",
		IsVerified:      true, // User is already verified
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	log.Println("Creating SUPERADMIN user...")
	if err := db.Create(&superadminUser).Error; err != nil {
		log.Printf("Failed to create SUPERADMIN user: %v", err)
		return err
	}
	log.Println("SUPERADMIN user created successfully")

	log.Println("Successfully created verified SUPERADMIN user")
	return nil
}

// SeedDefaultRoles creates the default system roles if they don't exist
func SeedDefaultRoles(db *gorm.DB) error {
	log.Println("Seeding default roles...")

	// Define default roles with fixed UUIDs for consistency
	defaultRoles := []struct {
		ID          uuid.UUID
		Name        string
		Description string
	}{
		{
			ID:          uuid.MustParse("45d15e0b-f09f-43d3-8354-971ec594ad6e"), // SUPERADMIN
			Name:        "SUPERADMIN",
			Description: "Super administrator with full access",
		},
		{
			ID:          uuid.MustParse("96070191-cf4f-497f-a670-dddf5eddce7c"), // ADMIN
			Name:        "ADMIN",
			Description: "Tenant administrator",
		},
		{
			ID:          uuid.MustParse("be75dfe6-1a52-4c17-a65a-3dc484be86fe"), // USER
			Name:        "USER",
			Description: "Regular user",
		},
	}

	// Create roles if they don't exist
	for _, role := range defaultRoles {
		var existingRole model.Role
		result := db.Where("id = ? OR name = ?", role.ID, role.Name).First(&existingRole)

		if result.Error == gorm.ErrRecordNotFound {
			newRole := model.Role{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				IsSystem:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			if err := db.Create(&newRole).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
				return err
			}

			log.Printf("Created role: %s with ID: %s", role.Name, role.ID)
		} else if result.Error != nil {
			log.Printf("Error checking for role %s: %v", role.Name, result.Error)
			return result.Error
		} else {
			// Update existing role to ensure consistency
			existingRole.Description = role.Description
			existingRole.IsSystem = true
			existingRole.UpdatedAt = time.Now()

			if err := db.Save(&existingRole).Error; err != nil {
				log.Printf("Failed to update role %s: %v", role.Name, err)
				return err
			}
			log.Printf("Updated existing role: %s with ID: %s", role.Name, existingRole.ID)
		}
	}

	log.Println("Default roles seeded successfully")
	return nil
}
