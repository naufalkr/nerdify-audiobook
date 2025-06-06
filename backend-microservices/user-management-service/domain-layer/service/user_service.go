package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"microservice/user/data-layer/entity"
	"microservice/user/data-layer/repository"
	"microservice/user/helpers/dto"
	"microservice/user/helpers/utils"
)

type UserService struct {
	repo        repository.UserRepository
	roleRepo    repository.RoleRepository
	tokenMaker  utils.TokenMaker
	emailSender utils.EmailSender
}

func NewUserService(
	repo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tokenMaker utils.TokenMaker,
	emailSender utils.EmailSender,
) *UserService {
	return &UserService{
		repo:        repo,
		roleRepo:    roleRepo,
		tokenMaker:  tokenMaker,
		emailSender: emailSender,
	}
}

// UpdateUserProfile updates a user's profile
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req dto.UserProfileUpdateRequest) error {
	// Get user with transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			s.repo.RollbackTx(ctx, tx)
		}
	}()

	// Get user with lock to prevent concurrent updates
	user, err := s.repo.FindUserByIdWithLock(ctx, tx, userID.String())
	if err != nil {
		return err
	}

	// Update basic info
	user.FullName = req.FullName
	user.Alamat = req.Alamat
	user.Latitude = req.Latitude
	user.Longitude = req.Longitude

	// Handle profile image upload if provided
	if req.ProfileImage != nil {
		// Delete old image if exists
		if user.ProfileImageURL != "" {
			oldPath := strings.TrimPrefix(user.ProfileImageURL, "/uploads/")
			utils.DeleteFile(oldPath)
		}

		// Upload new image
		filePath, err := utils.UploadFile(req.ProfileImage, "profiles")
		if err != nil {
			return fmt.Errorf("failed to upload profile image: %v", err)
		}

		user.ProfileImageURL = utils.GetFileURL(filePath)
	}

	// Update user
	_, err = s.repo.Update(ctx, tx, user)
	if err != nil {
		return err
	}

	// Commit transaction
	return s.repo.CommitTx(ctx, tx)
}

// ListUsers returns a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, req dto.UserListRequest) ([]dto.UserListResponse, int64, error) {
	// Get users with pagination
	users, total, err := s.repo.FindAllUser(ctx, (req.Page-1)*req.PageSize, req.PageSize, nil)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response DTO
	response := make([]dto.UserListResponse, len(users))
	for i, user := range users {
		// Handle DeletedAt field which is a gorm.DeletedAt
		var deletedAt time.Time
		isDeleted := false
		if user.DeletedAt.Valid {
			deletedAt = user.DeletedAt.Time
			isDeleted = true
		}

		response[i] = dto.UserListResponse{
			ID:                         user.ID.String(),
			Email:                      user.Email,
			UserName:                   user.UserName,
			FullName:                   user.FullName,
			Alamat:                     user.Alamat,
			Latitude:                   user.Latitude,
			Longitude:                  user.Longitude,
			ProfileImageURL:            user.ProfileImageURL,
			RoleID:                     user.RoleID.String(),
			Role:                       user.Role.Name,
			IsVerified:                 user.IsVerified,
			AccessToken:                user.AuthToken,
			TokenExpiry:                user.TokenExpiry,
			RefreshToken:               user.RefreshToken,
			TokenCreatedAt:             user.TokenCreatedAt,
			OTPCode:                    user.OTPCode,
			OTPCreatedAt:               user.OTPCreatedAt,
			OTPAttemptCount:            user.OTPAttemptCount,
			ResendCount:                user.ResendCount,
			LastResendAt:               user.LastResendAt,
			CooldownStartedAt:          user.CooldownStartedAt,
			VerificationToken:          user.VerificationToken,
			VerificationTokenCreatedAt: user.VerificationTokenCreatedAt,
			PendingEmail:               user.PendingEmail,
			Status:                     user.Status,
			CreatedAt:                  user.CreatedAt,
			UpdatedAt:                  user.UpdatedAt,
			DeletedAt:                  deletedAt,
			IsDeleted:                  isDeleted,
		}
	}

	return response, total, nil
}

// Register creates a new user, assigns USER role, and sends verification email with OTP and link
func (s *UserService) Register(ctx context.Context, req dto.UserRegisterRequest) (dto.UserProfileResponse, error) {
	if !utils.IsValidEmail(req.Email) {
		return dto.UserProfileResponse{}, dto.ErrInvalidEmail
	}

	exists, _ := s.repo.ExistsByEmailOrUsername(ctx, nil, req.Email, req.UserName)
	if exists {
		return dto.UserProfileResponse{}, dto.ErrUserAlreadyExist
	}

	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.UserProfileResponse{}, err
	}

	role, err := s.roleRepo.FindByName(ctx, "USER")
	if err != nil {
		return dto.UserProfileResponse{}, err
	}

	// Initialize with current time to avoid nil pointers
	now := time.Now()

	user := entity.User{
		ID:                         uuid.New(),
		Email:                      req.Email,
		UserName:                   req.UserName,
		Password:                   hashedPwd,
		FullName:                   req.FullName,
		Alamat:                     req.Alamat,
		Latitude:                   req.Latitude,
		Longitude:                  req.Longitude,
		ProfileImageURL:            "",       // Set to empty string as ProfileImageURL doesn't exist in dto.UserRegisterRequest
		RoleID:                     &role.ID, // Convert uuid.UUID to *uuid.UUID
		IsVerified:                 false,
		OTPCreatedAt:               &now, // Initialize to avoid nil pointer dereference
		VerificationTokenCreatedAt: &now, // Initialize to avoid nil pointer dereference
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
		Status:                     "pending", // Set default status
	}

	createdUser, err := s.repo.CreateUser(ctx, nil, user)
	if err != nil {
		return dto.UserProfileResponse{}, dto.ErrCreateUserFailed
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		// Continue with user creation even if OTP generation fails
		return dto.UserProfileResponse{
			ID:              createdUser.ID.String(),
			Username:        createdUser.UserName,
			Email:           createdUser.Email,
			FullName:        createdUser.FullName,
			Alamat:          createdUser.Alamat,
			Latitude:        createdUser.Latitude,
			Longitude:       createdUser.Longitude,
			ProfileImageURL: createdUser.ProfileImageURL,
			RoleID:          createdUser.RoleID.String(),
			Role:            role.Name,
			IsVerified:      createdUser.IsVerified,
			Status:          createdUser.Status,
			CreatedAt:       createdUser.CreatedAt,
			UpdatedAt:       createdUser.UpdatedAt,
		}, nil
	}

	// Generate verification token for the link method
	verificationToken, err := utils.GenerateVerificationToken()
	// if err != nil {
	// 	log.Printf("Failed to generate verification token: %v", err)
	// 	// Continue with user creation even if token generation fails
	// 	return dto.UserProfileResponse{
	// 		ID:              createdUser.ID.String(),
	// 		Username:        createdUser.UserName,
	// 		Email:           createdUser.Email,
	// 		FullName:        createdUser.FullName,
	// 		Alamat:          createdUser.Alamat,
	// 		Latitude:        createdUser.Latitude,
	// 		Longitude:       createdUser.Longitude,
	// 		ProfileImageURL: createdUser.ProfileImageURL,
	// 		RoleID:          createdUser.RoleID.String(),
	// 		Role:            role.Name,
	// 		IsVerified:      createdUser.IsVerified,
	// 		Status:          createdUser.Status,
	// 		CreatedAt:       createdUser.CreatedAt,
	// 		UpdatedAt:       createdUser.UpdatedAt,
	// 	}, nil
	// }

	// Save OTP and verification token to user record
	createdUser.OTPCode = otp
	// now variable already defined above
	createdUser.OTPCreatedAt = &now
	createdUser.OTPAttemptCount = 0
	createdUser.VerificationToken = verificationToken
	createdUser.VerificationTokenCreatedAt = &now

	// Try to update user with OTP info, but don't fail registration if this fails
	updatedUser, updateErr := s.repo.Update(ctx, nil, createdUser)
	if updateErr != nil {
		log.Printf("Failed to update user with OTP information: %v", updateErr)
		// Continue anyway with the original user
		updatedUser = createdUser
	}

	// Use production frontend URL for verification link
	verificationLink := utils.BuildVerificationLink("https://lecsens-iot.erplabiim.com", verificationToken, createdUser.Email)

	// Try to send email, but don't fail registration if email sending fails
	emailBody := utils.BuildOTPAndLinkVerificationEmail(createdUser.FullName, otp, verificationLink)
	emailErr := s.emailSender.Send(createdUser.Email, "Verifikasi Email", emailBody)
	if emailErr != nil {
		log.Printf("Failed to send verification email: %v", emailErr)
		// Return success with a warning that email wasn't sent
		return dto.UserProfileResponse{
			ID:                         updatedUser.ID.String(),
			Username:                   updatedUser.UserName,
			Email:                      updatedUser.Email,
			FullName:                   updatedUser.FullName,
			Alamat:                     updatedUser.Alamat,
			Latitude:                   updatedUser.Latitude,
			Longitude:                  updatedUser.Longitude,
			ProfileImageURL:            updatedUser.ProfileImageURL,
			RoleID:                     updatedUser.RoleID.String(),
			Role:                       role.Name,
			IsVerified:                 updatedUser.IsVerified,
			OTPCode:                    updatedUser.OTPCode,
			OTPCreatedAt:               updatedUser.OTPCreatedAt,
			OTPAttemptCount:            updatedUser.OTPAttemptCount,
			VerificationToken:          updatedUser.VerificationToken,
			VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
			Status:                     updatedUser.Status,
			CreatedAt:                  updatedUser.CreatedAt,
			UpdatedAt:                  updatedUser.UpdatedAt,
		}, nil
	}

	return dto.UserProfileResponse{
		ID:                         updatedUser.ID.String(),
		Username:                   updatedUser.UserName,
		Email:                      updatedUser.Email,
		FullName:                   updatedUser.FullName,
		Alamat:                     updatedUser.Alamat,
		Latitude:                   updatedUser.Latitude,
		Longitude:                  updatedUser.Longitude,
		ProfileImageURL:            updatedUser.ProfileImageURL,
		RoleID:                     updatedUser.RoleID.String(),
		Role:                       role.Name,
		IsVerified:                 updatedUser.IsVerified,
		OTPCode:                    updatedUser.OTPCode,
		OTPCreatedAt:               updatedUser.OTPCreatedAt,
		OTPAttemptCount:            updatedUser.OTPAttemptCount,
		VerificationToken:          updatedUser.VerificationToken,
		VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
		Status:                     updatedUser.Status,
		CreatedAt:                  updatedUser.CreatedAt,
		UpdatedAt:                  updatedUser.UpdatedAt,
	}, nil
}

