package models

import (
	"time"

	"gorm.io/gorm"
)

type Otp struct {
	gorm.Model
	User      uint      `json:"user"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}
