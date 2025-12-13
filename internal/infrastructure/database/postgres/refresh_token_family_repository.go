package postgres

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
	"github.com/dhanarrizky/Golang-template/internal/repository"
	"gorm.io/gorm"
)

type refreshTokenFamilyRepository struct {
	db *gorm.DB
}

func NewRefreshTokenFamilyRepository(db *gorm.DB) repository.RefreshTokenFamilyRepository {
	return &refreshTokenFamilyRepository{db: db}
}

func (r *refreshTokenFamilyRepository) CreateFamily(ctx context.Context, family *entities.RefreshTokenFamily) error {
	return r.db.WithContext(ctx).Create(family).Error
}

func (r *refreshTokenFamilyRepository) GetFamilyByID(ctx context.Context, id uint) (*entities.RefreshTokenFamily, error) {
	var fam entities.RefreshTokenFamily
	err := r.db.WithContext(ctx).First(&fam, id).Error
	return &fam, err
}

func (r *refreshTokenFamilyRepository) RevokeFamily(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entities.RefreshTokenFamily{}).
		Where("id = ?", id).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}
