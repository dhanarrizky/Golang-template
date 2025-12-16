package repositories

import (
	"context"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	mapper "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/mappers/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) ports.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(
	ctx context.Context,
	token *domain.PasswordResetToken,
) error {

	m := mapper.ToModelPasswordResetToken(token)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *passwordResetTokenRepository) GetByTokenHash(
	ctx context.Context,
	hash string,
) (*domain.PasswordResetToken, error) {

	var m model.PasswordResetToken

	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainPasswordResetToken(&m), nil
}

func (r *passwordResetTokenRepository) MarkUsed(
	ctx context.Context,
	id uint64,
) error {

	return r.db.WithContext(ctx).
		Model(&model.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used", true).Error
}
