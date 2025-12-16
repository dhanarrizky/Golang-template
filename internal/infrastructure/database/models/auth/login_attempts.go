package auth

import (
	"time"

	"gorm.io/gorm"
)

type LoginAttempt struct {
	ID uint `gorm:"primaryKey"`

	Email     string `gorm:"size:255;index"`
	IPAddress string `gorm:"size:45;index"`
	Success   bool   `gorm:"not null"`
	UserAgent string `gorm:"size:512"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
