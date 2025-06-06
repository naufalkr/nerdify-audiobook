package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email           string    `gorm:"size:100;index:idx_users_email"`
	UserName        string    `gorm:"size:100;index:idx_users_user_name"`
	Password        string    `gorm:"size:255"`
	FullName        string    `gorm:"size:100"`
	Alamat          string    `gorm:"size:255"`
	Latitude        float64
	Longitude       float64
	ProfileImageURL string     `gorm:"size:255"`
	RoleID          *uuid.UUID `gorm:"type:uuid"`
	Role            Role       `gorm:"foreignKey:RoleID"`
	IsVerified      bool       `gorm:"default:false"`

	// Authentication token fields
	AuthToken      string    `gorm:"size:500;column:auth_token"`
	TokenExpiry    time.Time `gorm:"column:token_expiry"`
	RefreshToken   string    `gorm:"size:500;column:refresh_token"`
	TokenCreatedAt time.Time `gorm:"column:token_created_at"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// OTP-related fields
	OTPCode         string     `gorm:"column:otp_code"`
	OTPCreatedAt    *time.Time `gorm:"column:otp_created_at"`
	OTPAttemptCount int        `gorm:"column:otp_attempt_count"`

	// Resend OTP tracking
	ResendCount       int        `gorm:"column:resend_count;default:0"`
	LastResendAt      *time.Time `gorm:"column:last_resend_at"`
	CooldownStartedAt *time.Time `gorm:"column:cooldown_started_at"`

	// New fields for email link verification
	VerificationToken          string     `gorm:"column:verification_token"`
	VerificationTokenCreatedAt *time.Time `gorm:"column:verification_token_created_at"`

	// Field for email change
	PendingEmail string `gorm:"column:pending_email"`

	Status string `gorm:"size:20;default:active"` // active, suspended, inactive

	Tenants []UserTenant `gorm:"foreignKey:UserID"`
}
