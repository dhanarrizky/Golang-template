package repository

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type RefreshTokenRepository interface {
	Create(
		ctx context.Context,
		token *entities.RefreshToken,
	) error

	GetByTokenHash(
		ctx context.Context,
		hash string,
	) (*entities.RefreshToken, error)

	Revoke(
		ctx context.Context,
		id uint,
	) error

	RevokeByFamily(
		ctx context.Context,
		familyID uint,
	) error

	DeleteExpired(
		ctx context.Context,
	) error
}
