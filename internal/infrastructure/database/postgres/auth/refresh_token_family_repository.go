package repositories

import (
	"context"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	mapper "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/mappers/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	"gorm.io/gorm"
)

type refreshTokenFamilyRepository struct {
	db *gorm.DB
}

func NewRefreshTokenFamilyRepository(db *gorm.DB) ports.RefreshTokenFamilyRepository {
	return &refreshTokenFamilyRepository{db: db}
}

func (r *refreshTokenFamilyRepository) Create(
	ctx context.Context,
	family *domain.RefreshTokenFamily,
) error {

	m := mapper.ToModelRefreshTokenFamily(family)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *refreshTokenFamilyRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*domain.RefreshTokenFamily, error) {

	var m model.RefreshTokenFamily

	err := r.db.WithContext(ctx).
		First(&m, id).Error
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainRefreshTokenFamily(&m), nil
}

func (r *refreshTokenFamilyRepository) Revoke(
	ctx context.Context,
	id uint64,
) error {

	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.RefreshTokenFamily{}).
		Where("id = ?", id).
		Update("revoked_at", &now).Error
}
