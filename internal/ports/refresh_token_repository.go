package ports

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *auth.RefreshToken) error
	GetByTokenHash(ctx context.Context, hash string) (*auth.RefreshToken, error)

	Revoke(ctx context.Context, id uint64) error
	RevokeByFamily(ctx context.Context, familyID uint64) error

	DeleteExpired(ctx context.Context) error
}
