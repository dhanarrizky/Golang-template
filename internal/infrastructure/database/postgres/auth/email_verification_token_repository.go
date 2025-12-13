package auth

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/repository"
	"gorm.io/gorm"
)

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

func NewEmailVerificationTokenRepository(db *gorm.DB) repository.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(ctx context.Context, token *entities.EmailVerificationToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *emailVerificationTokenRepository) GetByTokenHash(ctx context.Context, hash string) (*entities.EmailVerificationToken, error) {
	var token entities.EmailVerificationToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&token).Error
	return &token, err
}

func (r *emailVerificationTokenRepository) DeleteByUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.EmailVerificationToken{}).Error
}
