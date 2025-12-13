package postgres

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
	"github.com/dhanarrizky/Golang-template/internal/repository"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) repository.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, token *entities.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *passwordResetTokenRepository) GetByTokenHash(ctx context.Context, hash string) (*entities.PasswordResetToken, error) {
	var token entities.PasswordResetToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&token).Error
	return &token, err
}

func (r *passwordResetTokenRepository) MarkUsed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entities.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used", true).Error
}
