# Asset Management Service - Integration Example

Contoh implementasi Asset Management Service yang menggunakan User Management External API.

## Project Structure

```
asset_management/
├── main.go
├── go.mod
├── config/
│   └── config.go
├── models/
│   ├── asset.go
│   └── rental.go
├── services/
│   ├── asset_service.go
│   ├── rental_service.go
│   └── user_client.go
├── handlers/
│   ├── asset_handler.go
│   └── rental_handler.go
├── middleware/
│   └── auth_middleware.go
└── .env
```

## Environment Configuration (.env)

```env
PORT=8081
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=asset_management

# User Management Service
USER_MANAGEMENT_URL=http://localhost:8080
USER_MANAGEMENT_API_KEY=asset-management-key

# JWT Secret (should match User Management Service)
JWT_SECRET=your-jwt-secret-key
```

## Dependencies (go.mod)

```go
module asset_management

go 1.22

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/golang-jwt/jwt v3.2.2+incompatible
    github.com/google/uuid v1.6.0
    github.com/joho/godotenv v1.5.1
    gorm.io/gorm v1.25.11
    gorm.io/driver/postgres v1.5.9
)
```

## Models

### Asset Model (models/asset.go)

```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Asset struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    Name        string    `gorm:"size:255;not null" json:"name"`
    Category    string    `gorm:"size:100" json:"category"`
    Description string    `gorm:"size:1000" json:"description"`
    PricePerDay float64   `json:"pricePerDay"`
    IsAvailable bool      `gorm:"default:true" json:"isAvailable"`
    Location    string    `gorm:"size:255" json:"location"`
    ImageURL    string    `gorm:"size:500" json:"imageUrl"`
    
    CreatedAt   time.Time      `json:"createdAt"`
    UpdatedAt   time.Time      `json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Asset) BeforeCreate(tx *gorm.DB) error {
    a.ID = uuid.New()
    return nil
}
```

### Rental Model (models/rental.go)

```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Rental struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    AssetID   uuid.UUID `gorm:"type:uuid;not null" json:"assetId"`
    Asset     Asset     `gorm:"foreignKey:AssetID" json:"asset"`
    
    TenantID  uuid.UUID `gorm:"type:uuid;not null" json:"tenantId"`
    UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
    
    StartDate time.Time `gorm:"not null" json:"startDate"`
    EndDate   time.Time `gorm:"not null" json:"endDate"`
    
    Status    string    `gorm:"size:20;default:active" json:"status"` // active, completed, cancelled
    TotalCost float64   `json:"totalCost"`
    
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Rental) BeforeCreate(tx *gorm.DB) error {
    r.ID = uuid.New()
    return nil
}
```

## User Management Client (services/user_client.go)

```go
package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type UserClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

type TokenValidationResponse struct {
    Valid    bool   `json:"valid"`
    UserID   string `json:"userID"`
    UserRole string `json:"userRole"`
    Email    string `json:"email"`
}

type TenantLimitsResponse struct {
    TenantID string `json:"tenantID"`
    Limits   struct {
        MaxUsers    int    `json:"maxUsers"`
        MaxAssets   int    `json:"maxAssets"`
        MaxRentals  int    `json:"maxRentals"`
        Plan        string `json:"subscriptionPlan"`
    } `json:"limits"`
}

type UserTenantAccessResponse struct {
    UserID    string `json:"userID"`
    TenantID  string `json:"tenantID"`
    HasAccess bool   `json:"hasAccess"`
}

