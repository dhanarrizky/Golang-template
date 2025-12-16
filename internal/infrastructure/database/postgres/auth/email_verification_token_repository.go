package repositories

import (
	"context"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/database/mappers/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
	"gorm.io/gorm"
)

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

func NewEmailVerificationTokenRepository(db *gorm.DB) ports.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(
	ctx context.Context,
	token *domain.EmailVerificationToken,
) error {

	modelToken := auth.ToModelEmailVerificationToken(token)

	return r.db.WithContext(ctx).Create(modelToken).Error
}

func (r *emailVerificationTokenRepository) GetByTokenHash(
	ctx context.Context,
	hash string,
) (*domain.EmailVerificationToken, error) {

	var modelToken model.EmailVerificationToken

	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&modelToken).Error
	if err != nil {
		return nil, err
	}

	return auth.ToDomainEmailVerificationToken(&modelToken), nil
}

func (r *emailVerificationTokenRepository) DeleteByUser(
	ctx context.Context,
	userID uint64,
) error {

	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.EmailVerificationToken{}).Error
}
