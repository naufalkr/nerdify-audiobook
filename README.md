# Nerdify Audiobook Platform

A modern audiobook streaming platform built with Go microservice backend and React frontend, providing access to free audiobooks from LibriVox.

## Academic Context

This project is developed as a final project for Software Engineering (PPL) course, demonstrating:
- Full-stack web development
- Microservice architecture
- Design pattern implementation

## Architecture

```
├── Backend Microservice (Go + PostgreSQL)
├── Frontend (React)
└── Docker Support
```

## Design Patterns Implementation

This project implements several design patterns as part of Software Engineering (PPL) coursework:

### Repository Pattern
- **Purpose**: Separates data access logic from UI components and provides a centralized way to handle data operations
- **Implementation**: Applied across all data access layers with base repository class and specialized repositories

### Singleton Pattern
- **Purpose**: Ensures single instance creation and provides global access point for logging utilities
- **Implementation**: Applied to logging and utility classes that need to maintain state across the application

### Factory Pattern
- **Purpose**: Creates repository instances with custom configurations without exposing implementation details
- **Implementation**: Applied to repository creation and user session management


## Key Features

- JWT Authentication (register/login/logout)
- Audio streaming with progress tracking
- User profile management
- Search and browse audiobooks
- Responsive design
- Mock data support for standalone development

## Quick Start

### Prerequisites
- Go 1.19+
- Node.js 14+
- PostgreSQL (or Docker)

### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=user-management, catalog-services (tergantung service)
DB_USER=postgres
DB_PASSWORD=your_password
JWT_SECRET=your_secret
PORT=3160, 3161, 3162, 3163 (tergantung service)
```

### Frontend (.env)
```env
REACT_APP_USE_REAL_API=true
REACT_APP_API_BASE_URL=http://localhost:3160
```

### Backend Setup
```bash
cd backend-microservice
go mod download
cp .env.example .env  # Configure database settings

# Kalo baru pertama kali run, seed ke db nya:
go run main.go -seed

# Run biasa
go run main.go
```
Backend runs on `http://localhost:3160`

### Frontend Setup
```bash
cd frontend
npm install
cp .env.example .env  # Set REACT_APP_USE_REAL_API=true
npm start
```

## Authentication
### Superadmin
Username: superadmin@gmail.com
Password: superadmin123

### User
Username: user3@example.com
password: password123


Frontend runs on `http://localhost:3000`

### Docker (Alternative)
```bash
docker-compose up -d
```

## License

MIT License