package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UserManagementService handles external API calls to auth service
type UserManagementService struct {
	baseURL    string
	httpClient *http.Client
}

// SuperAdminValidationResponse represents the response from validate-superadmin endpoint
type SuperAdminValidationResponse struct {
	IsSuperAdmin bool   `json:"isSuperAdmin"`
	UserID       string `json:"userID"`
	UserRole     string `json:"userRole"`
	Valid        bool   `json:"valid"`
}

// NewUserManagementService creates a new user management service
func NewUserManagementService(baseURL string) *UserManagementService {
	return &UserManagementService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidateSuperAdmin validates if a user has SuperAdmin privileges
func (s *UserManagementService) ValidateSuperAdmin(ctx context.Context, token string) (*SuperAdminValidationResponse, error) {
	url := fmt.Sprintf("%s/api/external/auth/validate-superadmin", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validation failed with status: %d", resp.StatusCode)
	}

	var validationResponse SuperAdminValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &validationResponse, nil
}
