package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

// ================= TO DOMAIN =================
func ToDomainLoginAttempt(m *model.LoginAttempt) *domain.LoginAttempt {
	if m == nil {
		return nil
	}

	return &domain.LoginAttempt{
		ID:        m.ID,
		Identity:  m.Identifier,
		UserID:    m.UserID,
		IPAddress: m.IPAddress,
		UserAgent: m.UserAgent,
		Success:   m.Success,
		Reason:    m.FailureReason,
		CreatedAt: m.CreatedAt,
	}
}

// ================= TO MODEL =================

// ================= TO MODEL =================
func ToModelLoginAttempt(d *domain.LoginAttempt) *model.LoginAttempt {
	if d == nil {
		return nil
	}

	return &model.LoginAttempt{
		// ID biasanya auto-increment → isi hanya jika perlu (misal import data)
		ID:            d.ID,
		Identifier:    d.Identity,
		UserID:        d.UserID,
		IPAddress:     d.IPAddress,
		UserAgent:     d.UserAgent,
		Success:       d.Success,
		FailureReason: d.Reason,
		// CreatedAt → biarkan GORM
	}
}
