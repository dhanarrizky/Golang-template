package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

func ToDomainRefreshToken(m *model.RefreshToken) *domain.RefreshToken {
	if m == nil {
		return nil
	}

	return &domain.RefreshToken{
		ID:        m.ID,
		UserID:    m.UserID,
		FamilyID:  m.FamilyID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
		RevokedAt: m.RevokedAt,
		IPAddress: m.IPAddress,
		UserAgent: m.UserAgent,
	}
}

func ToModelRefreshToken(d *domain.RefreshToken) *model.RefreshToken {
	if d == nil {
		return nil
	}

	return &model.RefreshToken{
		ID:        d.ID,
		UserID:    d.UserID,
		FamilyID:  d.FamilyID,
		TokenHash: d.TokenHash,
		ExpiresAt: d.ExpiresAt,
		CreatedAt: d.CreatedAt,
		RevokedAt: d.RevokedAt,
		IPAddress: d.IPAddress,
		UserAgent: d.UserAgent,
	}
}