// Login verifies credentials and returns access token
func (s *UserService) Login(ctx context.Context, req dto.UserLoginRequest) (dto.UserProfileResponse, error) {
	var user entity.User
	var err error

	// Check if login is via email or username
	if req.Email != "" {
		// Validate email format
		if !utils.IsValidEmail(req.Email) {
			return dto.UserProfileResponse{}, dto.ErrInvalidEmail
		}

		// Find user by email
		user, err = s.repo.FindUserByEmail(ctx, nil, req.Email)
		if err != nil {
			log.Printf("Login failed: User not found with email %s: %v", req.Email, err)
			return dto.UserProfileResponse{}, dto.ErrUserNotFound
		}
	} else if req.Username != "" {
		// Find user by username
		user, err = s.repo.FindUserByUsername(ctx, nil, req.Username)
		if err != nil {
			log.Printf("Login failed: User not found with username %s: %v", req.Username, err)
			return dto.UserProfileResponse{}, dto.ErrUserNotFound
		}
	} else {
		return dto.UserProfileResponse{}, fmt.Errorf("username or email is required")
	}

	// Check if password matches
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		log.Printf("Login failed: Password mismatch for user %s", req.Email)
		return dto.UserProfileResponse{}, dto.ErrPasswordNotMatch
	}

	// Check if user is verified
	if !user.IsVerified {
		log.Printf("Login failed: User %s not verified", req.Email)
		return dto.UserProfileResponse{
			Username:   user.UserName,
			Email:      user.Email,
			IsVerified: false,
		}, dto.ErrUserNotVerified
	}

	// Create access token with role name
	now := time.Now()
	log.Printf("Creating access token for user %s at %v", req.Email, now)

	accessToken, err := s.tokenMaker.CreateAccessToken(user.ID.String(), user.RoleID.String(), user.Role.Name, 24*time.Hour)
	if err != nil {
		log.Printf("Login failed: Error creating access token for user %s: %v", req.Email, err)
		return dto.UserProfileResponse{}, dto.ErrLoginFailed
	}
	log.Printf("Access token created successfully for user %s", req.Email)

	// Create refresh token
	refreshToken, err := s.tokenMaker.CreateRefreshToken(user.ID.String(), 7*24*time.Hour)
	if err != nil {
		log.Printf("Login failed: Error creating refresh token for user %s: %v", req.Email, err)
		return dto.UserProfileResponse{}, dto.ErrLoginFailed
	}
	log.Printf("Refresh token created successfully for user %s", req.Email)

	// Start transaction to update user with tokens
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("Login failed: Error starting transaction: %v", err)
		return dto.UserProfileResponse{}, dto.ErrLoginFailed
	}

	// Update user with tokens
	user.AuthToken = accessToken
	user.TokenExpiry = now.Add(24 * time.Hour)
	user.RefreshToken = refreshToken
	user.TokenCreatedAt = now
	user.UpdatedAt = now

	log.Printf("Updating user %s with tokens - TokenExpiry: %v, TokenCreatedAt: %v",
		req.Email, user.TokenExpiry, user.TokenCreatedAt)

	// Save updated user with tokens
	updatedUser, err := s.repo.Update(ctx, tx, user)
	if err != nil {
		s.repo.RollbackTx(ctx, tx)
		log.Printf("Login failed: Error saving tokens for user %s: %v", req.Email, err)
		return dto.UserProfileResponse{}, dto.ErrLoginFailed
	}

	// Commit transaction
	if err := s.repo.CommitTx(ctx, tx); err != nil {
		s.repo.RollbackTx(ctx, tx)
		log.Printf("Login failed: Error committing transaction: %v", err)
		return dto.UserProfileResponse{}, dto.ErrLoginFailed
	}

	log.Printf("Login successful for user %s, token expires at %v", req.Email, user.TokenExpiry)

	// Handle DeletedAt field which is a gorm.DeletedAt
	var deletedAt time.Time
	isDeleted := false
	if updatedUser.DeletedAt.Valid {
		deletedAt = updatedUser.DeletedAt.Time
		isDeleted = true
	}

	return dto.UserProfileResponse{
		ID:                         updatedUser.ID.String(),
		Username:                   updatedUser.UserName,
		Email:                      updatedUser.Email,
		FullName:                   updatedUser.FullName,
		Alamat:                     updatedUser.Alamat,
		Latitude:                   updatedUser.Latitude,
		Longitude:                  updatedUser.Longitude,
		ProfileImageURL:            updatedUser.ProfileImageURL,
		RoleID:                     updatedUser.RoleID.String(),
		Role:                       updatedUser.Role.Name,
		IsVerified:                 updatedUser.IsVerified,
		AccessToken:                updatedUser.AuthToken,
		TokenExpiry:                updatedUser.TokenExpiry,
		RefreshToken:               updatedUser.RefreshToken,
		TokenCreatedAt:             updatedUser.TokenCreatedAt,
		OTPCode:                    updatedUser.OTPCode,
		OTPCreatedAt:               updatedUser.OTPCreatedAt,
		OTPAttemptCount:            updatedUser.OTPAttemptCount,
		ResendCount:                updatedUser.ResendCount,
		LastResendAt:               updatedUser.LastResendAt,
		CooldownStartedAt:          updatedUser.CooldownStartedAt,
		VerificationToken:          updatedUser.VerificationToken,
		VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
		PendingEmail:               updatedUser.PendingEmail,
		Status:                     updatedUser.Status,
		CreatedAt:                  updatedUser.CreatedAt,
		UpdatedAt:                  updatedUser.UpdatedAt,
		DeletedAt:                  deletedAt,
		IsDeleted:                  isDeleted,
	}, nil
}

