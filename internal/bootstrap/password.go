package bootstrap

import (
	"github.com/dhanarrizky/Golang-template/internal/infrastructure/security"
	"github.com/dhanarrizky/Golang-template/pkg/utils/password"
)

func InitPasswordHasher(cfg *config.Config) ports.PasswordHasher {
	return security.NewArgon2PasswordHasher(cfg)
}

func InitTokenHasher(cfg *config.Config) ports.TokenHasher {
	secret := cof.SecretToken
	if secret == "" {
		log.Fatal("EMAIL_VERIFICATION_TOKEN_SECRET is not set")
	}

	return security.NewHMACTokenHasher(secret)
}