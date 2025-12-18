package mailer

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type EmailVerificationTokenRepository interface {
	Create(ctx context.Context, token *auth.EmailVerificationToken) error
	GetByTokenHash(ctx context.Context, hash string) (*auth.EmailVerificationToken, error)
	DeleteByUser(ctx context.Context, userID uint64) error
}
