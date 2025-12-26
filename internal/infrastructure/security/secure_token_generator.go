package security

import (
	"crypto/rand"
	"encoding/base64"

	ports "github.com/dhanarrizky/Golang-template/internal/ports/others"
)

type secureTokenGenerator struct {
	verifier ports.TokenVerifier
}

func NewSecureTokenGenerator(verifier ports.TokenVerifier) ports.TokenGenerator {
	return &secureTokenGenerator{
		verifier: verifier,
	}
}

func (g *secureTokenGenerator) Generate() (plain string, hash string, err error) {
	raw := make([]byte, 48) // 384 bit
	if _, err = rand.Read(raw); err != nil {
		return "", "", err
	}

	plain = base64.RawURLEncoding.EncodeToString(raw)
	hash = g.verifier.Hash(plain)

	return plain, hash, nil
}
