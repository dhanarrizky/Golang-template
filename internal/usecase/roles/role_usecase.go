package roles

import (
	"context"
	"errors"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	otherPorts "github.com/dhanarrizky/Golang-template/internal/ports/others"
	rolePorts "github.com/dhanarrizky/Golang-template/internal/ports/roles"
	userPorts "github.com/dhanarrizky/Golang-template/internal/ports/users"
)

var (
	ErrDecode         = errors.New("internal server decode")
	ErrRoleNotFound   = errors.New("role not found")
	ErrRoleNameExists = errors.New("role name already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrRoleInUse      = errors.New("role is still assigned to users")
)

type RoleUsecase interface {
	List(ctx context.Context) ([]domain.Role, error)
	Create(ctx context.Context, name string) error
	Update(ctx context.Context, roleID, name string) error
	// AssignToUser(ctx context.Context, userID, roleID string) error
}

type roleUsecase struct {
	roleRepo rolePorts.RoleRepository
	userRepo userPorts.UserRepository
	idCodec  otherPorts.PublicIDCodec
}

func NewRoleUsecase(
	roleRepo rolePorts.RoleRepository,
	userRepo userPorts.UserRepository,
	idCodec otherPorts.PublicIDCodec,
) RoleUsecase {
	return &roleUsecase{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// =============== LIST =================

func (u *roleUsecase) List(ctx context.Context) ([]domain.Role, error) {
	roles, err := u.roleRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Role, 0, len(roles))
	for _, r := range roles {
		if r == nil {
			continue
		}
		result = append(result, *r)
	}

	return result, nil
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
	id, err := u.idCodec.Decode(roleID)
	if err != nil {
		return ErrDecode
	}

	role, _ := u.roleRepo.GetByID(ctx, id)
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

// func (u *roleUsecase) AssignToUser(ctx context.Context, userID, roleID string) error {
// 	id, err := u.idCodec.Decode(roleID)
// 	if err != nil {
// 		return ErrDecode
// 	}

// 	user, _ := u.userRepo.GetByID(ctx, id)
// 	if user == nil {
// 		return ErrUserNotFound
// 	}

// 	role, _ := u.roleRepo.GetByID(ctx, id)
// 	if role == nil {
// 		return ErrRoleNotFound
// 	}

// 	return u.roleRepo.AssignToUser(ctx, userID, roleID)
// }

// =============== DELETE =================

func (u *roleUsecase) Delete(ctx context.Context, roleID string) error {
	id, err := u.idCodec.Decode(roleID)
	if err != nil {
		return ErrDecode
	}

	role, err := u.roleRepo.GetByID(ctx, id)
	if err != nil || role == nil {
		return ErrRoleNotFound
	}

	// OPTIONAL: guard jika role masih dipakai user
	used, err := u.roleRepo.IsRoleUsed(ctx, id)
	if err != nil {
		return err
	}
	if used {
		return ErrRoleInUse
	}

	return u.roleRepo.Delete(ctx, id)
}
