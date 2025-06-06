package utils

import (
	"crypto/rand"
	"math/big"
	"time"
)

const (
	// OTP length
	OTPLength = 6

	// OTP expiration time in minutes
	OTPExpirationMinutes = 15

	// Maximum verification attempts
	MaxOTPAttempts = 3

	// Maximum resend OTP attempts
	MaxResendAttempts = 5

	// Cooldown time in minutes after max resend attempts
	ResendCooldownMinutes = 30

	// Allowed characters for OTP (excluding easily confused characters)
	OTPChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

// GenerateOTP creates a cryptographically secure OTP
func GenerateOTP() (string, error) {
	otp := make([]byte, OTPLength)

	for i := 0; i < OTPLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(OTPChars))))
		if err != nil {
			return "", err
		}
		otp[i] = OTPChars[num.Int64()]
	}

	return string(otp), nil
}

// IsOTPExpired checks if the OTP has expired
func IsOTPExpired(createdAt time.Time) bool {
	return time.Since(createdAt) > time.Minute*OTPExpirationMinutes
}

// IsInResendCooldown checks if user is in cooldown period after max resend attempts
func IsInResendCooldown(lastResendAt time.Time) bool {
	return time.Since(lastResendAt) < time.Duration(ResendCooldownMinutes)*time.Minute
}

// GetRemainingCooldownMinutes returns the remaining minutes in cooldown
func GetRemainingCooldownMinutes(lastResendAt time.Time) int {
	elapsed := time.Since(lastResendAt)
	cooldownDuration := time.Duration(ResendCooldownMinutes) * time.Minute

	if elapsed >= cooldownDuration {
		return 0
	}

	remaining := cooldownDuration - elapsed
	return int(remaining.Minutes()) + 1 // Round up to the next minute
}
