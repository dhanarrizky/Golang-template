package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

func ToDomainRefreshTokenFamily(m *model.RefreshTokenFamily) *domain.RefreshTokenFamily {
	if m == nil {
		return nil
	}

	return &domain.RefreshTokenFamily{
		ID:        m.ID,
		UserID:    m.UserID,
		RevokedAt: m.RevokedAt,
		CreatedAt: m.CreatedAt,
	}
}

func ToModelRefreshTokenFamily(d *domain.RefreshTokenFamily) *model.RefreshTokenFamily {
	if d == nil {
		return nil
	}

	return &model.RefreshTokenFamily{
		ID:        d.ID,
		UserID:    d.UserID,
		RevokedAt: d.RevokedAt,
		CreatedAt: d.CreatedAt,
	}
}
