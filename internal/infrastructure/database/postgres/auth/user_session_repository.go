package repositories

import (
	"context"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	mapper "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/mappers/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/users"
	"gorm.io/gorm"
)

type userSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) ports.UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r *userSessionRepository) Create(
	ctx context.Context,
	session *domain.UserSession,
) error {

	m := mapper.ToModelUserSession(session)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *userSessionRepository) UpdateLastSeen(
	ctx context.Context,
	id uint64,
) error {

	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("id = ?", id).
		Update("last_seen_at", &now).Error
}

func (r *userSessionRepository) Logout(
	ctx context.Context,
	id uint64,
) error {

	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("id = ?", id).
		Update("logout_at", &now).Error
}

func (r *userSessionRepository) GetActiveSessions(
	ctx context.Context,
	userID uint64,
) ([]*domain.UserSession, error) {

	var models []model.UserSession

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND logout_at IS NULL", userID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	sessions := make([]*domain.UserSession, 0, len(models))
	for i := range models {
		sessions = append(sessions, mapper.ToDomainUserSession(&models[i]))
	}

	return sessions, nil
}
