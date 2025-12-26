package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	authPorts "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	userPorts "github.com/dhanarrizky/Golang-template/internal/ports/users"
)

var (
	ErrResetTokenInvalid    = errors.New("invalid or expired reset token")
	ErrUsernameNotFound     = errors.New("username not found")
	ErrHashingPass          = errors.New("internal server error hash password")
	ErrResetTokenUsed       = errors.New("reset token already used")
	ErrCurrentPasswordWrong = errors.New("current password is incorrect")
	ErrPasswordSameAsOld    = errors.New("new password cannot be the same as old password")
)

type PasswordUsecase interface {
	Forgot(ctx context.Context, email string) error
	Reset(ctx context.Context, token, newPassword string) error
	Change(ctx context.Context, userID, currentPassword, newPassword string) error
}

type passwordUsecase struct {
	userRepo          userPorts.UserRepository
	resetTokenRepo    authPorts.PasswordResetTokenRepository
	refreshRepo       authPorts.RefreshTokenRepository
	refreshFamilyRepo authPorts.RefreshTokenFamilyRepository
	sessionRepo       userPorts.UserSessionRepository
	passwordHasher    userPorts.PasswordHasher
	resetTokenExp     time.Duration
}

func NewPasswordUsecase(
	userRepo userPorts.UserRepository,
	resetTokenRepo authPorts.PasswordResetTokenRepository,
	refreshRepo authPorts.RefreshTokenRepository,
	refreshFamilyRepo authPorts.RefreshTokenFamilyRepository,
	sessionRepo userPorts.UserSessionRepository,
	passwordHasher userPorts.PasswordHasher,
	resetTokenExp time.Duration,
) PasswordUsecase {
	return &passwordUsecase{
		userRepo:       userRepo,
		resetTokenRepo: resetTokenRepo,
		refreshRepo:    refreshRepo,
		sessionRepo:    sessionRepo,
		passwordHasher: passwordHasher,
		resetTokenExp:  resetTokenExp,
	}
}

// ================= FORGOT PASSWORD =================
func (u *passwordUsecase) Forgot(ctx context.Context, email string) error {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil // silent untuk security
	}

	// Generate secure random token
	raw := make([]byte, 48)
	if _, err := rand.Read(raw); err != nil {
		return err
	}
	plainToken := base64.URLEncoding.EncodeToString(raw)
	hashedToken := u.passwordHasher.has(plainToken)

	// ID:        uuid.NewString(),
	// UserID:    user.ID,
	// Token:     hashedToken,
	// ExpiresAt: time.Now().Add(u.resetTokenExp),
	// CreatedAt: time.Now(),
	// Used:      false,
	resetToken := domain.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(u.resetTokenExp),
		Used:      false,
		CreatedAt: time.Now(),
	}

	if err := u.resetTokenRepo.Create(ctx, &resetToken); err != nil {
		return err
	}

	return nil
}

// ================= RESET PASSWORD =================
func (u *passwordUsecase) Reset(ctx context.Context, username, newPassword string) error {
	user, err := u.userRepo.GetByEmailOrUsername(ctx, username)
	if err != nil || user == nil {
		return ErrUsernameNotFound
	}

	// Hash new password
	hashedNew, err := u.passwordHasher.HashPassword([]byte(newPassword))
	if err != nil {
		return ErrHashingPass
	}

	if err := u.userRepo.UpdatePassword(ctx, user.ID, hashedNew); err != nil {
		return err
	}

	// Revoke all sessions & refresh tokens (force logout everywhere)
	families, err := u.refreshFamilyRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, family := range families {
		// revoke semua refresh token dalam family
		_ = u.refreshRepo.RevokeByFamily(ctx, family.ID)

		// revoke family itu sendiri
		_ = u.refreshFamilyRepo.Revoke(ctx, family.ID)
	}

	return nil
}

// ================= CHANGE PASSWORD (logged in user) =================
func (u *passwordUsecase) Change(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !u.passwordHasher.VerifyPassword(currentPassword, user.Password) {
		return ErrCurrentPasswordWrong
	}

	// Prevent same password
	if u.passwordHasher.VerifyPassword(newPassword, user.Password) {
		return ErrPasswordSameAsOld
	}

	hashedNew, err := u.passwordHasher.HashPassword([]byte(newPassword))
	if err != nil {
		return ErrHashingPass
	}
	return u.userRepo.UpdatePassword(ctx, userID, hashedNew)
}
