package jwt

import (
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID int, ueserEmail string) (string, error) {
	claims := jwtlib.MapClaims{
		"user_id": userID,
		"email":   ueserEmail,
		"exp":     time.Now().Add(24 * time.Hour),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	return token.SignedString(secret)
}
