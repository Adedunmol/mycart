package util

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const TOKEN_EXPIRATION = 15 * time.Minute

func GenerateToken(username string) (string, error) {

	claims := jwt.MapClaims{
		"username":   username,
		"Expiration": time.Now().Add(TOKEN_EXPIRATION).Unix(),
		"IssuedAt":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