// VerifyEmail verifies user's email using OTP
func (s *UserService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		// Get role information
		role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
		if err != nil {
			log.Printf("Role not found for user %s: %v", user.ID.String(), err)
			return dto.VerifyEmailResponse{}, dto.ErrRoleNotFound
		}

		// Handle DeletedAt field which is a gorm.DeletedAt
		var deletedAt time.Time
		isDeleted := false
		if user.DeletedAt.Valid {
			deletedAt = user.DeletedAt.Time
			isDeleted = true
		}

		return dto.VerifyEmailResponse{
			Email:      user.Email,
			IsVerified: true,
			User: dto.UserProfileResponse{
				ID:                         user.ID.String(),
				Username:                   user.UserName,
				Email:                      user.Email,
				FullName:                   user.FullName,
				Alamat:                     user.Alamat,
				Latitude:                   user.Latitude,
				Longitude:                  user.Longitude,
				ProfileImageURL:            user.ProfileImageURL,
				RoleID:                     user.RoleID.String(),
				Role:                       role.Name,
				IsVerified:                 user.IsVerified,
				AccessToken:                user.AuthToken,
				TokenExpiry:                user.TokenExpiry,
				RefreshToken:               user.RefreshToken,
				TokenCreatedAt:             user.TokenCreatedAt,
				OTPCode:                    user.OTPCode,
				OTPCreatedAt:               user.OTPCreatedAt,
				OTPAttemptCount:            user.OTPAttemptCount,
				ResendCount:                user.ResendCount,
				LastResendAt:               user.LastResendAt,
				CooldownStartedAt:          user.CooldownStartedAt,
				VerificationToken:          user.VerificationToken,
				VerificationTokenCreatedAt: user.VerificationTokenCreatedAt,
				PendingEmail:               user.PendingEmail,
				Status:                     user.Status,
				CreatedAt:                  user.CreatedAt,
				UpdatedAt:                  user.UpdatedAt,
				DeletedAt:                  deletedAt,
				IsDeleted:                  isDeleted,
			},
		}, nil
	}

	// Check attempt count to prevent brute force
	if user.OTPAttemptCount >= utils.MaxOTPAttempts {
		return dto.VerifyEmailResponse{}, dto.ErrTooManyAttempts
	}

	// Update attempt count
	user.OTPAttemptCount++

	// Check if OTP is expired
	if user.OTPCreatedAt == nil || utils.IsOTPExpired(*user.OTPCreatedAt) {
		_, updateErr := s.repo.Update(ctx, nil, user) // Save updated attempt count
		if updateErr != nil {
			log.Printf("Failed to update attempt count: %v", updateErr)
		}
		return dto.VerifyEmailResponse{}, dto.ErrOTPExpired
	}

	// Check if OTP matches
	if user.OTPCode != req.OTP {
		_, updateErr := s.repo.Update(ctx, nil, user) // Save updated attempt count
		if updateErr != nil {
			log.Printf("Failed to update attempt count: %v", updateErr)
		}
		return dto.VerifyEmailResponse{}, dto.ErrOTPNotMatch
	}

	// OTP is valid, verify user
	user.IsVerified = true
	user.OTPCode = "" // Clear OTP after successful verification
	user.UpdatedAt = time.Now()

	updatedUser, err := s.repo.Update(ctx, nil, user)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUserFailed
	}

	// Get role information
	role, err := s.roleRepo.FindByID(ctx, *updatedUser.RoleID)
	if err != nil {
		log.Printf("Role not found for user %s: %v", updatedUser.ID.String(), err)
		return dto.VerifyEmailResponse{}, dto.ErrRoleNotFound
	}

	// Handle DeletedAt field which is a gorm.DeletedAt
	var deletedAt time.Time
	isDeleted := false
	if updatedUser.DeletedAt.Valid {
		deletedAt = updatedUser.DeletedAt.Time
		isDeleted = true
	}

	return dto.VerifyEmailResponse{
		Email:      updatedUser.Email,
		IsVerified: updatedUser.IsVerified,
		User: dto.UserProfileResponse{
			ID:                         updatedUser.ID.String(),
			Username:                   updatedUser.UserName,
			Email:                      updatedUser.Email,
			FullName:                   updatedUser.FullName,
			Alamat:                     updatedUser.Alamat,
			Latitude:                   updatedUser.Latitude,
			Longitude:                  updatedUser.Longitude,
			ProfileImageURL:            updatedUser.ProfileImageURL,
			RoleID:                     updatedUser.RoleID.String(),
			Role:                       role.Name,
			IsVerified:                 updatedUser.IsVerified,
			AccessToken:                updatedUser.AuthToken,
			TokenExpiry:                updatedUser.TokenExpiry,
			RefreshToken:               updatedUser.RefreshToken,
			TokenCreatedAt:             updatedUser.TokenCreatedAt,
			OTPCode:                    updatedUser.OTPCode,
			OTPCreatedAt:               updatedUser.OTPCreatedAt,
			OTPAttemptCount:            updatedUser.OTPAttemptCount,
			ResendCount:                updatedUser.ResendCount,
			LastResendAt:               updatedUser.LastResendAt,
			CooldownStartedAt:          updatedUser.CooldownStartedAt,
			VerificationToken:          updatedUser.VerificationToken,
			VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
			PendingEmail:               updatedUser.PendingEmail,
			Status:                     updatedUser.Status,
			CreatedAt:                  updatedUser.CreatedAt,
			UpdatedAt:                  updatedUser.UpdatedAt,
			DeletedAt:                  deletedAt,
			IsDeleted:                  isDeleted,
		},
	}, nil
}

