package repository

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
)

type UserRepository interface {
	GetByID(
		ctx context.Context,
		id uint,
	) (*entities.User, error)

	GetByEmail(
		ctx context.Context,
		email string,
	) (*entities.User, error)

	Create(
		ctx context.Context,
		user *entities.User,
	) error

	Update(
		ctx context.Context,
		user *entities.User,
	) error

	SoftDelete(
		ctx context.Context,
		id uint,
	) error
}
