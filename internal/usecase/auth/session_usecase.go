package auth

import (
	"context"
	"errors"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrSessionNotFound        = errors.New("session not found")
	ErrCannotRevokeOwnSession = errors.New("cannot revoke current session")
)

type SessionUsecase interface {
	List(ctx context.Context, userID string) ([]domain.UserSession, error)
	Revoke(ctx context.Context, userID, sessionID, currentFamilyID string) error
}

type sessionUsecase struct {
	sessionRepo ports.UserSessionRepository
	refreshRepo ports.RefreshTokenRepository
}

func NewSessionUsecase(
	sessionRepo ports.UserSessionRepository,
	refreshRepo ports.RefreshTokenRepository,
) SessionUsecase {
	return &sessionUsecase{
		sessionRepo: sessionRepo,
		refreshRepo: refreshRepo,
	}
}

// ================= LIST =================

func (u *sessionUsecase) List(
	ctx context.Context,
	userID string,
) ([]domain.UserSession, error) {

	return u.sessionRepo.FindByUser(ctx, userID)
}

// ================= REVOKE =================

func (u *sessionUsecase) Revoke(
	ctx context.Context,
	userID, sessionID, currentFamilyID string,
) error {

	session, err := u.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if session == nil || session.UserID != userID {
		return ErrSessionNotFound
	}

	// Tidak boleh revoke session sendiri
	if session.FamilyID == currentFamilyID {
		return ErrCannotRevokeOwnSession
	}

	// 1. Revoke refresh token family
	if err := u.refreshRepo.RevokeFamily(ctx, session.FamilyID); err != nil {
		return err
	}

	// 2. Revoke session
	return u.sessionRepo.RevokeByFamily(ctx, session.FamilyID)
}
