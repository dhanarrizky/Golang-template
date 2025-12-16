package auth

import "time"

type UserSession struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement;type:bigserial"`
	UserID uint64 `gorm:"not null;index:idx_us_user_id"`

	IPAddress string    `gorm:"size:60"`
	UserAgent string    `gorm:"type:text"`
	LoginAt   time.Time `gorm:"not null"`
	LastSeenAt *time.Time `gorm:"index"`
	LogoutAt   *time.Time
}
