package auth

import (
	"time"

	"gorm.io/gorm"
)

type RefreshTokenFamily struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;type:bigserial"`
	UserID    uint64     `gorm:"not null;index:idx_rtf_user_id"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	RevokedAt *time.Time
	CreatedAt time.Time

	RefreshTokens []RefreshToken
}
