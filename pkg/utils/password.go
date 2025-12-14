package password

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

/* ============================================================
   CONFIG
============================================================ */

type Config struct {
	Memory                uint32
	Iterations            uint32
	Parallelism           uint8
	SaltLength            uint32
	KeyLength             uint32
	Peppers               map[int]string
	CurrentPepperVersion  int
}

func DefaultConfig() *Config {
	return &Config{
		Memory:       64 * 1024, // 64 MB
		Iterations:   3,
		Parallelism:  2,
		SaltLength:   16,
		KeyLength:    32,
		Peppers: map[int]string{
			1: "pepper-v1-secret",
		},
		CurrentPepperVersion: 1,
	}
}

/* ============================================================
   UTIL
============================================================ */

func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func randomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func deriveInput(password []byte, pepper string) []byte {
	if pepper == "" {
		out := make([]byte, len(password))
		copy(out, password)
		return out
	}

	mac := hmac.New(sha256.New, []byte(pepper))
	mac.Write(password)
	return mac.Sum(nil)
}

/* ============================================================
   HASH PARSER
============================================================ */

var ErrInvalidHash = errors.New("invalid argon2id hash format")

type parsedHash struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	PepperVer   *int
	Salt        []byte
	Hash        []byte
}

// Format:
// $argon2id$v=19$m=...,t=...,p=...$pepver=<n>$<salt_b64>$<hash_b64>
func parseHash(encoded string) (*parsedHash, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 && len(parts) != 7 {
		return nil, ErrInvalidHash
	}

	if parts[1] != "argon2id" || parts[2] != "v=19" {
		return nil, ErrInvalidHash
	}

	var mem uint32
	var it uint32
	var par uint8

	for _, p := range strings.Split(parts[3], ",") {
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			return nil, ErrInvalidHash
		}

		switch kv[0] {
		case "m":
			v, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil { return nil, err }
			mem = uint32(v)
		case "t":
			v, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil { return nil, err }
			it = uint32(v)
		case "p":
			v, err := strconv.ParseUint(kv[1], 10, 8)
			if err != nil { return nil, err }
			par = uint8(v)
		}
	}

	idx := 4
	var pepver *int

	if len(parts) == 7 && strings.HasPrefix(parts[4], "pepver=") {
		v, err := strconv.Atoi(strings.TrimPrefix(parts[4], "pepver="))
		if err != nil {
			return nil, ErrInvalidHash
		}
		pepver = &v
		idx = 5
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[idx])
	if err != nil {
		return nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[idx+1])
	if err != nil {
		return nil, err
	}

	return &parsedHash{
		Memory:      mem,
		Iterations:  it,
		Parallelism: par,
		PepperVer:   pepver,
		Salt:        salt,
		Hash:        hash,
	}, nil
}

/* ============================================================
   PUBLIC API
============================================================ */

// HashPassword creates argon2id hash with pepper version embedded
func HashPassword(password []byte, cfg *Config) (string, int, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	salt, err := randomBytes(cfg.SaltLength)
	if err != nil {
		return "", 0, err
	}

	pepver := cfg.CurrentPepperVersion
	pepper := cfg.Peppers[pepver]

	input := deriveInput(password, pepper)
	defer zeroBytes(input)

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

	zeroBytes(hash)
	zeroBytes(salt)

	return encoded, pepver, nil
}

// VerifyPassword verifies hash and signals if rehash is recommended
func VerifyPassword(password []byte, encoded string, cfg *Config) (bool, bool, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	parsed, err := parseHash(encoded)
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

	input := deriveInput(password, pepper)
	defer zeroBytes(input)

	computed := argon2.IDKey(
		input,
		parsed.Salt,
		parsed.Iterations,
		parsed.Memory,
		parsed.Parallelism,
		uint32(len(parsed.Hash)),
	)
	defer zeroBytes(computed)

	match := subtle.ConstantTimeCompare(computed, parsed.Hash) == 1
	shouldRehash := match && cfg.CurrentPepperVersion > usedPepper

	return match, shouldRehash, nil
}



// example implementasion

// package main

// import (
// 	"fmt"
// 	"log"

// 	"example/password"
// )

// func main() {
// 	cfg := password.DefaultConfig()

// 	plain := []byte("SuperSecretPassword123!")

// 	// HASH
// 	hash, pepver, err := password.HashPassword(plain, cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("HASHED PASSWORD:")
// 	fmt.Println(hash)
// 	fmt.Println("Pepper version:", pepver)

// 	// VERIFY (correct)
// 	ok, rehash, err := password.VerifyPassword([]byte("SuperSecretPassword123!"), hash, cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("\nVERIFY CORRECT PASSWORD:")
// 	fmt.Println("Match:", ok)
// 	fmt.Println("Should Rehash:", rehash)

// 	// VERIFY (wrong)
// 	ok, rehash, _ = password.VerifyPassword([]byte("WrongPassword"), hash, cfg)

// 	fmt.Println("\nVERIFY WRONG PASSWORD:")
// 	fmt.Println("Match:", ok)
// 	fmt.Println("Should Rehash:", rehash)

// 	// SIMULATE PEPPER ROTATION
// 	cfg.Peppers[2] = "pepper-v2-rotated"
// 	cfg.CurrentPepperVersion = 2

// 	ok, rehash, _ = password.VerifyPassword(plain, hash, cfg)

// 	fmt.Println("\nAFTER PEPPER ROTATION:")
// 	fmt.Println("Match:", ok)
// 	fmt.Println("Should Rehash:", rehash)
// }
