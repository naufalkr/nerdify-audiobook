package dto

import "errors"

// Additional error messages for user interactions
const (
	MESSAGE_UPDATE_USER_DATA_SUCCESS = "data pengguna berhasil diperbarui"
)

var (
	// Additional errors related to user data update
	ErrEmailAlreadyInUse    = errors.New("email sudah digunakan oleh pengguna lain")
	ErrUsernameAlreadyInUse = errors.New("username sudah digunakan oleh pengguna lain")

	// Tenant and role related errors
	ErrRoleNotFound = errors.New("role tidak ditemukan")
	ErrUnauthorized = errors.New("tidak memiliki akses")
	ErrInvalidID    = errors.New("ID tidak valid")

	// Tenant specific errors
	ErrInvalidTenantID        = errors.New("ID tenant tidak valid")
	ErrInvalidUserID          = errors.New("ID user tidak valid")
	ErrCannotInviteSuperadmin = errors.New("tidak dapat mengundang superadmin ke tenant")
	ErrTenantMaxUsersReached  = errors.New("tenant telah mencapai batas maksimum user")
	ErrUserInSameTenant       = errors.New("User sudah berada dalam tenant ini")
	ErrUserInOtherTenant      = errors.New("User sudah berada dalam tenant lain")
)