// VerifyEmailByLink verifies user's email using the verification link
func (s *UserService) VerifyEmailByLink(ctx context.Context, req dto.VerifyEmailLinkRequest) (dto.VerifyEmailResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		// Get role information
		role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
		if err != nil {
			log.Printf("Role not found for user %s: %v", user.ID.String(), err)
			return dto.VerifyEmailResponse{}, dto.ErrRoleNotFound
		}

		// Handle DeletedAt field which is a gorm.DeletedAt
		var deletedAt time.Time
		isDeleted := false
		if user.DeletedAt.Valid {
			deletedAt = user.DeletedAt.Time
			isDeleted = true
		}

		return dto.VerifyEmailResponse{
			Email:      user.Email,
			IsVerified: true,
			User: dto.UserProfileResponse{
				ID:                         user.ID.String(),
				Username:                   user.UserName,
				Email:                      user.Email,
				FullName:                   user.FullName,
				Alamat:                     user.Alamat,
				Latitude:                   user.Latitude,
				Longitude:                  user.Longitude,
				ProfileImageURL:            user.ProfileImageURL,
				RoleID:                     user.RoleID.String(),
				Role:                       role.Name,
				IsVerified:                 user.IsVerified,
				AccessToken:                user.AuthToken,
				TokenExpiry:                user.TokenExpiry,
				RefreshToken:               user.RefreshToken,
				TokenCreatedAt:             user.TokenCreatedAt,
				OTPCode:                    user.OTPCode,
				OTPCreatedAt:               user.OTPCreatedAt,
				OTPAttemptCount:            user.OTPAttemptCount,
				ResendCount:                user.ResendCount,
				LastResendAt:               user.LastResendAt,
				CooldownStartedAt:          user.CooldownStartedAt,
				VerificationToken:          user.VerificationToken,
				VerificationTokenCreatedAt: user.VerificationTokenCreatedAt,
				PendingEmail:               user.PendingEmail,
				Status:                     user.Status,
				CreatedAt:                  user.CreatedAt,
				UpdatedAt:                  user.UpdatedAt,
				DeletedAt:                  deletedAt,
				IsDeleted:                  isDeleted,
			},
		}, nil
	}

	// Check if token is valid
	if user.VerificationToken != req.Token {
		return dto.VerifyEmailResponse{}, dto.ErrVerificationTokenInvalid
	}

	// Check if token is expired (24 hours)
	if user.VerificationTokenCreatedAt == nil {
		return dto.VerifyEmailResponse{}, dto.ErrVerificationTokenExpired
	}
	tokenAge := time.Since(*user.VerificationTokenCreatedAt)
	if tokenAge > utils.VerificationTokenExpiry {
		return dto.VerifyEmailResponse{}, dto.ErrVerificationTokenExpired
	}

	// Token is valid, verify user
	user.IsVerified = true
	user.VerificationToken = "" // Clear token after successful verification
	user.UpdatedAt = time.Now()

	updatedUser, err := s.repo.Update(ctx, nil, user)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUserFailed
	}

	// Get role information
	role, err := s.roleRepo.FindByID(ctx, *updatedUser.RoleID)
	if err != nil {
		log.Printf("Role not found for user %s: %v", updatedUser.ID.String(), err)
		return dto.VerifyEmailResponse{}, dto.ErrRoleNotFound
	}

	// Handle DeletedAt field which is a gorm.DeletedAt
	var deletedAt time.Time
	isDeleted := false
	if updatedUser.DeletedAt.Valid {
		deletedAt = updatedUser.DeletedAt.Time
		isDeleted = true
	}

	return dto.VerifyEmailResponse{
		Email:      updatedUser.Email,
		IsVerified: updatedUser.IsVerified,
		User: dto.UserProfileResponse{
			ID:                         updatedUser.ID.String(),
			Username:                   updatedUser.UserName,
			Email:                      updatedUser.Email,
			FullName:                   updatedUser.FullName,
			Alamat:                     updatedUser.Alamat,
			Latitude:                   updatedUser.Latitude,
			Longitude:                  updatedUser.Longitude,
			ProfileImageURL:            updatedUser.ProfileImageURL,
			RoleID:                     updatedUser.RoleID.String(),
			Role:                       role.Name,
			IsVerified:                 updatedUser.IsVerified,
			AccessToken:                updatedUser.AuthToken,
			TokenExpiry:                updatedUser.TokenExpiry,
			RefreshToken:               updatedUser.RefreshToken,
			TokenCreatedAt:             updatedUser.TokenCreatedAt,
			OTPCode:                    updatedUser.OTPCode,
			OTPCreatedAt:               updatedUser.OTPCreatedAt,
			OTPAttemptCount:            updatedUser.OTPAttemptCount,
			ResendCount:                updatedUser.ResendCount,
			LastResendAt:               updatedUser.LastResendAt,
			CooldownStartedAt:          updatedUser.CooldownStartedAt,
			VerificationToken:          updatedUser.VerificationToken,
			VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
			PendingEmail:               updatedUser.PendingEmail,
			Status:                     updatedUser.Status,
			CreatedAt:                  updatedUser.CreatedAt,
			UpdatedAt:                  updatedUser.UpdatedAt,
			DeletedAt:                  deletedAt,
			IsDeleted:                  isDeleted,
		},
	}, nil
}

// VerifyEmailById verifies a user's email by the admin directly using user ID
func (s *UserService) VerifyEmailById(ctx context.Context, userID string) (dto.VerifyEmailResponse, error) {
	// Parse the user ID to validate it's a valid UUID
	_, err := uuid.Parse(userID)
	if err != nil {
		return dto.VerifyEmailResponse{}, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get the user by ID
	user, err := s.repo.FindUserById(ctx, nil, userID)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	// Check if the user is already verified
	if user.IsVerified {
		return dto.VerifyEmailResponse{
			Email:      user.Email,
			IsVerified: true,
		}, nil
	}

	// Mark the email as verified
	user.IsVerified = true
	user.Status = "active" // Use the string constant directly as model.UserStatusActive is not defined
	user.UpdatedAt = time.Now()

	// Update the user in the database
	updatedUser, err := s.repo.Update(ctx, nil, user)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUserFailed
	}

	// Get role information
	role, err := s.roleRepo.FindByID(ctx, *updatedUser.RoleID)
	if err != nil {
		log.Printf("Role not found for user %s: %v", userID, err)
		return dto.VerifyEmailResponse{}, dto.ErrRoleNotFound
	}

	// Handle DeletedAt field which is a gorm.DeletedAt
	var deletedAt time.Time
	isDeleted := false
	if updatedUser.DeletedAt.Valid {
		deletedAt = updatedUser.DeletedAt.Time
		isDeleted = true
	}

	// Create a full response with all user fields
	return dto.VerifyEmailResponse{
		Email:      updatedUser.Email,
		IsVerified: updatedUser.IsVerified,
		User: dto.UserProfileResponse{
			ID:                         updatedUser.ID.String(),
			Username:                   updatedUser.UserName,
			Email:                      updatedUser.Email,
			FullName:                   updatedUser.FullName,
			Alamat:                     updatedUser.Alamat,
			Latitude:                   updatedUser.Latitude,
			Longitude:                  updatedUser.Longitude,
			ProfileImageURL:            updatedUser.ProfileImageURL,
			RoleID:                     updatedUser.RoleID.String(),
			Role:                       role.Name,
			IsVerified:                 updatedUser.IsVerified,
			AccessToken:                updatedUser.AuthToken,
			TokenExpiry:                updatedUser.TokenExpiry,
			RefreshToken:               updatedUser.RefreshToken,
			TokenCreatedAt:             updatedUser.TokenCreatedAt,
			OTPCode:                    updatedUser.OTPCode,
			OTPCreatedAt:               updatedUser.OTPCreatedAt,
			OTPAttemptCount:            updatedUser.OTPAttemptCount,
			ResendCount:                updatedUser.ResendCount,
			LastResendAt:               updatedUser.LastResendAt,
			CooldownStartedAt:          updatedUser.CooldownStartedAt,
			VerificationToken:          updatedUser.VerificationToken,
			VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
			PendingEmail:               updatedUser.PendingEmail,
			Status:                     updatedUser.Status,
			CreatedAt:                  updatedUser.CreatedAt,
			UpdatedAt:                  updatedUser.UpdatedAt,
			DeletedAt:                  deletedAt,
			IsDeleted:                  isDeleted,
		},
	}, nil
}

