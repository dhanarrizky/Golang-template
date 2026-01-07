package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

// ================= TO DOMAIN =================
func ToDomainEmailOTP(m *model.EmailOTP) *domain.EmailOTP {
	if m == nil {
		return nil
	}

	return &domain.EmailOTP{
		ID:        m.ID,
		Email:     m.Email,
		OTPHash:   m.OTPHash,
		ExpiredAt: m.ExpiredAt,
		IPAddress: m.IPAddress,
		UserAgent: m.UserAgent,
		CreatedAt: m.CreatedAt,
	}
}

// ================= TO MODEL =================
func ToModelEmailOTP(d *domain.EmailOTP) *model.EmailOTP {
	if d == nil {
		return nil
	}

	return &model.EmailOTP{
		ID:        d.ID,
		Email:     d.Email,
		OTPHash:   d.OTPHash,
		ExpiredAt: d.ExpiredAt,
		IPAddress: d.IPAddress,
		UserAgent: d.UserAgent,
		// CreatedAt biarkan GORM
	}
}
