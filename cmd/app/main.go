// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// 	"github.com/dhanarrizky/Golang-template/internal/config"
// 	httpdelivery "github.com/dhanarrizky/Golang-template/internal/delivery"
// 	"github.com/dhanarrizky/Golang-template/internal/infrastructure/cache"
// 	"github.com/dhanarrizky/Golang-template/internal/usecase"
// 	"github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres"
// 	"github.com/dhanarrizky/Golang-template/internal/usecase/user"
	
// 	"github.com/redis/go-redis/v9"
// )

// func main() {
// 	_ = godotenv.Load()

// 	cfg, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	if cfg.IsProduction() {
// 		gin.SetMode(gin.ReleaseMode)
// 	}

// 	log.Println("Starting application...")

// 	// === Init Database ===
// 	db, err := postgres.NewPostgresDB(cfg.DatabaseURL)
// 	if err != nil {
// 		log.Fatalf("DB error: %v", err)
// 	}

// 	// Connection pool
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		log.Fatalf("DB raw error: %v", err)
// 	}

// 	sqlDB.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)
// 	sqlDB.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)

// 	lifetime, err := time.ParseDuration(cfg.DatabaseConnMaxLifetime)
// 	if err != nil {
// 		lifetime = 30 * time.Minute
// 	}
// 	sqlDB.SetConnMaxLifetime(lifetime)

// 	// === Init Redis ===
// 	var redisClient *redis.Client

	
// 	if cfg.RedisEnable {
// 		redisClient, err = cache.NewRedisClient(cfg.RedisHost, cfg.RedisPassword)
// 		if err != nil {
// 			log.Fatalf("Redis error: %v", err)
// 		}
// 	}


// 	// === Repository Layer ===
// 	userRepo := postgres.NewUserRepository(db)

// 	// === Usecase Layer ===
// 	getUserUC := user.NewGetUserUsecase(userRepo)
// 	createUserUC := user.NewCreateUserUsecase(userRepo)
// 	deleteUserUC := user.NewDeleteUserUsecase(userRepo)

// 	// === HTTP Server ===
// 	router := gin.New()
// 	router.Use(gin.Recovery())

// 	authSvc := usecase.NewUserAuthUsecase(userRepo)

// 	httpdelivery.RegisterRoutes(
// 		router,
// 		getUserUC,
// 		createUserUC,
// 		deleteUserUC,
// 		authSvc,      // ← WAJIB ditambahkan
// 		redisClient, // ← pastikan redis/v9
// 		*cfg,
// 	)

// 	srv := &http.Server{
// 		Addr:         ":" + cfg.AppPort,
// 		Handler:      router,
// 		ReadTimeout:  15 * time.Second,
// 		WriteTimeout: 15 * time.Second,
// 		IdleTimeout:  60 * time.Second,
// 	}

// 	go func() {
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Server error: %v", err)
// 		}
// 	}()

// 	log.Printf("Server running on port %s", cfg.AppPort)

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	log.Println("Shutting down gracefully...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Fatalf("Forced shutdown: %v", err)
// 	}

// 	log.Println("Server stopped.")
// }


package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/dhanarrizky/Golang-template/internal/bootstrap"
)

func main() {
	_ = godotenv.Load()

	if err := bootstrap.RunHTTPServer(); err != nil {
		log.Fatalf("app stopped: %v", err)
	}
}
