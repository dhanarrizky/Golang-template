package auth

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;type:bigserial"`
	UserID    uint64    `gorm:"not null;index:idx_prt_user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	TokenHash string    `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index:idx_prt_expires_at"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}