// DeleteUserByID performs a soft delete on a user by ID
func (s *UserService) DeleteUserByID(ctx context.Context, userID string) error {
	log.Printf("DeleteUserByID: Starting soft delete for user ID: %s", userID)

	// Start a transaction to ensure data consistency
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("DeleteUserByID: Error beginning transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure proper cleanup of the transaction
	defer func() {
		if r := recover(); r != nil {
			s.repo.RollbackTx(ctx, tx)
			log.Printf("DeleteUserByID: Panic recovered: %v", r)
		}
	}()

	// Perform the soft delete
	err = s.repo.SoftDeleteUserByID(ctx, userID, tx)
	if err != nil {
		s.repo.RollbackTx(ctx, tx)
		log.Printf("DeleteUserByID: Error during soft delete: %v", err)
		return err
	}

	// Commit the transaction
	if err := s.repo.CommitTx(ctx, tx); err != nil {
		s.repo.RollbackTx(ctx, tx)
		log.Printf("DeleteUserByID: Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("DeleteUserByID: Successfully completed soft delete for user ID: %s", userID)
	return nil
}

// HardDeleteUserByID performs a permanent delete on a user by ID
func (s *UserService) HardDeleteUserByID(ctx context.Context, userID string) error {
	log.Printf("HardDeleteUserByID: Starting hard delete for user ID: %s", userID)

	// Start a transaction to ensure data consistency
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("HardDeleteUserByID: Error beginning transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure proper cleanup of the transaction
	defer func() {
		if r := recover(); r != nil {
			s.repo.RollbackTx(ctx, tx)
			log.Printf("HardDeleteUserByID: Panic recovered: %v", r)
		}
	}()

	// Perform the hard delete
	err = s.repo.HardDeleteUserByID(ctx, userID, tx)
	if err != nil {
		s.repo.RollbackTx(ctx, tx)
		log.Printf("HardDeleteUserByID: Error during hard delete: %v", err)
		return err
	}

	// Commit the transaction
	if err := s.repo.CommitTx(ctx, tx); err != nil {
		s.repo.RollbackTx(ctx, tx)
		log.Printf("HardDeleteUserByID: Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("HardDeleteUserByID: Successfully completed hard delete for user ID: %s", userID)
	return nil
}

// ResendVerificationEmail with enhanced functionality for both OTP and link
func (s *UserService) ResendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	user, err := s.repo.FindUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrUserNotFound
	}

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
	return s.emailSender.Send(user.Email, "Verifikasi Email", emailBody)
}

// ForgotPassword sends an email with password reset token
func (s *UserService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	user, err := s.repo.FindUserByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.ErrUserNotFound
	}

	token, err := s.tokenMaker.CreatePasswordResetToken(user.Email, 1*time.Hour)
	if err != nil {
		return dto.ErrSendEmailFailed
	}

	emailBody := utils.BuildPasswordResetEmail(user.FullName, token)
	return s.emailSender.Send(user.Email, "Reset Password", emailBody)
}

// ResetPassword changes password based on reset token
func (s *UserService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	email, err := s.tokenMaker.ParsePasswordResetToken(req.Token)
	if err != nil {
		return dto.ErrTokenInvalid
	}

	user, err := s.repo.FindUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Check if the new password is the same as the current password
	if utils.CheckPasswordHash(req.NewPassword, user.Password) {
		return dto.ErrSamePassword
	}

	hashedPwd, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPwd
	user.UpdatedAt = time.Now()

	_, err = s.repo.Update(ctx, nil, user)
	if err != nil {
		return dto.ErrUpdateUserFailed
	}

	return nil
}

// RefreshToken extends login session with new token
func (s *UserService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error) {
	// Validate refresh token
	claims, err := s.tokenMaker.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return dto.TokenResponse{}, dto.ErrTokenInvalid
	}

	user, err := s.repo.FindUserById(ctx, nil, claims.UserID)
	if err != nil {
		return dto.TokenResponse{}, dto.ErrUserNotFound
	}

	// Create new tokens with role name
	accessToken, err := s.tokenMaker.CreateAccessToken(user.ID.String(), user.RoleID.String(), user.Role.Name, 24*time.Hour)
	if err != nil {
		return dto.TokenResponse{}, dto.ErrLoginFailed
	}

	refreshToken, err := s.tokenMaker.CreateRefreshToken(user.ID.String(), 7*24*time.Hour)
	if err != nil {
		return dto.TokenResponse{}, dto.ErrLoginFailed
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout marks user token as invalid
func (s *UserService) Logout(ctx context.Context, userID string) error {
	// Parse userID to ensure it's valid
	_, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format in logout attempt: %s", userID)
		return dto.ErrUserNotFound
	}

	// Check if user exists
	user, err := s.repo.FindUserById(ctx, nil, userID)
	if err != nil {
		log.Printf("User not found during logout: %s", userID)
		return dto.ErrUserNotFound
	}

	// Check if the BlacklistUserTokens method is implemented
	// This is a type assertion to check if the method exists at runtime
	blacklister, ok := s.tokenMaker.(interface {
		BlacklistUserTokens(context.Context, string) error
	})

	if ok {
		// If the method is implemented, use it
		err = blacklister.BlacklistUserTokens(ctx, userID)
		if err != nil {
			log.Printf("Failed to blacklist tokens for user %s: %v", userID, err)
			return fmt.Errorf("logout failed: %w", err)
		}
	} else {
		// If not implemented, log a warning but don't fail the logout
		log.Printf("WARNING: TokenMaker does not implement BlacklistUserTokens - tokens will remain valid")
	}

	log.Printf("Successfully logged out user: %s (%s)", userID, user.Email)
	return nil
}

// ChangeUserRole changes a user's role (only SuperAdmin can execute)
func (s *UserService) ChangeUserRole(ctx context.Context, userID string, roleName string, changedBy string) (dto.UserResponse, error) {
	log.Printf("DEBUG: Starting role change process - UserID: %s, New Role: %s, Changed By: %s", userID, roleName, changedBy)

	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("DEBUG: Failed to begin transaction: %v", err)
		return dto.UserResponse{}, dto.ErrUpdateUserFailed
	}

	// Get user with row lock
	userToUpdate, err := s.repo.FindUserByIdWithLock(ctx, tx, userID)
	if err != nil {
		log.Printf("DEBUG: Failed to find user with lock: %v", err)
		s.repo.RollbackTx(ctx, tx)
		return dto.UserResponse{}, dto.ErrUserNotFound
	}
	log.Printf("DEBUG: Found user to update - ID: %s, Current RoleID: %s", userToUpdate.ID, userToUpdate.RoleID)

	// Find new role
	newRole, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		log.Printf("DEBUG: Error finding role: %v", err)
		s.repo.RollbackTx(ctx, tx)
		return dto.UserResponse{}, dto.ErrRoleNotFound
	}
	log.Printf("DEBUG: Found new role - ID: %s, Name: %s", newRole.ID, newRole.Name)

	// Log current role before change
	oldRoleID := *userToUpdate.RoleID
	oldRole, _ := s.roleRepo.FindByID(ctx, oldRoleID)
	if oldRole != nil {
		log.Printf("DEBUG: Current role before update - ID: %s, Name: %s", oldRoleID, oldRole.Name)
	}

	// Update the role using direct SQL
	err = s.repo.UpdateUserRoleDirectSQL(ctx, tx, userID, newRole.ID.String())
	if err != nil {
		log.Printf("DEBUG: Failed to update user role: %v", err)
		s.repo.RollbackTx(ctx, tx)
		return dto.UserResponse{}, dto.ErrUpdateUserFailed
	}
	log.Printf("DEBUG: Successfully updated user role in transaction")

	// Commit the transaction
	if err := s.repo.CommitTx(ctx, tx); err != nil {
		log.Printf("DEBUG: Failed to commit transaction: %v", err)
		return dto.UserResponse{}, dto.ErrUpdateUserFailed
	}
	log.Printf("DEBUG: Successfully committed transaction")

	// Verify the update after commit
	verifyUser, err := s.repo.FindUserById(ctx, nil, userID)
	if err != nil {
		log.Printf("DEBUG: Error verifying user after update: %v", err)
		return dto.UserResponse{}, dto.ErrUpdateUserFailed
	}
	log.Printf("DEBUG: Verification after commit - User ID: %s, RoleID: %s", verifyUser.ID, verifyUser.RoleID)

	if verifyUser.RoleID.String() != newRole.ID.String() {
		log.Printf("DEBUG: ⚠️ Role verification failed! Expected %s but found %s", newRole.ID, verifyUser.RoleID)
		return dto.UserResponse{}, dto.ErrUpdateUserFailed
	}
	log.Printf("DEBUG: ✓ Role verification successful")

	// Create audit log
	go func() {
		auditCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if auditRepo, ok := s.repo.(interface {
			LogAuditEvent(ctx context.Context, event string, userID string, details string) error
		}); ok {
			details := fmt.Sprintf("User role changed from %s to %s by %s", oldRole.Name, newRole.Name, changedBy)
			_ = auditRepo.LogAuditEvent(auditCtx, "ROLE_CHANGE", userID, details)
		}
	}()

	log.Printf("DEBUG: Role change process completed successfully")
	return dto.UserResponse{
		ID:         verifyUser.ID.String(),
		Username:   verifyUser.UserName,
		Email:      verifyUser.Email,
		UserRole:   newRole.Name,
		IsVerified: verifyUser.IsVerified,
	}, nil
}

// GetAllUsers gets all users with pagination (for SuperAdmin)
func (s *UserService) GetAllUsers(ctx context.Context, page, limit int) ([]dto.UserListResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, totalCount, err := s.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	response := make([]dto.UserListResponse, len(users))
	for i, user := range users {
		// Handle DeletedAt field which is a gorm.DeletedAt
		var deletedAt time.Time
		isDeleted := false
		if user.DeletedAt.Valid {
			deletedAt = user.DeletedAt.Time
			isDeleted = true
		}

		response[i] = dto.UserListResponse{
			ID:                         user.ID.String(),
			Email:                      user.Email,
			UserName:                   user.UserName,
			FullName:                   user.FullName,
			Alamat:                     user.Alamat,
			Latitude:                   user.Latitude,
			Longitude:                  user.Longitude,
			ProfileImageURL:            user.ProfileImageURL,
			RoleID:                     user.RoleID.String(),
			Role:                       user.Role.Name,
			IsVerified:                 user.IsVerified,
			AccessToken:                user.AuthToken,
			TokenExpiry:                user.TokenExpiry,
			RefreshToken:               user.RefreshToken,
			TokenCreatedAt:             user.TokenCreatedAt,
			OTPCode:                    user.OTPCode,
			OTPCreatedAt:               user.OTPCreatedAt,
			OTPAttemptCount:            user.OTPAttemptCount,
			ResendCount:                user.ResendCount,
			LastResendAt:               user.LastResendAt,
			CooldownStartedAt:          user.CooldownStartedAt,
			VerificationToken:          user.VerificationToken,
			VerificationTokenCreatedAt: user.VerificationTokenCreatedAt,
			PendingEmail:               user.PendingEmail,
			Status:                     user.Status,
			CreatedAt:                  user.CreatedAt,
			UpdatedAt:                  user.UpdatedAt,
			DeletedAt:                  deletedAt,
			IsDeleted:                  isDeleted,
		}
	}

	return response, totalCount, nil
}

// UpdateProfileImage updates the user's profile image URL
func (s *UserService) UpdateProfileImage(ctx context.Context, userID string, imageURL string) (dto.UserProfileResponse, error) {
	// Find user by ID
	user, err := s.repo.FindUserById(ctx, nil, userID)
	if err != nil {
		return dto.UserProfileResponse{}, dto.ErrUserNotFound
	}

	// Update profile image URL
	user.ProfileImageURL = imageURL
	user.UpdatedAt = time.Now()

	// Save changes
	updatedUser, err := s.repo.Update(ctx, nil, user)
	if err != nil {
		return dto.UserProfileResponse{}, dto.ErrUpdateUserFailed
	}

	// Get user role
	role, err := s.roleRepo.FindByID(ctx, *updatedUser.RoleID)
	if err != nil {
		return dto.UserProfileResponse{}, dto.ErrRoleNotFound
	}

	// Handle DeletedAt field which is a gorm.DeletedAt
	var deletedAt time.Time
	isDeleted := false
	if updatedUser.DeletedAt.Valid {
		deletedAt = updatedUser.DeletedAt.Time
		isDeleted = true
	}

	return dto.UserProfileResponse{
		ID:                         updatedUser.ID.String(),
		Username:                   updatedUser.UserName,
		Email:                      updatedUser.Email,
		FullName:                   updatedUser.FullName,
		Alamat:                     updatedUser.Alamat,
		Latitude:                   updatedUser.Latitude,
		Longitude:                  updatedUser.Longitude,
		ProfileImageURL:            updatedUser.ProfileImageURL,
		RoleID:                     updatedUser.RoleID.String(),
		Role:                       role.Name,
		IsVerified:                 updatedUser.IsVerified,
		AccessToken:                updatedUser.AuthToken,
		TokenExpiry:                updatedUser.TokenExpiry,
		RefreshToken:               updatedUser.RefreshToken,
		TokenCreatedAt:             updatedUser.TokenCreatedAt,
		OTPCode:                    updatedUser.OTPCode,
		OTPCreatedAt:               updatedUser.OTPCreatedAt,
		OTPAttemptCount:            updatedUser.OTPAttemptCount,
		ResendCount:                updatedUser.ResendCount,
		LastResendAt:               updatedUser.LastResendAt,
		CooldownStartedAt:          updatedUser.CooldownStartedAt,
		VerificationToken:          updatedUser.VerificationToken,
		VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
		PendingEmail:               updatedUser.PendingEmail,
		Status:                     updatedUser.Status,
		CreatedAt:                  updatedUser.CreatedAt,
		UpdatedAt:                  updatedUser.UpdatedAt,
		DeletedAt:                  deletedAt,
		IsDeleted:                  isDeleted,
	}, nil
}

// GetUserProfileFromToken gets the profile of the currently authenticated user
func (s *UserService) GetUserProfileFromToken(ctx context.Context, tokenUserID string) (dto.UserProfileResponse, error) {
	// Validate user ID format
	_, err := uuid.Parse(tokenUserID)
	if err != nil {
		log.Printf("Invalid user ID format in profile request: %s", tokenUserID)
		return dto.UserProfileResponse{}, dto.ErrUserNotFound
	}

	// Get user by ID
	user, err := s.repo.FindUserById(ctx, nil, tokenUserID)
	if err != nil {
		log.Printf("User not found in profile request: %s", tokenUserID)
		return dto.UserProfileResponse{}, dto.ErrUserNotFound
	}

	// Get role information
	role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		log.Printf("Role not found for user %s: %v", tokenUserID, err)
		return dto.UserProfileResponse{}, dto.ErrRoleNotFound
	}

	// Create audit log for profile access
	go func() {
		// Use a separate context for the goroutine
		auditCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Try to log this action but don't fail if logging fails
		if auditRepo, ok := s.repo.(interface {
			LogAuditEvent(ctx context.Context, event string, userID string, details string) error
		}); ok {
			_ = auditRepo.LogAuditEvent(
				auditCtx,
				"PROFILE_ACCESS",
				tokenUserID,
				"User accessed their own profile",
			)
		}
	}()

	// Handle DeletedAt field which is a gorm.DeletedAt
	var deletedAt time.Time
	isDeleted := false
	if user.DeletedAt.Valid {
		deletedAt = user.DeletedAt.Time
		isDeleted = true
	}

	// Return the user profile with all fields
	return dto.UserProfileResponse{
		ID:                         user.ID.String(),
		Username:                   user.UserName,
		Email:                      user.Email,
		FullName:                   user.FullName,
		Alamat:                     user.Alamat,
		Latitude:                   user.Latitude,
		Longitude:                  user.Longitude,
		ProfileImageURL:            user.ProfileImageURL,
		RoleID:                     user.RoleID.String(),
		Role:                       role.Name,
		IsVerified:                 user.IsVerified,
		AccessToken:                user.AuthToken,
		TokenExpiry:                user.TokenExpiry,
		RefreshToken:               user.RefreshToken,
		TokenCreatedAt:             user.TokenCreatedAt,
		OTPCode:                    user.OTPCode,
		OTPCreatedAt:               user.OTPCreatedAt,
		OTPAttemptCount:            user.OTPAttemptCount,
		ResendCount:                user.ResendCount,
		LastResendAt:               user.LastResendAt,
		CooldownStartedAt:          user.CooldownStartedAt,
		VerificationToken:          user.VerificationToken,
		VerificationTokenCreatedAt: user.VerificationTokenCreatedAt,
		PendingEmail:               user.PendingEmail,
		Status:                     user.Status,
		CreatedAt:                  user.CreatedAt,
		UpdatedAt:                  user.UpdatedAt,
		DeletedAt:                  deletedAt,
		IsDeleted:                  isDeleted,
	}, nil
}

