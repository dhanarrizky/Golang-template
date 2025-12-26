package users

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

// type UserRepository interface {
// 	GetByID(ctx context.Context, id uint64) (*auth.User, error)
// 	GetByEmail(ctx context.Context, email string) (*auth.User, error)

// 	GetByEmailOrUsername(ctx context.Context, identifier string) (*auth.User, error)

// 	Create(ctx context.Context, user *auth.User) error
// 	Update(ctx context.Context, user *auth.User) error

// 	UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error

// 	SoftDelete(ctx context.Context, id uint64) error
// }

type UserRepository interface {
	GetByID(ctx context.Context, id uint64) (*auth.User, error)
	GetByEmail(ctx context.Context, email string) (*auth.User, error)
	GetByEmailOrUsername(ctx context.Context, identifier string) (*auth.User, error)

	Create(ctx context.Context, user *auth.User) error
	Update(ctx context.Context, user *auth.User) error
	UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error
	UpdateUsername(ctx context.Context, id uint64, hashedPassword string) error

	ExistsByUsernameExceptID(ctx context.Context, username string, exceptID uint64) (bool, error)
	ExistsByEmailExceptID(ctx context.Context, email string, exceptID uint64) (bool, error)

	SoftDelete(ctx context.Context, id uint64) error
}
