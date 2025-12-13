package postgres

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
	"github.com/dhanarrizky/Golang-template/internal/repository"
	"gorm.io/gorm"
)

type userSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) repository.UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r *userSessionRepository) Create(ctx context.Context, s *entities.UserSession) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *userSessionRepository) UpdateLastSeen(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserSession{}).
		Where("id = ?", id).
		Update("last_seen_at", gorm.Expr("NOW()")).Error
}

func (r *userSessionRepository) Logout(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserSession{}).
		Where("id = ?", id).
		Update("logout_at", gorm.Expr("NOW()")).Error
}

func (r *userSessionRepository) GetActiveSessions(ctx context.Context, userID uint) ([]*entities.UserSession, error) {
	var sessions []*entities.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND logout_at IS NULL", userID).
		Find(&sessions).Error
	return sessions, err
}
