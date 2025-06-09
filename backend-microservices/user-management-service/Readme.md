# Nerdify Backend

Backend user-management-service for Nerdify application built with Go and PostgreSQL.

## Prerequisites

Before you begin, ensure you have the following installed:
- Go (version 1.22 or later)
- Docker & Docker Compose
- PostgreSQL (if running without Docker)
- Git

## Getting Started


### Environment Setup
Create .env file in root directory:
```bash
PORT=3160
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=nerdify-user-management
JWT_SECRET=your-secret-key
```

### Run Locally (Without Docker)

1. Install dependencies:
    ```bash
    go mod download
    ```
2. Update database connection in `.env`:
    ```bash
    DB_HOST=localhost
    ```
3. Run the application:
    ```bash
    # Kalo baru pertama kali run, seed ke db nya:
    go run main.go -seed

    # Run biasa
    go run main.go
    ```

## Service Integration

### External API for Role Validation

This service provides external APIs for validating user roles (including SuperAdmin) from other microservices:

#### Validating SuperAdmin Role

```bash
GET /api/external/auth/validate-superadmin
```

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
X-API-Key: your-api-key
```

**Response:**
```json
{
  "valid": true,
  "userID": "user-uuid",
  "userRole": "SUPERADMIN",
  "isSuperAdmin": true
}
```

For complete documentation on external APIs, see:
- [External API Documentation](/docs/external_api_guide.md)
- [SuperAdmin Validation API](/docs/external_api_superadmin_validation.md)
- [Role-Based API Routes](/docs/role_based_api_routes.md)