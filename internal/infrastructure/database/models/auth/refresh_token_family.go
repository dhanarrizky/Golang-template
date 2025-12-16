package auth

import "time"

type RefreshTokenFamily struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;type:bigserial"`
	UserID    uint64     `gorm:"not null;index:idx_rtf_user_id"`
	RevokedAt *time.Time `gorm:"index"`
	CreatedAt time.Time

	// Relasi hanya untuk ORM convenience
	RefreshTokens []RefreshToken `gorm:"foreignKey:FamilyID"`
}
