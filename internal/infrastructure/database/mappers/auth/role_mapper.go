package auth

import (
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
)

func ToDomainRole(m *model.Role) *domain.Role {
	if m == nil {
		return nil
	}

	return &domain.Role{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}

func ToModelRole(d *domain.Role) *model.Role {
	if d == nil {
		return nil
	}

	return &model.Role{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
	}
}
