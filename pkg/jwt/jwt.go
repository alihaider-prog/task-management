package jwt

import (
	"errors"
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwtlib.RegisteredClaims
}

var secret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID int64, userEmail string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  userEmail,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwtlib.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwtlib.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return secret, nil
		})

	if err != nil {
		return nil, err
	}

	Claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return Claims, nil
}
