package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Username     string    `json:"username" gorm:"unique"`
	Email        string    `json:"email" gorm:"unique"`
	Password     string    `json:"-"`
	Verified     bool      `json:"verified"`
	RefreshToken string    `json:"refresh_token"`
	Otp          Otp       `json:"otp"`
	RoleID       uint      `json:"role_id"`
	Products     []Product `gorm:"foreignKey:Vendor"`
	Orders       []Order   `gorm:"foreignKey:BuyerID"`
}
