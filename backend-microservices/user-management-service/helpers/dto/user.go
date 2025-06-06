package dto

import (
	"errors"
	"mime/multipart"
	"time"
)

const (
	MESSAGE_REGISTER_SUCCESS                  = "berhasil mendaftarkan pengguna dan mengirim email verifikasi"
	MESSAGE_LOGIN_SUCCESS                     = "berhasil login"
	MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS   = "berhasil mengirim email verifikasi"
	MESSAGE_VERIFY_EMAIL_SUCCESS              = "berhasil verifikasi email dengan OTP"
	MESSAGE_LOGOUT_SUCCESS                    = "berhasil logout"
	MESSAGE_GET_ALL_USER_SUCCESS              = "berhasil mendapatkan semua pengguna"
	MESSAGE_DELETE_USER_SUCCESS               = "berhasil menghapus pengguna"
	MESSAGE_HARD_DELETE_USER_SUCCESS          = "berhasil menghapus pengguna secara permanen"
	MESSAGE_EDIT_USER_SUCCESS                 = "berhasil mengedit pengguna"
	MESSAGE_RESEND_VERIFICATION_EMAIL_SUCCESS = "berhasil mengirim ulang email verifikasi"
	MESSAGE_FORGOT_PASSWORD_SUCCESS           = "berhasil mengirim email reset password"
	MESSAGE_RESET_PASSWORD_SUCCESS            = "berhasil mereset password"
	MESSAGE_REFRESH_TOKEN_SUCCESS             = "berhasil memperpanjang token"
	MESSAGE_OTP_SENT_SUCCESS                  = "berhasil mengirim kode OTP"
	MESSAGE_OTP_VERIFY_FAILED                 = "kode OTP tidak valid atau sudah kadaluarsa"
	MESSAGE_OTP_MAX_ATTEMPTS                  = "terlalu banyak percobaan, silakan minta kode baru"
	MESSAGE_VERIFY_EMAIL_LINK_SUCCESS         = "berhasil verifikasi email dengan link"
)

var (
	ErrInvalidUsername          = errors.New("username tidak valid")
	ErrUserAlreadyExist         = errors.New("pengguna sudah terdaftar")
	ErrCreateUserFailed         = errors.New("gagal membuat pengguna")
	ErrSendEmailFailed          = errors.New("gagal mengirim email")
	ErrTokenExpired             = errors.New("token kadaluarsa")
	ErrUserNotVerified          = errors.New("pengguna belum diverifikasi")
	ErrUserNotFound             = errors.New("pengguna tidak ditemukan")
	ErrPasswordNotMatch         = errors.New("password tidak cocok")
	ErrLoginFailed              = errors.New("gagal login")
	ErrUserAlreadyVerified      = errors.New("pengguna sudah diverifikasi")
	ErrGetAllEmailLimit         = errors.New("gagal mendapatkan semua email")
	ErrDeleteEmailLimit         = errors.New("gagal menghapus email")
	ErrEmailLimitReached        = errors.New("batas email tercapai")
	ErrCreateEmailLimit         = errors.New("gagal membuat email")
	ErrTokenInvalid             = errors.New("token tidak valid")
	ErrUpdateUserFailed         = errors.New("gagal memperbarui pengguna")
	ErrInvalidEmail             = errors.New("email tidak valid")
	ErrPasswordResetFailed      = errors.New("gagal reset password")
	ErrRefreshTokenFailed       = errors.New("gagal refresh token")
	ErrSamePassword             = errors.New("password baru tidak boleh sama dengan password lama")
	ErrUpdateRoleFailed         = errors.New("gagal mengubah role user")
	ErrOTPNotMatch              = errors.New("kode OTP tidak sesuai")
	ErrOTPExpired               = errors.New("kode OTP sudah kadaluarsa")
	ErrTooManyAttempts          = errors.New("terlalu banyak percobaan verifikasi")
	ErrTooManyResendAttempts    = errors.New("terlalu banyak permintaan pengiriman ulang OTP")
	ErrInCooldownPeriod         = errors.New("mohon tunggu beberapa saat sebelum meminta OTP baru")
	ErrVerificationTokenInvalid = errors.New("token verifikasi tidak valid")
	ErrVerificationTokenExpired = errors.New("token verifikasi sudah kadaluarsa")
)

