package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	BuyerID  uint8
	Products []Product
}
