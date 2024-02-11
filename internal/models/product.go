package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name      string `json:"name"`
	Details   string `json:"details"`
	Price     int    `json:"price"`
	Quantity  uint   `json:"quantity"`
	Category  string `json:"category"`
	Vendor    uint
	CartItems []CartItem `gorm:"foreignKey:ProductID"`
	// review (rating, comment)
}