type (
	UserRegisterRequest struct {
		UserName  string  `json:"user_name" validate:"required"`
		Email     string  `json:"email" validate:"required,email"`
		Password  string  `json:"password" validate:"required,min=6"`
		FullName  string  `json:"full_name" validate:"required"`
		Alamat    string  `json:"alamat" validate:"required"`
		Latitude  float64 `json:"latitude" validate:"required"`
		Longitude float64 `json:"longitude" validate:"required"`
		RoleName  string  `json:"role_name,omitempty"` // Optional, for admin creation
		Status    string  `json:"status,omitempty"`    // Optional, for admin creation
	}

	UserEditRequest struct {
		UserName  string  `json:"user_name" validate:"required"`
		Email     string  `json:"email" validate:"required,email"`
		FullName  string  `json:"full_name" validate:"required"`
		Alamat    string  `json:"alamat" validate:"required"`
		Latitude  float64 `json:"latitude" validate:"required"`
		Longitude float64 `json:"longitude" validate:"required"`
	}

	UserSearchKeywordRequest struct {
		Keyword string `json:"keyword" validate:"required"`
	}

	UserLoginRequest struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Username string `json:"username" binding:"omitempty"`
		Password string `json:"password" binding:"required,min=6"`
	}
	UserResponse struct {
		ID              string `json:"id"`
		Username        string `json:"username"`
		Email           string `json:"email"`
		UserRole        string `json:"user_role"`
		IsVerified      bool   `json:"is_verified"`
		AccessToken     string `json:"access_token,omitempty"`
		RefreshToken    string `json:"refresh_token,omitempty"`
		ProfileImageURL string `json:"profile_image_url,omitempty"`
		Status          string `json:"status,omitempty"`
	}

	DetailUserResponse struct {
		ID        string                 `json:"id,omitempty"`
		UserName  string                 `json:"user_name,omitempty"`
		FullName  string                 `json:"full_name,omitempty"`
		Email     string                 `json:"email,omitempty"`
		Alamat    string                 `json:"alamat,omitempty"`
		Latitude  float64                `json:"latitude,omitempty"`
		Longitude float64                `json:"longitude,omitempty"`
		Alat      []DetailUserAccessData `json:"detail_user_access_data,omitempty"`
	}

	DetailUserAccessData struct {
		AlatID     string `json:"alat_id,omitempty"`
		AccessRole string `json:"access_role,omitempty"`
		AlatName   string `json:"alat_name,omitempty"`
	}

	SendVerificationEmailRequest struct {
		Email string `json:"email" binding:"required"`
	}

	VerifyEmailRequest struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required"`
	}

	VerifyEmailLinkRequest struct {
		Email string `json:"email" binding:"required,email"`
		Token string `json:"token" binding:"required"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	ResetPasswordRequest struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	SendOTPRequest struct {
		Email string `json:"email" binding:"required,email"`
	} // UserProfileResponse is the response format for user profile data
	UserProfileResponse struct {
		ID                         string     `json:"id"`
		Username                   string     `json:"username"`
		Email                      string     `json:"email"`
		FullName                   string     `json:"full_name"`
		Alamat                     string     `json:"alamat,omitempty"`
		Latitude                   float64    `json:"latitude,omitempty"`
		Longitude                  float64    `json:"longitude,omitempty"`
		ProfileImageURL            string     `json:"profile_image_url,omitempty"`
		RoleID                     string     `json:"role_id,omitempty"`
		Role                       string     `json:"role"`
		IsVerified                 bool       `json:"is_verified"`
		OTPCode                    string     `json:"otp_code,omitempty"`
		OTPCreatedAt               *time.Time `json:"otp_created_at,omitempty"`
		OTPAttemptCount            int        `json:"otp_attempt_count,omitempty"`
		ResendCount                int        `json:"resend_count,omitempty"`
		LastResendAt               *time.Time `json:"last_resend_at,omitempty"`
		CooldownStartedAt          *time.Time `json:"cooldown_started_at,omitempty"`
		VerificationToken          string     `json:"verification_token,omitempty"`
		VerificationTokenCreatedAt *time.Time `json:"verification_token_created_at,omitempty"`
		PendingEmail               string     `json:"pending_email,omitempty"`
		AccessToken                string     `json:"access_token,omitempty"`
		RefreshToken               string     `json:"refresh_token,omitempty"`
		TokenExpiry                time.Time  `json:"token_expiry,omitempty"`
		TokenCreatedAt             time.Time  `json:"token_created_at,omitempty"`
		Status                     string     `json:"status,omitempty"`
		CreatedAt                  time.Time  `json:"created_at"`
		UpdatedAt                  time.Time  `json:"updated_at"`
		DeletedAt                  time.Time  `json:"deleted_at,omitempty"`
		IsDeleted                  bool       `json:"is_deleted,omitempty"`
	}

	// UserDataUpdateRequest contains fields for updating user data by admin
	UserDataUpdateRequest struct {
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		FullName  string  `json:"full_name"`
		Alamat    string  `json:"alamat"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Status    string  `json:"status"`
	}

	// External API DTOs for other services
	ExternalUserInfoResponse struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		RoleID   string `json:"role_id"`
		Status   string `json:"status"`
		IsActive bool   `json:"is_active"`
	}

	// TokenValidationRequest for external API token validation
	TokenValidationRequest struct {
		Token string `json:"token" binding:"required"`
	}

	// TokenValidationResponse for external API token validation response
	TokenValidationResponse struct {
		IsValid   bool                      `json:"is_valid"`
		UserInfo  *ExternalUserInfoResponse `json:"user_info,omitempty"`
		ExpiresAt *time.Time                `json:"expires_at,omitempty"`
		Error     string                    `json:"error,omitempty"`
	}
)

