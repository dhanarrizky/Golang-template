package roles

import (
	"context"
	"errors"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrRoleNotFound    = errors.New("role not found")
	ErrRoleNameExists  = errors.New("role name already exists")
	ErrUserNotFound    = errors.New("user not found")
)

type RoleUsecase interface {
	List(ctx context.Context) ([]domain.Role, error)
	Create(ctx context.Context, name string) error
	Update(ctx context.Context, roleID, name string) error
	AssignToUser(ctx context.Context, userID, roleID string) error
}

type roleUsecase struct {
	roleRepo ports.RoleRepository
	userRepo ports.UserRepository
}

func NewRoleUsecase(
	roleRepo ports.RoleRepository,
	userRepo ports.UserRepository,
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
	exists, _ := u.roleRepo.ExistsByName(ctx, name)
	if exists {
		return ErrRoleNameExists
	}
	return u.roleRepo.Create(ctx, name)
}

// =============== UPDATE =================

func (u *roleUsecase) Update(ctx context.Context, roleID, name string) error {
	role, _ := u.roleRepo.FindByID(ctx, roleID)
	if role == nil {
		return ErrRoleNotFound
	}

	exists, _ := u.roleRepo.ExistsByName(ctx, name)
	if exists && role.Name != name {
		return ErrRoleNameExists
	}

	return u.roleRepo.UpdateName(ctx, roleID, name)
}

// =============== ASSIGN =================

func (u *roleUsecase) AssignToUser(ctx context.Context, userID, roleID string) error {
	user, _ := u.userRepo.FindByID(ctx, userID)
	if user == nil {
		return ErrUserNotFound
	}

	role, _ := u.roleRepo.FindByID(ctx, roleID)
	if role == nil {
		return ErrRoleNotFound
	}

	return u.roleRepo.AssignToUser(ctx, userID, roleID)
}
