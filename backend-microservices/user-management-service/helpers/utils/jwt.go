package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID       string
	UserRole string
	jwt.StandardClaims
}

var JWTKey = []byte(os.Getenv("JWT_SECRET"))

func JWTGenerateToken(userID string, userRole string) (string, error) {
	expTime := time.Now().Add(time.Hour * 24).Unix()

	claims := &Claims{
		ID:       userID,
		UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}
