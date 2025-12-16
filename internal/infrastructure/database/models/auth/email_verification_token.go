package auth

import "time"

type EmailVerificationToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;type:bigserial"`
	UserID    uint64    `gorm:"not null;index:idx_evt_user_id"`

	TokenHash string    `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index:idx_evt_expires_at"`
	CreatedAt time.Time
}
