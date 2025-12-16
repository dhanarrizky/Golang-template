package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/dhanarrizky/Golang-template/internal/config"
	"github.com/dhanarrizky/Golang-template/internal/delivery/http"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/cache"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres"
	authinfra "github.com/dhanarrizky/Golang-template/internal/infrastructure/auth"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"
	"github.com/dhanarrizky/Golang-template/internal/usecase/user"
)

func RunHTTPServer() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// =====================
	// Infrastructure
	// =====================
	db, err := postgres.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}

	var redisClient *redis.Client
	if cfg.RedisEnable {
		redisClient, err = cache.NewRedisClient(cfg.RedisHost, cfg.RedisPassword)
		if err != nil {
			return err
		}
	}

	signer := authinfra.NewJWTSigner(cfg)

	// =====================
	// Repository
	// =====================
	userRepo := postgres.NewUserRepository(db)

	// =====================
	// Usecase
	// =====================
	getUserUC := user.NewGetUserUsecase(userRepo)
	createUserUC := user.NewCreateUserUsecase(userRepo)
	deleteUserUC := user.NewDeleteUserUsecase(userRepo)

	authUC := auth.NewUserAuthUsecase(
		userRepo,
		signer,
	)

	// =====================
	// HTTP
	// =====================
	router := gin.New()

	http.RegisterRoutes(
		router,
		http.RouteDeps{
			GetUserUC:    getUserUC,
			CreateUserUC: createUserUC,
			DeleteUserUC: deleteUserUC,
			AuthUC:       authUC,
			Redis:        redisClient,
			Config:       *cfg,
		},
	)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Println("HTTP server running on", cfg.AppPort)
		_ = srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
