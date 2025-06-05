package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims represents the claims in a JWT
type TokenClaims struct {
	UserID   string `json:"user_id"`
	RoleID   string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

// TokenMaker is an interface for managing tokens
type TokenMaker interface {
	// Access token methods
	CreateAccessToken(userID, roleID, roleName string, duration time.Duration) (string, error)
	ParseAccessToken(token string) (*TokenClaims, error)

	// Email verification token methods
	CreateEmailToken(email string, duration time.Duration) (string, error)
	ParseEmailToken(token string) (string, error)

	// Password reset token methods
	CreatePasswordResetToken(email string, duration time.Duration) (string, error)
	ParsePasswordResetToken(token string) (string, error)

	// Refresh token methods
	CreateRefreshToken(userID string, duration time.Duration) (string, error)
	ParseRefreshToken(token string) (*TokenClaims, error)

	// Token revocation methods
	BlacklistUserTokens(ctx context.Context, userID string) error

	// Email validation methods
	ValidateEmailToken(token string, email string) bool
}

// JWTMaker implements the TokenMaker interface using the JWT token
type JWTMaker struct {
	accessTokenSecret   string
	refreshTokenSecret  string
	emailTokenSecret    string
	passwordResetSecret string
}

// Minimum secret key size in characters
const MinSecretKeySize = 32

// JWTClaims contains the claims data for JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	RoleID   string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	Email    string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

// NewTokenMaker creates a new TokenMaker instance using environment variables
func NewTokenMaker() (TokenMaker, error) {
	// Get JWT secrets from environment variables
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	emailSecret := os.Getenv("JWT_EMAIL_SECRET")
	passwordResetSecret := os.Getenv("JWT_PASSWORD_RESET_SECRET")

	// Create a new JWTMaker
	return NewJWTMaker(accessSecret, refreshSecret, emailSecret, passwordResetSecret)
}

// NewJWTMaker creates a new JWTMaker with the given secret keys
func NewJWTMaker(accessSecret, refreshSecret, emailSecret, passwordResetSecret string) (TokenMaker, error) {
	if len(accessSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: access token secret must be at least %d characters", MinSecretKeySize)
	}
	if len(refreshSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: refresh token secret must be at least %d characters", MinSecretKeySize)
	}
	if len(emailSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: email token secret must be at least %d characters", MinSecretKeySize)
	}
	if len(passwordResetSecret) < MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: password reset secret must be at least %d characters", MinSecretKeySize)
	}

	return &JWTMaker{
		accessTokenSecret:   accessSecret,
		refreshTokenSecret:  refreshSecret,
		emailTokenSecret:    emailSecret,
		passwordResetSecret: passwordResetSecret,
	}, nil
}

// CreateToken implements the TokenMaker interface (legacy method kept for compatibility)
func (maker *JWTMaker) CreateToken(userID string, duration int) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(duration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.accessTokenSecret))
}

// CreateAccessToken implements the TokenMaker interface
func (maker *JWTMaker) CreateAccessToken(userID string, roleID string, roleName string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		RoleID:   roleID,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.accessTokenSecret))
}

// CreateEmailToken implements the TokenMaker interface
func (maker *JWTMaker) CreateEmailToken(email string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.emailTokenSecret))
}

// ParseEmailToken implements the TokenMaker interface
func (maker *JWTMaker) ParseEmailToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.emailTokenSecret), nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || claims.Email == "" {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.Email, nil
}

// VerifyToken implements the TokenMaker interface
func (maker *JWTMaker) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.accessTokenSecret), nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, nil
}

// ParseAccessToken implements the TokenMaker interface
func (maker *JWTMaker) ParseAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.accessTokenSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &TokenClaims{
		UserID:   claims.UserID,
		RoleID:   claims.RoleID,
		RoleName: claims.RoleName,
		Email:    claims.Email,
	}, nil
}

// CreatePasswordResetToken implements the TokenMaker interface
func (maker *JWTMaker) CreatePasswordResetToken(email string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.passwordResetSecret))
}

// ParsePasswordResetToken implements the TokenMaker interface
func (maker *JWTMaker) ParsePasswordResetToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.passwordResetSecret), nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || claims.Email == "" {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.Email, nil
}

// CreateRefreshToken implements the TokenMaker interface
func (maker *JWTMaker) CreateRefreshToken(userID string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.refreshTokenSecret))
}

// ParseRefreshToken implements the TokenMaker interface
func (maker *JWTMaker) ParseRefreshToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.refreshTokenSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &TokenClaims{
		UserID:   claims.UserID,
		RoleID:   claims.RoleID,
		RoleName: claims.RoleName,
		Email:    claims.Email,
	}, nil
}

// BlacklistUserTokens implements the TokenMaker interface
func (maker *JWTMaker) BlacklistUserTokens(ctx context.Context, userID string) error {
	// Implementation for blacklisting user tokens
	// This is a placeholder and should be replaced with actual logic
	return nil
}

// ValidateEmailToken checks if the email token is valid
func (maker *JWTMaker) ValidateEmailToken(token string, email string) bool {
	claims, err := maker.ParseEmailToken(token)
	if err != nil {
		return false
	}
	return claims == email
}
