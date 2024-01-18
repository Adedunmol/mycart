package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	RefreshToken string `json:"refresh_token"`
	RoleID       uint   `json:"role_id"`
}
