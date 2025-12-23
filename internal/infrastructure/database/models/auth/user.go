package auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID uint64 `gorm:"primaryKey;autoIncrement;type:bigserial"`

	Email         string  `gorm:"size:255;uniqueIndex;not null"`
	EmailVerified bool    `gorm:"default:false"`
	PasswordHash  string  `gorm:"type:text;not null"`
	Name          *string `gorm:"size:255"`

	RoleID uint64 `gorm:"not null;index"`
	Role   Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations (ORM only)
	RefreshTokenFamilies []RefreshTokenFamily
	RefreshTokens        []RefreshToken
	Sessions             []UserSession
}
