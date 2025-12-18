package valueobjects

import "time"

type TokenPayload struct {
	UserID    string
	TokenID   string
	ExpiresAt time.Time
}
