package auth

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type RefreshTokenFamilyRepository interface {
	Create(ctx context.Context, family *auth.RefreshTokenFamily) error
	GetByID(ctx context.Context, id uint64) (*auth.RefreshTokenFamily, error)
	GetByUserID(ctx context.Context, userID uint64) ([]*auth.RefreshTokenFamily, error)
	Revoke(ctx context.Context, id uint64) error
}