// UpdateUserData updates a user's data by an admin, including email and username validation
func (s *UserService) UpdateUserData(ctx context.Context, userID string, adminID string, req dto.UserDataUpdateRequest) (dto.UserProfileResponse, error) {
	// Validate admin ID and permissions
	admin, err := s.repo.FindUserById(ctx, nil, adminID)
	if err != nil {
		log.Printf("Admin not found: %s", adminID)
		return dto.UserProfileResponse{}, dto.ErrUnauthorized
	}

	// Get admin's role
	adminRole, err := s.roleRepo.FindByID(ctx, *admin.RoleID)
	if err != nil {
		log.Printf("Admin role not found: %s", *admin.RoleID)
		return dto.UserProfileResponse{}, dto.ErrUnauthorized
	}

	// Verify admin has SUPERADMIN role
	if adminRole.Name != "SUPERADMIN" {
		log.Printf("Unauthorized: Admin role is %s, not SUPERADMIN", adminRole.Name)
		return dto.UserProfileResponse{}, dto.ErrUnauthorized
	}

	// Create a transaction for atomicity
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return dto.UserProfileResponse{}, dto.ErrUpdateUserFailed
	}

	// Ensure proper cleanup of the transaction
	defer func() {
		if r := recover(); r != nil {
			s.repo.RollbackTx(ctx, tx)
			log.Printf("Panic recovered in UpdateUserData: %v", r)
		}
	}()

	// Find user with row lock to prevent concurrent updates
	user, err := s.repo.FindUserByIdWithLock(ctx, tx, userID)
	if err != nil {
		log.Printf("Error finding user with ID %s: %v", userID, err)
		s.repo.RollbackTx(ctx, tx)
		return dto.UserProfileResponse{}, dto.ErrUserNotFound
	}

	// Check if email exists if it's being changed
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.repo.FindUserByEmail(ctx, tx, req.Email)
		if err == nil && existingUser.ID != uuid.Nil {
			// Found a user with this email
			s.repo.RollbackTx(ctx, tx)
			return dto.UserProfileResponse{}, dto.ErrEmailAlreadyInUse
		}
	}

	// Check if username exists if it's being changed
	if req.Username != "" && req.Username != user.UserName {
		existingUser, err := s.repo.FindUserByUsername(ctx, tx, req.Username)
		if err == nil && existingUser.ID != uuid.Nil {
			// Found a user with this username
			s.repo.RollbackTx(ctx, tx)
			return dto.UserProfileResponse{}, dto.ErrUsernameAlreadyInUse
		}
	}

	// Update user fields
	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Username != "" {
		user.UserName = req.Username
	}

	// Always update fullName regardless if it's empty or not,
	// allowing users to clear their name if desired
	user.FullName = req.FullName

	if req.Alamat != "" {
		user.Alamat = req.Alamat
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	// Update coordinates only if they are provided (not zero)
	if req.Latitude != 0 {
		user.Latitude = req.Latitude
	}

	if req.Longitude != 0 {
		user.Longitude = req.Longitude
	}

	user.UpdatedAt = time.Now()

	// Save changes
	updatedUser, err := s.repo.Update(ctx, tx, user)
	if err != nil {
		log.Printf("Failed to update user data: %v", err)
		s.repo.RollbackTx(ctx, tx)
		return dto.UserProfileResponse{}, dto.ErrUpdateUserFailed
	}

	// Commit transaction
	if err := s.repo.CommitTx(ctx, tx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return dto.UserProfileResponse{}, dto.ErrUpdateUserFailed
	}

	// Get user's role
	role, err := s.roleRepo.FindByID(ctx, *updatedUser.RoleID)
	if err != nil {
		log.Printf("Role not found for user %s: %v", userID, err)
		return dto.UserProfileResponse{}, dto.ErrRoleNotFound
	}

	// Create audit log
	go func() {
		auditCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if auditRepo, ok := s.repo.(interface {
			LogAuditEvent(ctx context.Context, event string, userID string, details string) error
		}); ok {
			details := fmt.Sprintf("Admin %s updated user %s data", adminID, userID)
			_ = auditRepo.LogAuditEvent(auditCtx, "USER_DATA_UPDATE", adminID, details)
		}
	}()

	// Handle DeletedAt field which is a gorm.DeletedAt
	var deletedAt time.Time
	isDeleted := false
	if updatedUser.DeletedAt.Valid {
		deletedAt = updatedUser.DeletedAt.Time
		isDeleted = true
	}

	// Return the updated user profile with all fields
	return dto.UserProfileResponse{
		ID:                         updatedUser.ID.String(),
		Username:                   updatedUser.UserName,
		Email:                      updatedUser.Email,
		FullName:                   updatedUser.FullName,
		Alamat:                     updatedUser.Alamat,
		Latitude:                   updatedUser.Latitude,
		Longitude:                  updatedUser.Longitude,
		ProfileImageURL:            updatedUser.ProfileImageURL,
		RoleID:                     updatedUser.RoleID.String(),
		Role:                       role.Name,
		IsVerified:                 updatedUser.IsVerified,
		AccessToken:                updatedUser.AuthToken,
		TokenExpiry:                updatedUser.TokenExpiry,
		RefreshToken:               updatedUser.RefreshToken,
		TokenCreatedAt:             updatedUser.TokenCreatedAt,
		OTPCode:                    updatedUser.OTPCode,
		OTPCreatedAt:               updatedUser.OTPCreatedAt,
		OTPAttemptCount:            updatedUser.OTPAttemptCount,
		ResendCount:                updatedUser.ResendCount,
		LastResendAt:               updatedUser.LastResendAt,
		CooldownStartedAt:          updatedUser.CooldownStartedAt,
		VerificationToken:          updatedUser.VerificationToken,
		VerificationTokenCreatedAt: updatedUser.VerificationTokenCreatedAt,
		PendingEmail:               updatedUser.PendingEmail,
		Status:                     updatedUser.Status,
		CreatedAt:                  updatedUser.CreatedAt,
		UpdatedAt:                  updatedUser.UpdatedAt,
		DeletedAt:                  deletedAt,
		IsDeleted:                  isDeleted,
	}, nil
}

