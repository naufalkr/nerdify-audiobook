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

	// Check if superadmin already exists with the new credentials
	var existingUser model.User
	result := db.Where("user_name = ? OR email = ?", "SUPERADMIN", "superadmin@gmail.com").First(&existingUser)
	if result.Error == nil {
		log.Printf("SUPERADMIN user already exists with ID: %s", existingUser.ID)
		log.Printf("Existing user details:")
		log.Printf("- Username: %s", existingUser.UserName)
		log.Printf("- Email: %s", existingUser.Email)
		log.Printf("- Full Name: %s", existingUser.FullName)
		log.Printf("- Is Verified: %t", existingUser.IsVerified)
		return nil
	} else if result.Error != gorm.ErrRecordNotFound {
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
		log.Printf("Found existing SUPERADMIN role with ID: %s", superadminRole.ID)
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword("superadmin123")
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}

	// Create the superadmin user with custom data
	superadminUser := model.User{
		ID:              uuid.New(),
		RoleID:          &superadminRole.ID,
		UserName:        "SUPERADMIN",
		Email:           "superadmin@gmail.com",
		Password:        hashedPassword,
		FullName:        "SUPERADMIN",
		Alamat:          "SUPERADMIN",
		Latitude:        -6.2088,
		Longitude:       106.8456,
		ProfileImageURL: "",
		IsVerified:      true, // User is already verified
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	log.Println("Creating SUPERADMIN user with custom credentials...")
	if err := db.Create(&superadminUser).Error; err != nil {
		log.Printf("Failed to create SUPERADMIN user: %v", err)
		return err
	}

	log.Println("✅ SUPERADMIN user created successfully!")
	log.Printf("📋 SUPERADMIN Details:")
	log.Printf("- ID: %s", superadminUser.ID)
	log.Printf("- Username: %s", superadminUser.UserName)
	log.Printf("- Email: %s", superadminUser.Email)
	log.Printf("- Full Name: %s", superadminUser.FullName)
	log.Printf("- Address: %s", superadminUser.Alamat)
	log.Printf("- Location: %.4f, %.4f", superadminUser.Latitude, superadminUser.Longitude)
	log.Printf("- Is Verified: %t", superadminUser.IsVerified)
	log.Printf("- Role ID: %s", *superadminUser.RoleID)

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
