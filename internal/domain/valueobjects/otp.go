package valueobjects

import "time"

type OTP struct {
	Hash      string
	ExpiredAt time.Time
}

func (o OTP) IsExpired(now time.Time) bool {
	return now.After(o.ExpiredAt)
}
