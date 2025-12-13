package repository

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
)

type LoginAttemptRepository interface {
	LogAttempt(
		ctx context.Context,
		attempt *entities.LoginAttempt,
	) error

	CountFailedByIP(
		ctx context.Context,
		ip string,
		minutes int,
	) (int, error)

	CountFailedByEmail(
		ctx context.Context,
		email string,
		minutes int,
	) (int, error)
}
