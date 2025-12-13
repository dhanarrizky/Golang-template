package repository

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type RefreshTokenFamilyRepository interface {
	CreateFamily(
		ctx context.Context,
		family *entities.RefreshTokenFamily,
	) error

	GetFamilyByID(
		ctx context.Context,
		id uint,
	) (*entities.RefreshTokenFamily, error)

	RevokeFamily(
		ctx context.Context,
		id uint,
	) error
}
