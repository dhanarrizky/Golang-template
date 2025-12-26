package auth

import (
	"time"
)

// type LoginAttempt struct {
// 	ID uint `gorm:"primaryKey"`

// 	Email     string `gorm:"size:255;index"`
// 	IPAddress string `gorm:"size:45;index"`
// 	Success   bool   `gorm:"not null"`
// 	UserAgent string `gorm:"size:512"`

// 	CreatedAt time.Time      `gorm:"autoCreateTime"`
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// }

type LoginAttempt struct {
	ID uint64 `gorm:"primaryKey"`

	Identifier string `gorm:"size:255;index;not null"`

	UserID *uint64 `gorm:"index"`

	IPAddress string `gorm:"size:45;index;not null"`
	UserAgent string `gorm:"size:512"`

	Success bool `gorm:"not null"`

	FailureReason string    `gorm:"size:64"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}
