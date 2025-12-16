package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

func ToDomainUserSession(m *model.UserSession) *domain.UserSession {
	if m == nil {
		return nil
	}

	return &domain.UserSession{
		ID:         m.ID,
		UserID:     m.UserID,
		IPAddress:  m.IPAddress,
		UserAgent:  m.UserAgent,
		LoginAt:    m.LoginAt,
		LastSeenAt: m.LastSeenAt,
		LogoutAt:   m.LogoutAt,
	}
}

func ToModelUserSession(d *domain.UserSession) *model.UserSession {
	if d == nil {
		return nil
	}

	return &model.UserSession{
		ID:         d.ID,
		UserID:     d.UserID,
		IPAddress:  d.IPAddress,
		UserAgent:  d.UserAgent,
		LoginAt:    d.LoginAt,
		LastSeenAt: d.LastSeenAt,
		LogoutAt:   d.LogoutAt,
	}
}
