package roles

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

// type RoleRepository interface {
// 	GetByID(ctx context.Context, id uint64) (*auth.Role, error)
// 	GetByName(ctx context.Context, name string) (*auth.Role, error)

// 	List(ctx context.Context) ([]*auth.Role, error)

// 	Create(ctx context.Context, role *auth.Role) error
// 	Update(ctx context.Context, role *auth.Role) error
// }

type RoleRepository interface {
	GetByID(ctx context.Context, id uint64) (*auth.Role, error)
	GetByName(ctx context.Context, name string) (*auth.Role, error)

	List(ctx context.Context) ([]*auth.Role, error)

	Create(ctx context.Context, role *auth.Role) error
	Update(ctx context.Context, role *auth.Role) error

	Delete(ctx context.Context, id uint64) error
	IsRoleUsed(ctx context.Context, id uint64) (bool, error)
}
