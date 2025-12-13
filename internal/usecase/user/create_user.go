package user

import (
	"context"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/repository"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
)

type CreateUserUsecase struct {
	repo repository.UserRepository
}

func NewCreateUserUsecase(repo repository.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{repo: repo}
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (uc *CreateUserUsecase) Execute(ctx context.Context, input CreateUserInput) (*entities.User, error) {
	user := &entities.User{
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, utils.WrapError(err, "validation failed")
	}

	// Additional business logic

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, utils.WrapError(err, "failed to create user")
	}

	return user, nil
}