package config

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabaseConnection creates and configures the database connection
func SetupDatabaseConnection() *gorm.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database configuration from environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Debug output to verify values are loaded
	fmt.Println("==== ENV DEBUG ====")
	fmt.Println("DB_HOST =", dbHost)
	fmt.Println("DB_PORT =", dbPort)
	fmt.Println("DB_USER =", dbUser)
	fmt.Println("DB_PASSWORD =", dbPassword)
	fmt.Println("DB_NAME =", dbName)
	fmt.Println("===================")

	// --- CREATE DATABASE IF NOT EXISTS ---
	if dbName == "user_service" {
		adminDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
			dbHost, dbUser, dbPassword, dbPort)
		adminDb, err := sql.Open("postgres", adminDsn)
		if err != nil {
			log.Fatalf("Failed to connect to postgres for db creation: %v", err)
		}
		defer adminDb.Close()
		var exists bool
		checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
		err = adminDb.QueryRow(checkQuery).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check if database exists: %v", err)
		}
		if !exists {
			createQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
			_, err = adminDb.Exec(createQuery)
			if err != nil {
				log.Fatalf("Failed to create database %s: %v", dbName, err)
			}
			log.Printf("Database %s created successfully", dbName)
		}
	}

	// Format DSN for PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db
}

// CloseDatabaseConnection closes the database connection
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	err = dbSQL.Close()
	if err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}

	log.Println("PostgreSQL database connection closed")
}
