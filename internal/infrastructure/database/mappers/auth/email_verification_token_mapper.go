package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

func ToDomainEmailVerificationToken(m *model.EmailVerificationToken) *domain.EmailVerificationToken {
	if m == nil {
		return nil
	}

	return &domain.EmailVerificationToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}

func ToModelEmailVerificationToken(e *domain.EmailVerificationToken) *model.EmailVerificationToken {
	if e == nil {
		return nil
	}

	return &model.EmailVerificationToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		ExpiresAt: e.ExpiresAt,
		CreatedAt: e.CreatedAt,
	}
}