// EmailUpdateRequest adalah DTO untuk permintaan update email
type EmailUpdateRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// EmailUpdateVerificationRequest adalah DTO untuk verifikasi update email
type EmailUpdateVerificationRequest struct {
	OTP string `json:"otp" binding:"required"`
}

// UserProfileUpdateRequest untuk pembaruan profil pengguna
type UserProfileUpdateRequest struct {
	UserName     string                `json:"user_name,omitempty"`
	Email        string                `json:"email,omitempty"`
	FullName     string                `json:"full_name,omitempty"`
	Alamat       string                `json:"alamat,omitempty"`
	Latitude     float64               `json:"latitude,omitempty"`
	Longitude    float64               `json:"longitude,omitempty"`
	ProfileImage *multipart.FileHeader `json:"profile_image,omitempty"`
}

// UserListRequest adalah DTO untuk permintaan daftar pengguna dengan pagination
type UserListRequest struct {
	Page     int    `json:"page" binding:"required,min=1"`
	PageSize int    `json:"page_size" binding:"required,min=1,max=100"`
	Search   string `json:"search,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   string `json:"status,omitempty"`
}

// UserListResponse adalah DTO untuk respons daftar pengguna
type UserListResponse struct {
	ID                         string     `json:"id"`
	Username                   string     `json:"username"`
	Email                      string     `json:"email"`
	UserName                   string     `json:"user_name"`
	FullName                   string     `json:"full_name"`
	Alamat                     string     `json:"alamat,omitempty"`
	Latitude                   float64    `json:"latitude,omitempty"`
	Longitude                  float64    `json:"longitude,omitempty"`
	ProfileImageURL            string     `json:"profile_image_url,omitempty"`
	RoleID                     string     `json:"role_id,omitempty"`
	Role                       string     `json:"role"`
	IsVerified                 bool       `json:"is_verified"`
	OTPCode                    string     `json:"otp_code,omitempty"`
	OTPCreatedAt               *time.Time `json:"otp_created_at,omitempty"`
	OTPAttemptCount            int        `json:"otp_attempt_count,omitempty"`
	ResendCount                int        `json:"resend_count,omitempty"`
	LastResendAt               *time.Time `json:"last_resend_at,omitempty"`
	CooldownStartedAt          *time.Time `json:"cooldown_started_at,omitempty"`
	VerificationToken          string     `json:"verification_token,omitempty"`
	VerificationTokenCreatedAt *time.Time `json:"verification_token_created_at,omitempty"`
	PendingEmail               string     `json:"pending_email,omitempty"`
	AccessToken                string     `json:"access_token,omitempty"`
	RefreshToken               string     `json:"refresh_token,omitempty"`
	TokenExpiry                time.Time  `json:"token_expiry,omitempty"`
	TokenCreatedAt             time.Time  `json:"token_created_at,omitempty"`
	Status                     string     `json:"status"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
	DeletedAt                  time.Time  `json:"deleted_at,omitempty"`
	IsDeleted                  bool       `json:"is_deleted,omitempty"`
}

// VerifyEmailResponse untuk respons verifikasi email
type VerifyEmailResponse struct {
	Email      string              `json:"email"`
	IsVerified bool                `json:"is_verified"`
	User       UserProfileResponse `json:"user,omitempty"`
}
