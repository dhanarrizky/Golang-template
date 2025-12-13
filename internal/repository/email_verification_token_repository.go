package repository

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
)

type EmailVerificationTokenRepository interface {
	Create(
		ctx context.Context,
		token *entities.EmailVerificationToken,
	) error

	GetByTokenHash(
		ctx context.Context,
		hash string,
	) (*entities.EmailVerificationToken, error)

	DeleteByUser(
		ctx context.Context,
		userID uint,
	) error
}
