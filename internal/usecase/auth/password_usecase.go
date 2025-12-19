package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrResetTokenInvalid    = errors.New("invalid or expired reset token")
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
	userRepo       ports.UserRepository
	resetTokenRepo ports.PasswordResetTokenRepository
	refreshRepo    ports.RefreshTokenRepository
	sessionRepo    ports.UserSessionRepository
	passwordHasher ports.PasswordHasher // interface untuk hash & compare (bisa bcrypt wrapper)
	mailer         ports.Mailer
	resetTokenExp  time.Duration
}

func NewPasswordUsecase(
	userRepo ports.UserRepository,
	resetTokenRepo ports.PasswordResetTokenRepository,
	refreshRepo ports.RefreshTokenRepository,
	sessionRepo ports.UserSessionRepository,
	passwordHasher ports.PasswordHasher,
	mailer ports.Mailer,
	resetTokenExp time.Duration,
) PasswordUsecase {
	return &passwordUsecase{
		userRepo:       userRepo,
		resetTokenRepo: resetTokenRepo,
		refreshRepo:    refreshRepo,
		sessionRepo:    sessionRepo,
		passwordHasher: passwordHasher,
		mailer:         mailer,
		resetTokenExp:  resetTokenExp,
	}
}

// ================= FORGOT PASSWORD =================
func (u *passwordUsecase) Forgot(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil // silent untuk security
	}

	// Generate secure random token
	raw := make([]byte, 48)
	if _, err := rand.Read(raw); err != nil {
		return err
	}
	plainToken := base64.URLEncoding.EncodeToString(raw)
	hashedToken := u.passwordHasher.Hash(plainToken)

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

	// Buat link reset (ganti dengan domain app kamu)
	resetLink := "https://yourapp.com/reset-password?token=" + plainToken

	// Kirim email
	if err := u.mailer.Send(ctx, user.Email, "Reset Password", "Klik link untuk reset password: "+resetLink); err != nil {
		// Log error tapi jangan fail request
	}

	return nil
}

// ================= RESET PASSWORD =================
func (u *passwordUsecase) Reset(ctx context.Context, plainToken, newPassword string) error {
	hashedToken := u.passwordHasher.Hash(plainToken)

	token, err := u.resetTokenRepo.FindByToken(ctx, hashedToken)
	if err != nil || token == nil || token.Used || time.Now().After(token.ExpiresAt) {
		return ErrResetTokenInvalid
	}

	// Hash new password
	hashedNew := u.passwordHasher.Hash(newPassword)

	if err := u.userRepo.UpdatePassword(ctx, token.UserID, hashedNew); err != nil {
		return err
	}

	// Mark token as used
	if err := u.resetTokenRepo.MarkAsUsed(ctx, token.ID); err != nil {
		// log only
	}

	// Revoke all sessions & refresh tokens (force logout everywhere)
	u.refreshRepo.RevokeAllByUser(ctx, token.UserID)
	u.sessionRepo.RevokeAllByUser(ctx, token.UserID)

	return nil
}

// ================= CHANGE PASSWORD (logged in user) =================
func (u *passwordUsecase) Change(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !u.passwordHasher.Compare(currentPassword, user.Password) {
		return ErrCurrentPasswordWrong
	}

	// Prevent same password
	if u.passwordHasher.Compare(newPassword, user.Password) {
		return ErrPasswordSameAsOld
	}

	hashedNew := u.passwordHasher.Hash(newPassword)
	return u.userRepo.UpdatePassword(ctx, userID, hashedNew)
}
