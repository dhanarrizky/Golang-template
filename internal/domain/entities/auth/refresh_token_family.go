package auth

import "time"

// RefreshTokenFamily adalah security aggregate
// Satu family = satu login session group
type RefreshTokenFamily struct {
	ID        uint64
	UserID    uint64
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (f *RefreshTokenFamily) IsRevoked() bool {
	return f.RevokedAt != nil
}
