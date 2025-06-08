package config

import (
	"fmt"
	"os"
)

// GetUserManagementBaseURL returns the base URL for user management service
func GetUserManagementBaseURL() string {
	host := os.Getenv("USER_MANAGEMENT_HOST")
	if host == "" {
		host = "localhost" // default
	}

	port := os.Getenv("USER_MANAGEMENT_PORT")
	if port == "" {
		port = "3120" // default
	}

	return fmt.Sprintf("http://%s:%s", host, port)
}

// GetUserManagementValidateURL returns the validation endpoint
func GetUserManagementValidateURL() string {
	url := os.Getenv("USER_MANAGEMENT_VALIDATE_URL")
	if url == "" {
		url = "/api/external/auth/validate-superadmin" // default
	}
	return url
}
