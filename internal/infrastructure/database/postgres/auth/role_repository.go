package repositories

import (
	"context"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	mapper "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/mappers/auth"
	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/roles"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) ports.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*domain.Role, error) {

	var m model.Role

	err := r.db.WithContext(ctx).
		First(&m, id).Error
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainRole(&m), nil
}

func (r *roleRepository) GetByName(
	ctx context.Context,
	name string,
) (*domain.Role, error) {

	var m model.Role

	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&m).Error
	if err != nil {
		return nil, err
	}

	return mapper.ToDomainRole(&m), nil
}

func (r *roleRepository) List(
	ctx context.Context,
) ([]*domain.Role, error) {

	var models []model.Role

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	roles := make([]*domain.Role, 0, len(models))
	for i := range models {
		roles = append(roles, mapper.ToDomainRole(&models[i]))
	}

	return roles, nil
}

func (r *roleRepository) Create(
	ctx context.Context,
	role *domain.Role,
) error {

	m := mapper.ToModelRole(role)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *roleRepository) Update(
	ctx context.Context,
	role *domain.Role,
) error {

	m := mapper.ToModelRole(role)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Role{}).Error
}

func (r *roleRepository) IsRoleUsed(ctx context.Context, id uint64) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("role_id = ?", id).
		Count(&count).Error

	return count > 0, err
}
