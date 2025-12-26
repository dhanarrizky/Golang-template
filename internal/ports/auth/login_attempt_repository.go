package auth

import (
	"context"
)

// type LoginAttemptRepository interface {
// 	LogAttempt(ctx context.Context, attempt *auth.LoginAttempt) error
// 	CountFailedByIP(ctx context.Context, ip string, sinceMinutes int) (int, error)
// 	CountFailedByEmail(ctx context.Context, email string, sinceMinutes int) (int, error)
// }

type LoginAttemptRepository interface {
	IsRateLimited(ctx context.Context, identifier string) bool
	RecordFailedAttempt(ctx context.Context, identifier string) error
	ResetAttempts(ctx context.Context, identifier string) error
}
