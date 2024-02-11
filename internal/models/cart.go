package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	BuyerID    uint       `json:"buyer_id"`
	TotalPrice uint       `json:"total_price"`
	CartItems  []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	CartID      uint   `json:"cart_id"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	UnitPrice   uint   `json:"unit_price"`
	Quantity    uint   `json:"quantity"`
	TotalPrice  uint   `json:"total_price"`
}
