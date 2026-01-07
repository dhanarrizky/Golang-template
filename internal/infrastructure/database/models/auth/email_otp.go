package auth

import "time"

type EmailOTP struct {
	ID uint64 `gorm:"primaryKey"`

	Email string `gorm:"size:255;index;not null"`

	OTPHash string `gorm:"size:255;not null"`

	ExpiredAt time.Time `gorm:"index;not null"`

	IPAddress string `gorm:"size:45"`
	UserAgent string `gorm:"size:512"`
	Purpose   string `gorm:"size:512"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}
