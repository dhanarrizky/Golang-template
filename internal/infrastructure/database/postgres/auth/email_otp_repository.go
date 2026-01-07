package auth

import (
	"context"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	mapper "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/mappers/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/email"
	"gorm.io/gorm"
)

type emailOTPRepository struct {
	db *gorm.DB
}

func NewEmailOTPRepository(db *gorm.DB) ports.EmailOTPRepository {
	return &emailOTPRepository{db: db}
}

func (r *emailOTPRepository) Save(
	ctx context.Context,
	otp *domain.EmailOTP,
) error {

	m := mapper.ToModelEmailOTP(otp)

	return r.db.WithContext(ctx).Create(m).Error
}

func (r *emailOTPRepository) FindActiveByEmail(
	ctx context.Context,
	email string,
) (*domain.EmailOTP, error) {

	var m model.EmailOTP

	err := r.db.WithContext(ctx).
		Where(
			"email = ? AND expired_at > ?",
			email,
			time.Now(),
		).
		Order("created_at DESC").
		First(&m).
		Error

	if err != nil {
		return nil, err
	}

	return mapper.ToDomainEmailOTP(&m), nil
}

func (r *emailOTPRepository) DeleteByEmail(
	ctx context.Context,
	email string,
) error {

	return r.db.WithContext(ctx).
		Where("email = ?", email).
		Delete(&model.EmailOTP{}).
		Error
}