// UpdateUserEmail handles email update requests with verification
func (s *UserService) UpdateUserEmail(ctx context.Context, userID uuid.UUID, newEmail string) error {
	// Find user by ID
	user, err := s.repo.FindUserById(ctx, nil, userID.String())
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Check if the new email is the same as the current email
	if user.Email == newEmail {
		return fmt.Errorf("email baru sama dengan email saat ini")
	}

	// Check if new email already exists for OTHER users
	existingUser, err := s.repo.FindUserByEmail(ctx, nil, newEmail)
	if err == nil && existingUser.ID != uuid.Nil && existingUser.ID != userID {
		// Email exists and belongs to another user
		return dto.ErrUserAlreadyExist
	}

	// Check role of the user
	// Admin users can bypass OTP verification
	role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err == nil && role.Name != "USER" {
		// For non-USER roles, directly update the email
		oldEmail := user.Email
		user.Email = newEmail
		user.UpdatedAt = time.Now()

		// Save changes
		_, err = s.repo.Update(ctx, nil, user)
		if err != nil {
			return dto.ErrUpdateUserFailed
		}

		// Optionally, send confirmation email
		emailBody := fmt.Sprintf("Your email has been updated successfully from %s to %s", oldEmail, newEmail)
		s.emailSender.Send(newEmail, "Email Updated Successfully", emailBody)

		return nil
	}

	// For USER role, generate OTP and send verification email
	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		return dto.ErrSendEmailFailed
	}

	// Store the new email and OTP in user model temporarily
	now := time.Now()
	user.OTPCode = otp
	user.OTPCreatedAt = &now
	user.OTPAttemptCount = 0
	user.PendingEmail = newEmail

	// Save changes
	_, err = s.repo.Update(ctx, nil, user)
	if err != nil {
		return dto.ErrUpdateUserFailed
	}

	// Send verification email with OTP to the NEW email address
	emailBody := utils.BuildEmailUpdateOTPVerificationEmail(user.FullName, user.Email, newEmail, otp)
	err = s.emailSender.Send(newEmail, "Verifikasi Perubahan Email", emailBody)
	if err != nil {
		return dto.ErrSendEmailFailed
	}

	return nil
}

