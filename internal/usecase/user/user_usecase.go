package user

import (
	"context"
	"errors"

	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUsernameTaken = errors.New("username already taken")
)

type UserUsecase interface {
	GetMe(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID, username string) error
	SoftDelete(ctx context.Context, userID string) error
}

type userUsecase struct {
	userRepo    ports.UserRepository
	sessionRepo ports.UserSessionRepository
	refreshRepo ports.RefreshTokenRepository
}

func NewUserUsecase(
	userRepo ports.UserRepository,
	sessionRepo ports.UserSessionRepository,
	refreshRepo ports.RefreshTokenRepository,
) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		refreshRepo: refreshRepo,
	}
}

// ================= GET ME =================

func (u *userUsecase) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// ================= UPDATE =================

func (u *userUsecase) UpdateProfile(ctx context.Context, userID, username string) error {
	if username == "" {
		return nil
	}

	exists, _ := u.userRepo.ExistsByUsername(ctx, username)
	if exists {
		return ErrUsernameTaken
	}

	return u.userRepo.UpdateUsername(ctx, userID, username)
}

// ================= SOFT DELETE =================

func (u *userUsecase) SoftDelete(ctx context.Context, userID string) error {
	// Soft delete user
	if err := u.userRepo.SoftDelete(ctx, userID); err != nil {
		return err
	}

	// Revoke ALL session & token (critical)
	u.refreshRepo.RevokeAllByUser(ctx, userID)
	u.sessionRepo.RevokeAllByUser(ctx, userID)

	return nil
}
