package auth

import "time"

type PasswordResetToken struct {
	ID        uint64
	UserID    uint64

	TokenHash string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}
