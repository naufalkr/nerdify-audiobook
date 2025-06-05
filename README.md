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

### Factory Pattern
- **Purpose**: Creates instance objects without exposing class implementation
- **Implementation**: Applied to `AudiobookPlayer` and `UserSessionManager` objects with specific configurations
- **Location**: Used in audio player initialization and user session management

### Observer Pattern
- **Purpose**: Monitors state changes across the application
- **Implementation**: Tracks audiobook playback status changes and user progress updates
- **Location**: Audio player state management and progress tracking system

### Repository Pattern
- **Purpose**: Separates data access logic from UI components
- **Implementation**: Handles data operations from remote API, database, and cache
- **Location**: API service layer, data persistence, and state management


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
DB_NAME=audiobook
DB_USER=postgres
DB_PASSWORD=your_password
JWT_SECRET=your_secret
PORT=3160
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
Frontend runs on `http://localhost:3000`

### Docker (Alternative)
```bash
docker-compose up -d
```

## License

MIT License