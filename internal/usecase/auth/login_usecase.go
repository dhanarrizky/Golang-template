package auth

import (
	"context"
	"errors"

	// "github.com/dhanarrizky/Golang-template/internal/ports"
	authPorts "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	rolePorts "github.com/dhanarrizky/Golang-template/internal/ports/roles"
	userPorts "github.com/dhanarrizky/Golang-template/internal/ports/users"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTooManyAttempts    = errors.New("too many login attempts")
	ErrRoleNotFound       = errors.New("invalid role access")
	ErrAccountLocked      = errors.New("account locked")
)

type LoginResult struct {
	UserID        uint64
	Email         string
	Username      string
	Roles         string
	RolesID       uint64
	EmailVerified bool
}

type LoginUsecase interface {
	Login(ctx context.Context, identifier, password string) (*LoginResult, error)
	Logout(ctx context.Context, refreshToken string) error
}

type loginUsecase struct {
	userRepo         userPorts.UserRepository
	loginAttemptRepo authPorts.LoginAttemptRepository
	passwordHasher   userPorts.PasswordHasher
	roleRepo         rolePorts.RoleRepository
	tokenUsecase     TokenUsecase
}

func NewLoginUsecase(
	userRepo userPorts.UserRepository,
	loginAttemptRepo authPorts.LoginAttemptRepository,
	passwordHasher userPorts.PasswordHasher,
	roleRepo rolePorts.RoleRepository,
	tokenUsecase TokenUsecase,
) LoginUsecase {
	return &loginUsecase{
		userRepo:         userRepo,
		loginAttemptRepo: loginAttemptRepo,
		passwordHasher:   passwordHasher,
		tokenUsecase:     tokenUsecase,
	}
}

// ================= LOGIN =================

func (u *loginUsecase) Login(
	ctx context.Context,
	identifier, password string,
) (*LoginResult, error) {

	if u.loginAttemptRepo.IsRateLimited(ctx, identifier) {
		return nil, ErrTooManyAttempts
	}

	user, err := u.userRepo.GetByEmailOrUsername(ctx, identifier)
	if err != nil || user == nil {
		u.loginAttemptRepo.RecordFailedAttempt(ctx, identifier)
		return nil, ErrInvalidCredentials
	}

	if user.Locked {
		return nil, ErrAccountLocked
	}

	matched, shouldRehash, err :=
		u.passwordHasher.VerifyPassword([]byte(password), user.PasswordHash)

	if err != nil || !matched {
		u.loginAttemptRepo.RecordFailedAttempt(ctx, identifier)
		return nil, ErrInvalidCredentials
	}

	if shouldRehash {
		newHash, err := u.passwordHasher.HashPassword([]byte(password))
		if err == nil {
			_ = u.userRepo.UpdatePassword(ctx, user.ID, newHash)
		}
	}

	u.loginAttemptRepo.ResetAttempts(ctx, identifier)

	role, err := u.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil || role == nil {
		return nil, ErrRoleNotFound
	}

	return &LoginResult{
		UserID:        user.ID,
		Email:         user.Email,
		Username:      user.Username,
		Roles:         role.Name,
		RolesID:       role.ID,
		EmailVerified: user.EmailVerified,
	}, nil
}

// ================= LOGOUT =================

func (u *loginUsecase) Logout(ctx context.Context, refreshToken string) error {
	return u.tokenUsecase.Revoke(ctx, refreshToken)
}
