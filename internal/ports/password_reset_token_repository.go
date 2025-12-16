package ports

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *auth.PasswordResetToken) error
	GetByTokenHash(ctx context.Context, hash string) (*auth.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id uint64) error
}
