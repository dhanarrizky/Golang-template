package auth

import (
	"context"
	"errors"

	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTooManyAttempts    = errors.New("too many login attempts")
	ErrAccountLocked      = errors.New("account locked")
)

type LoginResult struct {
	UserID        string
	Email         string
	Username      string
	Roles         []string
	EmailVerified bool
}

type LoginUsecase interface {
	Login(ctx context.Context, identifier, password string) (*LoginResult, error)
	Logout(ctx context.Context, refreshToken string) error
}

type loginUsecase struct {
	userRepo         ports.UserRepository
	loginAttemptRepo ports.LoginAttemptRepository
	passwordHasher   ports.PasswordHasher
	tokenUsecase     TokenUsecase
}

func NewLoginUsecase(
	userRepo ports.UserRepository,
	loginAttemptRepo ports.LoginAttemptRepository,
	passwordHasher ports.PasswordHasher,
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

	user, err := u.userRepo.FindByEmailOrUsername(ctx, identifier)
	if err != nil || user == nil {
		u.loginAttemptRepo.RecordFailedAttempt(ctx, identifier)
		return nil, ErrInvalidCredentials
	}

	if user.Locked {
		return nil, ErrAccountLocked
	}

	if !u.passwordHasher.Compare(password, user.HashedPassword) {
		u.loginAttemptRepo.RecordFailedAttempt(ctx, identifier)
		return nil, ErrInvalidCredentials
	}

	u.loginAttemptRepo.ResetAttempts(ctx, identifier)

	return &LoginResult{
		UserID:        user.ID,
		Email:         user.Email,
		Username:      user.Username,
		Roles:         user.Roles,
		EmailVerified: user.EmailVerified,
	}, nil
}

// ================= LOGOUT =================

func (u *loginUsecase) Logout(ctx context.Context, refreshToken string) error {
	return u.tokenUsecase.Revoke(ctx, refreshToken)
}
