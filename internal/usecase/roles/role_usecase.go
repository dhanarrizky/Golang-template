package roles

import (
	"context"
	"errors"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	rolePorts "github.com/dhanarrizky/Golang-template/internal/ports/roles"
	userPorts "github.com/dhanarrizky/Golang-template/internal/ports/users"
)

var (
	ErrRoleNotFound   = errors.New("role not found")
	ErrRoleNameExists = errors.New("role name already exists")
	ErrUserNotFound   = errors.New("user not found")
)

type RoleUsecase interface {
	List(ctx context.Context) ([]domain.Role, error)
	Create(ctx context.Context, name string) error
	Update(ctx context.Context, roleID, name string) error
	AssignToUser(ctx context.Context, userID, roleID string) error
}

type roleUsecase struct {
	roleRepo rolePorts.RoleRepository
	userRepo userPorts.UserRepository
}

func NewRoleUsecase(
	roleRepo rolePorts.RoleRepository,
	userRepo userPorts.UserRepository,
) RoleUsecase {
	return &roleUsecase{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// =============== LIST =================

func (u *roleUsecase) List(ctx context.Context) ([]domain.Role, error) {
	return u.roleRepo.FindAll(ctx)
}

// =============== CREATE =================

func (u *roleUsecase) Create(ctx context.Context, name string) error {
	exists, _ := u.roleRepo.GetByName(ctx, name)
	if exists == nil {
		return ErrRoleNameExists
	}

	newRole := domain.Role{
		Name:        name,
		Description: nil,

		CreatedAt: time.Now(),
	}

	return u.roleRepo.Create(ctx, &newRole)
}

// =============== UPDATE =================

func (u *roleUsecase) Update(ctx context.Context, roleID, name string) error {
	role, _ := u.roleRepo.GetByID(ctx, roleID)
	if role == nil {
		return ErrRoleNotFound
	}

	exists, _ := u.roleRepo.GetByName(ctx, name)
	if exists != nil && role.Name != name {
		return ErrRoleNameExists
	}

	exists.Name = name

	return u.roleRepo.Update(ctx, exists)
}

// =============== ASSIGN =================

func (u *roleUsecase) AssignToUser(ctx context.Context, userID, roleID string) error {
	user, _ := u.userRepo.GetByID(ctx, userID)
	if user == nil {
		return ErrUserNotFound
	}

	role, _ := u.roleRepo.GetByID(ctx, roleID)
	if role == nil {
		return ErrRoleNotFound
	}

	return u.roleRepo.AssignToUser(ctx, userID, roleID)
}
