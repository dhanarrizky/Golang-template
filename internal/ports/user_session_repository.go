package ports

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type UserSessionRepository interface {
	Create(ctx context.Context, session *auth.UserSession) error

	UpdateLastSeen(ctx context.Context, id uint64) error
	Logout(ctx context.Context, id uint64) error

	GetActiveSessions(ctx context.Context, userID uint64) ([]*auth.UserSession, error)
}
