package postgres

import (
	authModels "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&authModels.User{},
		&authModels.Role{},
		// &authModels.PasswordResetToken{},
		&authModels.RefreshTokenFamily{},
		&authModels.RefreshToken{},
		&authModels.UserSession{},
		&authModels.LoginAttempt{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
