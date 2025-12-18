package http

import (
	"github.com/gin-gonic/gin"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/auth"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/middleware"
)

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

	// tokenHandler := auth.NewEmailHandler(
	// 	d.EmailUC,
	// 	d.Validator,
	// )

	// passwordHandler := auth.NewPasswordHandler(
	// 	d.PasswordUC,
	// 	d.Validator,
	// )

	// emailHandler := auth.NewEmailHandler(
	// 	d.EmailUC,
	// 	d.Validator,
	// )

	// sessionHandler := sessions.NewSessionHandler(
	// 	d.SessionUC,
	// )

	// userHandler := users.NewUserHandler(
	// 	d.UserUC,
	// 	d.Validator,
	// )

	// roleHandler := roles.NewRoleHandler(
	// 	d.RoleUC,
	// 	d.Validator,
	// )

	// =====================================================
	// PUBLIC ROUTES
	// =====================================================
	public := r.Group("/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		// public.POST("/auth/refresh", tokenHandler.Refresh)

		// public.POST("/auth/password/forgot", passwordHandler.Forgot)
		// public.POST("/auth/password/reset", passwordHandler.Reset)

		// public.POST("/auth/email/verify", emailHandler.Verify)
		// public.POST("/auth/email/resend", emailHandler.Resend)
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

		// // password
		// protected.POST("/auth/password/change", passwordHandler.Change)

		// // sessions
		// protected.GET("/sessions", sessionHandler.List)
		// protected.DELETE("/sessions/:id", sessionHandler.Revoke)

		// // user
		// protected.GET("/users/me", userHandler.Me)
		// protected.PUT("/users/me", userHandler.Update)
		// protected.DELETE("/users/me", userHandler.Delete)

		// // role (biasanya admin only)
		// protected.GET("/roles", roleHandler.List)
		// protected.POST("/roles", roleHandler.Create)
		// protected.PUT("/roles/:id", roleHandler.Update)
		// protected.POST("/roles/assign", roleHandler.Assign)
	}
}
