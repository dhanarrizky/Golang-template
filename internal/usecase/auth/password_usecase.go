package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/domain"
	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrResetTokenInvalid   = errors.New("invalid or expired reset token")
	ErrResetTokenUsed      = errors.New("reset token already used")
	ErrCurrentPasswordWrong = errors.New("current password is incorrect")
	ErrPasswordSameAsOld   = errors.New("new password cannot be the same as old")
)

type PasswordUsecase interface {
	Forgot(ctx context.Context, email string) error
	Reset(ctx context.Context, token, newPassword string) error
	Change(ctx context.Context, userID, currentPassword, newPassword string) error
}

type passwordUsecase struct {
	userRepo          ports.UserRepository
	resetTokenRepo    ports.PasswordResetTokenRepository
	refreshRepo       ports.RefreshTokenRepository
	sessionRepo       ports.UserSessionRepository
	passwordHasher    ports.PasswordHasher
	mailer            ports.Mailer // interface untuk kirim email
	resetTokenExp     time.Duration
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

// Forgot: buat token + kirim email
func (u *passwordUsecase) Forgot(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		// Silent: jangan kasih tahu kalau email tidak ada (security)
		return nil
	}

	// Generate secure random token
	rawToken := make([]byte, 48) // 48 byte → ~64 char base64
	_, err = rand.Read(rawToken)
	if err != nil {
		return err
	}
	plainToken := base64.URLEncoding.EncodeToString(rawToken)
	hashedToken := u.passwordHasher.Hash(plainToken) // hash sebelum simpan

	resetToken := domain.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(u.resetTokenExp),
		CreatedAt: time.Now(),
	}

	if err := u.resetTokenRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	// Kirim email (di real app: gunakan template + link frontend)
	resetLink := "https://yourapp.com/reset-password?token=" + plainToken
	err = u.mailer.Send(ctx, user.Email, "Password Reset", "Click here to reset: "+resetLink)
	if err != nil {
		// Log error, tapi jangan fail request (user experience)
	}

	return nil
}

// Reset: validasi token + update password + revoke all sessions
func (u *passwordUsecase) Reset(ctx context.Context, plainToken, newPassword string) error {
	hashedToken := u.passwordHasher.Hash(plainToken)

	token, err := u.resetTokenRepo.FindByToken(ctx, hashedToken)
	if err != nil || token == nil {
		return ErrResetTokenInvalid
	}

	if token.Used {
		return ErrResetTokenUsed
	}

	if time.Now().After(token.ExpiresAt) {
		return ErrResetTokenInvalid
	}

	// Update password
	hashedNew := u.passwordHasher.Hash(newPassword)
	if err := u.userRepo.UpdatePassword(ctx, token.UserID, hashedNew); err != nil {
		return err
	}

	// Mark token used
	u.resetTokenRepo.MarkAsUsed(ctx, hashedToken)

	// CRITICAL: Revoke all sessions & refresh tokens (paksa logout semua device)
	u.refreshRepo.RevokeAllByUser(ctx, token.UserID)
	u.sessionRepo.RevokeAllByUser(ctx, token.UserID)

	return nil
}

// Change: saat user sudah login
func (u *passwordUsecase) Change(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if !u.passwordHasher.Compare(currentPassword, user.HashedPassword) {
		return ErrCurrentPasswordWrong
	}

	if u.passwordHasher.Compare(newPassword, user.HashedPassword) {
		return ErrPasswordSameAsOld
	}

	hashedNew := u.passwordHasher.Hash(newPassword)
	return u.userRepo.UpdatePassword(ctx, userID, hashedNew)
}