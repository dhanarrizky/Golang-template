package user

import (
	"context"

	"github.com/dhanarrizky/Golang-template/internal/repository"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
)

type DeleteUserUsecase struct {
	repo repository.UserRepository
}

func NewDeleteUserUsecase(repo repository.UserRepository) *DeleteUserUsecase {
	return &DeleteUserUsecase{repo: repo}
}

func (uc *DeleteUserUsecase) Execute(ctx context.Context, id uint) error {
	// Business logic if any, e.g., check if exists first
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return utils.WrapError(err, "user not found")
	}

	return uc.repo.Delete(ctx, id)
}