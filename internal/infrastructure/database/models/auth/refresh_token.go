package auth

import "time"

type RefreshToken struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;type:bigserial"`
	UserID    uint64     `gorm:"not null;index:idx_rt_user_family,priority:1"`
	FamilyID  uint64     `gorm:"not null;index:idx_rt_user_family,priority:2"`

	TokenHash string     `gorm:"size:255;uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"index"`
	CreatedAt time.Time
	RevokedAt *time.Time `gorm:"index"`

	IPAddress string `gorm:"size:60"`
	UserAgent string `gorm:"type:text"`
}
