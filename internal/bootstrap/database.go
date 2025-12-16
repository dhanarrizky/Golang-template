package bootstrap

import (
	"log"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/config"
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/database/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfg *config.Config) *gorm.DB {
	db, err := postgres.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql DB:", err)
	}

	sqlDB.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)

	lifetime, err := time.ParseDuration(cfg.DatabaseConnMaxLifetime)
	if err != nil {
		lifetime = 30 * time.Minute
	}
	sqlDB.SetConnMaxLifetime(lifetime)

	return db
}
