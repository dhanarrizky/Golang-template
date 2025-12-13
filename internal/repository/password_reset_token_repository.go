package repository

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
)

type PasswordResetTokenRepository interface {
	Create(
		ctx context.Context,
		token *entities.PasswordResetToken,
	) error

	GetByTokenHash(
		ctx context.Context,
		hash string,
	) (*entities.PasswordResetToken, error)

	MarkUsed(
		ctx context.Context,
		id uint,
	) error
}
