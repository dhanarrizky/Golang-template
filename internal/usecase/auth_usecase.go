package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/repository"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
)

type UserAuthUsecase struct {
	UserRepo repository.UserRepository
}

func NewUserAuthUsecase(repo repository.UserRepository) *UserAuthUsecase {
	return &UserAuthUsecase{
		UserRepo: repo,
	}
}

func (u *UserAuthUsecase) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := u.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// Convert uint → string
	return strconv.FormatUint(uint64(user.ID), 10), nil
}

func (u *UserAuthUsecase) GetByID(ctx context.Context, id string) (*entities.User, error) {
	// Convert string → uint
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	return u.UserRepo.GetByID(ctx, uint(parsedID))
}
