package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/dhanarrizky/Golang-template/internal/config"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/handlers"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http/middleware"
	"github.com/dhanarrizky/Golang-template/internal/usecase/user"
	auth "github.com/dhanarrizky/Golang-template/pkg/auth"
)


func RegisterRoutes(
	r *gin.Engine,
	getUC *user.GetUserUsecase,
	createUC *user.CreateUserUsecase,
	deleteUC *user.DeleteUserUsecase,
	authUC handlers.UserAuthService,              // ⭐ Auth with Access+Refresh Token
	redisClient *redis.Client,             // 🟢 optional
	cfg config.Config,
) {

	// ============================================================
	// GLOBAL MIDDLEWARE
	// ============================================================

	r.Use(middleware.RecoveryWithLogging())

	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	if cfg.IsDevelopment() {
		r.Use(middleware.LoggingMiddleware())
	}

	// ============================================================
	// PUBLIC ROUTES (NO AUTH)
	// ============================================================

	signer := auth.NewSigner(
		cfg.JWTSecret,
		"github.com/dhanarrizky/Golang-template",        // issuer
		"github.com/dhanarrizky/Golang-template_client", // audience
		time.Hour,      // access TTL
		24*time.Hour,   // refresh TTL
	)

	// authHandlers := &handlers.AuthHandlers{
	// 	Signer:      signer,
	// 	RedisClient: redisClient,
	// 	Config:      &cfg,
	// 	UserSvc:     authUC.UserService, // pastikan ada interface UserAuthService
	// }
	authHandlers := &handlers.AuthHandlers{
		UserSvc:     authUC,       // langsung implementasi UserAuthService
		RedisClient: redisClient,
		Signer:      signer,
		Config:      &cfg,
	}

	public := r.Group("/v1")
	{
		public.GET("/health", authHandlers.HealthHandler) // buat HealthHandler method jika perlu

		public.POST("/auth/login", authHandlers.LoginHandler())
		public.POST("/auth/refresh", authHandlers.RefreshHandler())
	}

	protected := r.Group("/v1")

	// rate limit
	if cfg.IsProduction() && redisClient != nil && cfg.RateLimiterEnable {
		protected.Use(middleware.RateLimitMiddleware(
			60,
			time.Minute,
			func(c *gin.Context) string {
				if apiKey := c.GetHeader("X-Api-Key"); apiKey != "" {
					return apiKey
				}
				return c.ClientIP()
			},
		))
	}

	// JWT auth middleware
	protected.Use(middleware.AuthMiddleware(signer))

	// user routes
	protected.GET("/users/:id", handlers.GetUserHandler(getUC))
	protected.POST("/users", handlers.CreateUserHandler(createUC))
	protected.DELETE("/users/:id", handlers.DeleteUserHandler(deleteUC))

}
