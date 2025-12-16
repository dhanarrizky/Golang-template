package bootstrap

import (
	"github.com/gin-gonic/gin"

	http "github.com/dhanarrizky/Golang-template/internal/delivery/http"
	"github.com/dhanarrizky/Golang-template/internal/config"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres"
	authRepo "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres/auth"
)

func InitHTTPApp(cfg *config.Config) *gin.Engine {
	// =====================
	// Infrastructure
	// =====================
	db := InitDatabase(cfg)
	redis := InitRedis(cfg)
	passwordHasher := InitTokenHasher(cfg)
	// jwtSigner := authinfra.NewJWTSigner(cfg)

	// =====================
	// Repository
	// =====================

	// Auth
	emailRepo := authRepo.NewEmailVerificationTokenRepository(db)
	loginAttemptRepo := authRepo.NewLoginAttemptRepository(db)
	passwordResetTokenRepo := authRepo.NewPasswordResetTokenRepository(db)
	refreshTokenFamilyRepo := authRepo.NewRefreshTokenFamilyRepository(db)
	refreshTokenRepo := authRepo.NewRefreshTokenRepository(db)
	roleRepo := authRepo.NewRoleRepository(db)
	userRepo := authRepo.NewUserRepository(db)
	sessionRepo := authRepo.NewUserSessionRepository(db)

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
