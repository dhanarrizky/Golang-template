package auth

import "time"

type LoginAttempt struct {
	ID uint

	Email     string
	IPAddress string
	Success   bool
	UserAgent string

	CreatedAt time.Time
}
