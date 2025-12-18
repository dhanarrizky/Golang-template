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

type loginAttemptRepository struct {
	db *gorm.DB
}

func NewLoginAttemptRepository(db *gorm.DB) ports.LoginAttemptRepository {
	return &loginAttemptRepository{db: db}
}

func (r *loginAttemptRepository) LogAttempt(
	ctx context.Context,
	attempt *domain.LoginAttempt,
) error {

	m := mapper.ToModelLoginAttempt(attempt)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *loginAttemptRepository) CountFailedByIP(
	ctx context.Context,
	ip string,
	minutes int,
) (int, error) {

	var count int64
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)

	err := r.db.WithContext(ctx).
		Model(&model.LoginAttempt{}).
		Where("ip_address = ? AND success = FALSE AND created_at >= ?", ip, since).
		Count(&count).Error

	return int(count), err
}

func (r *loginAttemptRepository) CountFailedByEmail(
	ctx context.Context,
	email string,
	minutes int,
) (int, error) {

	var count int64
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)

	err := r.db.WithContext(ctx).
		Model(&model.LoginAttempt{}).
		Where("email = ? AND success = FALSE AND created_at >= ?", email, since).
		Count(&count).Error

	return int(count), err
}
