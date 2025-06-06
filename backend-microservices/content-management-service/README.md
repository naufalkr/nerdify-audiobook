# Content Management Service

Content Management Service for managing audiobook data including categories, audiobooks, metadata, and analytics.

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin (HTTP Router)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: External Auth Service (localhost:3160)

## Database Schema

Main tables:
- `categories` - Audiobook categories
- `audiobooks` - Main audiobook data
- `metadata` - Audiobook metadata
- `audiobook_analytics` - User activity tracking

## Quick Start

### 1. Navigate to Service Directory
```bash
cd d:\Kampus\SEMESTER 6\PPL\nerdify-audiobook\backend-microservices\content-management-service
```

### 2. Install Dependencies
```bash
go mod download
# or
go mod tidy
```

### 3. Database Setup
1. **Ensure PostgreSQL is running**
2. **Create database**: `nerdify-content-management`
3. **Setup .env file**:
```env
PORT=3163
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=nerdify-content-management
JWT_SECRET=your-secret-key
```

### 4. Run Database Migration & Seeding
```bash
go run cmd/seeder/main.go
```

### 5. Start Service
```bash
go run main.go
```

Service runs on **port 3163**.

## API Endpoints

### Public Endpoints
- `GET /api/health` - Health check
- `GET /api/public/audiobooks` - List audiobooks
- `GET /api/public/audiobooks/:id` - Get audiobook
- `GET /api/public/audiobooks/search?q=query` - Search audiobooks
- `GET /api/public/categories` - List categories

### Admin Endpoints (Requires SUPERADMIN)
- `POST /api/admin/audiobooks` - Create audiobook
- `PUT /api/admin/audiobooks/:id` - Update audiobook
- `DELETE /api/admin/audiobooks/:id` - Delete audiobook
- `POST /api/admin/categories` - Create category
- `PUT /api/admin/categories/:id` - Update category

### Query Parameters
- `?page=1` - Page number
- `?limit=20` - Items per page (max: 100)
- `?q=search` - Search query

## Authentication

### Headers
```bash
Authorization: Bearer <jwt_token>
```

### Role Requirements
- **Public endpoints**: No auth required
- **Admin endpoints**: SUPERADMIN role required

## Response Format

### Success
```json
{
  "message": "Success",
  "data": { ... }
}
```

### List with Pagination
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error
```json
{
  "error": "Error message"
}
```

## Seeder Commands

```bash
# Run migration + seeding
go run cmd/seeder/main.go

# Migration only
go run cmd/seeder/main.go -migrate

# Seeding only
go run cmd/seeder/main.go -seed

# Show statistics
go run cmd/seeder/main.go -stats

# Clear data
go run cmd/seeder/main.go -clear
```

## Build & Deploy

```bash
# Build
go build -o content-service main.go

# Run built binary
./content-service
```
