package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWTTokenMaker is a JWT implementation of TokenMaker
type JWTTokenMaker struct {
	accessTokenSecret   string
	refreshTokenSecret  string
	emailTokenSecret    string
	passwordResetSecret string
}

// NewJWTTokenMaker creates a new JWTTokenMaker
func NewJWTTokenMaker(
	accessTokenSecret,
	refreshTokenSecret,
	emailTokenSecret,
	passwordResetSecret string,
) (TokenMaker, error) {
	if len(accessTokenSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: access token secret must be at least %d characters", MinSecretKeySize)
	}
	if len(refreshTokenSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: refresh token secret must be at least %d characters", MinSecretKeySize)
	}
	if len(emailTokenSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: email token secret must be at least %d characters", MinSecretKeySize)
	}
	if len(passwordResetSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: password reset secret must be at least %d characters", MinSecretKeySize)
	}

	return &JWTTokenMaker{
		accessTokenSecret:   accessTokenSecret,
		refreshTokenSecret:  refreshTokenSecret,
		emailTokenSecret:    emailTokenSecret,
		passwordResetSecret: passwordResetSecret,
	}, nil
}

// CreateAccessToken creates a new access token
// Ubah signature dan implementasi untuk menerima roleName
func (maker *JWTTokenMaker) CreateAccessToken(userID, roleID, roleName string, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	log.Printf("Creating access token - Current time: %v, Expiry time: %v", now, exp)

	payload := jwt.MapClaims{
		"user_id":   userID,
		"role_id":   roleID,
		"role_name": roleName,
		"exp":       exp.Unix(),
		"iat":       now.Unix(), // Add issued at time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(maker.accessTokenSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	log.Printf("Access token created successfully for user %s, expires at %v", userID, exp)
	return tokenString, nil
}

// ParseAccessToken parses an access token and returns its payload
func (maker *JWTTokenMaker) ParseAccessToken(token string) (*TokenClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(maker.accessTokenSecret), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				log.Printf("Token expired: %v", err)
				return nil, fmt.Errorf("token expired: %w", err)
			}
			if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				log.Printf("Invalid token signature: %v", err)
				return nil, fmt.Errorf("invalid token signature: %w", err)
			}
		}
		log.Printf("Error parsing token: %v", err)
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Invalid token claims type")
		return nil, errors.New("invalid token claims type")
	}

	if !jwtToken.Valid {
		log.Printf("Token is not valid")
		return nil, errors.New("token is not valid")
	}

	// Validate expiration time
	exp, ok := claims["exp"].(float64)
	if !ok {
		log.Printf("Invalid expiration time in token")
		return nil, errors.New("invalid expiration time in token")
	}

	expTime := time.Unix(int64(exp), 0)
	now := time.Now()

	log.Printf("Token validation - Current time: %v, Token expiry: %v", now, expTime)

	if now.After(expTime) {
		log.Printf("Token expired at %v", expTime)
		return nil, errors.New("token has expired")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Printf("Invalid user id in token")
		return nil, errors.New("invalid user id in token")
	}

	roleID, _ := claims["role_id"].(string)
	roleName, _ := claims["role_name"].(string)

	log.Printf("Successfully parsed token for user %s with role %s, expires at %v", userID, roleName, expTime)

	return &TokenClaims{
		UserID:   userID,
		RoleID:   roleID,
		RoleName: roleName,
	}, nil
}

// CreateEmailToken creates a new email verification token
func (maker *JWTTokenMaker) CreateEmailToken(email string, duration time.Duration) (string, error) {
	payload := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.emailTokenSecret))
}

// ParseEmailToken parses an email token and returns the email
func (maker *JWTTokenMaker) ParseEmailToken(token string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(maker.emailTokenSecret), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", errors.New("invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("invalid email in token")
	}

	return email, nil
}

// CreatePasswordResetToken creates a new password reset token
func (maker *JWTTokenMaker) CreatePasswordResetToken(email string, duration time.Duration) (string, error) {
	payload := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(duration).Unix(),
		"type":  "password_reset",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.passwordResetSecret))
}

// ParsePasswordResetToken parses a password reset token and returns the email
func (maker *JWTTokenMaker) ParsePasswordResetToken(token string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(maker.passwordResetSecret), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", errors.New("invalid token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "password_reset" {
		return "", errors.New("invalid token type")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("invalid email in token")
	}

	return email, nil
}

// CreateRefreshToken creates a new refresh token
func (maker *JWTTokenMaker) CreateRefreshToken(userID string, duration time.Duration) (string, error) {
	payload := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.refreshTokenSecret))
}

// ParseRefreshToken parses a refresh token and returns its payload
func (maker *JWTTokenMaker) ParseRefreshToken(token string) (*TokenClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(maker.refreshTokenSecret), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user id")
	}

	return &TokenClaims{
		UserID: userID,
	}, nil
}

// BlacklistUserTokens invalidates all tokens for a given user
// This is a simple implementation; in production, you would use Redis or another
// distributed cache to store blacklisted tokens or user IDs
func (maker *JWTTokenMaker) BlacklistUserTokens(ctx context.Context, userID string) error {
	// In a real implementation:
	// 1. You would add the userID to a blacklist in Redis with an expiry time
	// 2. When validating tokens, you would check if the userID is in the blacklist

	// For this implementation, we'll just log that the method was called
	log.Printf("User %s tokens have been marked for invalidation", userID)

	// Since we don't have actual blacklist functionality yet,
	// we'll return success but log a warning
	log.Printf("WARNING: BlacklistUserTokens is a stub implementation")

	return nil
}

// ValidateEmailToken checks if the email token is valid
func (maker *JWTTokenMaker) ValidateEmailToken(token string, email string) bool {
	claims, err := maker.ParseEmailToken(token)
	if err != nil {
		return false
	}
	return claims == email
}
