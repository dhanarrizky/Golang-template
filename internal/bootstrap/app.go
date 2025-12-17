package bootstrap

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/dhanarrizky/Golang-template/internal/config"
	http "github.com/dhanarrizky/Golang-template/internal/delivery/http"
	authRepo "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres/auth"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/mailer"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/security"
	authUC "github.com/dhanarrizky/Golang-template/internal/usecase/auth"
	roleUC "github.com/dhanarrizky/Golang-template/internal/usecase/roles"
	userUC "github.com/dhanarrizky/Golang-template/internal/usecase/user"
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

	// =====================
	// Repository
	// =====================

	// Auth
	emailRepo := authRepo.NewEmailVerificationTokenRepository(db)
	loginAttemptRepo := authRepo.NewLoginAttemptRepository(db)
	passwordResetTokenRepo := authRepo.NewPasswordResetTokenRepository(db)
	// refreshTokenFamilyRepo := authRepo.NewRefreshTokenFamilyRepository(db)
	refreshTokenRepo := authRepo.NewRefreshTokenRepository(db)
	roleRepo := authRepo.NewRoleRepository(db)
	userRepo := authRepo.NewUserRepository(db)
	sessionRepo := authRepo.NewUserSessionRepository(db)

	// =====================
	// Mailer
	// =====================
	mailer, err := mailer.NewSMTPMailer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFromName,
		cfg.SMTPFromAddress,
	)

	if err != nil {
		log.Fatal(err)
	}

	// =====================
	// Usecases
	// =====================

	// Auth

	accessExp, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		log.Fatalf("invalid JWT_EXPIRES_IN: %v", err)
	}

	refreshExp, err := time.ParseDuration(cfg.JWTRefreshExpiresIn)
	if err != nil {
		log.Fatalf("invalid JWT_REFRESH_EXPIRES_IN: %v", err)
	}

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

	emailUC := authUC.NewEmailUsecase(
		userRepo,
		emailRepo,
		tokenHasher,
		mailer,
		accessExp,
	)

	passwordUC := authUC.NewPasswordUsecase(
		userRepo,
		passwordResetTokenRepo,
		refreshTokenRepo,
		sessionRepo,
		passwordHasher,
		mailer,
		accessExp, // harusnya bukan access exp
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
			Validator: validator.New(),
			Config:    cfg,

			EmailUC:    emailUC,
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
