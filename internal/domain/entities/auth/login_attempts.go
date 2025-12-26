package auth

import "time"

type LoginAttempt struct {
	ID       uint64
	Identity string
	UserID   *uint64

	IPAddress string
	UserAgent string

	Success bool
	Reason  string

	CreatedAt time.Time
}
