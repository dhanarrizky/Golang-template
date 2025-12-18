package bootstrap

import (
	"log"

	"github.com/dhanarrizky/Golang-template/internal/config"

	"github.com/dhanarrizky/Golang-template/internal/infrastructure/security"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/others"
)

// func InitPasswordHasher(cfg *config.Config) ports.PasswordHasher {
// 	return security.NewArgon2PasswordHasher(cfg)
// }

func InitTokenHasher(cfg *config.Config) ports.TokenHasher {
	secret := cfg.SecretToken
	if secret == "" {
		log.Fatal("EMAIL_VERIFICATION_TOKEN_SECRET is not set")
	}

	return security.NewHMACTokenHasher(secret)
}
