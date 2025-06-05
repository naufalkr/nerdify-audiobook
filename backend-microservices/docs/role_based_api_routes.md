# Role-Based API Routes Documentation

## Overview
API ini menggunakan system role-based access control (RBAC) dengan 3 level utama:
1. **User** - User biasa
2. **Admin** - Administrator tenant
3. **SuperAdmin** - Administrator sistem

## Route Structure

### 1. Public Routes (No Authentication Required)
```
POST /api/auth/register
POST /api/auth/login
POST /api/auth/verify-email
POST /api/auth/resend-verification-email
POST /api/auth/forgot-password
POST /api/auth/reset-password
POST /api/auth/refresh-token
```

### 2. User Routes (Authentication Required)
**Base Path**: `/api/users/*` dan `/api/tenants/*`

#### Profile Management
```
GET    /api/users/profile
PUT    /api/users/profile
PATCH  /api/users/profile
DELETE /api/users/profile
POST   /api/users/profile/upload-image
PUT    /api/users/email
POST   /api/users/email/verify
POST   /api/users/logout
```

#### Tenant Context Management
```
GET  /api/v1/user-tenant/current           # Get current tenant
PUT  /api/v1/user-tenant/current           # Set current tenant
GET  /api/v1/user-tenant/tenants           # Get user tenants
POST /api/v1/user-tenant/switch            # Switch tenant
POST /api/v1/user-tenant/validate-access   # Validate user tenant access
```

#### User Tenant Information
```
GET /api/tenants/user-tenants              # Get tenants user belongs to
GET /api/tenants/detail/tenant             # Get current tenant details
GET /api/tenants/users                     # Get users in current tenant
```

### 3. Admin Routes (Admin + SuperAdmin Access)
**Base Path**: `/api/v1/admin/*`

#### Tenant Management
```
GET   /api/v1/admin/tenant                 # Get current tenant details
PUT   /api/v1/admin/tenant                 # Update current tenant
PATCH /api/v1/admin/tenant/contact         # Update tenant contact
PATCH /api/v1/admin/tenant/logo            # Update tenant logo
```

#### User Management in Tenant
```
GET    /api/v1/admin/users                 # Get users in tenant
POST   /api/v1/admin/users/invite          # Invite user to tenant
DELETE /api/v1/admin/users/:userID         # Remove user from tenant
GET    /api/v1/admin/users/:userID/profile # Get user profile
```

#### Subscription Management
```
GET  /api/v1/admin/subscription            # Get subscription info
POST /api/v1/admin/subscription            # Update subscription
```

#### Access Control
```
POST /api/v1/admin/validate-user-access    # Validate user access
```

#### Audit & Monitoring
```
GET /api/v1/admin/audit-logs               # Get audit logs for tenant
GET /api/v1/admin/stats/users              # User statistics
GET /api/v1/admin/stats/usage              # Usage statistics
```

### 4. SuperAdmin Routes (SuperAdmin Only)
**Base Path**: `/api/v1/superadmin/*`

#### System-Wide User Management
```
GET    /api/v1/superadmin/users            # List all users
POST   /api/v1/superadmin/users            # Create user
GET    /api/v1/superadmin/users/:id        # Get user by ID
PUT    /api/v1/superadmin/users/:id        # Update user
DELETE /api/v1/superadmin/users/:id        # Soft delete user
DELETE /api/v1/superadmin/users/:id/permanent # Hard delete user
POST   /api/v1/superadmin/users/:id/verify-email # Verify user email
PUT    /api/v1/superadmin/users/:id/role   # Change user role
```

#### System-Wide Tenant Management
```
GET    /api/v1/superadmin/tenants          # List all tenants
POST   /api/v1/superadmin/tenants          # Create tenant
GET    /api/v1/superadmin/tenants/:tenantID # Get tenant details
PUT    /api/v1/superadmin/tenants/:tenantID # Update tenant
DELETE /api/v1/superadmin/tenants/:tenantID # Delete tenant
```

#### Tenant-User Relationship Management
```
GET    /api/v1/superadmin/tenants/:tenantID/users              # Get tenant users
POST   /api/v1/superadmin/tenants/:tenantID/users/:userID/invite # Invite user
DELETE /api/v1/superadmin/tenants/:tenantID/users/:userID      # Remove user
POST   /api/v1/superadmin/tenants/:tenantID/users/:userID/promote # Promote to admin
POST   /api/v1/superadmin/tenants/:tenantID/users/:userID/demote  # Demote from admin
```

