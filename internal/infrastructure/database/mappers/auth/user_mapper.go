package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

func ToDomainUser(m *model.User) *domain.User {
	if m == nil {
		return nil
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &domain.User{
		ID:            m.ID,
		Email:         m.Email,
		EmailVerified: m.EmailVerified,
		PasswordHash:  m.PasswordHash,
		Name:          m.Name,
		RoleID:        m.RoleID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func ToModelUser(d *domain.User) *model.User {
	if d == nil {
		return nil
	}

	m := &model.User{
		ID:            d.ID,
		Email:         d.Email,
		EmailVerified: d.EmailVerified,
		PasswordHash:  d.PasswordHash,
		Name:          d.Name,
		RoleID:        d.RoleID,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}

	if d.DeletedAt != nil {
		m.DeletedAt = gorm.DeletedAt{
			Time:  *d.DeletedAt,
			Valid: true,
		}
	}

	return m
}
