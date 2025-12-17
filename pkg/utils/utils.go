package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func RandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func DeriveInput(password []byte, pepper string) []byte {
	if pepper == "" {
		out := make([]byte, len(password))
		copy(out, password)
		return out
	}

	mac := hmac.New(sha256.New, []byte(pepper))
	mac.Write(password)
	return mac.Sum(nil)
}

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
func ParseHash(encoded string) (*parsedHash, error) {
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
			if err != nil {
				return nil, err
			}
			mem = uint32(v)
		case "t":
			v, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return nil, err
			}
			it = uint32(v)
		case "p":
			v, err := strconv.ParseUint(kv[1], 10, 8)
			if err != nil {
				return nil, err
			}
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
