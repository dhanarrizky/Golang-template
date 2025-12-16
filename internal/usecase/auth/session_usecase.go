package auth

import (
	"context"
	"errors"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrSessionNotFound = errors.New("session not found")
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

func (u *sessionUsecase) List(ctx context.Context, userID string) ([]domain.UserSession, error) {
	return u.sessionRepo.FindByUser(ctx, userID)
}

// ================= REVOKE =================

func (u *sessionUsecase) Revoke(ctx context.Context, userID, sessionID, currentFamilyID string) error {
	session, err := u.sessionRepo.FindByID(ctx, sessionID)
	if err != nil || session == nil {
		return ErrSessionNotFound
	}

	if session.UserID != userID {
		return ErrSessionNotFound
	}

	// Jangan revoke session sendiri lewat endpoint ini
	if session.FamilyID == currentFamilyID {
		return ErrCannotRevokeOwnSession
	}

	// Revoke refresh token family + session
	u.refreshRepo.RevokeFamily(ctx, session.FamilyID)
	return u.sessionRepo.RevokeByFamily(ctx, session.FamilyID)
}
