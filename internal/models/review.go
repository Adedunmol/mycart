package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Comment   string `json:"comment"`
	Rating    uint   `json:"rating"`
	ProductID uint   `json:"product_id"`
	UserID    uint   `json:"user_id"`
}
