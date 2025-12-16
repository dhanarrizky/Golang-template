package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

func ToDomainLoginAttempt(m *model.LoginAttempt) *domain.LoginAttempt {
	if m == nil {
		return nil
	}

	return &domain.LoginAttempt{
		ID:        m.ID,
		Email:     m.Email,
		IPAddress: m.IPAddress,
		Success:   m.Success,
		UserAgent: m.UserAgent,
		CreatedAt: m.CreatedAt,
	}
}

func ToModelLoginAttempt(d *domain.LoginAttempt) *model.LoginAttempt {
	if d == nil {
		return nil
	}

	return &model.LoginAttempt{
		ID:        d.ID,
		Email:     d.Email,
		IPAddress: d.IPAddress,
		Success:   d.Success,
		UserAgent: d.UserAgent,
		CreatedAt: d.CreatedAt,
	}
}
