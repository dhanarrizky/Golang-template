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

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) ports.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(
	ctx context.Context,
	token *domain.RefreshToken,
) error {

	m := mapper.ToModelRefreshToken(token)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *refreshTokenRepository) GetByTokenHash(
	ctx context.Context,
	hash string,
) (*domain.RefreshToken, error) {

	var m model.RefreshToken

	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainRefreshToken(&m), nil
}

func (r *refreshTokenRepository) Revoke(
	ctx context.Context,
	id uint64,
) error {

	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", &now).Error
}

func (r *refreshTokenRepository) RevokeByFamily(
	ctx context.Context,
	familyID uint64,
) error {

	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("family_id = ?", familyID).
		Update("revoked_at", &now).Error
}

func (r *refreshTokenRepository) DeleteExpired(
	ctx context.Context,
) error {

	now := time.Now()

	return r.db.WithContext(ctx).
		Where("expires_at < ?", now).
		Delete(&model.RefreshToken{}).Error
}
