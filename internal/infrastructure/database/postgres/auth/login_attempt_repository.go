package auth

import (
	"context"
	"time"

	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	"gorm.io/gorm"
)

type loginAttemptRepository struct {
	db *gorm.DB
}

const (
	maxAttempts   = 5
	windowMinutes = 15
)

func NewLoginAttemptRepository(db *gorm.DB) ports.LoginAttemptRepository {
	return &loginAttemptRepository{db: db}
}

func (r *loginAttemptRepository) IsRateLimited(
	ctx context.Context,
	identifier string,
) bool {

	var count int64
	since := time.Now().Add(-windowMinutes * time.Minute)

	err := r.db.WithContext(ctx).
		Model(&model.LoginAttempt{}).
		Where(
			"identifier = ? AND success = FALSE AND created_at >= ?",
			identifier,
			since,
		).
		Count(&count).Error

	if err != nil {
		// Fail-safe: kalau DB error, anggap limited
		return true
	}

	return count >= maxAttempts
}

func (r *loginAttemptRepository) RecordFailedAttempt(
	ctx context.Context,
	identifier string,
) error {

	attempt := &model.LoginAttempt{
		Identifier: identifier,
		Success:    false,
		CreatedAt:  time.Now(),
	}

	return r.db.WithContext(ctx).Create(attempt).Error
}

func (r *loginAttemptRepository) ResetAttempts(
	ctx context.Context,
	identifier string,
) error {

	return r.db.WithContext(ctx).
		Where("identifier = ?", identifier).
		Delete(&model.LoginAttempt{}).
		Error
}
