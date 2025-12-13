package user

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/repository"
)

type GetUserUsecase struct {
	repo repository.UserRepository
}

func NewGetUserUsecase(repo repository.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{repo: repo}
}

func (uc *GetUserUsecase) Execute(ctx context.Context, id uint) (*entities.User, error) {
	// Business logic if any
	return uc.repo.GetByID(ctx, id)
}