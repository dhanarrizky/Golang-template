package auth

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;type:bigserial"`
	Name        string    `gorm:"size:50;uniqueIndex;not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time

	// Relations
	Users []User
}
