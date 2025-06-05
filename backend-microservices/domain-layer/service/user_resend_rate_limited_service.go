package service

import (
	"context"
	"fmt"
	"microservice/user/helpers/dto"
	"microservice/user/helpers/utils"
	"time"
)

// ResendVerificationEmailWithRateLimiting adds rate limiting to the email verification resend process
func (s *UserService) ResendVerificationEmailWithRateLimiting(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	// Find user by email
	user, err := s.repo.FindUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Check if user is already verified
	if user.IsVerified {
		return dto.ErrUserAlreadyVerified
	}

	// Check if user is in cooldown period
	if user.CooldownStartedAt != nil && utils.IsInResendCooldown(*user.CooldownStartedAt) {
		// Return error with remaining cooldown time
		remainingMinutes := utils.GetRemainingCooldownMinutes(*user.CooldownStartedAt)
		return fmt.Errorf("%s, silakan coba lagi dalam %d menit", dto.ErrInCooldownPeriod.Error(), remainingMinutes)
	}

	// Check if resend count exceeded
	if user.ResendCount >= utils.MaxResendAttempts {
		// Start cooldown period
		now := time.Now()
		user.CooldownStartedAt = &now
		user.ResendCount = 0 // Reset counter

		_, err = s.repo.Update(ctx, nil, user)
		if err != nil {
			return dto.ErrUpdateUserFailed
		}

		return dto.ErrTooManyResendAttempts
	}

	// Generate new OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return dto.ErrSendEmailFailed
	}

	// Generate new verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return dto.ErrSendEmailFailed
	}

	// Update user with new OTP and verification token
	user.OTPCode = otp
	now := time.Now()
	user.OTPCreatedAt = &now
	user.OTPAttemptCount = 0
	user.VerificationToken = verificationToken
	user.VerificationTokenCreatedAt = &now

	// Increment resend count and update last resend timestamp
	user.ResendCount++
	user.LastResendAt = &now

	_, err = s.repo.Update(ctx, nil, user)
	if err != nil {
		return dto.ErrUpdateUserFailed
	}

	// Use production frontend URL for verification link
	verificationLink := utils.BuildVerificationLink("https://lecsens-iot.erplabiim.com", verificationToken, user.Email)

	// Send email with both OTP and verification link
	emailBody := utils.BuildOTPAndLinkVerificationEmail(user.FullName, otp, verificationLink)
	err = s.emailSender.Send(user.Email, "Verifikasi Email", emailBody)
	if err != nil {
		return dto.ErrSendEmailFailed
	}

	return nil
}
