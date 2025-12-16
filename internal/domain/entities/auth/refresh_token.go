package auth

import "time"

// RefreshToken adalah credential (bukan session)
type RefreshToken struct {
	ID        uint64
	UserID    uint64
	FamilyID  uint64

	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time

	IPAddress string
	UserAgent string
}

func (t *RefreshToken) IsExpired(now time.Time) bool {
	return now.After(t.ExpiresAt)
}

func (t *RefreshToken) IsRevoked() bool {
	return t.RevokedAt != nil
}
