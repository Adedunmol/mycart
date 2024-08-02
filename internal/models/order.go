package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	BuyerID uint8
	CartID  uint8
}
