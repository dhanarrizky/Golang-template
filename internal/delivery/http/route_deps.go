package http

import (
	"github.com/go-playground/validator/v10"

	"github.com/dhanarrizky/Golang-template/internal/config"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	authUC "github.com/dhanarrizky/Golang-template/internal/usecase/auth"
	emailUC "github.com/dhanarrizky/Golang-template/internal/usecase/email"
	roleUC "github.com/dhanarrizky/Golang-template/internal/usecase/roles"
	userUC "github.com/dhanarrizky/Golang-template/internal/usecase/user"
)

// type RouteDeps struct {
// 	JwtSigner *ports.TokenSigner
// 	Validator *validator.Validate
// 	Config    *config.Config

// 	// Auth usecases
// 	LoginUC    authUC.LoginUsecase
// 	PasswordUC authUC.PasswordUsecase
// 	SessionUC  authUC.SessionUsecase
// 	TokenUC    authUC.TokenUsecase
// 	RoleUC     roleUC.RoleUsecase
// 	UserUC     userUC.UserUsecase
// }

type RouteDeps struct {
	JwtSigner *ports.TokenSigner  // JWT signer
	Validator *validator.Validate // Validator (e.g., go-playground/validator)
	Config    *config.Config      // Config struct dengan JWTSecret, CORSAllowedOrigins, IsDevelopment()
	// EmailSender emailUC.OTPUsecase  // Tambahan: Interface untuk send email (e.g., gomail)

	LoginUC    authUC.LoginUsecase    // UseCase untuk login
	TokenUC    authUC.TokenUsecase    // UseCase untuk token
	PasswordUC authUC.PasswordUsecase // UseCase untuk password
	UserUC     userUC.UserUsecase     // UseCase untuk user
	RoleUC     roleUC.RoleUsecase     // UseCase untuk role
	OTPUC      emailUC.OTPUsecase     // Tambahan: UseCase untuk OTP (generate, verify, resend)
	// ForgotPasswordUC                        // Tambahan: UseCase untuk forgot password (send OTP, reset)
	SessionUC authUC.SessionUsecase // Tambahan: UseCase untuk session management
	// Tambah lain jika perlu, seperti RateLimiter untuk OTP/resend
}
