package auth

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type LoginAttemptRepository interface {
	LogAttempt(ctx context.Context, attempt *auth.LoginAttempt) error
	CountFailedByIP(ctx context.Context, ip string, sinceMinutes int) (int, error)
	CountFailedByEmail(ctx context.Context, email string, sinceMinutes int) (int, error)
}
