package repositories

import (
	"context"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	mapper "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/mappers/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*domain.User, error) {

	var m model.User

	err := r.db.WithContext(ctx).
		First(&m, id).Error
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainUser(&m), nil
}

func (r *userRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {

	var m model.User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainUser(&m), nil
}

func (r *userRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {

	m := mapper.ToModelUser(user)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *userRepository) Update(
	ctx context.Context,
	user *domain.User,
) error {

	m := mapper.ToModelUser(user)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *userRepository) SoftDelete(
	ctx context.Context,
	id uint64,
) error {

	now := time.Now()

	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("deleted_at", &now).Error
}
