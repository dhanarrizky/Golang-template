package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	authPorts "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	otherPorts "github.com/dhanarrizky/Golang-template/internal/ports/others"
	userPorts "github.com/dhanarrizky/Golang-template/internal/ports/users"
)

var (
	ErrDecode               = errors.New("internal server decode")
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
	codecHasher       otherPorts.PublicIDCodec
	tokenGenerator    otherPorts.TokenGenerator
	hmacTokenVerifier otherPorts.TokenVerifier
	resetTokenExp     time.Duration
	idCodec           otherPorts.PublicIDCodec
}

func NewPasswordUsecase(
	userRepo userPorts.UserRepository,
	resetTokenRepo authPorts.PasswordResetTokenRepository,
	refreshRepo authPorts.RefreshTokenRepository,
	refreshFamilyRepo authPorts.RefreshTokenFamilyRepository,
	sessionRepo userPorts.UserSessionRepository,
	passwordHasher userPorts.PasswordHasher,
	resetTokenExp time.Duration,
	idCodec otherPorts.PublicIDCodec,
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

	_, hashedToken, err := u.tokenGenerator.Generate()
	if err != nil {
		return err
	}

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
	id, err := u.idCodec.Decode(userID)
	if err != nil {
		return ErrDecode
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// Verify current password
	match, _, err := u.passwordHasher.VerifyPassword(
		[]byte(currentPassword),
		user.PasswordHash,
	)
	if err != nil {
		return err
	}

	if !match {
		return ErrCurrentPasswordWrong
	}

	// Prevent same password
	same, _, err := u.passwordHasher.VerifyPassword(
		[]byte(newPassword),
		user.PasswordHash,
	)
	if err != nil {
		return err
	}

	if same {
		return ErrPasswordSameAsOld
	}

	hashedNew, err := u.passwordHasher.HashPassword([]byte(newPassword))
	if err != nil {
		return ErrHashingPass
	}
	return u.userRepo.UpdatePassword(ctx, id, hashedNew)
}
