package models

import (
	"gorm.io/gorm"
)

type Purchase struct {
	gorm.Model
	BuyerID   uint8
	ProductID uint8
}
