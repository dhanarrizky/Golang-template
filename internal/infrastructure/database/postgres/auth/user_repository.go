package auth

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/repository"

	"gorm.io/gorm"
	dbctx "github.com/dhanarrizky/Golang-template/pkg/database"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	db := dbctx.GetDB(ctx, r.db)

	err := db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	db := dbctx.GetDB(ctx, r.db)

	err := db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	return &user, err
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	db := dbctx.GetDB(ctx, r.db)
	return db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	db := dbctx.GetDB(ctx, r.db)
	return db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) SoftDelete(ctx context.Context, id uint) error {
	db := dbctx.GetDB(ctx, r.db)

	return db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
