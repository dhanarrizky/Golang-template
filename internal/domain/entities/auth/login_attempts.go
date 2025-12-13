package auth

import (
	"time"

	"gorm.io/gorm"
)

type LoginAttempt struct {
	ID uint `gorm:"primaryKey"`

	// Identitas percobaan login
	Email     string `gorm:"size:255;index"`
	IPAddress string `gorm:"size:45;index"` // IPv6 safe

	// Status
	Success bool `gorm:"not null"`

	// Metadata (optional tapi berguna)
	UserAgent string `gorm:"size:512"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Soft delete optional (biasanya tidak perlu)
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
