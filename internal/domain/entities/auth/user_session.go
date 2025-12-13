package auth

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;type:bigserial"`
	UserID    uint64     `gorm:"not null;index:idx_us_user_id"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	IPAddress  string
	UserAgent  string
	LoginAt    time.Time  `gorm:"not null;default:now()"`
	LastSeenAt *time.Time `gorm:"index:idx_us_last_seen"`
	LogoutAt   *time.Time
}
