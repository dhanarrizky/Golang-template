package auth

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/repository"
	"gorm.io/gorm"
)

type loginAttemptRepository struct {
	db *gorm.DB
}

func NewLoginAttemptRepository(db *gorm.DB) repository.LoginAttemptRepository {
	return &loginAttemptRepository{db: db}
}

func (r *loginAttemptRepository) LogAttempt(ctx context.Context, attempt *entities.LoginAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

func (r *loginAttemptRepository) CountFailedByIP(ctx context.Context, ip string, minutes int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.LoginAttempt{}).
		Where("ip_address = ? AND success = FALSE AND created_at > NOW() - INTERVAL '? minutes'", ip, minutes).
		Count(&count).Error
	return int(count), err
}

func (r *loginAttemptRepository) CountFailedByEmail(ctx context.Context, email string, minutes int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.LoginAttempt{}).
		Where("email = ? AND success = FALSE AND created_at > NOW() - INTERVAL '? minutes'", email, minutes).
		Count(&count).Error
	return int(count), err
}
