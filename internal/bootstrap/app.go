package bootstrap

import (
	"log"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/config"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http"
	authRepo "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres/auth"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/security"
	authUC "github.com/dhanarrizky/Golang-template/internal/usecase/auth"
	roleUC "github.com/dhanarrizky/Golang-template/internal/usecase/roles"
	userUC "github.com/dhanarrizky/Golang-template/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func InitHTTPApp(cfg *config.Config) *gin.Engine {
	// =====================
	// Infrastructure
	// =====================
	db := InitDatabase(cfg)
	// redis := InitRedis(cfg)
	tokenHasher := InitTokenHasher(cfg)
	passwordHasher := security.NewPasswordHasher(&security.PasswordConfig{
		Memory:               cfg.Password.Memory,
		Iterations:           cfg.Password.Iterations,
		Parallelism:          cfg.Password.Parallelism,
		SaltLength:           cfg.Password.SaltLength,
		KeyLength:            cfg.Password.KeyLength,
		Peppers:              cfg.Password.Peppers,
		CurrentPepperVersion: cfg.Password.CurrentPepperVersion,
	})
	// jwtSigner := authinfra.NewJWTSigner(cfg)
	accessExp, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		log.Fatalf("invalid JWT_EXPIRES_IN: %v", err)
	}

	refreshExp, err := time.ParseDuration(cfg.JWTRefreshExpiresIn)
	if err != nil {
		log.Fatalf("invalid JWT_REFRESH_EXPIRES_IN: %v", err)
	}

	jwtSigner := security.NewJWTSigner(
		cfg.JWTSecret,
		cfg.JWTSecret,
		accessExp,
		refreshExp,
	)

	// =====================
	// Repository
	// =====================

	// Auth
	loginAttemptRepo := authRepo.NewLoginAttemptRepository(db)
	passwordResetTokenRepo := authRepo.NewPasswordResetTokenRepository(db)
	// refreshTokenFamilyRepo := authRepo.NewRefreshTokenFamilyRepository(db)
	refreshTokenRepo := authRepo.NewRefreshTokenRepository(db)
	roleRepo := authRepo.NewRoleRepository(db)
	userRepo := authRepo.NewUserRepository(db)
	sessionRepo := authRepo.NewUserSessionRepository(db)

	// =====================
	// Usecases
	// =====================

	loginUC := authUC.NewLoginUsecase(
		userRepo,
		loginAttemptRepo,
		sessionRepo,
		refreshTokenRepo,
		passwordHasher,
		cfg.JWTSecret,
		accessExp,
		refreshExp,
	)

	passwordUC := authUC.NewPasswordUsecase(
		userRepo,
		passwordResetTokenRepo,
		refreshTokenRepo,
		sessionRepo,
		passwordHasher,
		mailer,
		accessExp,
	)

	sessionUC := authUC.NewSessionUsecase(
		sessionRepo,
		refreshTokenRepo,
	)

	roleUC := roleUC.NewRoleUsecase(
		roleRepo,
		userRepo,
	)

	tokenUC := authUC.NewTokenUsecase(
		refreshTokenRepo,
		sessionRepo,
		cfg.JWTSecret,
		accessExp,
		refreshExp,
	)

	userUC := userUC.NewUserUsecase(
		userRepo,
		sessionRepo,
		refreshTokenRepo)

	// =====================
	// HTTP Router
	// =====================
	router := gin.New()

	http.RegisterRoutes(
		router,
		http.RouteDeps{
			JwtSigner: &jwtSigner,
			Validator: validator.New(),
			Config:    cfg,

			LoginUC:    loginUC,
			PasswordUC: passwordUC,
			SessionUC:  sessionUC,
			TokenUC:    tokenUC,
			RoleUC:     roleUC,
			UserUC:     userUC,
		},
	)

	return router
}
