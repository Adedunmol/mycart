package util

import (
	"fmt"
	"time"

	"github.com/Adedunmol/mycart/internal/config"
	jwt "github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_EXPIRATION = 15 * time.Minute
const REFRESH_TOKEN_EXPIRATION = 1 * time.Hour

func GenerateToken(username string, expiration time.Duration) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      jwt.NewNumericDate(time.Now().Add(expiration)),
		"IssuedAt": time.Now(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.EnvConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeToken(tokenString string) (string, error) {
	var err error

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println(config.EnvConfig.SecretKey)
		return []byte(config.EnvConfig.SecretKey), nil
	})

	if err != nil {
		fmt.Println("error1: ", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return string(claims["username"].(string)), nil
	}

	return "", err
}
