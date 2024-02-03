package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	BuyerID    uint8
	TotalPrice uint8
}

type CartItem struct {
	gorm.Model
	CartID     uint8
	ProductID  uint8
	Quantity   uint8
	TotalPrice uint8
}
