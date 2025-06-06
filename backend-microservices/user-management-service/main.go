// filepath: /home/alvn/Documents/playground/kp/clean_architecture/user_management/main.go
package main

import (
	"flag"
	"log"
	"microservice/user/data-layer/config"
	"microservice/user/data-layer/migration"
	"microservice/user/data-layer/repository"
	"microservice/user/domain-layer/middleware"
	"microservice/user/domain-layer/service"
	"microservice/user/helpers/cmd"
	"microservice/user/helpers/utils"
	"microservice/user/presentation-layer/controller"
	"microservice/user/presentation-layer/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// NoOpEmailSender is a no-operation email sender that doesn't actually send emails
type NoOpEmailSender struct{}

func (n *NoOpEmailSender) Send(to, subject, body string) error {
	log.Printf("Email would be sent to: %s, Subject: %s", to, subject)
	// Log the email content so you can see the verification token
	log.Printf("Email content: %s", body)
	return nil
}

func main() {
	// Check for command line flags
	migrate := flag.Bool("migrate", false, "Run database migrations")
	partialIndex := flag.Bool("partial-index", false, "Run partial unique index migration only")
	seed := flag.Bool("seed", false, "Run database seeders")
	seedUsers := flag.Bool("seed-users", false, "Run only the verified users seeder")
	reset := flag.Bool("reset", false, "Reset database (drop all tables and re-migrate)")
	resetAll := flag.Bool("reset-all", false, "Completely reset database (drops everything including columns)")

	// Parse flags before doing anything else
	flag.Parse()

	// If any database-related flag is specified, run the migration command from cmd
	if *migrate || *partialIndex || *seed || *seedUsers || *reset || *resetAll {
		// Create FlagOptions struct with the flag values
		opts := cmd.FlagOptions{
			Migrate:      *migrate,
			PartialIndex: *partialIndex,
			Seed:         *seed,
			SeedUsers:    *seedUsers,
			Reset:        *reset,
			ResetAll:     *resetAll,
		}
		cmd.RunMigrationCommand(opts)
		return
	}

	// Continue with normal application startup if no migration flag is provided

	// Set up database connection using config package
	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Run database migrations
	log.Println("Starting database migration...")
	err := migration.RunMigration(db)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database migration completed successfully")

	// Apply partial unique index migration for soft-deleted records
	log.Println("Applying partial unique indexes...")
	err = migration.PartialUniqueIndexMigration(db)
	if err != nil {
		log.Fatalf("Partial unique index migration failed: %v", err)
	}
	log.Println("Partial unique indexes created successfully")

	// Seed default roles
	log.Println("Seeding default roles...")
	err = migration.SeedDefaultRoles(db)
	if err != nil {
		log.Fatalf("Failed to seed default roles: %v", err)
	}
	log.Println("Default roles seeded successfully")

	// Initialize repository
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	userTenantRepo := repository.NewUserTenantRepository(db)

	// Initialize other dependencies
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	emailSecret := os.Getenv("EMAIL_TOKEN_SECRET")
	passwordResetSecret := os.Getenv("PASSWORD_RESET_SECRET")

	// Set default values if not provided
	if accessSecret == "" {
		accessSecret = "default_access_secret_key_at_least_32_chars"
	}
	if refreshSecret == "" {
		refreshSecret = "default_refresh_secret_key_at_least_32_chars"
	}
	if emailSecret == "" {
		emailSecret = "default_email_secret_key_at_least_32_chars"
	}
	if passwordResetSecret == "" {
		passwordResetSecret = "default_pwd_reset_secret_key_at_least_32_chars"
	}

	tokenMaker, err := utils.NewJWTMaker(accessSecret, refreshSecret, emailSecret, passwordResetSecret)
	if err != nil {
		log.Fatalf("Cannot create token maker: %v", err)
	}

	// Initialize email configuration
	emailConfig, err := config.NewEmailConfig()
	if err != nil {
		log.Printf("Failed to initialize email config: %v", err)
		log.Println("Using NoOp email sender (emails will be logged but not sent)")
		emailSender := &NoOpEmailSender{}
		initializeServices(db, roleRepo, userRepo, tenantRepo, userTenantRepo, tokenMaker, emailSender)
	} else {
		// Create real email sender
		emailSender := utils.NewSMTPEmailSender(
			emailConfig.Host,
			emailConfig.Port,
			emailConfig.AuthEmail,
			emailConfig.AuthPassword,
			emailConfig.SenderName,
		)
		initializeServices(db, roleRepo, userRepo, tenantRepo, userTenantRepo, tokenMaker, emailSender)
	}
}

func initializeServices(
	db *gorm.DB,
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
	tenantRepo repository.TenantRepository,
	userTenantRepo repository.UserTenantRepository,
	tokenMaker utils.TokenMaker,
	emailSender utils.EmailSender,
) {
	// Initialize audit log repository and service
	auditLogRepo := repository.NewAuditLogRepository(db)
	auditSvc := service.NewAuditService(auditLogRepo)

	// Initialize Cloudinary service
	cloudinaryService, err := config.NewCloudinaryService()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary service: %v", err)
	}

	// Initialize service
	uService := service.NewUserService(userRepo, roleRepo, tokenMaker, emailSender)
	tService := service.NewTenantService(tenantRepo, userRepo, roleRepo, emailSender)
	rService := service.NewRoleService(roleRepo)
	userTenantContextService := service.NewUserTenantContextService(userTenantRepo, userRepo, tenantRepo, db)

	// Initialize controller
	userController := controller.NewUserController(uService, tokenMaker, cloudinaryService)
	tenantController := controller.NewTenantController(tService, tokenMaker, cloudinaryService)
	tenantAPIController := controller.NewTenantAPIController(tService, userTenantContextService, tokenMaker)
	roleController := controller.NewRoleController(rService, uService)
	auditController := controller.NewAuditController(auditSvc)

	// Setup router Gin
	router := gin.Default()

	// Root endpoint to check API status
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User Management API is running",
		})
	})

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add audit middleware
	router.Use(middleware.AuditMiddlewareFunc(auditSvc))

	// Register routes
	routes.SetupRoutes(router, userController, roleController, tenantController, tenantAPIController, auditController, tokenMaker)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
