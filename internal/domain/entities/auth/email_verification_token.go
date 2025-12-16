package auth

import "time"

type EmailVerificationToken struct {
	ID        uint64
	UserID    uint64
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}
