package postgres

import (
	authEntities "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&authEntities.User{},
		&authEntities.Role{},
		&authEntities.EmailVerificationToken{},
		&authEntities.PasswordResetToken{},
		&authEntities.RefreshTokenFamily{},
		&authEntities.RefreshToken{},
		&authEntities.UserSession{},
		&authEntities.LoginAttempt{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
