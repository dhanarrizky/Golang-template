package security

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/dhanarrizky/Golang-template/internal/ports"
	utils "github.com/dhanarrizky/Golang-template/pkg/utils"
	"golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
	Memory               uint32
	Iterations           uint32
	Parallelism          uint8
	SaltLength           uint32
	KeyLength            uint32
	Peppers              map[int]string
	CurrentPepperVersion int
}

func DefaultConfig() *PasswordConfig {
	return &PasswordConfig{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
		Peppers: map[int]string{
			1: "pepper-v1-secret",
		},
		CurrentPepperVersion: 1,
	}
}

// argon2Hasher adalah implementasi konkret dari PasswordHasher
type argon2Hasher struct {
	config *PasswordConfig
}

// NewPasswordHasher membuat instance baru yang bisa di-inject
func NewPasswordHasher(config *PasswordConfig) ports.PasswordHasher {
	if config == nil {
		config = DefaultConfig()
	}
	return &argon2Hasher{config: config}
}

func (a *argon2Hasher) HashPassword(password []byte) (string, int, error) {
	cfg := a.config

	salt, err := utils.RandomBytes(cfg.SaltLength)
	if err != nil {
		return "", 0, err
	}

	pepver := cfg.CurrentPepperVersion
	pepper := cfg.Peppers[pepver]

	input := utils.DeriveInput(password, pepper)
	defer utils.ZeroBytes(input)

	hash := argon2.IDKey(
		input,
		salt,
		cfg.Iterations,
		cfg.Memory,
		cfg.Parallelism,
		cfg.KeyLength,
	)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$pepver=%d$%s$%s",
		cfg.Memory,
		cfg.Iterations,
		cfg.Parallelism,
		pepver,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	utils.ZeroBytes(hash)
	utils.ZeroBytes(salt)

	return encoded, pepver, nil
}

func (a *argon2Hasher) VerifyPassword(password []byte, encoded string) (bool, bool, error) {
	cfg := a.config

	parsed, err := utils.ParseHash(encoded)
	if err != nil {
		return false, false, err
	}

	usedPepper := cfg.CurrentPepperVersion
	if parsed.PepperVer != nil {
		usedPepper = *parsed.PepperVer
	}

	pepper, ok := cfg.Peppers[usedPepper]
	if !ok {
		return false, false, fmt.Errorf("pepper version %d not available", usedPepper)
	}

	input := utils.DeriveInput(password, pepper)
	defer utils.ZeroBytes(input)

	computed := argon2.IDKey(
		input,
		parsed.Salt,
		parsed.Iterations,
		parsed.Memory,
		parsed.Parallelism,
		uint32(len(parsed.Hash)),
	)
	defer utils.ZeroBytes(computed)

	match := subtle.ConstantTimeCompare(computed, parsed.Hash) == 1
	shouldRehash := match && cfg.CurrentPepperVersion > usedPepper

	return match, shouldRehash, nil
}
