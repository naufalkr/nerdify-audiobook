package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
}

// GetDatabaseConfig returns database configuration from environment variables
func GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_DATABASE", "content_nerdify"),
	}
}

// InitDatabase initializes and returns database connection
func InitDatabase() (*gorm.DB, error) {
	config := GetDatabaseConfig()

	// First, try to create the database if it doesn't exist
	if err := createDatabaseIfNotExists(config); err != nil {
		log.Printf("Warning: Failed to create database: %v", err)
	}

	// Connect to the target database
	dsn := "host=" + config.Host + " user=" + config.Username + " password=" + config.Password + " dbname=" + config.Database + " port=" + config.Port + " sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to database: %s", config.Database)
	return db, nil
}

// createDatabaseIfNotExists creates the database if it doesn't exist
func createDatabaseIfNotExists(config DatabaseConfig) error {
	// Connect to postgres database (default database) first
	adminDSN := "host=" + config.Host + " user=" + config.Username + " password=" + config.Password + " dbname=postgres port=" + config.Port + " sslmode=disable TimeZone=Asia/Jakarta"

	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	// Get underlying sql.DB to execute raw SQL
	sqlDB, err := adminDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	err = sqlDB.QueryRow(query, config.Database).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		// Create database
		createDBQuery := "CREATE DATABASE " + config.Database
		_, err = sqlDB.Exec(createDBQuery)
		if err != nil {
			return err
		}
		log.Printf("Database '%s' created successfully", config.Database)
	} else {
		log.Printf("Database '%s' already exists", config.Database)
	}

	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
