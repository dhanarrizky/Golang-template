package auth

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/repository"

	"gorm.io/gorm"
	dbctx "github.com/dhanarrizky/Golang-template/pkg/database"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(
	ctx context.Context,
	token *entities.RefreshToken,
) error {
	db := dbctx.GetDB(ctx, r.db)
	return db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepository) GetByTokenHash(
	ctx context.Context,
	hash string,
) (*entities.RefreshToken, error) {

	var token entities.RefreshToken
	db := dbctx.GetDB(ctx, r.db)

	err := db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&token).Error

	return &token, err
}

func (r *refreshTokenRepository) Revoke(
	ctx context.Context,
	id uint,
) error {

	db := dbctx.GetDB(ctx, r.db)

	return db.WithContext(ctx).
		Model(&entities.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}

func (r *refreshTokenRepository) RevokeByFamily(
	ctx context.Context,
	familyID uint,
) error {

	db := dbctx.GetDB(ctx, r.db)

	return db.WithContext(ctx).
		Model(&entities.RefreshToken{}).
		Where("family_id = ?", familyID).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	db := dbctx.GetDB(ctx, r.db)

	return db.WithContext(ctx).
		Where("expires_at < NOW()").
		Delete(&entities.RefreshToken{}).Error
}
