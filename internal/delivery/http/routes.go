package http

import (
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/auth"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/otp"      // Tambahan untuk OTP handler
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/password" // Adjust untuk forgot password
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/roles"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/session" // Tambahan untuk session handler
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/users"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

// RouteDeps adalah struct untuk dependencies (asumsikan ini sudah ada di package Anda)
// Tambah fields baru untuk OTP, Forgot Password, dan Session jika belum ada
// type RouteDeps struct {
// 	LoginUC          // UseCase untuk login
// 	TokenUC          // UseCase untuk token
// 	PasswordUC       // UseCase untuk password
// 	UserUC           // UseCase untuk user
// 	RoleUC           // UseCase untuk role
// 	Validator        // Validator (e.g., go-playground/validator)
// 	Config           // Config struct dengan JWTSecret, CORSAllowedOrigins, IsDevelopment()
// 	JwtSigner        // JWT signer
// 	OTPUC            // Tambahan: UseCase untuk OTP (generate, verify, resend)
// 	ForgotPasswordUC // Tambahan: UseCase untuk forgot password (send OTP, reset)
// 	EmailSender      // Tambahan: Interface untuk send email (e.g., gomail)
// 	SessionUC        // Tambahan: UseCase untuk session management
// 	// Tambah lain jika perlu, seperti RateLimiter untuk OTP/resend
// }

func RegisterRoutes(r *gin.Engine, d RouteDeps) {
	// =====================================================
	// GLOBAL MIDDLEWARE
	// =====================================================
	r.Use(middleware.RecoveryWithLogging())
	r.Use(middleware.CORSMiddleware(d.Config.CORSAllowedOrigins))
	if d.Config.IsDevelopment() {
		r.Use(middleware.LoggingMiddleware())
	}

	// =====================================================
	// HANDLERS (CONSTRUCTION)
	// =====================================================
	authHandler := auth.NewAuthHandler(
		d.LoginUC,
		d.Validator,
	)
	tokenHandler := auth.NewTokenHandler(
		d.TokenUC,
		d.Config.JWTSecret,
	)
	passwordHandler := auth.NewPasswordHandler( // Ini untuk change password
		d.PasswordUC,
		d.Validator,
	)
	userHandler := users.NewUserHandler(
		d.UserUC,
		d.Validator,
	)
	roleHandler := roles.NewRoleHandler(
		d.RoleUC,
		d.Validator,
	)

	// Tambahan Handlers
	otpHandler := otp.NewOTPHandler( // Asumsikan package otp dengan NewOTPHandler
		d.OTPUC,
		d.Validator,
		d.EmailSender, // Untuk send OTP via email
	)
	forgotHandler := password.NewForgotPasswordHandler( // Asumsikan di package password atau auth
		d.ForgotPasswordUC,
		d.Validator,
		d.EmailSender,
	)
	sessionHandler := session.NewSessionHandler(d.SessionUC) // Tambahan dari session management

	// =====================================================
	// PUBLIC ROUTES
	// =====================================================
	public := r.Group("/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", tokenHandler.Refresh)
		// register
		public.POST("/users", userHandler.Create) // Setelah create, trigger send OTP di use case

		// Tambahan untuk OTP dan Forgot Password
		public.POST("/auth/verify-otp", otpHandler.Verify)          // Verify OTP untuk aktivasi
		public.POST("/auth/resend-otp", otpHandler.Resend)          // Resend OTP (rate limited di middleware jika perlu)
		public.POST("/auth/forgot-password", forgotHandler.Request) // Send OTP untuk reset
		public.POST("/auth/reset-password", forgotHandler.Reset)    // Reset password setelah verify OTP
	}

	// =====================================================
	// PROTECTED ROUTES
	// =====================================================
	protected := r.Group("/v1")
	protected.Use(middleware.AuthMiddleware(*d.JwtSigner))
	{
		// auth
		protected.POST("/auth/logout", authHandler.Logout)
		protected.POST("/auth/logout-all", authHandler.LogoutAll)
		protected.GET("/auth/me", authHandler.Me)
		// password
		protected.POST("/auth/password/change", passwordHandler.Change)
		// user (self)
		protected.GET("/users/me", userHandler.Me)
		protected.PUT("/users/me", userHandler.Update)
		protected.DELETE("/users/me", userHandler.Delete) // Soft delete

		// Tambahan untuk verify email jika change email
		protected.POST("/users/me/verify-email", userHandler.VerifyEmail) // Asumsikan method baru di UserHandler

		// Tambahan untuk session management
		protected.GET("/auth/sessions", sessionHandler.List)
		protected.DELETE("/auth/sessions/:session_id", sessionHandler.Revoke)
		protected.DELETE("/auth/sessions/all", sessionHandler.RevokeAll) // Optional, kalau ingin pakai session juga
	}

	// =====================================================
	// ADMIN ROUTES
	// =====================================================
	admin := r.Group("/v1")
	admin.Use(
		middleware.AuthMiddleware(*d.JwtSigner),
		middleware.RequireRole("admin"),
	)
	{
		// user management
		admin.GET("/users", userHandler.List) // optional, dengan pagination
		admin.GET("/users/:id", userHandler.GetByID)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id/permanent", userHandler.PermanentDelete)

		// Tambahan untuk assign role (jika belum include di update)
		admin.POST("/users/:id/assign-role", roleHandler.AssignRole) // Method baru

		// role
		admin.GET("/roles", roleHandler.List)
		admin.POST("/roles", roleHandler.Create)
		admin.PUT("/roles/:id", roleHandler.Update)
		admin.DELETE("/roles/:id", roleHandler.Delete)
	}

	// =====================================================
	// OPTIONAL ROUTES UNTUK TEMPLATE LEBIH LENGKAP
	// =====================================================
	// Health check (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
