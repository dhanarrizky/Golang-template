package postgres

import (
	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Auto migrate
	db.AutoMigrate(&entities.User{})
	return db, nil
}