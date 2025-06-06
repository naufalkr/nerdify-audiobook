# User Management System API Documentation

## Table of Contents
1. [Authentication Routes](#authentication-routes)
2. [User Routes](#user-routes)
3. [Admin Routes](#admin-routes)
4. [SuperAdmin Routes](#superadmin-routes)
5. [Tenant Routes](#tenant-routes)
6. [Role Routes](#role-routes)
7. [Audit Routes](#audit-routes)
8. [External API Routes](#external-api-routes)

## Authentication Routes
Base path: `/api/auth`

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/register` | Register a new user | No |
| POST | `/login` | Login user | No |
| POST | `/verify-email` | Verify user email | No |
| POST | `/resend-verification-email` | Resend verification email | No |
| POST | `/forgot-password` | Request password reset | No |
| POST | `/reset-password` | Reset password | No |
| POST | `/refresh-token` | Refresh access token | No |

## User Routes
Base path: `/api/users`

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/logout` | Logout user | Yes |
| GET | `/profile` | Get user profile | Yes |
| PUT/PATCH | `/profile` | Update user profile | Yes |
| DELETE | `/profile` | Delete own account | Yes |
| POST | `/profile/upload-image` | Upload profile image | Yes |
| PUT | `/email` | Update email | Yes |
| POST | `/email/verify` | Verify email update | Yes |

## Admin Routes
Base path: `/api/admin/users`

| Method | Endpoint | Description | Authentication Required | Role Required |
|--------|----------|-------------|------------------------|---------------|
| GET | `/` | List all users (with pagination) | Yes | SUPERADMIN |
| POST | `/` | Create new user | Yes | SUPERADMIN |
| GET | `/:id` | Get user profile by ID | Yes | SUPERADMIN |
| PUT | `/:id` | Update user data | Yes | SUPERADMIN |
| PUT | `/:id/role` | Change user role | Yes | SUPERADMIN |
| POST | `/:id/verify-email` | Verify user email | Yes | SUPERADMIN |
| DELETE | `/:id` | Soft delete user | Yes | SUPERADMIN |
| DELETE | `/:id/permanent` | Hard delete user | Yes | SUPERADMIN |

## SuperAdmin Routes
Base path: `/api/superadmin`

| Method | Endpoint | Description | Authentication Required | Role Required |
|--------|----------|-------------|------------------------|---------------|
| GET | `/users` | List all users | Yes | SUPERADMIN |
| GET | `/tenants` | List all tenants | Yes | SUPERADMIN |
| GET | `/roles` | List all roles | Yes | SUPERADMIN |
| GET | `/audit-logs` | Get audit logs | Yes | SUPERADMIN |

## Tenant Routes
Base path: `/api/tenants`

| Method | Endpoint | Description | Authentication Required | Role Required |
|--------|----------|-------------|------------------------|---------------|
| GET | `/` | List tenants | Yes | ADMIN/SUPERADMIN |
| POST | `/` | Create tenant | Yes | SUPERADMIN |
| GET | `/:id` | Get tenant details | Yes | ADMIN/SUPERADMIN |
| PUT | `/:id` | Update tenant | Yes | SUPERADMIN |
| DELETE | `/:id` | Delete tenant | Yes | SUPERADMIN |

## Role Routes
Base path: `/api/roles`

| Method | Endpoint | Description | Authentication Required | Role Required |
|--------|----------|-------------|------------------------|---------------|
| GET | `/` | List roles | Yes | ADMIN/SUPERADMIN |
| POST | `/` | Create role | Yes | SUPERADMIN |
| GET | `/:id` | Get role details | Yes | ADMIN/SUPERADMIN |
| PUT | `/:id` | Update role | Yes | SUPERADMIN |
| DELETE | `/:id` | Delete role | Yes | SUPERADMIN |

## Audit Routes
Base path: `/api/audit`

| Method | Endpoint | Description | Authentication Required | Role Required |
|--------|----------|-------------|------------------------|---------------|
| GET | `/logs` | Get audit logs | Yes | ADMIN/SUPERADMIN |
| GET | `/logs/:id` | Get specific audit log | Yes | ADMIN/SUPERADMIN |

## External API Routes
Base path: `/api/external`

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| GET | `/tenant/:id` | Get tenant information | Yes (Service Token) |
| POST | `/tenant/validate` | Validate tenant | Yes (Service Token) |

## Notes
1. All authenticated routes require a valid JWT token in the Authorization header
2. Role-based access control is implemented for admin and superadmin routes
3. Service tokens are required for external API routes
4. Pagination is supported for list endpoints using query parameters `page` and `limit`
5. File upload endpoints support multipart/form-data
6. All responses are in JSON format
