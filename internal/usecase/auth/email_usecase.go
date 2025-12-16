package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/dhanarrizky/Golang-template/internal/domain"
	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrVerificationTokenInvalid = errors.New("invalid or expired verification token")
	ErrEmailAlreadyVerified     = errors.New("email already verified")
)

type EmailUsecase interface {
	Verify(ctx context.Context, plainToken string) error
	Resend(ctx context.Context, email string) error
}

type emailUsecase struct {
	userRepo        ports.UserRepository
	tokenRepo       ports.EmailVerificationTokenRepository
	passwordHasher  ports.PasswordHasher
	mailer          ports.Mailer
	tokenExpiry     time.Duration
}

func NewEmailUsecase(
	userRepo ports.UserRepository,
	tokenRepo ports.EmailVerificationTokenRepository,
	passwordHasher ports.PasswordHasher,
	mailer ports.Mailer,
	tokenExpiry time.Duration,
) EmailUsecase {
	return &emailUsecase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		passwordHasher: passwordHasher,
		mailer:         mailer,
		tokenExpiry:    tokenExpiry,
	}
}

// ================= VERIFY =================

func (u *emailUsecase) Verify(ctx context.Context, plainToken string) error {
	hashed := u.passwordHasher.Hash(plainToken)

	token, err := u.tokenRepo.FindByToken(ctx, hashed)
	if err != nil || token == nil {
		return ErrVerificationTokenInvalid
	}

	if token.Used || time.Now().After(token.ExpiresAt) {
		return ErrVerificationTokenInvalid
	}

	user, err := u.userRepo.FindByID(ctx, token.UserID)
	if err != nil || user == nil {
		return ErrVerificationTokenInvalid
	}

	if user.EmailVerified {
		return ErrEmailAlreadyVerified
	}

	// Mark verified
	if err := u.userRepo.MarkEmailVerified(ctx, user.ID); err != nil {
		return err
	}

	// Mark token used + cleanup old tokens
	u.tokenRepo.MarkAsUsed(ctx, hashed)
	u.tokenRepo.DeleteAllByUser(ctx, user.ID)

	return nil
}

// ================= RESEND =================

func (u *emailUsecase) Resend(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil // silent
	}

	if user.EmailVerified {
		return nil // silent
	}

	// Generate secure token
	raw := make([]byte, 48)
	rand.Read(raw)

	plain := base64.URLEncoding.EncodeToString(raw)
	hashed := u.passwordHasher.Hash(plain)

	token := domain.EmailVerificationToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     hashed,
		ExpiresAt: time.Now().Add(u.tokenExpiry),
		CreatedAt: time.Now(),
	}

	u.tokenRepo.DeleteAllByUser(ctx, user.ID)
	u.tokenRepo.Create(ctx, token)

	link := "https://yourapp.com/verify-email?token=" + plain
	u.mailer.Send(ctx, user.Email, "Verify your email", "Click here: "+link)

	return nil
}
