package auth

import "time"

type EmailOTP struct {
	ID uint64

	Email string

	OTPHash   string
	ExpiredAt time.Time

	IPAddress string
	UserAgent string

	CreatedAt time.Time
}
