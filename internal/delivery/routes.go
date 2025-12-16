// package http

// import (
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/redis/go-redis/v9"

// 	"github.com/dhanarrizky/Golang-template/internal/config"
// 	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers"
// 	"github.com/dhanarrizky/Golang-template/internal/delivery/http/middleware"
// 	"github.com/dhanarrizky/Golang-template/internal/usecase/user"
// 	auth "github.com/dhanarrizky/Golang-template/pkg/auth"
// )


// func RegisterRoutes(
// 	r *gin.Engine,
// 	getUC *user.GetUserUsecase,
// 	createUC *user.CreateUserUsecase,
// 	deleteUC *user.DeleteUserUsecase,
// 	authUC handlers.UserAuthService,              // ⭐ Auth with Access+Refresh Token
// 	redisClient *redis.Client,             // 🟢 optional
// 	cfg config.Config,
// ) {

// 	// ============================================================
// 	// GLOBAL MIDDLEWARE
// 	// ============================================================

// 	r.Use(middleware.RecoveryWithLogging())

// 	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

// 	if cfg.IsDevelopment() {
// 		r.Use(middleware.LoggingMiddleware())
// 	}

// 	// ============================================================
// 	// PUBLIC ROUTES (NO AUTH)
// 	// ============================================================

// 	signer := auth.NewSigner(
// 		cfg.JWTSecret,
// 		"github.com/dhanarrizky/Golang-template",        // issuer
// 		"github.com/dhanarrizky/Golang-template_client", // audience
// 		time.Hour,      // access TTL
// 		24*time.Hour,   // refresh TTL
// 	)

// 	// authHandlers := &handlers.AuthHandlers{
// 	// 	Signer:      signer,
// 	// 	RedisClient: redisClient,
// 	// 	Config:      &cfg,
// 	// 	UserSvc:     authUC.UserService, // pastikan ada interface UserAuthService
// 	// }
// 	authHandlers := &handlers.AuthHandlers{
// 		UserSvc:     authUC,       // langsung implementasi UserAuthService
// 		RedisClient: redisClient,
// 		Signer:      signer,
// 		Config:      &cfg,
// 	}

// 	public := r.Group("/v1")
// 	{
// 		public.GET("/health", authHandlers.HealthHandler) // buat HealthHandler method jika perlu

// 		public.POST("/auth/login", authHandlers.LoginHandler())
// 		public.POST("/auth/refresh", authHandlers.RefreshHandler())
// 	}

// 	protected := r.Group("/v1")

// 	// rate limit
// 	if cfg.IsProduction() && redisClient != nil && cfg.RateLimiterEnable {
// 		protected.Use(middleware.RateLimitMiddleware(
// 			60,
// 			time.Minute,
// 			func(c *gin.Context) string {
// 				if apiKey := c.GetHeader("X-Api-Key"); apiKey != "" {
// 					return apiKey
// 				}
// 				return c.ClientIP()
// 			},
// 		))
// 	}

// 	// JWT auth middleware
// 	protected.Use(middleware.AuthMiddleware(signer))

// 	// user routes
// 	protected.GET("/users/:id", handlers.GetUserHandler(getUC))
// 	protected.POST("/users", handlers.CreateUserHandler(createUC))
// 	protected.DELETE("/users/:id", handlers.DeleteUserHandler(deleteUC))

// }


package http

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/auth"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/roles"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers/sessions"
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
		d.Config.JWTSecret,
		time.Minute*time.Duration(d.Config.JWTAccessTTL),
	)

	tokenHandler := auth.NewTokenHandler(
		d.TokenUC,
		d.Config.JWTSecret,
	)

	passwordHandler := auth.NewPasswordHandler(
		d.PasswordUC,
		d.Validator,
	)

	emailHandler := auth.NewEmailHandler(
		d.EmailUC,
		d.Validator,
	)

	sessionHandler := sessions.NewSessionHandler(
		d.SessionUC,
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

		public.POST("/auth/password/forgot", passwordHandler.Forgot)
		public.POST("/auth/password/reset", passwordHandler.Reset)

		public.POST("/auth/email/verify", emailHandler.Verify)
		public.POST("/auth/email/resend", emailHandler.Resend)
	}

	// =====================================================
	// PROTECTED ROUTES
	// =====================================================
	protected := r.Group("/v1")
	protected.Use(middleware.JWTAuthMiddleware(d.Config.JWTSecret))

	{
		// auth
		protected.POST("/auth/logout", authHandler.Logout)
		protected.POST("/auth/logout-all", authHandler.LogoutAll)
		protected.GET("/auth/me", authHandler.Me)

		// password
		protected.POST("/auth/password/change", passwordHandler.Change)

		// sessions
		protected.GET("/sessions", sessionHandler.List)
		protected.DELETE("/sessions/:id", sessionHandler.Revoke)

		// user
		protected.GET("/users/me", userHandler.Me)
		protected.PUT("/users/me", userHandler.Update)
		protected.DELETE("/users/me", userHandler.Delete)

		// role (biasanya admin only)
		protected.GET("/roles", roleHandler.List)
		protected.POST("/roles", roleHandler.Create)
		protected.PUT("/roles/:id", roleHandler.Update)
		protected.POST("/roles/assign", roleHandler.Assign)
	}
}
