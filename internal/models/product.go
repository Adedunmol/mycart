package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name     string `json:"name"`
	Details  string `json:"details"`
	Price    int    `json:"price"`
	Category string `json:"category"`
	Vendor   uint
	// review (rating, comment)
}
