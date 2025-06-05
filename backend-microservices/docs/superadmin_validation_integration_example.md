# Integration Example: SuperAdmin Validation in Asset Management Service

This example shows how the Asset Management service can validate if a user has SuperAdmin role by communicating with the User Management service.

## Go Implementation

```go
package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	// UserManagementServiceURL is the base URL of the User Management service
	UserManagementServiceURL = "http://user-management-service:3120"
)

// SuperAdminResponse represents the response from the validate-superadmin endpoint
type SuperAdminResponse struct {
	Valid        bool   `json:"valid"`
	UserID       string `json:"userID"`
	UserRole     string `json:"userRole"`
	IsSuperAdmin bool   `json:"isSuperAdmin"`
}

// IsSuperAdmin checks if the user identified by the JWT token is a SuperAdmin
func IsSuperAdmin(token string) (bool, error) {
	// Create request
	req, err := http.NewRequest(
		"GET",
		UserManagementServiceURL+"/api/external/auth/validate-superadmin",
		nil
	)
	if err != nil {
		return false, err
	}

	// Set headers - only Authorization is needed, X-API-Key is not required for this endpoint
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("failed to validate superadmin: " + string(body))
	}

	// Parse response
	var response SuperAdminResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return false, err
	}

	return response.Valid, nil
}
```

## Example Usage in HTTP Handler

```go
// RequireSuperAdmin is a middleware that checks if the user is a SuperAdmin
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from header
		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		// Check if user is SuperAdmin
		isSuperAdmin, err := auth.IsSuperAdmin(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate SuperAdmin role: " + err.Error()})
			c.Abort()
			return
		}

		// Return error if user is not SuperAdmin
		if !isSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "SuperAdmin role required"})
			c.Abort()
			return
		}

		// Continue if user is SuperAdmin
		c.Next()
	}
}

// Example protected route
func SetupSuperAdminRoutes(router *gin.Engine) {
	superAdminRoutes := router.Group("/api/v1/assets/system")
	superAdminRoutes.Use(RequireSuperAdmin())
	{
		superAdminRoutes.GET("/statistics", GetSystemStatistics)
		superAdminRoutes.POST("/maintenance", ToggleSystemMaintenance)
	}
}
```

## Error Handling

The service should handle different error scenarios:

1. Network errors - when User Management service is unreachable
2. Authentication errors - when API key is invalid
3. Validation errors - when token is expired or invalid
4. Permission errors - when user doesn't have SuperAdmin role

```go
// Example error handler
func handleSuperAdminValidationError(err error) {
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "connection refused"):
			log.Error("User Management service is down")
			// Implement fallback mechanism or return service unavailable
			
		case strings.Contains(err.Error(), "invalid API key"):
			log.Error("Invalid API key for service communication")
			// Alert operations team about configuration issue
			
		case strings.Contains(err.Error(), "token is expired"):
			// Return appropriate response to client
			
		default:
			log.Error("Unexpected error during SuperAdmin validation: " + err.Error())
		}
	}
}
```

## Performance Considerations

For high-traffic services, consider:

1. Implementing a local caching mechanism for validated tokens
2. Using circuit breakers to handle User Management service outages
3. Setting appropriate timeouts for HTTP requests
