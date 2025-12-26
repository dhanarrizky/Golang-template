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

func (r *userRepository) GetByEmailOrUsername(
	ctx context.Context,
	identifier string,
) (*domain.User, error) {

	var m model.User

	err := r.db.WithContext(ctx).
		Where(
			"email = ? OR username = ?",
			identifier,
			identifier,
		).
		First(&m).Error

	if err != nil {
		return nil, err
	}

	return mapper.ToDomainUser(&m), nil
}

func (r *userRepository) GetList(
	ctx context.Context,
) ([]*domain.User, error) {

	var models []model.User

	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, 0, len(models))
	for _, m := range models {
		users = append(users, mapper.ToDomainUser(&m))
	}

	return users, nil
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

func (r *userRepository) UpdatePassword(
	ctx context.Context,
	id uint64,
	hashedPassword string,
) error {

	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_hash": hashedPassword,
			"updated_at":    time.Now(),
		}).Error
}

func (r *userRepository) UpdateUsername(
	ctx context.Context,
	id uint64,
	username string,
) error {

	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_hash": username,
			"updated_at":    time.Now(),
		}).Error
}

func (r *userRepository) ExistsByUsernameExceptID(
	ctx context.Context,
	username string,
	exceptID uint64,
) (bool, error) {

	var exists bool

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("1").
		Where("username = ? AND id <> ?", username, exceptID).
		Limit(1).
		Scan(&exists).Error

	return exists, err
}

func (r *userRepository) ExistsByEmailExceptID(
	ctx context.Context,
	email string,
	exceptID uint64,
) (bool, error) {

	var exists bool

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("1").
		Where("email = ? AND id <> ?", email, exceptID).
		Limit(1).
		Scan(&exists).Error

	return exists, err
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
