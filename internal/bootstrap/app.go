package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/dhanarrizky/Golang-template/internal/config"

	authinfra "github.com/dhanarrizky/Golang-template/internal/infrastructure/auth"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"
	"github.com/dhanarrizky/Golang-template/internal/usecase/user"
)

func InitHTTPApp(cfg *config.Config) *gin.Engine {
	// =====================
	// Infrastructure
	// =====================
	db := InitDatabase(cfg)
	redis := InitRedis(cfg)
	passwordHasher := InitPasswordHasher(cfg)
	jwtSigner := authinfra.NewJWTSigner(cfg)

	// =====================
	// Repository
	// =====================

	// Auth
	emailRepo := postgres.NewEmailVerificationTokenRepository(db)
	loginAttemptRepo := postgres.NewLoginAttemptRepository(db)
	passwordResetTokenRepo := postgres.NewPasswordResetTokenRepository(db)
	refreshTokenFamilyRepo := postgres.NewRefreshTokenFamilyRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewUserSessionRepository(db)

	// =====================
	// Mailer
	// =====================
	mailer, err := mailerinfra.NewSMTPMailer(
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

	// loginUC := auth.NewLoginUsecase(
	// 	userRepo,
	// 	sessionRepo,
	// 	jwtSigner,
	// 	redis,
	// )

	// tokenUC := auth.NewTokenUsecase(
	// 	sessionRepo,
	// 	jwtSigner,
	// )

	// userUC := user.NewUserUsecase(userRepo)

	// =====================
	// HTTP Router
	// =====================
	router := gin.New()

	http.RegisterRoutes(
		router,
		http.RouteDeps{
			Validator: validator.New(),
			Config:    http.ConfigFrom(cfg),

			LoginUC: loginUC,
			TokenUC: tokenUC,
			UserUC:  userUC,
		},
	)

	return router
}
