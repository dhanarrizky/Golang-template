package auth

import "time"

// UserSession merepresentasikan login session (bukan token)
type UserSession struct {
	ID     uint64
	UserID uint64

	IPAddress string
	UserAgent string

	LoginAt    time.Time
	LastSeenAt *time.Time
	LogoutAt   *time.Time
}

func (s *UserSession) IsActive() bool {
	return s.LogoutAt == nil
}

func (s *UserSession) Touch(now time.Time) {
	s.LastSeenAt = &now
}

func (s *UserSession) Logout(now time.Time) {
	s.LogoutAt = &now
}
