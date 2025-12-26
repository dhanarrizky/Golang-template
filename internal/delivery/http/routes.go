package http

import (
	"github.com/gin-gonic/gin"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/auth"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/roles"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/users"
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

	tokenHandler := auth.NewTokenHandler(
		d.TokenUC,
		d.Config.JWTSecret,
	)

	passwordHandler := auth.NewPasswordHandler(
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

	// =====================================================
	// PUBLIC ROUTES
	// =====================================================
	public := r.Group("/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", tokenHandler.Refresh)

		// register
		public.POST("/users", userHandler.Create)
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
		protected.DELETE("/users/me", userHandler.Delete)
	}

	admin := r.Group("/v1")
	admin.Use(
		middleware.AuthMiddleware(*d.JwtSigner),
		middleware.RequireRole("admin"),
	)
	{
		// user management
		admin.GET("/users", userHandler.List) // optional
		admin.GET("/users/:id", userHandler.GetByID)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id/permanent", userHandler.PermanentDelete)

		// role
		admin.GET("/roles", roleHandler.List)
		admin.POST("/roles", roleHandler.Create)
		admin.PUT("/roles/:id", roleHandler.Update)
		admin.DELETE("/roles/:id", roleHandler.Delete)
	}

}
