package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Username     string    `json:"username" gorm:"unique"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	RefreshToken string    `json:"refresh_token"`
	RoleID       uint      `json:"role_id"`
	Products     []Product `gorm:"foreignKey:Vendor"`
	Orders       []Order   `gorm:"foreignKey:BuyerID"`
}