#### Subscription & Billing Management
```
GET  /api/v1/superadmin/tenants/:tenantID/subscription # Get subscription
POST /api/v1/superadmin/tenants/:tenantID/subscription # Update subscription
```

#### Tenant Customization
```
PATCH /api/v1/superadmin/tenants/:tenantID/contact # Update tenant contact
PATCH /api/v1/superadmin/tenants/:tenantID/logo    # Update tenant logo
```

#### Role Management
```
GET  /api/v1/superadmin/roles              # List all roles
POST /api/v1/superadmin/roles              # Create role
GET  /api/v1/superadmin/roles/:id          # Get role by ID
GET  /api/v1/superadmin/roles/name/:name   # Get role by name
PUT  /api/v1/superadmin/roles/:id          # Update role
GET  /api/v1/superadmin/roles/system       # Get system roles
POST /api/v1/superadmin/roles/seed         # Seed default roles
```

#### System Audit & Monitoring
```
GET /api/v1/superadmin/audit-logs                    # Get all audit logs
GET /api/v1/superadmin/audit-logs/users/:userID     # User-specific audit logs
GET /api/v1/superadmin/audit-logs/tenants/:tenantID # Tenant-specific audit logs
```

#### System Statistics & Analytics
```
GET /api/v1/superadmin/stats/overview      # System overview statistics
GET /api/v1/superadmin/stats/tenants       # Tenant statistics
GET /api/v1/superadmin/stats/users         # User statistics
```

#### System Maintenance
```
POST /api/v1/superadmin/system/maintenance/enable  # Enable maintenance mode
POST /api/v1/superadmin/system/maintenance/disable # Disable maintenance mode
```

#### External API Management
```
GET    /api/v1/superadmin/api-keys         # List API keys
POST   /api/v1/superadmin/api-keys         # Create API key
DELETE /api/v1/superadmin/api-keys/:keyID  # Revoke API key
```

### 5. External API Routes (For Microservice Communication)
**Base Path**: `/api/external/*`

#### Tenant Validation
```
GET  /api/external/tenants/:id/validate    # Validate tenant access
GET  /api/external/tenants                 # List tenants for external service
GET  /api/external/tenants/:id             # Get tenant by ID
GET  /api/external/tenants/:id/subscription # Get tenant subscription
GET  /api/external/tenants/:id/limits      # Get tenant limits
GET  /api/external/tenants/:id/users       # Get tenant users
POST /api/external/tenants/:id/validate-user-access # Validate user access
GET  /api/external/users/:userId/tenants   # Get user tenants
```

#### Authentication & Authorization
```
POST /api/external/auth/validate-token     # Validate JWT token
GET  /api/external/auth/user-info          # Get user info from token
POST /api/external/auth/validate-permissions # Validate user permissions
GET  /api/external/auth/validate-superadmin # Validate if user is SuperAdmin
```

## Authentication

### Headers Required
```
Authorization: Bearer <jwt_token>
```

### Token Claims
```json
{
  "userID": "uuid",
  "email": "user@example.com",
  "roleName": "Admin|SuperAdmin|User",
  "exp": timestamp,
  "iat": timestamp
}
```

## Response Format

### Success Response
```json
{
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "error": "Error message",
  "details": "Additional error details"
}
```

## Role Permissions Summary

| Feature | User | Admin | SuperAdmin |
|---------|------|-------|------------|
| Own Profile Management | ✅ | ✅ | ✅ |
| Tenant Context Management | ✅ | ✅ | ✅ |
| Manage Current Tenant | ❌ | ✅ | ✅ |
| Manage Users in Tenant | ❌ | ✅ | ✅ |
| View Tenant Audit Logs | ❌ | ✅ | ✅ |
| Manage Any Tenant | ❌ | ❌ | ✅ |
| Manage Any User | ❌ | ❌ | ✅ |
| System-wide Audit Logs | ❌ | ❌ | ✅ |
| Role Management | ❌ | ❌ | ✅ |
| System Maintenance | ❌ | ❌ | ✅ |
| External API Management | ❌ | ❌ | ✅ |

## Migration from Legacy Routes

Beberapa route lama masih tersedia untuk backward compatibility:
- `/api/admin/*` - Legacy admin routes
- `/api/superadmin/*` - Legacy superadmin routes (different from new `/api/v1/superadmin/*`)
- `/roles/*` - Legacy role management

Disarankan untuk menggunakan route baru dengan prefix `/api/v1/` untuk fitur yang lebih lengkap dan konsisten.
