package http

import (
	"github.com/go-playground/validator/v10"

	"github.com/dhanarrizky/Golang-template/internal/config"
	authUC "github.com/dhanarrizky/Golang-template/internal/usecase/auth"
	roleUC "github.com/dhanarrizky/Golang-template/internal/usecase/roles"
	userUC "github.com/dhanarrizky/Golang-template/internal/usecase/user"
)

type RouteDeps struct {
	Validator *validator.Validate
	Config    *config.Config

	// Auth usecases
	EmailUC    authUC.EmailUsecase
	LoginUC    authUC.LoginUsecase
	PasswordUC authUC.PasswordUsecase
	SessionUC  authUC.SessionUsecase
	TokenUC    authUC.TokenUsecase
	RoleUC     roleUC.RoleUsecase
	UserUC     userUC.UserUsecase
}
