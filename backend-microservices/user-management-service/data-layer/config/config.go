package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all non-database configuration for the application
type Config struct {
	// SMTP configuration
	SMTPHost         string
	SMTPPort         int
	SMTPSenderName   string
	SMTPAuthEmail    string
	SMTPAuthPassword string
	AppURL           string

	// JWT configuration
	AccessSecret        string
	RefreshSecret       string
	EmailSecret         string
	PasswordResetSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Parse SMTP port
	smtpPort := 587 // Default SMTP port
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			smtpPort = port
		}
	}

	config := &Config{
		SMTPHost:         os.Getenv("SMTP_HOST"),
		SMTPPort:         smtpPort,
		SMTPSenderName:   os.Getenv("SMTP_SENDER_NAME"),
		SMTPAuthEmail:    os.Getenv("SMTP_AUTH_EMAIL"),
		SMTPAuthPassword: os.Getenv("SMTP_AUTH_PASSWORD"),
		AppURL:           os.Getenv("APP_URL"),

		// JWT secrets
		AccessSecret:        os.Getenv("JWT_ACCESS_SECRET"),
		RefreshSecret:       os.Getenv("JWT_REFRESH_SECRET"),
		EmailSecret:         os.Getenv("JWT_EMAIL_SECRET"),
		PasswordResetSecret: os.Getenv("JWT_PASSWORD_RESET_SECRET"),
	}

	return config, nil
}
