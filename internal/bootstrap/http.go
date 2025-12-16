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

	"github.com/dhanarrizky/Golang-template/internal/config"
)

func RunHTTPServer() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// =====================
	// Gin Mode
	// =====================
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := InitHTTPApp(cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Printf(
			"[BOOTSTRAP] HTTP server running on :%s | env=%s | debug=%v",
			cfg.AppPort,
			cfg.Environment,
			cfg.AppDebug,
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error:", err)
		}
	}()

	// =====================
	// Graceful Shutdown
	// =====================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[BOOTSTRAP] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
