package util

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_EXPIRATION = 15 * time.Minute
const REFRESH_TOKEN_EXPIRATION = 1 * time.Hour

func GenerateToken(username string, expiration time.Duration) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration),
		"IssuedAt": time.Now(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
