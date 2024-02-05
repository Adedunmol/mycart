package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	BuyerID    uint
	TotalPrice uint
	CartItems  []CartItem `gorm:"foreignKey:Cart"`
}

type CartItem struct {
	gorm.Model
	Cart       uint
	ProductID  uint
	Quantity   uint
	TotalPrice uint
}