func NewUserClient() *UserClient {
    return &UserClient{
        baseURL: os.Getenv("USER_MANAGEMENT_URL"),
        apiKey:  os.Getenv("USER_MANAGEMENT_API_KEY"),
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (u *UserClient) ValidateToken(token string) (*TokenValidationResponse, error) {
    reqBody := map[string]string{"token": token}
    jsonData, _ := json.Marshal(reqBody)
    
    req, err := http.NewRequest("POST", u.baseURL+"/api/external/auth/validate-token", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", u.apiKey)
    
    resp, err := u.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("token validation failed: %d", resp.StatusCode)
    }
    
    var result TokenValidationResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

func (u *UserClient) GetTenantLimits(tenantID string) (*TenantLimitsResponse, error) {
    req, err := http.NewRequest("GET", u.baseURL+"/api/external/tenants/"+tenantID+"/limits", nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("X-API-Key", u.apiKey)
    
    resp, err := u.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get tenant limits: %d", resp.StatusCode)
    }
    
    var result TenantLimitsResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

func (u *UserClient) ValidateUserTenantAccess(userID, tenantID string) (*UserTenantAccessResponse, error) {
    reqBody := map[string]string{"userId": userID}
    jsonData, _ := json.Marshal(reqBody)
    
    req, err := http.NewRequest("POST", u.baseURL+"/api/external/tenants/"+tenantID+"/validate-user-access", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", u.apiKey)
    
    resp, err := u.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to validate user access: %d", resp.StatusCode)
    }
    
    var result UserTenantAccessResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

## Authentication Middleware (middleware/auth_middleware.go)

```go
package middleware

import (
    "asset_management/services"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware(userClient *services.UserClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Extract token
        token := authHeader
        if strings.HasPrefix(authHeader, "Bearer ") {
            token = authHeader[7:]
        }

        // Validate token with User Management Service
        userInfo, err := userClient.ValidateToken(token)
        if err != nil || !userInfo.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Store user info in context
        c.Set("userID", userInfo.UserID)
        c.Set("userRole", userInfo.UserRole)
        c.Set("email", userInfo.Email)
        
        c.Next()
    }
}

func TenantAccessMiddleware(userClient *services.UserClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("userID")
        tenantID := c.GetHeader("X-Tenant-ID")
        
        if tenantID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID header required"})
            c.Abort()
            return
        }

        // Validate user access to tenant
        access, err := userClient.ValidateUserTenantAccess(userID, tenantID)
        if err != nil || !access.HasAccess {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to tenant"})
            c.Abort()
            return
        }

        c.Set("tenantID", tenantID)
        c.Next()
    }
}
```

## Asset Handler (handlers/asset_handler.go)

```go
package handlers

import (
    "asset_management/models"
    "asset_management/services"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type AssetHandler struct {
    assetService *services.AssetService
    userClient   *services.UserClient
}

func NewAssetHandler(assetService *services.AssetService, userClient *services.UserClient) *AssetHandler {
    return &AssetHandler{
        assetService: assetService,
        userClient:   userClient,
    }
}

func (h *AssetHandler) GetAssets(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    category := c.Query("category")

    assets, total, err := h.assetService.GetAssets(page, limit, category)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "assets": assets,
        "total":  total,
        "page":   page,
        "limit":  limit,
    })
}

func (h *AssetHandler) RentAsset(c *gin.Context) {
    assetID := c.Param("id")
    userID := c.GetString("userID")
    tenantID := c.GetString("tenantID")

    var req struct {
        StartDate string `json:"startDate" binding:"required"`
        EndDate   string `json:"endDate" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check tenant limits
    limits, err := h.userClient.GetTenantLimits(tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check tenant limits"})
        return
    }

    // Check current rental count
    currentRentals, err := h.assetService.GetActiveTenantRentals(tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check current rentals"})
        return
    }

    if len(currentRentals) >= limits.Limits.MaxRentals {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Rental limit exceeded",
            "current": len(currentRentals),
            "limit": limits.Limits.MaxRentals,
        })
        return
    }

    // Create rental
    rental, err := h.assetService.CreateRental(assetID, userID, tenantID, req.StartDate, req.EndDate)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"rental": rental})
}

func (h *AssetHandler) GetTenantRentals(c *gin.Context) {
    tenantID := c.GetString("tenantID")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    rentals, total, err := h.assetService.GetTenantRentals(tenantID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "rentals": rentals,
        "total":   total,
        "page":    page,
        "limit":   limit,
    })
}
```

## Main Application (main.go)

```go
package main

import (
    "asset_management/handlers"
    "asset_management/middleware"
    "asset_management/services"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    // Initialize services
    userClient := services.NewUserClient()
    assetService := services.NewAssetService() // Implementation not shown
    
    // Initialize handlers
    assetHandler := handlers.NewAssetHandler(assetService, userClient)

    // Setup router
    router := gin.Default()

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Public routes
    api := router.Group("/api")
    {
        // Public asset listing
        api.GET("/assets", assetHandler.GetAssets)
        api.GET("/assets/:id", assetHandler.GetAssetByID)
    }

    // Protected routes
    protected := api.Group("/")
    protected.Use(middleware.AuthMiddleware(userClient))
    protected.Use(middleware.TenantAccessMiddleware(userClient))
    {
        // Asset rental management
        protected.POST("/assets/:id/rent", assetHandler.RentAsset)
        protected.GET("/rentals", assetHandler.GetTenantRentals)
        protected.PUT("/rentals/:id/return", assetHandler.ReturnAsset)
        protected.GET("/rentals/:id", assetHandler.GetRental)
    }

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }

    log.Printf("Asset Management Service starting on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

## Usage Example

### 1. Get Assets (Public)
```bash
curl -X GET "http://localhost:8081/api/assets?page=1&limit=10&category=laptop"
```

### 2. Rent Asset (Protected)
```bash
curl -X POST "http://localhost:8081/api/assets/123e4567-e89b-12d3-a456-426614174000/rent" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: YOUR_TENANT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "startDate": "2024-01-01",
    "endDate": "2024-01-15"
  }'
```

### 3. Get Tenant Rentals (Protected)
```bash
curl -X GET "http://localhost:8081/api/rentals?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-Tenant-ID: YOUR_TENANT_ID"
```

## Key Integration Points

1. **Token Validation**: Setiap request protected divalidasi melalui User Management Service
2. **Tenant Access Control**: Memastikan user memiliki akses ke tenant yang diminta
3. **Business Rules**: Menggunakan tenant limits untuk membatasi rental berdasarkan subscription
4. **Audit Trail**: Semua aktivitas akan tercatat di User Management Service audit logs

Implementasi ini menunjukkan bagaimana Asset Management Service dapat securely berinteraksi dengan User Management Service menggunakan External API yang telah kita buat.
