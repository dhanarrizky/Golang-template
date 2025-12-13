package repository

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type UserSessionRepository interface {
	Create(
		ctx context.Context,
		s *entities.UserSession,
	) error

	UpdateLastSeen(
		ctx context.Context,
		id uint,
	) error

	Logout(
		ctx context.Context,
		id uint,
	) error

	GetActiveSessions(
		ctx context.Context,
		userID uint,
	) ([]*entities.UserSession, error)
}
