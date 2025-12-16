package security

import (
	"github.com/dhanarrizky/Golang-template/internal/ports"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
)

type Argon2PasswordHasher struct {
	cfg *utils.Config
}

func NewArgon2PasswordHasher(cfg *utils.Config) ports.PasswordHasher {
	if cfg == nil {
		cfg = password.DefaultConfig()
	}
	return &Argon2PasswordHasher{cfg: cfg}
}

// =====================
// IMPLEMENT PORT
// =====================

func (h *Argon2PasswordHasher) Hash(plain string) (string, error) {
	hash, _, err := password.HashPassword([]byte(plain), h.cfg)
	return hash, err
}

func (h *Argon2PasswordHasher) Compare(plain, hashed string) (bool, error) {
	match, _, err := password.VerifyPassword([]byte(plain), hashed, h.cfg)
	return match, err
}