// VerifyEmailUpdate verifies a pending email update
func (s *UserService) VerifyEmailUpdate(ctx context.Context, userID uuid.UUID, otp string) error {
	// Find user by ID
	user, err := s.repo.FindUserById(ctx, nil, userID.String())
	if err != nil {
		return dto.ErrUserNotFound
	}

	// Check if there's a pending email update
	if user.PendingEmail == "" {
		return fmt.Errorf("tidak ada permintaan perubahan email")
	}

	// Check attempt count
	if user.OTPAttemptCount >= utils.MaxOTPAttempts {
		return dto.ErrTooManyAttempts
	}

	// Update attempt count
	user.OTPAttemptCount++

	// Check if OTP is expired
	if user.OTPCreatedAt == nil || utils.IsOTPExpired(*user.OTPCreatedAt) {
		_, err = s.repo.Update(ctx, nil, user) // Save updated attempt count
		return dto.ErrOTPExpired
	}

	// Check if OTP matches
	if user.OTPCode != otp {
		_, err = s.repo.Update(ctx, nil, user) // Save updated attempt count
		return dto.ErrOTPNotMatch
	}

	// OTP is valid, update email
	oldEmail := user.Email
	newEmail := user.PendingEmail

	// Store the old values before updating
	user.Email = newEmail
	user.PendingEmail = ""
	user.OTPCode = ""
	user.UpdatedAt = time.Now()

	// Start a transaction for updating the user
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			s.repo.RollbackTx(ctx, tx)
		}
	}()

	// Save changes
	_, err = s.repo.Update(ctx, tx, user)
	if err != nil {
		return dto.ErrUpdateUserFailed
	}

	// Commit transaction
	if err := s.repo.CommitTx(ctx, tx); err != nil {
		return dto.ErrUpdateUserFailed
	}

	// Send confirmation emails to both old and new email addresses
	// 1. Notification to the old email
	oldEmailBody := utils.BuildEmailChangeNotificationEmail(user.FullName, oldEmail, newEmail)
	go s.emailSender.Send(oldEmail, "Pemberitahuan Perubahan Email", oldEmailBody)

	// 2. Confirmation to the new email
	newEmailBody := fmt.Sprintf("Selamat! Email akun Anda telah berhasil diubah dari %s menjadi %s", oldEmail, newEmail)
	go s.emailSender.Send(newEmail, "Email Berhasil Diperbarui", newEmailBody)

	return nil
}

// CreateUserByAdmin creates a new user by an admin
func (s *UserService) CreateUserByAdmin(ctx context.Context, req dto.UserRegisterRequest, isVerified bool) (dto.UserResponse, error) {
	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		return dto.UserResponse{}, dto.ErrInvalidEmail
	}

	// Check if user already exists
	exists, _ := s.repo.ExistsByEmailOrUsername(ctx, nil, req.Email, req.UserName)
	if exists {
		return dto.UserResponse{}, dto.ErrUserAlreadyExist
	}

	// Hash password
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Get role by name (default to USER if not specified)
	roleName := "USER"
	if req.RoleName != "" {
		roleName = req.RoleName
	}

	role, err := s.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Initialize with current time to avoid nil pointers
	now := time.Now()

	user := entity.User{
		ID:                         uuid.New(),
		Email:                      req.Email,
		UserName:                   req.UserName,
		Password:                   hashedPwd,
		FullName:                   req.FullName,
		Alamat:                     req.Alamat,
		Latitude:                   req.Latitude,
		Longitude:                  req.Longitude,
		ProfileImageURL:            "",
		RoleID:                     &role.ID,   // Convert uuid.UUID to *uuid.UUID
		IsVerified:                 isVerified, // Set verified status based on admin preference
		OTPCreatedAt:               &now,
		VerificationTokenCreatedAt: &now,
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
		Status:                     "active", // Set status to active by default
	}

	// Create user in database
	createdUser, err := s.repo.CreateUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUserFailed
	}

	// If not verified, generate verification materials
	if !isVerified {
		// Generate OTP
		otp, err := utils.GenerateOTP()
		// Generate verification token for the link method
		verificationToken, err := utils.GenerateVerificationToken()
		if err != nil {
			log.Printf("Failed to generate verification token: %v", err)
		} else {
			// Save OTP and verification token to user record
			createdUser.OTPCode = otp
			createdUser.OTPCreatedAt = &now
			createdUser.OTPAttemptCount = 0
			createdUser.VerificationToken = verificationToken
			createdUser.VerificationTokenCreatedAt = &now

			// Update user with OTP info
			updatedUser, updateErr := s.repo.Update(ctx, nil, createdUser)
			if updateErr != nil {
				log.Printf("Failed to update user with OTP information: %v", updateErr)
			} else {
				createdUser = updatedUser
			}

			// Send verification email
			verificationLink := utils.BuildVerificationLink("https://lecsens-iot.erplabiim.com", verificationToken, createdUser.Email)
			emailBody := utils.BuildOTPAndLinkVerificationEmail(createdUser.FullName, otp, verificationLink)
			emailErr := s.emailSender.Send(createdUser.Email, "Verifikasi Email", emailBody)
			if emailErr != nil {
				log.Printf("Failed to send verification email: %v", emailErr)
			}
		}
	}

	return dto.UserResponse{
		ID:              createdUser.ID.String(),
		Username:        createdUser.UserName,
		Email:           createdUser.Email,
		UserRole:        role.Name,
		IsVerified:      createdUser.IsVerified,
		Status:          createdUser.Status,
		ProfileImageURL: createdUser.ProfileImageURL,
	}, nil
}

// GetUserByID retrieves user information by user ID for external API usage
func (s *UserService) GetUserByID(ctx context.Context, userID string) (dto.ExternalUserInfoResponse, error) {
	// Validate user ID format
	_, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %s", userID)
		return dto.ExternalUserInfoResponse{}, dto.ErrUserNotFound
	}

	// Get user by ID
	user, err := s.repo.FindUserById(ctx, nil, userID)
	if err != nil {
		log.Printf("User not found: %s", userID)
		return dto.ExternalUserInfoResponse{}, dto.ErrUserNotFound
	}

	// Get role information
	role, err := s.roleRepo.FindByID(ctx, *user.RoleID)
	if err != nil {
		log.Printf("Role not found for user %s: %v", userID, err)
		return dto.ExternalUserInfoResponse{}, fmt.Errorf("role not found")
	}

	// Return user info for external API
	return dto.ExternalUserInfoResponse{
		UserID:   user.ID.String(),
		Username: user.UserName,
		Email:    user.Email,
		Role:     role.Name,
		RoleID:   user.RoleID.String(),
		Status:   user.Status,
		IsActive: user.Status == "active",
	}, nil
}
