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

func InitPublicIdCodec(cfg *config.Config) ports.PublicIDCodec {
	secret := cfg.PublicIdAesKey
	if secret == "" {
		log.Fatal("PUBLIC_ID_AES_KEY is not set")
	}

	codec, err := security.NewPublicIDCodecFromBase64(secret)
	if err != nil {
		log.Fatal(err)
	}
	return codec
}

func InitTokenVerifier(cfg *config.Config) ports.TokenVerifier {
	secret := cfg.TokenHmacSecret
	if secret == "" {
		log.Fatal("TOKEN_HMAC_SECRET is not set")
	}

	return security.NewHMACTokenVerifier(secret)
}

func InitTokenGenerator(cfg *config.Config) ports.TokenGenerator {
	verifier := InitTokenVerifier(cfg)
	return security.NewSecureTokenGenerator(verifier)
}
