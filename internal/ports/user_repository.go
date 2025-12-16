package ports

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uint64) (*auth.User, error)
	GetByEmail(ctx context.Context, email string) (*auth.User, error)

	Create(ctx context.Context, user *auth.User) error
	Update(ctx context.Context, user *auth.User) error

	SoftDelete(ctx context.Context, id uint64) error
}
