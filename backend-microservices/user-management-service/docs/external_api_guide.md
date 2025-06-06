# External API Guide untuk Microservice Integration

## Overview

User Management Service menyediakan External API untuk komunikasi antar microservice, khususnya untuk service seperti Asset Management yang membutuhkan autentikasi dan otorisasi.

## Authentication

Semua External API endpoint memerlukan **API Key** yang dikirim via header:

```
X-API-Key: your-service-api-key
```

**Pengecualian:** Endpoint `/api/external/auth/validate-superadmin` tidak memerlukan API Key, hanya memerlukan Bearer token.

## Available Endpoints

### 1. Authentication & Authorization APIs

#### POST /api/external/auth/validate-token
Memvalidasi JWT token dari microservice lain.

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "isValid": true,
  "userInfo": {
    "userID": "user-uuid",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "Admin",
    "roleID": "role-uuid", 
    "status": "active",
    "isActive": true
  }
}
```

#### GET /api/external/auth/validate-superadmin
Memvalidasi apakah user memiliki role SuperAdmin.

**Header:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-API-Key: your-api-key-here
```

**Response (SuperAdmin):**
```json
{
  "valid": true,
  "userID": "user-uuid",
  "userRole": "SUPERADMIN",
  "isSuperAdmin": true
}
```

**Response (Non-SuperAdmin):**
```json
{
  "valid": false,
  "userID": "user-uuid",
  "userRole": "Admin",
  "isSuperAdmin": false
}
```

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "valid": true,
  "userID": "123e4567-e89b-12d3-a456-426614174000",
  "userRole": "USER",
  "email": "user@example.com"
}
```

#### GET /api/external/auth/user-info
Mendapatkan informasi user dari Authorization header.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**
```json
{
  "userID": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "userRole": "USER"
}
```

#### POST /api/external/auth/validate-user-permissions
Memvalidasi apakah user memiliki permission tertentu.

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tenantId": "123e4567-e89b-12d3-a456-426614174000",
  "requiredRole": "ADMIN",
  "permissions": ["read:assets", "write:assets"]
}
```

**Response:**
```json
{
  "valid": true,
  "hasRolePermission": true,
  "hasTenantAccess": true,
  "userID": "123e4567-e89b-12d3-a456-426614174000",
  "userRole": "ADMIN"
}
```

### 2. Tenant Management APIs

#### GET /api/external/tenants
Mendapatkan daftar semua tenant.

**Response:**
```json
{
  "tenants": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Company A",
      "subscriptionPlan": "premium",
      "isActive": true
    }
  ]
}
```

#### GET /api/external/tenants/:id
Mendapatkan detail tenant by ID.

**Response:**
```json
{
  "tenant": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Company A",
    "description": "Tech company",
    "subscriptionPlan": "premium",
    "maxUsers": 100,
    "isActive": true
  }
}
```

#### GET /api/external/tenants/:id/validate
Memvalidasi akses tenant.

**Response:**
```json
{
  "valid": true,
  "tenant": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Company A"
  }
}
```

### 3. Business Logic APIs

#### GET /api/external/tenants/:id/subscription
Mendapatkan informasi subscription tenant.

**Response:**
```json
{
  "tenantID": "123e4567-e89b-12d3-a456-426614174000",
  "subscriptionPlan": "premium",
  "subscriptionStartDate": "2024-01-01T00:00:00Z",
  "subscriptionEndDate": "2024-12-31T23:59:59Z",
  "isActive": true
}
```

#### GET /api/external/tenants/:id/limits
Mendapatkan limits tenant berdasarkan subscription.

**Response:**
```json
{
  "tenantID": "123e4567-e89b-12d3-a456-426614174000",
  "limits": {
    "maxUsers": 100,
    "maxAssets": 50,
    "maxRentals": 25,
    "subscriptionPlan": "premium"
  }
}
```

#### GET /api/external/tenants/:id/users
Mendapatkan daftar user dalam tenant.

**Response:**
```json
{
  "tenantID": "123e4567-e89b-12d3-a456-426614174000",
  "users": [
    {
      "id": "user-123",
      "email": "user@company.com",
      "fullName": "John Doe",
      "role": "USER"
    }
  ],
  "total": 15,
  "page": 1,
  "limit": 100
}
```

#### POST /api/external/tenants/:id/validate-user-access
Memvalidasi akses user ke tenant.

**Request Body:**
```json
{
  "userId": "123e4567-e89b-12d3-a456-426614174000"
}
```

**Response:**
```json
{
  "userID": "123e4567-e89b-12d3-a456-426614174000",
  "tenantID": "123e4567-e89b-12d3-a456-426614174000",
  "hasAccess": true
}
```

#### GET /api/external/users/:userId/tenants
Mendapatkan daftar tenant yang diikuti user.

**Response:**
```json
{
  "userID": "123e4567-e89b-12d3-a456-426614174000",
  "tenants": [
    {
      "id": "tenant-123",
      "name": "Company A",
      "role": "USER"
    }
  ]
}
```

## Environment Configuration

Tambahkan konfigurasi berikut ke `.env`:

```env
# API Keys untuk external services (comma-separated)
VALID_API_KEYS=asset-management-key,other-service-key

# Atau untuk development/testing
VALID_API_KEYS=alat-service-api-key
```

## Error Responses

Semua endpoint dapat mengembalikan error response dalam format:

```json
{
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `400` - Bad Request (invalid parameters)
- `401` - Unauthorized (missing/invalid API key or token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource not found)
- `500` - Internal Server Error

## Rate Limiting

External API endpoints tidak memiliki rate limiting khusus, namun direkomendasikan untuk mengimplementasikan caching di sisi client untuk mengurangi beban server.

## Security Notes

1. **API Keys**: Simpan API key dengan aman dan jangan expose di logs
2. **JWT Tokens**: Validate expiration dan signature sebelum menggunakan
3. **HTTPS**: Selalu gunakan HTTPS di production
4. **Audit Logs**: Semua akses external API akan tercatat di audit logs
