package util

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	jwt "github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_EXPIRATION = 15 * time.Minute
const REFRESH_TOKEN_EXPIRATION = 1 * time.Hour

func GenerateToken(username string, expiration time.Duration) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(), // jwt.NewNumericDate(time.Now().Add(expiration)),
		"iat":      time.Now(),
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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
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

func AuthMiddleware(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			RespondWithJSON(w, http.StatusUnauthorized, APIResponse{Message: "No auth token in the header", Data: nil, Status: "error"})
			return
		}

		tokenString := strings.Split(authHeader, " ")

		if len(tokenString) != 2 {
			RespondWithJSON(w, http.StatusUnauthorized, APIResponse{Message: "Malformed token", Data: nil, Status: "error"})
			return
		}

		username, err := DecodeToken(tokenString[1])
		if err != nil || username == "" {
			RespondWithJSON(w, http.StatusUnauthorized, APIResponse{Message: "Bad token or token is expired", Data: nil, Status: "error"})
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		newReq := r.WithContext(ctx)

		handler.ServeHTTP(w, newReq)

	})
}

func RoleAuthorization(permissions ...uint8) func(handler http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.Context().Value("username")

			if username == nil {
				RespondWithJSON(w, http.StatusUnauthorized, "Not authorized")
				return
			}
			var foundUser models.User

			result := database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

			if result.Error != nil {
				RespondWithJSON(w, http.StatusUnauthorized, APIResponse{Message: "Unauthorized", Data: nil, Status: "error"})
				return
			}

			var role models.Role

			result = database.DB.First(&role, foundUser.RoleID)
			if result.Error != nil {
				RespondWithJSON(w, http.StatusUnauthorized, APIResponse{Message: "Unauthorized", Data: nil, Status: "error"})
				return
			}

			for _, perm := range permissions {

				if !role.HasPermission(uint8(perm)) {
					RespondWithJSON(w, http.StatusForbidden, APIResponse{Message: "Forbidden", Data: nil, Status: "error"})
					return
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
