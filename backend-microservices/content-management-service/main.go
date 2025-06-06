package main

import (
	"content-management-service/data_layer/migration"
	"content-management-service/data_layer/repository"
	"content-management-service/domain_layer/service"
	"content-management-service/helpers/config"
	"content-management-service/presentation_layer/controller"
	"content-management-service/presentation_layer/route"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load environment variables
	config.LoadEnv()

	// Initialize database
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := migration.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	authorRepo := repository.NewAuthorRepository(db)
	readerRepo := repository.NewReaderRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	audiobookRepo := repository.NewAudiobookRepository(db)
	trackRepo := repository.NewTrackRepository(db)
	userRepo := repository.NewUserRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	// Initialize user management service for API validation
	userManagementBaseURL := config.GetUserManagementBaseURL()
	userManagementService := service.NewUserManagementService(userManagementBaseURL)
	log.Printf("User Management Service configured at: %s", userManagementBaseURL)

	// Initialize services
	authorService := service.NewAuthorService(authorRepo)
	readerService := service.NewReaderService(readerRepo)
	genreService := service.NewGenreService(genreRepo)
	audiobookService := service.NewAudiobookService(audiobookRepo, authorRepo, readerRepo, genreRepo)
	trackService := service.NewTrackService(trackRepo)
	userService := service.NewUserService(userRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)

	// Initialize controllers
	authorController := controller.NewAuthorController(authorService)
	readerController := controller.NewReaderController(readerService)
	genreController := controller.NewGenreController(genreService)
	audiobookController := controller.NewAudiobookController(audiobookService)
	trackController := controller.NewTrackController(trackService)
	userController := controller.NewUserController(userService)
	analyticsController := controller.NewAnalyticsController(analyticsService)

	// Setup Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes with user management service for middleware
	route.SetupRoutes(router, authorController, readerController, genreController, audiobookController, trackController, userController, analyticsController, userManagementService)

	// Get port from environment or use default
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3160" // default port
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
